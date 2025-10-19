from __future__ import annotations

import asyncio
from typing import Any, Callable, Dict, Iterable, List, Optional

from pydantic import BaseModel

from LLM.LLM_tools import REGISTRY, register_default_tools
from LLM.LLM_tools.base import BaseTool


class ToolSpec(BaseModel):
    name: str
    description: str
    input_schema: Dict[str, Any]
    output_type: str


class Tool:
    """Wrapper around callable or LLM Tool that exposes unified spec."""

    def __init__(
        self,
        func: Callable[..., Any],
        name: str,
        description: str,
        input_schema: Dict[str, Any],
        output_type: str = "string",
        base_tool: Optional[BaseTool] = None,
    ) -> None:
        self.func = func
        self.spec = ToolSpec(
            name=name,
            description=description,
            input_schema=input_schema,
            output_type=output_type,
        )
        self.base_tool = base_tool

    async def __call__(self, **kwargs: Any) -> Any:
        if asyncio.iscoroutinefunction(self.func):
            return await self.func(**kwargs)
        return self.func(**kwargs)

    @property
    def name(self) -> str:
        return self.spec.name

    @property
    def description(self) -> str:
        return self.spec.description

    @property
    def input_schema(self) -> Dict[str, Any]:
        return self.spec.input_schema

    @classmethod
    def from_llm_tool(cls, tool: BaseTool) -> "Tool":
        """Create a registry Tool wrapper from a BaseTool instance."""

        return cls(
            func=tool.__call__,
            name=tool.name,
            description=tool.description,
            input_schema=tool.input_schema.model_json_schema(),
            output_type="json",
            base_tool=tool,
        )


def load_llm_tools() -> List[Tool]:
    """Instantiate Tool wrappers for all registered LLM tools."""

    register_default_tools()
    tools: List[Tool] = []
    for meta in REGISTRY.list_tools():
        base_tool = REGISTRY.get_tool(meta["name"])
        tools.append(Tool.from_llm_tool(base_tool))
    return tools


class ToolRegistry:
    def __init__(
        self,
        initial_tools: Optional[Iterable[Tool]] = None,
        auto_register_llm_tools: bool = False,
    ) -> None:
        self._tools: Dict[str, Tool] = {}

        if auto_register_llm_tools:
            self.register_many(load_llm_tools())

        if initial_tools:
            self.register_many(initial_tools)

    def register(self, tool: Tool) -> None:
        self._tools[tool.name] = tool

    def register_many(self, tools: Iterable[Tool]) -> None:
        for tool in tools:
            self.register(tool)

    def get(self, name: str) -> Optional[Tool]:
        return self._tools.get(name)

    def list_specs(self) -> List[ToolSpec]:
        return [tool.spec for tool in self._tools.values()]

    async def execute(self, name: str, arguments: Dict[str, Any]) -> Any:
        tool = self.get(name)
        if not tool:
            raise ValueError(f"Tool '{name}' not found")
        return await tool(**arguments)