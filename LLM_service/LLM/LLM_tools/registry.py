"""Tool registry implementation."""

from __future__ import annotations

from typing import Dict, Iterable, List

from LLM_service.utils.logger import AutoClassFuncLogger

from .base import BaseTool


class ToolRegistry:
    """Central registry that stores tool instances by name."""

    def __init__(self, logger: AutoClassFuncLogger) -> None:
        self._logger = logger
        self._tools: Dict[str, BaseTool] = {}

    def register(self, tool: BaseTool) -> None:
        """Register a tool instance using its declared name."""

        if tool.name in self._tools:
            self._logger.warning(
                "Tool already registered",
                extra={"tool": tool.name},
            )
            return

        self._logger.info(
            "Registering tool",
            extra={"tool": tool.name},
        )
        self._tools[tool.name] = tool

    def get_tool(self, name: str) -> BaseTool:
        """Return a tool by name or raise a KeyError."""

        try:
            return self._tools[name]
        except KeyError as exc:  # pragma: no cover - defensive
            self._logger.error(
                "Tool not found",
                extra={"tool": name},
            )
            raise exc

    def list_tools(self) -> List[dict]:
        """Return metadata for all registered tools."""

        return [
            {
                "name": tool.name,
                "description": tool.description,
                "input_schema": tool.input_schema.model_json_schema(),
                "output_schema": tool.output_schema.model_json_schema(),
            }
            for tool in self._tools.values()
        ]

    def bulk_register(self, tools: Iterable[BaseTool]) -> None:
        """Register multiple tool instances."""

        for tool in tools:
            self.register(tool)


