"""Core abstractions for LLM tools.

This module defines the fundamental building blocks required to create
tooling that can be orchestrated by LLM agents. The abstractions mirror the
technical specification and are intentionally lightweight to ease testing
and extension.
"""

from __future__ import annotations

import asyncio
from abc import ABC, abstractmethod
from time import perf_counter
from typing import Any, Dict, Generic, Optional, Type, TypeVar

from pydantic import BaseModel, ValidationError

from utils.logger import AutoClassFuncLogger, get_logger


class ToolInput(BaseModel):
    """Base class for tool input schemas."""


class ToolOutput(BaseModel):
    """Base class for tool output schemas."""

    success: bool
    data: Optional[Dict[str, Any]] = None
    error: Optional[str] = None
    execution_time_ms: float = 0.0


InputType = TypeVar("InputType", bound=ToolInput)
OutputType = TypeVar("OutputType", bound=ToolOutput)


class BaseTool(Generic[InputType, OutputType], ABC):
    """Abstract base class for all tools.

    Every tool is asynchronous and uses pydantic models for input/output
    validation. Concrete implementations must define meta-information and the
    `execute` coroutine.
    """

    name: str
    description: str
    input_schema: Type[InputType]
    output_schema: Type[OutputType]

    def __init__(self, logger: AutoClassFuncLogger | None = None) -> None:
        self.logger = logger or get_logger("LLM-TOOLS-1.0")

    def validate_input(self, params: Dict[str, Any]) -> InputType:
        """Validate raw parameters using the declared input schema."""

        try:
            return self.input_schema(**params)
        except ValidationError as exc:  # pragma: no cover - exercised indirectly
            self.logger.error("Input validation error", extra={"error": str(exc)})
            raise

    @abstractmethod
    async def execute(self, params: InputType) -> OutputType:
        """Execute the tool logic."""

    async def __call__(self, **kwargs: Any) -> Dict[str, Any]:
        """Convenience wrapper to make tools directly callable."""

        validated = self.validate_input(kwargs)
        start = perf_counter()

        try:
            result = await self.execute(validated)
        except Exception as exc:  # pragma: no cover - exercised indirectly
            duration = (perf_counter() - start) * 1000
            self.logger.error(
                "Tool execution failed",
                extra={
                    "tool": self.name,
                    "error": str(exc),
                },
            )
            return self.output_schema(
                success=False,
                data=None,
                error=str(exc),
                execution_time_ms=duration,
            ).model_dump()

        result.execution_time_ms = (perf_counter() - start) * 1000
        return result.model_dump()

    async def run_in_thread(self, func, *args, **kwargs):
        """Helper to offload blocking work to a thread."""

        return await asyncio.to_thread(func, *args, **kwargs)


