"""Integration tests for registry usage."""

from LLM.LLM_tools import REGISTRY, register_default_tools


def test_registry_has_default_tools():
    register_default_tools()
    tools = {tool["name"] for tool in REGISTRY.list_tools()}
    assert "analyze_csv_structure" in tools
    assert "connect_to_postgres" in tools
    assert "generate_postgres_ddl" in tools


