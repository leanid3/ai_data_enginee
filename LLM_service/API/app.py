from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, File, UploadFile
from fastapi.responses import JSONResponse
from pathlib import Path
import aiofiles

from LLM.LLM_tools.connectors.minio_conn import MinIOConnector, MinIOConnectionConfig, QueryRequest
from utils.logger import AutoClassFuncLogger
from API.schemas import LLMRequest, LLMResponse
from API.core import DataProcessor
from API.core import DataEngineeringOrchestrator
from API.interfaces import OpenRouterClient
from API.schemas import SourceConfig
from API.schemas import TargetConfig
from LLM.core import LLMChainManager
from LLM.core import AgentRole

from utils.logger import get_logger

__VERSION__ = "0.1"

logger = get_logger(__VERSION__)

pipeline_storage = {}

@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("Starting LLM Service...")

    # Инициализация LLM
    llm_client = OpenRouterClient()
    chain_manager = LLMChainManager(llm_client)
    app.state.orchestrator = DataEngineeringOrchestrator(chain_manager)
    app.state.data_processor = DataProcessor()
    minio_logger = AutoClassFuncLogger("MinIOConnector")
    minio_config = MinIOConnectionConfig(
        host="localhost",
        port=9000,
        username="minioadmin",
        password="minioadmin",
        database="ai-data-engineer",  # placeholder, как обсуждалось
        use_ssl=False,
        timeout=30.0,
    )
    minio_connector = MinIOConnector(minio_logger)
    app.state.minio_connector = minio_connector
    app.state.minio_config = minio_config
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
async def analyze_file(user_id: str, target_system: str = "auto"):
    try:
        connector: MinIOConnector = app.state.minio_connector
        config: MinIOConnectionConfig = app.state.minio_config

        query_request = QueryRequest(operation="query", query=user_id)
        result = await connector.run(config, query_request)

        if result.connection_status != "success":
            raise HTTPException(status_code=500, detail=f"MinIO query failed: {result.error}")

        files = result.data.get("files", [])
        if not files:
            raise HTTPException(status_code=404, detail=f"No files found for user '{user_id}'")

        filename = files[0]["filename"]
        logger.info(f"Selected file for analysis: {filename} (user: {user_id})")

        read_request = QueryRequest(operation="read", query=f"{user_id}/{filename}")
        read_result = await connector.run(config, read_request)

        if read_result.connection_status != "success":
            raise HTTPException(status_code=500, detail=f"MinIO read failed: {read_result.error}")

        content_bytes = read_result.data["content_preview"]

        try:
            content = content_bytes.encode("utf-8")
        except UnicodeEncodeError:
            raise HTTPException(status_code=500, detail="File is binary and cannot be processed")

        # 4. Определяем тип
        if filename.endswith(".csv"):
            source_type = "csv"
        elif filename.endswith(".json"):
            source_type = "json"
        else:
            raise HTTPException(status_code=400, detail="Unsupported file type")

        data_sample = await app.state.data_processor(source_type, content, filename)

        req = LLMRequest(
            user_query=f"Проанализировать файл {filename} и создать ETL пайплайн",
            source_config=SourceConfig(type=source_type, file_path=f"{user_id}/{filename}"),
            target_config=TargetConfig(type=target_system),
            data_sample=data_sample,
            operation_type="full_process",
        )
        return await process_data_request(req)

    except Exception as e:
        logger.error(f"Error analyzing file for user {user_id}: {e}")
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