"""Utility functions for mapping data types between ecosystems."""

from __future__ import annotations

from collections import defaultdict
from typing import Dict

import pandas as pd


PANDAS_TO_POSTGRES: Dict[str, str] = defaultdict(
    lambda: "TEXT",
    {
        "int64": "BIGINT",
        "int32": "INTEGER",
        "float64": "DOUBLE PRECISION",
        "float32": "REAL",
        "bool": "BOOLEAN",
        "datetime64[ns]": "TIMESTAMP",
    },
)


def pandas_dtype_to_postgres(dtype: pd.api.types.ExtensionDtype | str) -> str:
    """Map pandas dtype to PostgreSQL type."""

    dtype_name = str(dtype)
    return PANDAS_TO_POSTGRES[dtype_name]


CLICKHOUSE_TYPE_MAP: Dict[str, str] = defaultdict(
    lambda: "String",
    {
        "int64": "Int64",
        "int32": "Int32",
        "float64": "Float64",
        "float32": "Float32",
        "bool": "UInt8",
        "datetime64[ns]": "DateTime",
    },
)


def pandas_dtype_to_clickhouse(dtype: pd.api.types.ExtensionDtype | str) -> str:
    """Map pandas dtype to ClickHouse type."""

    dtype_name = str(dtype)
    return CLICKHOUSE_TYPE_MAP[dtype_name]


