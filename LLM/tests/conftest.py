import pytest
from fastapi.testclient import TestClient

from main import app, DataProcessor
from core.agents import LLMChainManager


class FakeLLM:
    async def chat(
        self,
        messages,
        *,
        model=None,
        temperature=0.3,
        max_tokens=2048,
        response_format=None
    ):
        return {"choices": [{"message": {"content": '{"ok": true}'}}]}


@pytest.fixture()
def client():
    app.state.llm_client = FakeLLM()
    app.state.chain_manager = LLMChainManager(app.state.llm_client)
    app.state.data_processor = DataProcessor()
    with TestClient(app) as c:
        yield c
