"""Environment helpers for tool configuration."""

from __future__ import annotations

import os
from typing import Any, Callable, Dict


def _to_bool(value: str) -> bool:
    return value.strip().lower() in {"1", "true", "yes", "on"}


def _set_if_missing(
    data: Dict[str, Any],
    field: str,
    env_name: str,
    caster: Callable[[str], Any] | None = None,
) -> None:
    if data.get(field) not in (None, ""):
        return

    raw_value = os.getenv(env_name)
    if raw_value is None:
        return

    if caster:
        try:
            data[field] = caster(raw_value)
        except ValueError as exc:  # pragma: no cover - configuration issue
            raise ValueError(f"Invalid value for {env_name}: {raw_value}") from exc
    else:
        data[field] = raw_value


def apply_postgres_env_defaults(data: Dict[str, Any] | None) -> Dict[str, Any]:
    """Populate PostgreSQL connection fields from environment when missing."""

    result: Dict[str, Any] = dict(data or {})
    _set_if_missing(result, "host", "POSTGRES_HOST")
    _set_if_missing(result, "port", "POSTGRES_PORT", int)
    _set_if_missing(result, "database", "POSTGRES_DB")
    _set_if_missing(result, "username", "POSTGRES_USER")
    _set_if_missing(result, "password", "POSTGRES_PASSWORD")
    _set_if_missing(result, "use_ssl", "POSTGRES_USE_SSL", _to_bool)
    _set_if_missing(result, "timeout", "POSTGRES_TIMEOUT", float)
    return result


def apply_clickhouse_env_defaults(data: Dict[str, Any] | None) -> Dict[str, Any]:
    """Populate ClickHouse connection fields from environment when missing."""

    result: Dict[str, Any] = dict(data or {})
    _set_if_missing(result, "host", "CLICKHOUSE_HOST")
    _set_if_missing(result, "port", "CLICKHOUSE_PORT", int)
    _set_if_missing(result, "database", "CLICKHOUSE_DB")
    _set_if_missing(result, "username", "CLICKHOUSE_USER")
    _set_if_missing(result, "password", "CLICKHOUSE_PASSWORD")
    _set_if_missing(result, "use_ssl", "CLICKHOUSE_USE_SSL", _to_bool)
    _set_if_missing(result, "timeout", "CLICKHOUSE_TIMEOUT", float)
    return result


