package service

import (
	"ai-data-engineer-backend/internal/models"
	"io"
)

// AnalyzeFileRequest запрос на анализ файла
type AnalyzeFileRequest struct {
	File     io.Reader
	Filename string
	UserID   string
	FileType string
	TargetDB string
}

// AnalyzeFileResponse ответ на анализ файла
type AnalyzeFileResponse struct {
	FileID     string
	AnalysisID string
	Status     string
	Message    string
}

// CreatePipelineRequest запрос на создание пайплайна
type CreatePipelineRequest struct {
	AnalysisID string
	UserID     string
	Config     map[string]interface{}
}

// ExecutePipelineRequest запрос на выполнение пайплайна
type ExecutePipelineRequest struct {
	PipelineID string
	UserID     string
	Parameters map[string]interface{}
}

// StartAnalysisRequest запрос на запуск анализа
type StartAnalysisRequest struct {
	FileID   string
	UserID   string
	FilePath string
}

// DatabaseTestRequest запрос на тестирование БД
type DatabaseTestRequest struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Secure   bool
}

// DatabaseTestResponse ответ на тестирование БД
type DatabaseTestResponse struct {
	Status    string
	Message   string
	Connected bool
	TestedAt  string
	Details   map[string]interface{}
}

// LLMRequest запрос к LLM
type LLMRequest struct {
	UserQuery     string                 `json:"user_query"`
	SourceConfig  map[string]interface{} `json:"source_config"`
	TargetConfig  map[string]interface{} `json:"target_config"`
	OperationType string                 `json:"operation_type"`
	DataSample    string                 `json:"data_sample,omitempty"`
}

// LLMResponse ответ от LLM
type LLMResponse struct {
	PipelineID string    `json:"pipeline_id"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	Error      *LLMError `json:"error,omitempty"`
}

// LLMError ошибка LLM
type LLMError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// GenerateDDLRequest запрос на генерацию DDL
type GenerateDDLRequest struct {
	DataProfile models.DataProfile `json:"data_profile"`
	Target      models.DataTarget  `json:"target"`
}

// GenerateDDLResponse ответ на генерацию DDL
type GenerateDDLResponse struct {
	DDLScript   string `json:"ddl_script"`
	Status      string `json:"status"`
	Explanation string `json:"explanation"`
}

