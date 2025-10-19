"""Tool for generating PostgreSQL DDL scripts."""

from __future__ import annotations

from typing import Any, Dict, List, Optional

from utils.logger import AutoClassFuncLogger

from ..base import BaseTool, ToolInput, ToolOutput
from ..utils.type_mapping import pandas_dtype_to_postgres
from ..utils.validators import ensure_non_empty, ensure_unique


class PostgresColumnConfig(ToolInput):
    name: str
    data_type: str
    nullable: bool = True
    is_potential_key: bool = False


class PostgresDDLInput(ToolInput):
    table_name: str
    columns: List[PostgresColumnConfig]
    primary_key: List[str]
    indexes: List[Dict[str, Any]] = []
    partitioning: Optional[Dict[str, Any]] = None
    constraints: List[Dict[str, Any]] = []
    options: Dict[str, Any] = {}


class PostgresDDLOutput(ToolOutput):
    ddl_scripts: List[Dict[str, Any]]
    estimated_size: Optional[str] = None
    recommendations: List[Dict[str, Any]] = []
    migration_script: Optional[str] = None


class GeneratePostgresDDLTool(BaseTool[PostgresDDLInput, PostgresDDLOutput]):
    name = "generate_postgres_ddl"
    description = "Generate PostgreSQL DDL based on column definitions"
    input_schema = PostgresDDLInput
    output_schema = PostgresDDLOutput

    def __init__(self, logger: AutoClassFuncLogger | None = None) -> None:
        super().__init__(logger)

    async def execute(self, params: PostgresDDLInput) -> PostgresDDLOutput:
        ensure_non_empty(params.table_name, "table_name")
        ensure_unique((column.name for column in params.columns), "columns")

        create_table = self._render_create_table(params)
        indexes = [self._render_index(index, params.table_name) for index in params.indexes]

        scripts = [
            {
                "type": "TABLE",
                "name": params.table_name,
                "script": create_table,
                "description": "Create base table",
                "execution_order": 1,
            },
        ]

        for order, script in enumerate(filter(None, indexes), start=2):
            scripts.append(script | {"execution_order": order})

        return PostgresDDLOutput(
            success=True,
            ddl_scripts=scripts,
            recommendations=self._recommendations(params),
        )

    def _render_create_table(self, params: PostgresDDLInput) -> str:
        column_definitions = []
        for column in params.columns:
            pg_type = pandas_dtype_to_postgres(column.data_type)
            nullable = "" if column.nullable else " NOT NULL"
            column_definitions.append(f'    "{column.name}" {pg_type}{nullable}')

        primary_key = ""
        if params.primary_key:
            pk_cols = ", ".join(f'"{col}"' for col in params.primary_key)
            primary_key = f",\n    PRIMARY KEY ({pk_cols})"

        options = []
        if params.options.get("if_not_exists", True):
            options.append("IF NOT EXISTS")
        if params.options.get("temporary"):
            options.append("TEMPORARY")
        if params.options.get("unlogged"):
            options.append("UNLOGGED")

        option_prefix = f"{' '.join(options)} " if options else ""

        columns_ddl = ",\n".join(column_definitions) + primary_key
        return f"CREATE {option_prefix}TABLE \"{params.table_name}\" (\n{columns_ddl}\n);"

    def _render_index(self, index: Dict[str, Any], table_name: str) -> Dict[str, Any] | None:
        index_type = index.get("type", "btree").upper()
        columns = index.get("columns", [])
        ensure_unique(columns, "index columns")
        if not columns:
            return None

        columns_ddl = ", ".join(f'"{col}"' for col in columns)
        name = index.get("name", f"{table_name}_{'_'.join(columns)}_{index_type.lower()}")
        script = (
            f"CREATE INDEX IF NOT EXISTS \"{name}\" ON \"{table_name}\" USING {index_type} ({columns_ddl});"
        )
        return {
            "type": "INDEX",
            "name": name,
            "script": script,
            "description": index.get("description", "Create index"),
        }

    def _recommendations(self, params: PostgresDDLInput) -> List[Dict[str, Any]]:
        recommendations: List[Dict[str, Any]] = []
        if not params.primary_key:
            recommendations.append(
                {
                    "type": "constraint",
                    "reason": "No primary key defined",
                    "impact": "high",
                }
            )
        for column in params.columns:
            if column.is_potential_key and column.name not in params.primary_key:
                recommendations.append(
                    {
                        "type": "constraint",
                        "reason": f"Column {column.name} appears unique",
                        "impact": "medium",
                    }
                )
        return recommendations


