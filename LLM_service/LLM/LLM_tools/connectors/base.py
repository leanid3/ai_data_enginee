"""数据库连接器基础组件。
提供统一的数据库连接和操作接口，支持异步操作和连接池管理。
实现了连接生命周期管理和错误处理机制。
"""

from __future__ import annotations

import asyncio
from abc import ABC, abstractmethod
from contextlib import asynccontextmanager
from typing import Any, AsyncIterator, Dict, Optional

from pydantic import BaseModel, Field, ValidationInfo, field_validator

from LLM_service.utils.logger import AutoClassFuncLogger


# 数据库连接配置基类
# 定义所有数据库连接器共享的基本配置参数
# 包括：主机、端口、数据库名、用户认证等
class ConnectionConfig(BaseModel):
    """Base configuration for database connections."""

    host: str
    port: int = Field(..., ge=1, le=65535)
    database: str
    username: str
    password: str
    use_ssl: bool = False
    timeout: float = Field(10.0, gt=0)

    @field_validator("host")
    @classmethod
    def validate_host(cls, value: str) -> str:
        if not value:
            msg = "Host must not be empty"
            raise ValueError(msg)
        return value


class QueryRequest(BaseModel):
    """Representation of a query execution request."""

    operation: str = Field("test", pattern=r"^(test|query|schema|stats|read)$")
    query: Optional[str] = None

    @field_validator("query")
    @classmethod
    def validate_query(cls, query: Optional[str], info: ValidationInfo) -> Optional[str]:
        operation = info.data.get("operation")
        if operation == "query" and not query:
            msg = "Query must be provided for 'query' operation"
            raise ValueError(msg)
        return query


class ConnectionResult(BaseModel):
    """Base result returned by connectors."""

    connection_status: str
    server_version: Optional[str] = None
    error: Optional[str] = None
    data: Dict[str, Any] = Field(default_factory=dict)


# 数据库连接器抽象基类
# 定义统一的异步数据库操作接口
# 实现：
# - 连接生命周期管理
# - 查询执行
# - 错误处理
# - 资源清理
class BaseConnector(ABC):
    """Abstract asynchronous connector interface."""

    DRIVER_IMPORT_ERROR: str = ""

    def __init__(self, logger: AutoClassFuncLogger) -> None:
        self._logger = logger

    @property
    @abstractmethod
    def name(self) -> str:
        """Unique name of the connector."""

    @abstractmethod
    async def connect(self, config: ConnectionConfig) -> Any:
        """Establish a raw connection."""

    @abstractmethod
    async def close(self, connection: Any) -> None:
        """Close an open connection."""

    @abstractmethod
    async def fetch_server_version(self, connection: Any) -> Optional[str]:
        """Return server version information."""

    @abstractmethod
    async def execute_operation(
        self,
        connection: Any,
        request: QueryRequest,
    ) -> Dict[str, Any]:
        """Execute operation-specific logic."""

    async def test_connection(self, config: ConnectionConfig) -> ConnectionResult:
        """Default test implementation."""

        async with self.connection_scope(config) as conn:
            version = await self.fetch_server_version(conn)
            return ConnectionResult(connection_status="success", server_version=version)

    # 执行数据库操作的主要方法
    # 处理流程：
    # 1. 建立数据库连接
    # 2. 获取服务器版本信息
    # 3. 执行具体操作
    # 4. 处理可能的错误
    # 5. 返回执行结果
    async def run(self, config: ConnectionConfig, request: QueryRequest) -> ConnectionResult:
        """Execute connector workflow."""

        try:
            async with self.connection_scope(config) as connection:
                version = await self.fetch_server_version(connection)

                if request.operation == "test":
                    data = {}
                else:
                    data = await self.execute_operation(connection, request)

                return ConnectionResult(
                    connection_status="success",
                    server_version=version,
                    data=data,
                )
        except Exception as exc:  # pragma: no cover - defensive path
            self._logger.error(
                "Connector operation failed",
                extra={
                    "connector": self.name,
                    "error": str(exc),
                },
            )
            return ConnectionResult(connection_status="failed", error=str(exc))

    @asynccontextmanager
    async def connection_scope(self, config: ConnectionConfig) -> AsyncIterator[Any]:
        """Manage connection lifecycle."""

        connection = await self.connect(config)
        try:
            yield connection
        finally:
            await self.safe_close(connection)

    async def safe_close(self, connection: Any) -> None:
        """Close connection suppressing exceptions."""

        try:
            await self.close(connection)
        except Exception:  # pragma: no cover - best effort
            self._logger.warning(
                "Failed to close connection cleanly",
                extra={"connector": self.name},
            )

    async def run_with_timeout(
        self,
        config: ConnectionConfig,
        request: QueryRequest,
    ) -> ConnectionResult:
        """Execute run with timeout."""

        return await asyncio.wait_for(self.run(config, request), timeout=config.timeout)


