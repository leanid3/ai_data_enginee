package repository

import (
	"ai-data-engineer-backend/domain/models"
	"context"
)

// PipelineRepository интерфейс для работы с пайплайнами
type PipelineRepository interface {
	SavePipeline(ctx context.Context, pipeline *models.Pipeline) (string, error)
	GetPipeline(ctx context.Context, id string) (*models.Pipeline, error)
	GetPipelinesByUser(ctx context.Context, userID string, limit, offset int) ([]*models.Pipeline, error)
	UpdatePipeline(ctx context.Context, pipeline *models.Pipeline) error
	DeletePipeline(ctx context.Context, id string) error
	GetPipelinesByStatus(ctx context.Context, status models.PipelineStatus) ([]*models.Pipeline, error)
}

// FileRepository интерфейс для работы с файлами
type FileRepository interface {
	SaveFile(ctx context.Context, file *models.FileMetadata) error
	GetFile(ctx context.Context, id string) (*models.FileMetadata, error)
	GetFilesByUser(ctx context.Context, userID string, limit, offset int) ([]*models.FileMetadata, error)
	UpdateFile(ctx context.Context, file *models.FileMetadata) error
	DeleteFile(ctx context.Context, id string) error
	GetFilesByStatus(ctx context.Context, status models.FileStatus) ([]*models.FileMetadata, error)
}

// AnalysisRepository интерфейс для работы с анализами
type AnalysisRepository interface {
	SaveAnalysis(ctx context.Context, analysis *models.AnalysisResult) error
	GetAnalysis(ctx context.Context, id string) (*models.AnalysisResult, error)
	GetAnalysesByUser(ctx context.Context, userID string, limit, offset int) ([]*models.AnalysisResult, error)
	UpdateAnalysis(ctx context.Context, analysis *models.AnalysisResult) error
	DeleteAnalysis(ctx context.Context, id string) error
	GetAnalysesByStatus(ctx context.Context, status models.AnalysisStatus) ([]*models.AnalysisResult, error)
}

// ExecutionRepository интерфейс для работы с выполнениями пайплайнов
type ExecutionRepository interface {
	SaveExecution(ctx context.Context, execution *models.PipelineExecution) error
	GetExecution(ctx context.Context, id string) (*models.PipelineExecution, error)
	GetExecutionsByPipeline(ctx context.Context, pipelineID string, limit, offset int) ([]*models.PipelineExecution, error)
	UpdateExecution(ctx context.Context, execution *models.PipelineExecution) error
	GetExecutionsByStatus(ctx context.Context, status models.ExecutionStatus) ([]*models.PipelineExecution, error)
}

// DatabaseRepository интерфейс для работы с базами данных
type DatabaseRepository interface {
	TestConnection(ctx context.Context, config interface{}) error
	ExecuteQuery(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	GetTableSchema(ctx context.Context, tableName string) (*models.TableSchema, error)
	CreateTable(ctx context.Context, schema *models.TableSchema) error
	InsertData(ctx context.Context, tableName string, data []map[string]interface{}) error
}
