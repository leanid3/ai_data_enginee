from datetime import datetime
import uuid
from typing import Any, Dict, List

from LLM.core import AgentChain, AgentRole, LLMChainManager
from LLM.core.tools import load_llm_tools

from API.schemas import LLMRequest, LLMResponse


class DataEngineeringOrchestrator:
    def __init__(self, chain_manager: LLMChainManager):
        self.chain_manager = chain_manager
        self._available_tools = sorted({tool.name for tool in load_llm_tools()})
    async def execute_chain(self, request: LLMRequest) -> LLMResponse:
        start_time = datetime.now()
        chain = AgentChain.get_chain(request.operation_type)
        results: Dict[Any, Any] = {}
        agents_used: List[str] = []
        tools_used: List[str] = []

        for role in chain:
            context = self._build_context(role, request, results)
            result = await self.chain_manager.call_agent(role, context)
            results[role] = result
            agents_used.append(role.value)

            declared_tools: List[str] = []
            if isinstance(result, dict):
                candidate = result.get("tools_used") or result.get("tools_available")
                if isinstance(candidate, list):
                    declared_tools = [str(tool) for tool in candidate]

            tools_used.extend(declared_tools)

        processing_time = (datetime.now() - start_time).total_seconds()
        pipeline_id = str(uuid.uuid4())

        data_analysis = results.get(AgentRole.DATA_ANALYZER, {})
        storage_recommendation = results.get(AgentRole.DB_SELECTOR, {})
        ddl_results = results.get(AgentRole.DDL_GENERATOR, {})
        etl_results = results.get(AgentRole.ETL_BUILDER, {})
        optimization_results = results.get(AgentRole.QUERY_OPTIMIZER, {})
        report_results = results.get(AgentRole.REPORT_GENERATOR, {})

        # 1. Получаем список строк
        optimizations_list_of_strings = optimization_results.get("optimizations", [])

        # 2. Преобразуем список строк в список словарей
        formatted_optimized_queries = [{"query": text} for text in optimizations_list_of_strings]
        unique_tools = sorted(set(tools_used)) or self._available_tools

        return LLMResponse(
            success=True,
            pipeline_id=pipeline_id,
            data_analysis=data_analysis,
            storage_recommendation=storage_recommendation,
            ddl_scripts=ddl_results.get("ddl_scripts", []),
            dag_code=etl_results.get("python_code", ""),
            
            # 3. Передаем в модель уже отформатированный список
            optimized_queries=formatted_optimized_queries,

            visualization_config={"nodes": [], "edges": []},
            user_report=report_results.get("summary") or "Отчёт не сгенерирован.",
            processing_time=processing_time,
            agents_used=agents_used,
            tools_used=unique_tools,
            confidence_score=0.85,
            errors=[],
            warnings=[],
        )

    def _build_context(self, role: AgentRole, request: LLMRequest, results: Dict) -> Dict[str, Any]:
        base = {
            "user_request": request.user_query,
            "source_type": request.source_config.type,
            "target_system": request.target_config.type,
            "operation_type": request.operation_type,
        }

        if role == AgentRole.DATA_ANALYZER and request.data_sample:
            base["data_sample"] = {
                "headers": request.data_sample.headers,
                "sample_rows": request.data_sample.sample_rows[:5],
                "total_rows": request.data_sample.total_rows,
            }
            base["metadata"] = request.metadata or {}

        if role in (AgentRole.DB_SELECTOR, AgentRole.DDL_GENERATOR):
            base["analysis_results"] = results.get(AgentRole.DATA_ANALYZER, {})

        if role == AgentRole.DDL_GENERATOR:
            base["target_db"] = request.target_config.type
            base["data_structure"] = results.get(AgentRole.DATA_ANALYZER, {}).get("structure", {})
            base["storage_config"] = results.get(AgentRole.DB_SELECTOR, {}).get("config", {})

        if role == AgentRole.ETL_BUILDER:
            base["source_config"] = request.source_config.dict()
            base["target_config"] = request.target_config.dict()
            base["schedule"] = "0 */1 * * *"

        if role == AgentRole.QUERY_OPTIMIZER:
            base["db_type"] = request.target_config.type
            base["queries"] = []
            base["data_volume"] = (
                results.get(AgentRole.DATA_ANALYZER, {})
                .get("characteristics", {})
                .get("volume", "medium")
            )

        if role == AgentRole.REPORT_GENERATOR:
            base["all_results"] = results

        return base