"""MinIO connector for AI Data Engineer platform.
Supports listing and reading user files under users/{username}/files/
"""

from __future__ import annotations

from typing import Literal
import asyncio
from contextlib import asynccontextmanager
from typing import Any, AsyncIterator, Dict, List, Optional
from urllib.parse import quote
from urllib3 import PoolManager, Timeout

from minio import Minio
from minio.error import S3Error
from pydantic import Field, field_validator

from utils.logger import AutoClassFuncLogger
from LLM.LLM_tools.connectors.base import (
    BaseConnector,
    ConnectionConfig,
    ConnectionResult,
    QueryRequest,
)


class MinIOConnectionConfig(ConnectionConfig):
    database: str = "ai-data-engineer"  # ← строка, как того требует базовый класс
    bucket_name: Literal["ai-data-engineer"] = "ai-data-engineer"
    region: str = Field(default="us-east-1")
    secure: bool = Field(default=False)

    @field_validator("database")
    @classmethod
    def validate_database(cls, v: str) -> str:
        if v != "ai-data-engineer":
            raise ValueError("MinIO connector only supports 'ai-data-engineer' as database (placeholder)")
        return v

    @field_validator("bucket_name")
    @classmethod
    def validate_bucket_name(cls, v: str) -> str:
        if v != "ai-data-engineer":
            raise ValueError("Only 'ai-data-engineer' bucket is supported")
        return v


class MinIOConnector(BaseConnector):
    DRIVER_IMPORT_ERROR = ""

    def __init__(self, logger: AutoClassFuncLogger) -> None:
        super().__init__(logger)
        self._last_config: Optional[MinIOConnectionConfig] = None

    @property
    def name(self) -> str:
        return "minio"

    async def connect(self, config: ConnectionConfig) -> Minio:
        if not isinstance(config, MinIOConnectionConfig):
            config = MinIOConnectionConfig(**config.model_dump())

        # Настройка HTTP-клиента с таймаутами
        timeout = Timeout(connect=config.timeout, read=config.timeout)  # config.timeout из вашего ConnectionConfig
        http_client = PoolManager(
            timeout=timeout,
            retries=False,  # или настройте retries по желанию
        )

        loop = asyncio.get_event_loop()
        client = await loop.run_in_executor(
            None,
            lambda: Minio(
                endpoint=f"{config.host}:{config.port}",
                access_key=config.username,
                secret_key=config.password,
                secure=config.secure,
                region=config.region,
                http_client=http_client,  # ← передаём клиент с таймаутами
            ),
        )

        try:
            exists = await loop.run_in_executor(None, client.bucket_exists, config.bucket_name)
            if not exists:
                raise RuntimeError(f"Bucket '{config.bucket_name}' does not exist")
        except S3Error as e:
            self._logger.error(f"Bucket check failed: {e}")
            raise
        except Exception as e:
            self._logger.error(f"MinIO connection failed: {e}")
            raise RuntimeError(f"MinIO connection error: {e}")

        self._last_config = config
        return client

    async def close(self, connection: Minio) -> None:
        pass

    async def fetch_server_version(self, connection: Minio) -> Optional[str]:
        return "MinIO (ai-data-engineer platform)"

    async def execute_operation(
        self,
        connection: Minio,
        request: QueryRequest,
    ) -> Dict[str, Any]:
        if self._last_config is None:
            raise RuntimeError("Config not available")

        bucket = self._last_config.bucket_name
        loop = asyncio.get_event_loop()

        if request.operation == "query":
            username = request.query
            if not username:
                raise ValueError("Username is required for 'query' operation")
            prefix = f"users/{username}/files/"
            objects = await loop.run_in_executor(
                None,
                lambda: list(connection.list_objects(bucket, prefix=prefix, recursive=False)),
            )
            files = []
            for obj in objects:
                if obj.is_dir:
                    continue
                rel_name = obj.object_name[len(prefix):] if obj.object_name.startswith(prefix) else obj.object_name
                files.append({
                    "filename": rel_name,
                    "size_bytes": obj.size,
                    "last_modified": obj.last_modified.isoformat() if obj.last_modified else None,
                })
            return {"files": files, "username": username}

        elif request.operation == "read":
            if not request.query or "/" not in request.query:
                raise ValueError("For 'read' operation, query must be 'username/filename'")

            try:
                username, filename = request.query.split("/", 1)
            except ValueError as e:
                raise ValueError("Invalid query format for 'read': expected 'username/filename'") from e

            if not filename:
                raise ValueError("Filename cannot be empty")

            object_name = f"users/{username}/files/{filename}"
            self._logger.info(f"Reading object: {object_name} from bucket {bucket}")

            try:
                response = await loop.run_in_executor(
                    None,
                    lambda: connection.get_object(bucket, object_name),
                )
                # Читаем до 10 KB (или настраиваемый лимит) для preview
                MAX_READ_SIZE = 10 * 1024  # 10 KB
                content = response.read(MAX_READ_SIZE)
                response.close()
                response.release_conn()

                # Попытка декодировать как UTF-8 (для текстовых файлов)
                try:
                    text_content = content.decode("utf-8")
                except UnicodeDecodeError:
                    # Если бинарный — возвращаем base64 или просто длину
                    text_content = f"<binary data, {len(content)} bytes>"

                return {
                    "filename": filename,
                    "username": username,
                    "content_preview": text_content,
                    "bytes_read": len(content),
                    "truncated": len(content) == MAX_READ_SIZE,
                }

            except S3Error as e:
                self._logger.error(f"Failed to read {object_name}: {e}")
                raise RuntimeError(f"File not found or access denied: {e}")

        elif request.operation == "schema":
            return {
                "bucket": bucket,
                "path_template": "users/{username}/files/{filename}",
                "supported_formats": ["csv", "json", "txt", "parquet", "xlsx"],
                "read_operation_format": "username/filename",
            }

        elif request.operation == "stats":
            username = request.query or ""
            prefix = f"users/{username}/files/" if username else "users/"
            objects = await loop.run_in_executor(
                None,
                lambda: list(connection.list_objects(bucket, prefix=prefix, recursive=True)),
            )
            file_objects = [obj for obj in objects if not obj.is_dir]
            total_size = sum(obj.size for obj in file_objects)
            return {
                "file_count": len(file_objects),
                "total_size_bytes": total_size,
                "prefix": prefix,
            }

        else:
            raise ValueError(f"Unsupported operation: {request.operation}")

    @asynccontextmanager
    async def connection_scope(self, config: ConnectionConfig) -> AsyncIterator[Minio]:
        connection = await self.connect(config)
        try:
            yield connection
        finally:
            await self.safe_close(connection)
            self._last_config = None

    async def safe_close(self, connection: Minio) -> None:
        pass