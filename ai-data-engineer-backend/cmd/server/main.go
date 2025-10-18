package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-data-engineer-backend/internal/api"
	"ai-data-engineer-backend/internal/config"
	"ai-data-engineer-backend/internal/repository"
	"ai-data-engineer-backend/internal/service"
	"ai-data-engineer-backend/pkg/logger"
)

func main() {
	// Загружаем конфигурацию
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализируем логгер
	logger := logger.NewLogger(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.Output)
	logger.Info("Starting AI Data Engineer Backend")

	// Инициализируем репозитории
	repositories, err := initializeRepositories(cfg, logger)
	if err != nil {
		logger.Fatalf("Ошибка инициализации репозиториев: %v", err)
	}

	// Инициализируем сервисы
	services, err := initializeServices(cfg, logger, repositories)
	if err != nil {
		logger.Fatalf("Ошибка инициализации сервисов: %v", err)
	}

	// Настраиваем маршруты
	router := api.SetupRoutes(
		services.FileService,
		services.PipelineService,
		services.AnalyzeService,
		services.HealthService,
		logger,
	)

	// Проверяем переменную окружения PORT
	port := cfg.Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Запускаем сервер в горутине
	go func() {
		logger.Infof("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

// Repositories содержит все репозитории
type Repositories struct {
	Pipeline  repository.PipelineRepository
	File      repository.FileRepository
	Analysis  repository.AnalysisRepository
	Execution repository.ExecutionRepository
	Database  repository.DatabaseRepository
}

// Services содержит все сервисы
type Services struct {
	FileService     service.FileService
	PipelineService service.PipelineService
	AnalyzeService  service.AnalyzeService
	HealthService   service.HealthService
}

// initializeRepositories инициализирует репозитории
func initializeRepositories(cfg *config.Config, logger logger.Logger) (*Repositories, error) {
	// Заглушки для репозиториев - будут реализованы в следующих итерациях
	logger.Info("Initializing repositories (stub implementation)")

	return &Repositories{
		// Pipeline:  repository.NewPostgreSQLPipelineRepository(cfg, logger),
		// File:      repository.NewPostgreSQLFileRepository(cfg, logger),
		// Analysis:  repository.NewPostgreSQLAnalysisRepository(cfg, logger),
		// Execution: repository.NewPostgreSQLExecutionRepository(cfg, logger),
		// Database:  repository.NewDatabaseRepository(cfg, logger),
	}, nil
}

// initializeServices инициализирует сервисы
func initializeServices(cfg *config.Config, logger logger.Logger, repos *Repositories) (*Services, error) {
	logger.Info("Initializing services with real implementations")

	// Создаем LLM клиент
	llmClient := service.NewLLMClient(cfg.LLM.BaseURL, cfg.LLM.APIKey, logger)

	// Создаем MinIO клиент
	minioClient, err := service.NewMinIOClient(
		cfg.Storage.Endpoint,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.UseSSL,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Создаем файловый процессор
	fileProcessor := service.NewFileProcessor(logger)

	// Создаем анализатор данных
	dataAnalyzer := service.NewDataAnalyzer(logger)

	// Создаем сервисы с зависимостями
	fileService := service.NewFileService(fileProcessor, dataAnalyzer, llmClient, minioClient, logger)

	return &Services{
		FileService:     fileService,
		PipelineService: service.NewPipelineService(logger),
		AnalyzeService:  service.NewAnalyzeService(logger),
		HealthService:   service.NewHealthService(logger),
	}, nil
}
