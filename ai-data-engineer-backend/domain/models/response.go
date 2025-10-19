package models

import (
	"time"
)

// FileUploadResponse ответ на загрузку файла
type FileUploadResponse struct {
	FileID    string    `json:"file_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	UploadURL string    `json:"upload_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// AnalysisResponse ответ на анализ файла
type AnalysisResponse struct {
	AnalysisID string                 `json:"analysis_id"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Result     map[string]interface{} `json:"result,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// PipelineResponse ответ на создание пайплайна
type PipelineResponse struct {
	PipelineID string                 `json:"pipeline_id"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Config     map[string]interface{} `json:"config,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// ExecutePipelineResponse ответ на выполнение пайплайна
type ExecutePipelineResponse struct {
	ExecutionID string                 `json:"execution_id"`
	Status      string                 `json:"status"`
	Message     string                 `json:"message"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
}

// HealthResponse ответ на health check
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// DatabaseTestResponse ответ на тестирование БД
type DatabaseTestResponse struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Connected bool                   `json:"connected"`
	TestedAt  time.Time              `json:"tested_at"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// ErrorResponse стандартизированный ответ об ошибке
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// LLMRequest запрос к LLM
type LLMRequest struct {
	UserQuery     string                 `json:"user_query"`
	SourceConfig  map[string]interface{} `json:"source_config,omitempty"`
	TargetConfig  map[string]interface{} `json:"target_config,omitempty"`
	OperationType string                 `json:"operation_type"`
	Prompt        string                 `json:"prompt,omitempty"`
	Model         string                 `json:"model,omitempty"`
	MaxTokens     int                    `json:"max_tokens,omitempty"`
	Temperature   float64                `json:"temperature,omitempty"`
	Context       map[string]interface{} `json:"context,omitempty"`
	DataProfile   *DataProfile           `json:"data_profile,omitempty"`
}

// LLMResponse ответ от LLM
type LLMResponse struct {
	Content    string                 `json:"content"`
	Model      string                 `json:"model"`
	Tokens     int                    `json:"tokens,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	PipelineID string                 `json:"pipeline_id,omitempty"`
	Message    string                 `json:"message,omitempty"`
	Status     string                 `json:"status,omitempty"`
	UserReport string                 `json:"user_report,omitempty"`
	Error      *string                `json:"error,omitempty"`
}

// GenerateDDLRequest запрос на генерацию DDL
type GenerateDDLRequest struct {
	Schema      *DataSchema            `json:"schema"`
	Database    string                 `json:"database"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Target      *TargetConfig          `json:"target"`
	DataProfile *DataProfile           `json:"data_profile"`
}

// TargetConfig конфигурация целевой системы
type TargetConfig struct {
	Type      string `json:"type"`
	TableName string `json:"table_name"`
}

// GenerateDDLResponse ответ с DDL
type GenerateDDLResponse struct {
	DDL      string                 `json:"ddl"`
	Database string                 `json:"database"`
	Tables   []string               `json:"tables,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
