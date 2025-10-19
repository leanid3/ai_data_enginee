"""Tool for generating ClickHouse DDL scripts."""

from __future__ import annotations

from typing import Any, Dict, List, Optional

from LLM_service.utils.logger import AutoClassFuncLogger

from ..base import BaseTool, ToolInput, ToolOutput
from ..utils.type_mapping import pandas_dtype_to_clickhouse
from ..utils.validators import ensure_non_empty, ensure_unique


class ClickHouseColumnConfig(ToolInput):
    name: str
    data_type: str
    nullable: bool = True


class ClickHouseDDLInput(ToolInput):
    table_name: str
    columns: List[ClickHouseColumnConfig]
    engine: str = "MergeTree"
    partition_by: Optional[str] = None
    order_by: List[str]
    primary_key: List[str] = []
    ttl: Optional[str] = None
    settings: Dict[str, Any] = {}


class ClickHouseDDLOutput(ToolOutput):
    ddl_scripts: List[Dict[str, Any]]
    engine_config: Dict[str, Any]
    performance_tips: List[str] = []
    materialized_views: List[Dict[str, Any]] = []


class GenerateClickHouseDDLTool(
    BaseTool[ClickHouseDDLInput, ClickHouseDDLOutput]
):
    name = "generate_clickhouse_ddl"
    description = "Generate ClickHouse DDL based on column definitions"
    input_schema = ClickHouseDDLInput
    output_schema = ClickHouseDDLOutput

    def __init__(self, logger: AutoClassFuncLogger | None = None) -> None:
        super().__init__(logger)

    async def execute(self, params: ClickHouseDDLInput) -> ClickHouseDDLOutput:
        ensure_non_empty(params.table_name, "table_name")
        ensure_unique((column.name for column in params.columns), "columns")
        ensure_unique(params.order_by, "order_by")

        script = self._render_create_table(params)
        engine_config = self._build_engine_config(params)
        tips = self._performance_tips(params)

        return ClickHouseDDLOutput(
            success=True,
            ddl_scripts=[
                {
                    "type": "TABLE",
                    "script": script,
                    "on_cluster": params.settings.get("on_cluster", False),
                    "replicated": params.settings.get("replicated", False),
                }
            ],
            engine_config=engine_config,
            performance_tips=tips,
        )

    def _render_create_table(self, params: ClickHouseDDLInput) -> str:
        columns = []
        for column in params.columns:
            ch_type = pandas_dtype_to_clickhouse(column.data_type)
            if column.nullable:
                ch_type = f"Nullable({ch_type})"
            columns.append(f"    `{column.name}` {ch_type}")
        columns_ddl = ",\n".join(columns)

        engine = params.engine
        order_by = f"ORDER BY ({', '.join(params.order_by)})"
        partition_by = f"PARTITION BY {params.partition_by}" if params.partition_by else ""
        primary_key = f"PRIMARY KEY ({', '.join(params.primary_key)})" if params.primary_key else ""
        ttl = f"TTL {params.ttl}" if params.ttl else ""
        settings = self._render_settings(params.settings)

        clauses = "\n".join(filter(None, [partition_by, primary_key, ttl, settings]))

        return (
            f"CREATE TABLE `{params.table_name}` (\n{columns_ddl}\n)\n"
            f"ENGINE = {engine}\n{order_by}\n{clauses};"
        )

    def _render_settings(self, settings: Dict[str, Any]) -> str:
        ignore_keys = {"on_cluster", "replicated"}
        relevant = {k: v for k, v in settings.items() if k not in ignore_keys}
        if not relevant:
            return ""
        formatted = ", ".join(f"{key} = {value}" for key, value in relevant.items())
        return f"SETTINGS {formatted}"

    def _build_engine_config(self, params: ClickHouseDDLInput) -> Dict[str, Any]:
        return {
            "engine": params.engine,
            "partition_key": params.partition_by,
            "sorting_key": ", ".join(params.order_by),
            "primary_key": ", ".join(params.primary_key),
            "ttl": params.ttl,
            "settings": params.settings,
        }

    def _performance_tips(self, params: ClickHouseDDLInput) -> List[str]:
        tips: List[str] = []
        if not params.partition_by:
            tips.append("Consider setting PARTITION BY for large tables")
        if not params.primary_key:
            tips.append("Define PRIMARY KEY to improve query efficiency")
        return tips


