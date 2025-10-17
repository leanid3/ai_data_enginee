"""Tests for BaseTool behavior."""

import pytest

from LLM.LLM_tools.base import BaseTool, ToolInput, ToolOutput


class EchoInput(ToolInput):
    message: str


class EchoOutput(ToolOutput):
    result: str


class EchoTool(BaseTool[EchoInput, EchoOutput]):
    name = "echo"
    description = "Echo tool"
    input_schema = EchoInput
    output_schema = EchoOutput

    async def execute(self, params: EchoInput) -> EchoOutput:
        return EchoOutput(success=True, result=params.message)


@pytest.mark.asyncio
async def test_base_tool_call_executes_successfully():
    tool = EchoTool()
    result = await tool(message="hello")

    assert result["success"] is True
    assert result["result"] == "hello"


