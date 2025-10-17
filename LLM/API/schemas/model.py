from typing import Any, Dict, List, Optional
from pydantic import BaseModel


class SourceConfig(BaseModel):
    type: str
    connection_string: Optional[str] = None
    file_path: Optional[str] = None
    table_name: Optional[str] = None
    query: Optional[str] = None
    credentials: Optional[Dict[str, Any]] = None


class TargetConfig(BaseModel):
    type: str
    connection_string: Optional[str] = None
    table_name: Optional[str] = None
    database_name: Optional[str] = None
    credentials: Optional[Dict[str, Any]] = None


class DataSample(BaseModel):
    headers: List[str]
    sample_rows: List[Dict[str, Any]]
    total_rows: int
    file_size: Optional[int] = None


class LLMRequest(BaseModel):
    user_query: str
    source_config: SourceConfig
    target_config: TargetConfig
    data_sample: Optional[DataSample] = None
    metadata: Optional[Dict[str, Any]] = None
    operation_type: str
    context: Optional[Dict[str, Any]] = None
    conversation_id: Optional[str] = None


class LLMResponse(BaseModel):
    success: bool
    pipeline_id: str
    data_analysis: Optional[Dict[str, Any]] = None
    storage_recommendation: Optional[Dict[str, Any]] = None
    ddl_scripts: Optional[List[Dict[str, Any]]] = None
    dag_code: Optional[str] = None
    optimized_queries: Optional[List[Dict[str, Any]]] = None
    visualization_config: Optional[Dict[str, Any]] = None
    user_report: str
    processing_time: float
    agents_used: List[str]
    tools_used: List[str] = []
    confidence_score: float
    errors: List[str] = []
    warnings: List[str] = []
