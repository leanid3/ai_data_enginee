"""LLM tools package initialization.

This package provides a modular toolset that can be consumed by LLM agents.
It exposes a registry with all available tools and centralizes logger usage
across the package.
"""

from __future__ import annotations

from typing import Iterable

from LLM_service.utils.logger import AutoClassFuncLogger, get_logger

from .base import BaseTool
from .registry import ToolRegistry

DEFAULT_LOGGER_VERSION = "LLM-TOOLS-1.0"


def _initialize_logger() -> AutoClassFuncLogger:
    """Create a shared logger instance for the package."""

    return get_logger(DEFAULT_LOGGER_VERSION)


LOGGER: AutoClassFuncLogger = _initialize_logger()
REGISTRY: ToolRegistry = ToolRegistry(logger=LOGGER)

_DEFAULT_TOOLS_REGISTERED = False


def register_tools(tools: Iterable[type[BaseTool]]) -> None:
    """Register tools in the global registry.

    Parameters
    ----------
    tools:
        Iterable of tool classes to register.
    """

    for tool_cls in tools:
        REGISTRY.register(tool_cls())


def register_default_tools() -> None:
    """Register all built-in tools with the global registry."""

    global _DEFAULT_TOOLS_REGISTERED

    if _DEFAULT_TOOLS_REGISTERED:
        return

    from .analyzers import AnalyzeCsvStructureTool, SchemaAnalyzerTool
    from .connectors import ConnectToClickHouseTool, ConnectToPostgresTool
    from .formatters import FormatMarkdownReportTool
    from .generators import (
        CreateAirflowDagTool,
        GenerateClickHouseDDLTool,
        GeneratePostgresDDLTool,
    )

    register_tools(
        [
            AnalyzeCsvStructureTool,
            SchemaAnalyzerTool,
            ConnectToPostgresTool,
            ConnectToClickHouseTool,
            GeneratePostgresDDLTool,
            GenerateClickHouseDDLTool,
            CreateAirflowDagTool,
            FormatMarkdownReportTool,
        ]
    )

    _DEFAULT_TOOLS_REGISTERED = True


register_default_tools()


__all__ = [
    "BaseTool",
    "REGISTRY",
    "LOGGER",
    "register_default_tools",
    "register_tools",
]


