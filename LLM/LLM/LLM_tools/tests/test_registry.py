"""Tests for the tool registry."""

from unittest.mock import Mock

from LLM.LLM_tools import LOGGER
from LLM.LLM_tools.registry import ToolRegistry


def test_registry_register_and_get_tool():
    registry = ToolRegistry(LOGGER)
    dummy_tool = Mock()
    dummy_tool.name = "dummy"
    dummy_tool.description = "dummy tool"
    dummy_tool.input_schema = Mock()
    dummy_tool.output_schema = Mock()
    registry.register(dummy_tool)

    tool = registry.get_tool("dummy")
    assert tool is dummy_tool


