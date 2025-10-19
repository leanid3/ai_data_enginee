import os
from typing import Any, Dict, List, Optional
from dotenv import load_dotenv
from openai import AsyncOpenAI
from LLM.interfaces import LLMClientProtocol
from LLM.interfaces.client import ChatResponse

load_dotenv()

class OpenRouterClient(LLMClientProtocol):
    def __init__(self) -> None:
        api_key = os.getenv("OPENROUTER_API_KEY", "").strip()
        base_url = os.getenv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1").strip()
        if not api_key:
            raise ValueError("OPENROUTER_API_KEY отсутствует в .env")
        self._client = AsyncOpenAI(base_url=base_url, api_key=api_key)

    async def chat(
        self,
        messages: List[Dict[str, Any]],
        *,
        model: Optional[str] = None,
        temperature: float = 0.3,
        max_tokens: int = 2048,
        response_format: Optional[Dict[str, Any]] = None,
    ) -> ChatResponse:
        model = model or os.getenv("DEFAULT_MODEL", "deepseek/deepseek-chat-v3.1:free")
        response_format = response_format or {"type": "json_object"}

        response = await self._client.chat.completions.create(
            model=model,
            messages=messages,
            temperature=temperature,
            max_tokens=max_tokens,
            response_format=response_format,
        )
        return {
            "choices": [
                {
                    "message": {
                        "content": response.choices[0].message.content or "",
                    }
                }
            ]
        }