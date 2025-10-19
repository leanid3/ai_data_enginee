package models

import (
	"time"
)

// DataField представляет поле данных
type DataField struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Nullable    bool    `json:"nullable"`
	NullCount   int     `json:"null_count"`
	SampleValue string  `json:"sample_value"`
	MinValue    float64 `json:"min_value"`
	MaxValue    float64 `json:"max_value"`
	Description string  `json:"description"`
}

// DataProfile профиль данных
type DataProfile struct {
	DataType         string      `json:"data_type"`
	TotalRows        int         `json:"total_rows"`
	SampledRows      int         `json:"sampled_rows"`
	Fields           []DataField `json:"fields"`
	SampleData       string      `json:"sample_data"`
	DataQualityScore float64     `json:"data_quality_score"`
	FileSize         int64       `json:"file_size"`
	Encoding         string      `json:"encoding"`
	Delimiter        string      `json:"delimiter,omitempty"`
	HasHeaders       bool        `json:"has_headers"`
	CreatedAt        time.Time   `json:"created_at"`
}

// AnalysisResult результат анализа
type AnalysisResult struct {
	AnalysisID            string                `json:"analysis_id"`
	FileID                string                `json:"file_id"`
	UserID                string                `json:"user_id"`
	DataProfile           DataProfile           `json:"data_profile"`
	Recommendations       []string              `json:"recommendations"`
	StorageRecommendation StorageRecommendation `json:"storage_recommendation"`
	TableSchema           TableSchema           `json:"table_schema"`
	DDLMetadata           DDLMetadata           `json:"ddl_metadata"`
	LLMAnalysis           string                `json:"llm_analysis"`
	Status                AnalysisStatus        `json:"status"`
	CreatedAt             time.Time             `json:"created_at"`
	CompletedAt           *time.Time            `json:"completed_at,omitempty"`
	Error                 string                `json:"error,omitempty"`
}

// AnalysisStatus статус анализа
type AnalysisStatus string

const (
	AnalysisStatusPending   AnalysisStatus = "pending"
	AnalysisStatusRunning   AnalysisStatus = "running"
	AnalysisStatusCompleted AnalysisStatus = "completed"
	AnalysisStatusFailed    AnalysisStatus = "failed"
)

// StorageRecommendation рекомендация по хранилищу
type StorageRecommendation struct {
	PrimaryStorage   string                 `json:"primary_storage"`
	SecondaryStorage []string               `json:"secondary_storage"`
	Reasoning        map[string]interface{} `json:"reasoning"`
	StorageOptions   map[string]interface{} `json:"storage_options"`
}

// TableSchema схема таблицы
type TableSchema struct {
	TableName   string            `json:"table_name"`
	Fields      []TableField      `json:"fields"`
	PrimaryKey  []string          `json:"primary_key"`
	Indexes     []TableIndex      `json:"indexes"`
	Constraints []TableConstraint `json:"constraints"`
}

// TableField поле таблицы
type TableField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Nullable    bool   `json:"nullable"`
	Indexed     bool   `json:"indexed"`
	Description string `json:"description,omitempty"`
}

// TableIndex индекс таблицы
type TableIndex struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
	Type   string   `json:"type,omitempty"`
}

// TableConstraint ограничение таблицы
type TableConstraint struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Expression string `json:"expression,omitempty"`
}

// DDLMetadata метаданные DDL
type DDLMetadata struct {
	DDLGeneration       map[string]interface{} `json:"ddl_generation"`
	DataCharacteristics map[string]interface{} `json:"data_characteristics"`
}

// FileMetadata метаданные файла
type FileMetadata struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	UserID      string     `json:"user_id" gorm:"index"`
	Filename    string     `json:"filename"`
	ContentType string     `json:"content_type"`
	Size        int64      `json:"size"`
	Path        string     `json:"path"`
	Bucket      string     `json:"bucket"`
	Status      FileStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// FileStatus статус файла
type FileStatus string

const (
	FileStatusUploading  FileStatus = "uploading"
	FileStatusUploaded   FileStatus = "uploaded"
	FileStatusProcessing FileStatus = "processing"
	FileStatusProcessed  FileStatus = "processed"
	FileStatusError      FileStatus = "error"
	FileStatusDeleted    FileStatus = "deleted"
)
