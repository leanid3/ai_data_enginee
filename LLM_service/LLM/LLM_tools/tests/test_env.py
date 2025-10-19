"""Tests for environment helpers."""

import os

from LLM.LLM_tools.utils.env import (
    apply_clickhouse_env_defaults,
    apply_postgres_env_defaults,
)


def test_apply_postgres_env_defaults(monkeypatch):
    monkeypatch.setenv("POSTGRES_HOST", "db.local")
    monkeypatch.setenv("POSTGRES_PORT", "5439")
    monkeypatch.setenv("POSTGRES_USE_SSL", "true")
    config = apply_postgres_env_defaults({})
    assert config["host"] == "db.local"
    assert config["port"] == 5439
    assert config["use_ssl"] is True


def test_apply_clickhouse_env_defaults(monkeypatch):
    monkeypatch.setenv("CLICKHOUSE_HOST", "ch.local")
    monkeypatch.setenv("CLICKHOUSE_PORT", "9000")
    monkeypatch.setenv("CLICKHOUSE_TIMEOUT", "5")
    config = apply_clickhouse_env_defaults({})
    assert config["host"] == "ch.local"
    assert config["port"] == 9000
    assert config["timeout"] == 5.0


