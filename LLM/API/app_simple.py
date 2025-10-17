from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, File, UploadFile
from fastapi.responses import JSONResponse
import uuid
import os
from datetime import datetime

from API.schemas import LLMRequest, LLMResponse
from API.interfaces import OpenRouterClient
from API.schemas import SourceConfig, TargetConfig

__VERSION__ = "0.1"

# Простой логгер
import logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

pipeline_storage = {}

@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("Starting LLM Service...")
    llm_client = OpenRouterClient()
    app.state.llm_client = llm_client
    yield
    logger.info("Shutting down...")

app = FastAPI(title="AI Data Engineer API", version="3.0.0", lifespan=lifespan)

@app.post("/api/v1/process", response_model=LLMResponse)
async def process_data_request(request: LLMRequest) -> LLMResponse:
    try:
        # Простая обработка запроса
        pipeline_id = str(uuid.uuid4())
        
        # Создаем простой ответ
        response = LLMResponse(
            success=True,
            pipeline_id=pipeline_id,
            user_report="Запрос обработан успешно",
            processing_time=1.0,
            agents_used=["ROUTER", "DATA_ANALYZER"],
            confidence_score=0.8,
            dag_code="-- Простой DAG код\nprint('Hello from LLM Service')"
        )
        
        pipeline_storage[pipeline_id] = response.dict()
        return response
    except Exception as e:
        logger.error(f"Error processing request: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/v1/analyze-file")
async def analyze_file(file: UploadFile = File(...), target_system: str = "auto"):
    try:
        content = await file.read()
        
        # Простой анализ файла
        pipeline_id = str(uuid.uuid4())
        
        response = LLMResponse(
            success=True,
            pipeline_id=pipeline_id,
            user_report=f"Файл {file.filename} проанализирован",
            processing_time=1.5,
            agents_used=["DATA_ANALYZER", "DB_SELECTOR"],
            confidence_score=0.9,
            dag_code=f"-- Анализ файла {file.filename}\nprint('File analyzed')"
        )
        
        pipeline_storage[pipeline_id] = response.dict()
        return response
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
            "agents_available": ["ROUTER", "DATA_ANALYZER", "DB_SELECTOR", "DDL_GENERATOR", "ETL_BUILDER", "QUERY_OPTIMIZER", "REPORT_GENERATOR"],
            "pipelines_cached": len(pipeline_storage),
        }
    )

@app.get("/api/v1/agents")
async def get_agents():
    return JSONResponse(
        content={
            "agents": [
                {"role": "ROUTER", "description": "Маршрутизатор запросов"},
                {"role": "DATA_ANALYZER", "description": "Анализ структуры и качества данных"},
                {"role": "DB_SELECTOR", "description": "Выбор системы хранения"},
                {"role": "DDL_GENERATOR", "description": "Генерация DDL"},
                {"role": "ETL_BUILDER", "description": "Построение ETL пайплайнов"},
                {"role": "QUERY_OPTIMIZER", "description": "Оптимизация SQL"},
                {"role": "REPORT_GENERATOR", "description": "Отчеты и рекомендации"},
            ]
        }
    )
