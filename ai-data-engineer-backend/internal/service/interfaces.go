package service

import (
	"ai-data-engineer-backend/internal/models"
	"context"
)

// FileService интерфейс для работы с файлами
type FileService interface {
	AnalyzeFile(ctx context.Context, req *AnalyzeFileRequest) (*AnalyzeFileResponse, error)
	GetFileInfo(ctx context.Context, fileID string) (*models.FileMetadata, error)
	DeleteFile(ctx context.Context, fileID string) error
	ListFiles(ctx context.Context, userID string, limit, offset int) ([]*models.FileMetadata, error)
}

// PipelineService интерфейс для работы с пайплайнами
type PipelineService interface {
	CreatePipeline(ctx context.Context, req *CreatePipelineRequest) (*models.Pipeline, error)
	GetPipeline(ctx context.Context, pipelineID string) (*models.Pipeline, error)
	ExecutePipeline(ctx context.Context, req *ExecutePipelineRequest) (*models.PipelineExecution, error)
	DeletePipeline(ctx context.Context, pipelineID string) error
	ListPipelines(ctx context.Context, userID string, limit, offset int) ([]*models.Pipeline, error)
}

// AnalyzeService интерфейс для анализа данных
type AnalyzeService interface {
	StartAnalysis(ctx context.Context, req *StartAnalysisRequest) (*models.AnalysisResult, error)
	GetAnalysisStatus(ctx context.Context, analysisID string) (*models.AnalysisResult, error)
	GetAnalysisResult(ctx context.Context, analysisID string) (*models.AnalysisResult, error)
	ListAnalyses(ctx context.Context, userID string, limit, offset int) ([]*models.AnalysisResult, error)
}

// HealthService интерфейс для health checks
type HealthService interface {
	CheckHealth(ctx context.Context) (string, error)
	CheckDatabase(ctx context.Context) (string, error)
	CheckLLM(ctx context.Context) (string, error)
	TestDatabaseConnection(ctx context.Context, req *DatabaseTestRequest) (*DatabaseTestResponse, error)
}

// LLMClient интерфейс для работы с LLM
type LLMClient interface {
	ProcessRequest(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
	AnalyzeSchema(ctx context.Context, schema *models.DataSchema) (*models.AnalysisResult, error)
	GenerateDDL(ctx context.Context, req *GenerateDDLRequest) (*GenerateDDLResponse, error)
}

// FileProcessor интерфейс для обработки файлов
type FileProcessor interface {
	ParseCSV(ctx context.Context, content []byte) (*models.DataProfile, error)
	ParseJSON(ctx context.Context, content []byte) (*models.DataProfile, error)
	ParseXML(ctx context.Context, content []byte) (*models.DataProfile, error)
	DetectFileType(filename string) string
}

// DataAnalyzer интерфейс для анализа данных
type DataAnalyzer interface {
	AnalyzeDataStructure(ctx context.Context, profile *models.DataProfile) (*models.AnalysisResult, error)
	CalculateDataQuality(ctx context.Context, profile *models.DataProfile) float64
	GenerateRecommendations(ctx context.Context, profile *models.DataProfile) []string
}

// DDLGenerator интерфейс для генерации DDL
type DDLGenerator interface {
	GeneratePostgreSQLDDL(ctx context.Context, schema *models.TableSchema) (string, error)
	GenerateClickHouseDDL(ctx context.Context, schema *models.TableSchema) (string, error)
	GenerateHDFSSchema(ctx context.Context, schema *models.TableSchema) (string, error)
}

// PipelineBuilder интерфейс для построения пайплайнов
type PipelineBuilder interface {
	BuildETLPipeline(ctx context.Context, source *models.DataSource, target *models.DataTarget) (*models.Pipeline, error)
	ValidatePipeline(ctx context.Context, pipeline *models.Pipeline) error
	OptimizePipeline(ctx context.Context, pipeline *models.Pipeline) (*models.Pipeline, error)
}

// DAGGenerator интерфейс для генерации Airflow DAG
type DAGGenerator interface {
	GenerateDAG(ctx context.Context, pipeline *models.Pipeline) (string, error)
	SaveDAG(ctx context.Context, dagContent string, filename string) error
	ValidateDAG(ctx context.Context, dagContent string) error
}

// DatabaseConnector интерфейс для подключения к БД
type DatabaseConnector interface {
	TestConnection(ctx context.Context, config interface{}) error
	ExecuteQuery(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	GetTableSchema(ctx context.Context, tableName string) (*models.TableSchema, error)
	CreateTable(ctx context.Context, schema *models.TableSchema) error
	InsertData(ctx context.Context, tableName string, data []map[string]interface{}) error
}
