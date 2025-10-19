"""PostgreSQL数据库连接器实现。
基于asyncpg驱动的高性能异步PostgreSQL客户端。
支持连接池管理、事务处理和高级查询功能。
"""

from __future__ import annotations

import asyncio
from typing import Any, Dict, Optional

import asyncpg
from pydantic import Field

from LLM_service.utils.logger import AutoClassFuncLogger

from ..base import BaseTool, ToolOutput
from ..utils.env import apply_postgres_env_defaults
from .base import BaseConnector, ConnectionConfig, QueryRequest


class PostgresConfig(ConnectionConfig):
    """PostgreSQL-specific connection configuration."""

    port: int = Field(5432, ge=1, le=65535)


class PostgresResult(ToolOutput):
    """Tool output for PostgreSQL operations."""

    connection_status: str = "failed"
    server_version: Optional[str] = None
    query_result: Optional[Dict[str, Any]] = None


# PostgreSQL异步连接器实现
# 特点：
# - 基于asyncpg的异步操作
# - 高效的连接池管理
# - 支持prepared statements
# - 内置的类型转换系统
class AsyncpgConnector(BaseConnector):
    """Implementation of BaseConnector using asyncpg."""

    def __init__(self, logger: AutoClassFuncLogger) -> None:
        super().__init__(logger)

    @property
    def name(self) -> str:
        return "asyncpg"

    async def connect(self, config: ConnectionConfig) -> asyncpg.Connection:
        return await asyncpg.connect(
            host=config.host,
            port=config.port,
            user=config.username,
            password=config.password,
            database=config.database,
            timeout=config.timeout,
        )

    async def close(self, connection: asyncpg.Connection) -> None:
        await connection.close()

    async def fetch_server_version(self, connection: asyncpg.Connection) -> Optional[str]:
        version = connection.get_server_version()
        if version:
            return ".".join(str(part) for part in version)
        return None

    async def execute_operation(
        self,
        connection: asyncpg.Connection,
        request: QueryRequest,
    ) -> Dict[str, Any]:
        if request.operation == "schema":
            return await self._fetch_schema(connection)
        if request.operation == "stats":
            return await self._fetch_stats(connection)

        assert request.query is not None
        rows = await connection.fetch(request.query)
        return {
            "columns": [key for key in rows[0].keys()] if rows else [],
            "rows": [dict(row) for row in rows],
            "row_count": len(rows),
        }

    # 获取数据库模式信息
    # 查询information_schema视图
    # 返回：
    # - 模式名称
    # - 表名称
    # - 表类型信息
    async def _fetch_schema(self, connection: asyncpg.Connection) -> Dict[str, Any]:
        query = """
        SELECT table_schema, table_name
        FROM information_schema.tables
        WHERE table_type='BASE TABLE'
        ORDER BY table_schema, table_name;
        """
        rows = await connection.fetch(query)
        return {
            "tables": [dict(row) for row in rows],
        }

    # 获取数据库统计信息
    # 查询pg_stat_user_tables视图
    # 统计：
    # - 表行数
    # - 表大小
    # - 访问统计
    async def _fetch_stats(self, connection: asyncpg.Connection) -> Dict[str, Any]:
        query = """
        SELECT relname as table_name, n_live_tup AS row_count
        FROM pg_stat_user_tables;
        """
        rows = await connection.fetch(query)
        return {
            "stats": [dict(row) for row in rows],
        }


class ConnectToPostgresInput(PostgresConfig, QueryRequest):
    """Combined input model for the tool."""


class ConnectToPostgresTool(BaseTool[ConnectToPostgresInput, PostgresResult]):
    """LLM tool exposing PostgreSQL connector functionality."""

    name = "connect_to_postgres"
    description = "Connect to PostgreSQL and execute database operations"
    input_schema = ConnectToPostgresInput
    output_schema = PostgresResult

    def __init__(self, logger: AutoClassFuncLogger | None = None) -> None:
        super().__init__(logger)
        self._connector = AsyncpgConnector(self.logger)

    async def execute(self, params: ConnectToPostgresInput) -> PostgresResult:
        enriched = apply_postgres_env_defaults(params.model_dump())
        config = PostgresConfig(**enriched)
        request = QueryRequest(**enriched)

        result = await self._connector.run_with_timeout(config, request)

        return PostgresResult(
            success=result.connection_status == "success",
            connection_status=result.connection_status,
            server_version=result.server_version,
            data=result.data,
            error=result.error,
        )


