package models

import (
	"time"
)

// Pipeline представляет пайплайн данных
type Pipeline struct {
	ID          string                 `json:"id" gorm:"primaryKey"`
	UserID      string                 `json:"user_id" gorm:"index"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      PipelineStatus         `json:"status"`
	Config      map[string]interface{} `json:"config" gorm:"type:jsonb"`
	Source      DataSource             `json:"source" gorm:"type:jsonb"`
	Target      DataTarget             `json:"target" gorm:"type:jsonb"`
	Steps       []PipelineStep         `json:"steps" gorm:"type:jsonb"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	ExecutedAt  *time.Time             `json:"executed_at,omitempty"`
}

// PipelineStatus статус пайплайна
type PipelineStatus string

const (
	PipelineStatusDraft     PipelineStatus = "draft"
	PipelineStatusReady     PipelineStatus = "ready"
	PipelineStatusRunning   PipelineStatus = "running"
	PipelineStatusCompleted PipelineStatus = "completed"
	PipelineStatusFailed    PipelineStatus = "failed"
	PipelineStatusCancelled PipelineStatus = "cancelled"
)

// DataSource источник данных
type DataSource struct {
	Type   string                 `json:"type"`
	Path   string                 `json:"path,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
	Schema DataSchema             `json:"schema,omitempty"`
}

// DataTarget целевая система
type DataTarget struct {
	Type             string                 `json:"type"`
	ConnectionString string                 `json:"connection_string,omitempty"`
	TableName        string                 `json:"table_name,omitempty"`
	Schema           string                 `json:"schema,omitempty"`
	Config           map[string]interface{} `json:"config,omitempty"`
}

// PipelineStep шаг пайплайна
type PipelineStep struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        StepType               `json:"type"`
	Config      map[string]interface{} `json:"config"`
	DependsOn   []string               `json:"depends_on,omitempty"`
	Status      StepStatus             `json:"status"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// StepType тип шага пайплайна
type StepType string

const (
	StepTypeExtract   StepType = "extract"
	StepTypeTransform StepType = "transform"
	StepTypeLoad      StepType = "load"
	StepTypeValidate  StepType = "validate"
)

// StepStatus статус шага
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped"
)

// DataSchema схема данных
type DataSchema struct {
	Fields []DataField              `json:"fields"`
	Sample []map[string]interface{} `json:"sample,omitempty"`
}

// DataField поле данных
type DataField struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Nullable    bool    `json:"nullable"`
	NullCount   int     `json:"null_count,omitempty"`
	SampleValue string  `json:"sample_value,omitempty"`
	MinValue    float64 `json:"min_value,omitempty"`
	MaxValue    float64 `json:"max_value,omitempty"`
	Description string  `json:"description,omitempty"`
}

// PipelineExecution выполнение пайплайна
type PipelineExecution struct {
	ID          string                 `json:"id" gorm:"primaryKey"`
	PipelineID  string                 `json:"pipeline_id" gorm:"index"`
	UserID      string                 `json:"user_id" gorm:"index"`
	Status      ExecutionStatus        `json:"status"`
	Parameters  map[string]interface{} `json:"parameters" gorm:"type:jsonb"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Logs        []ExecutionLog         `json:"logs" gorm:"type:jsonb"`
}

// ExecutionStatus статус выполнения
type ExecutionStatus string

const (
	ExecutionStatusScheduled ExecutionStatus = "scheduled"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

// ExecutionLog лог выполнения
type ExecutionLog struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	StepID    string    `json:"step_id,omitempty"`
}

