"""ClickHouse数据库连接器实现。
基于clickhouse-connect驱动，提供高性能的异步数据访问接口。
支持查询执行、模式分析和统计信息收集。
"""

from __future__ import annotations

from typing import Any, Dict, Optional

from clickhouse_connect.driver import Client as CHClient
from clickhouse_connect.driver.exceptions import ClickHouseError
from pydantic import Field

from utils.logger import AutoClassFuncLogger

from ..base import BaseTool, ToolOutput
from ..utils.env import apply_clickhouse_env_defaults
from .base import BaseConnector, ConnectionConfig, QueryRequest


class ClickHouseConfig(ConnectionConfig):
    """Configuration model for ClickHouse."""

    port: int = Field(8123, ge=1, le=65535)
    database: str = Field("default", min_length=1)
    username: str = Field("default", min_length=1)
    password: str = ""


class ClickHouseResult(ToolOutput):
    """Result schema for ClickHouse tool."""

    connection_status: str = "failed"
    server_version: Optional[str] = None
    data: Dict[str, Any] = {}


# ClickHouse数据库连接器类
# 特点：
# - 使用官方clickhouse-connect驱动
# - 支持高性能数据传输
# - 提供系统表查询功能
# - 实现数据统计分析
class ClickHouseConnector(BaseConnector):
    """ClickHouse connector using clickhouse-connect driver."""

    def __init__(self, logger: AutoClassFuncLogger) -> None:
        super().__init__(logger)

    @property
    def name(self) -> str:
        return "clickhouse-connect"

    async def connect(self, config: ClickHouseConfig) -> CHClient:
        return CHClient(
            host=config.host,
            port=config.port,
            username=config.username,
            password=config.password,
            database=config.database,
            secure=config.use_ssl,
        )

    async def close(self, connection: CHClient) -> None:
        connection.close()

    async def fetch_server_version(self, connection: CHClient) -> Optional[str]:
        info = connection.server_info
        return getattr(info, "version", None)

    async def execute_operation(
        self,
        connection: CHClient,
        request: QueryRequest,
    ) -> Dict[str, Any]:
        if request.operation == "schema":
            return self._fetch_schema(connection)
        if request.operation == "stats":
            return self._fetch_stats(connection)

        assert request.query is not None
        result = connection.query(request.query)
        return {
            "columns": result.column_names,
            "rows": result.result_rows,
            "row_count": result.row_count,
        }

    # 获取数据库模式信息
    # 查询system.tables系统表
    # 返回：
    # - 数据库名称
    # - 表名称
    # - 表引擎类型
    def _fetch_schema(self, connection: CHClient) -> Dict[str, Any]:
        query = """
        SELECT database, table, engine
        FROM system.tables
        WHERE is_temporary = 0;
        """
        result = connection.query(query)
        return {
            "tables": [dict(zip(result.column_names, row)) for row in result.result_rows],
        }

    # 获取数据库统计信息
    # 查询system.parts系统表
    # 统计：
    # - 表大小
    # - 行数
    # - 分区信息
    def _fetch_stats(self, connection: CHClient) -> Dict[str, Any]:
        query = """
        SELECT database, table, total_rows AS row_count, total_bytes / (1024*1024) AS size_mb
        FROM system.parts
        GROUP BY database, table;
        """
        result = connection.query(query)
        return {
            "stats": [dict(zip(result.column_names, row)) for row in result.result_rows],
        }


class ConnectToClickHouseInput(ClickHouseConfig, QueryRequest):
    """Input schema for the ClickHouse tool."""


class ConnectToClickHouseTool(
    BaseTool[ConnectToClickHouseInput, ClickHouseResult]
):
    """LLM tool for ClickHouse interaction."""

    name = "connect_to_clickhouse"
    description = "Connect to ClickHouse and execute database operations"
    input_schema = ConnectToClickHouseInput
    output_schema = ClickHouseResult

    def __init__(self, logger: AutoClassFuncLogger | None = None) -> None:
        super().__init__(logger)
        self._connector = ClickHouseConnector(self.logger)

    async def execute(self, params: ConnectToClickHouseInput) -> ClickHouseResult:
        enriched = apply_clickhouse_env_defaults(params.model_dump())
        config = ClickHouseConfig(**enriched)
        request = QueryRequest(**enriched)

        result = await self._connector.run_with_timeout(config, request)

        return ClickHouseResult(
            success=result.connection_status == "success",
            connection_status=result.connection_status,
            server_version=result.server_version,
            data=result.data,
            error=result.error,
        )


