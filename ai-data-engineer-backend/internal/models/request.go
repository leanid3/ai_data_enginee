package models

// FileUploadRequest запрос на загрузку файла
type FileUploadRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	FileType string `json:"file_type" binding:"required,oneof=csv json xml"`
	TargetDB string `json:"target_db" binding:"required,oneof=postgres clickhouse hdfs"`
}

// AnalysisRequest запрос на анализ файла
type AnalysisRequest struct {
	FileID   string `json:"file_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	FilePath string `json:"file_path" binding:"required"`
}

// PipelineRequest запрос на создание пайплайна
type PipelineRequest struct {
	AnalysisID string                 `json:"analysis_id" binding:"required"`
	UserID     string                 `json:"user_id" binding:"required"`
	Config     map[string]interface{} `json:"config"`
}

// ExecutePipelineRequest запрос на выполнение пайплайна
type ExecutePipelineRequest struct {
	PipelineID string                 `json:"pipeline_id" binding:"required"`
	UserID     string                 `json:"user_id" binding:"required"`
	Parameters map[string]interface{} `json:"parameters"`
}

// DatabaseTestRequest запрос на тестирование подключения к БД
type DatabaseTestRequest struct {
	Type     string `json:"type" binding:"required,oneof=postgres clickhouse"`
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password"`
	DBName   string `json:"dbname" binding:"required"`
	SSLMode  string `json:"sslmode,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
}
