from typing import Any, Dict, List, Protocol, Optional, TypedDict


class ChatResponse(TypedDict):
    choices: List[dict]


class LLMClientProtocol(Protocol):
    async def chat(
        self,
        messages: List[Dict[str, Any]],
        *,
        model: Optional[str] = None,
        temperature: float = 0.3,
        max_tokens: int = 2048,
        response_format: Optional[Dict[str, Any]] = None,
    ) -> ChatResponse: ...
