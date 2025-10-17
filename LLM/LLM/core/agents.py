import json
from typing import Any, Dict, Optional
from LLM.interfaces import LLMClientProtocol
from .roles import AgentRole
from .prompts import PromptTemplates
from .tools import ToolRegistry


class BaseAgent:
    def __init__(self, prompt_template: str, role: AgentRole, tools: Optional[ToolRegistry]):
        self.prompt_template = prompt_template
        self.role = role
        self.tools = tools

    async def execute(
        self, context: Dict[str, Any], llm_client: LLMClientProtocol
    ) -> Dict[str, Any]:
        prompt = self.prompt_template
        for key, value in context.items():
            placeholder = "{" + str(key) + "}"
            if placeholder in prompt:
                prompt = prompt.replace(
                    placeholder,
                    (
                        json.dumps(value, ensure_ascii=False)
                        if isinstance(value, (dict, list))
                        else str(value)
                    ),
                )

        if self.tools:
            tools_desc = "\n".join(
                f"- {t.name}: {t.description} (вход: {t.input_schema})"
                for t in self.tools.list_specs()
            )
            prompt += f"\n\nAvailable instruments:\n{tools_desc}"

        response = await llm_client.chat(
            [
                {
                    "role": "system",
                    "content": "Ты - агент. Отвечай только в формате JSON.",
                },
                {"role": "user", "content": prompt},
            ],
            model=None,
            temperature=float(context.get("temperature", 0.3)),
            max_tokens=int(context.get("max_tokens", 2048)),
            response_format={"type": "json_object"},
        )
        return response


class LLMChainManager:
    def __init__(self, llm_client: LLMClientProtocol):
        self.llm_client = llm_client
        self.tool_registry = ToolRegistry(auto_register_llm_tools=True)
        self.agents: Dict[AgentRole, BaseAgent] = {
            AgentRole.ROUTER: BaseAgent(PromptTemplates.router(), AgentRole.ROUTER, None),
            AgentRole.DATA_ANALYZER: BaseAgent(
                PromptTemplates.analyzer(), AgentRole.DATA_ANALYZER, self.tool_registry
            ),
            AgentRole.DB_SELECTOR: BaseAgent(
                PromptTemplates.db_selector(), AgentRole.DB_SELECTOR, self.tool_registry
            ),
            AgentRole.DDL_GENERATOR: BaseAgent(
                PromptTemplates.ddl_generator(), AgentRole.DDL_GENERATOR, self.tool_registry
            ),
            AgentRole.ETL_BUILDER: BaseAgent(
                PromptTemplates.etl_builder(), AgentRole.ETL_BUILDER, self.tool_registry
            ),
            AgentRole.QUERY_OPTIMIZER: BaseAgent(
                PromptTemplates.optimizer(), AgentRole.QUERY_OPTIMIZER, self.tool_registry
            ),
            AgentRole.REPORT_GENERATOR: BaseAgent(
                PromptTemplates.reporter(), AgentRole.REPORT_GENERATOR, self.tool_registry
            ),
        }

    async def call_agent(
        self, role: AgentRole, context: Dict[str, Any]
    ) -> Dict[str, Any]:
        agent = self.agents[role]
        raw = await agent.execute(context, self.llm_client)
        try:
            if "choices" in raw:
                content = raw["choices"][0]["message"]["content"]
                return json.loads(content)
            else:
                return raw
        except Exception as e:
            return {"error": str(e), "raw_response": raw}