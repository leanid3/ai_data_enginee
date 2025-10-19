from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, File, UploadFile
from fastapi.responses import JSONResponse
from pathlib import Path
import aiofiles

from API.schemas import LLMRequest, LLMResponse
from API.core import DataProcessor
from API.core import DataEngineeringOrchestrator
from API.interfaces import OpenRouterClient
from API.schemas import SourceConfig
from API.schemas import TargetConfig
from LLM.core import LLMChainManager
from LLM.core import AgentRole

from utils import get_logger

__VERSION__ = "0.1"

logger = get_logger(__VERSION__)

pipeline_storage = {}

@asynccontextmanager
async def lifespan(app: FastAPI):

    logger.info("Starting LLM Service...")
    llm_client = OpenRouterClient()
    chain_manager = LLMChainManager(llm_client)
    app.state.orchestrator = DataEngineeringOrchestrator(chain_manager)
    app.state.data_processor = DataProcessor()
    yield
    logger.info("Shutting down...")

app = FastAPI(title="AI Data Engineer API", version="3.0.0", lifespan=lifespan)

@app.post("/api/v1/process", response_model=LLMResponse)
async def process_data_request(request: LLMRequest) -> LLMResponse:
    try:
        response = await app.state.orchestrator.execute_chain(request)
        pipeline_storage[response.pipeline_id] = response.dict()
        return response
    except Exception as e:
        logger.error(f"Error processing request: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/v1/analyze-file")
async def analyze_file(file: UploadFile = File(...), target_system: str = "auto"):
    try:
        content = await file.read()
        filename = file.filename
        
        if filename.endswith(".csv"):
            source_type = "csv"
        elif filename.endswith(".json"):
            source_type = "json"
        else:
            raise HTTPException(status_code=400, detail="Unsupported file type")

        data_sample = await app.state.data_processor(source_type, content, filename)

        req = LLMRequest(
            user_query=f"Проанализировать файл {filename} и создать ETL пайплайн",
            source_config=SourceConfig(type=source_type, file_path=filename),
            target_config=TargetConfig(type=target_system),
            data_sample=data_sample,
            operation_type="full_process",
        )
        return await process_data_request(req)
    except Exception as e:
        logger.error(f"Error analyzing file: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/v1/pipeline/{pipeline_id}")
async def get_pipeline(pipeline_id: str):
    if pipeline_id not in pipeline_storage:
        raise HTTPException(status_code=404, detail="Pipeline not found")
    return pipeline_storage[pipeline_id]

@app.post("/api/v1/pipeline/{pipeline_id}/execute")
async def execute_pipeline(pipeline_id: str):
    if pipeline_id not in pipeline_storage:
        raise HTTPException(status_code=404, detail="Pipeline not found")
    pipeline = pipeline_storage[pipeline_id]
    return {
        "status": "scheduled",
        "pipeline_id": pipeline_id,
        "message": "Pipeline scheduled for execution",
        "dag_code": pipeline.get("dag_code", ""),
    }

@app.get("/api/v1/health")
async def health():
    return JSONResponse(
        content={
            "status": "healthy",
            "service": "AI Data Engineer",
            "version": "3.0.0",
            "agents_available": [role.value for role in AgentRole],
            "pipelines_cached": len(pipeline_storage),
        }
    )

@app.get("/api/v1/agents")
async def get_agents():
    return JSONResponse(
        content={
            "agents": [
                {"role": AgentRole.ROUTER, "description": "Маршрутизатор запросов"},
                {"role": AgentRole.DATA_ANALYZER, "description": "Анализ структуры и качества данных"},
                {"role": AgentRole.DB_SELECTOR, "description": "Выбор системы хранения"},
                {"role": AgentRole.DDL_GENERATOR, "description": "Генерация DDL"},
                {"role": AgentRole.ETL_BUILDER, "description": "Построение ETL пайплайнов"},
                {"role": AgentRole.QUERY_OPTIMIZER, "description": "Оптимизация SQL"},
                {"role": AgentRole.REPORT_GENERATOR, "description": "Отчеты и рекомендации"},
            ]
        }
    )