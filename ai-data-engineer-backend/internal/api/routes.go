package api

import (
	"ai-data-engineer-backend/internal/api/handlers"
	"ai-data-engineer-backend/internal/api/middleware"
	"ai-data-engineer-backend/internal/service"
	"ai-data-engineer-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты API
func SetupRoutes(
	fileService service.FileService,
	pipelineService service.PipelineService,
	analyzeService service.AnalyzeService,
	healthService service.HealthService,
	log logger.Logger,
) *gin.Engine {
	// Создаем Gin router
	r := gin.New()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.LoggingContext(log))
	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.Recovery(log))
	r.Use(middleware.ErrorHandler(log))

	// 404 и 405 handlers
	r.NoRoute(middleware.NotFoundHandler())
	r.NoMethod(middleware.MethodNotAllowedHandler())

	// Создаем handlers
	fileHandler := handlers.NewFileHandler(fileService, log)
	pipelineHandler := handlers.NewPipelineHandler(pipelineService, log)
	analyzeHandler := handlers.NewAnalyzeHandler(analyzeService, log)
	healthHandler := handlers.NewHealthHandler(healthService, log)

	// API v1 группа
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", healthHandler.HealthCheck)
		v1.POST("/databases/test", healthHandler.DatabaseTest)

		// File operations
		files := v1.Group("/files")
		{
			files.POST("/upload", fileHandler.UploadFile)
			files.GET("/:id", fileHandler.GetFileInfo)
			files.DELETE("/:id", fileHandler.DeleteFile)
			files.GET("", fileHandler.ListFiles)
		}

		// Analysis operations
		analysis := v1.Group("/analysis")
		{
			analysis.POST("/start", analyzeHandler.StartAnalysis)
			analysis.GET("/:id/status", analyzeHandler.GetAnalysisStatus)
			analysis.GET("/:id/result", analyzeHandler.GetAnalysisResult)
			analysis.GET("", analyzeHandler.ListAnalyses)
		}

		// Pipeline operations
		pipelines := v1.Group("/pipelines")
		{
			pipelines.POST("", pipelineHandler.CreatePipeline)
			pipelines.GET("/:id", pipelineHandler.GetPipeline)
			pipelines.POST("/:id/execute", pipelineHandler.ExecutePipeline)
			pipelines.DELETE("/:id", pipelineHandler.DeletePipeline)
			pipelines.GET("", pipelineHandler.ListPipelines)
		}
	}

	return r
}

