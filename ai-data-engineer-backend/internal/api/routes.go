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
	fileService handlers.FileService,
	dataAnalyzer *service.DataAnalyzer,
	pipelineService *service.PipelineService,
	healthService handlers.HealthService,
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
