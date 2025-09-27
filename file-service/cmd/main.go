package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"file-service/config"
	"file-service/internal/handlers"
	"file-service/internal/models"
	"file-service/internal/storage"
	"file-service/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	db, err := connectToDatabase(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Автомиграция моделей
	if err := db.AutoMigrate(&models.File{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	// Создание хранилища
	storageConfig := storage.StorageConfig{
		Type:      cfg.Storage.Type,
		Endpoint:  cfg.Storage.Endpoint,
		AccessKey: cfg.Storage.AccessKey,
		SecretKey: cfg.Storage.SecretKey,
		Bucket:    cfg.Storage.Bucket,
		UseSSL:    cfg.Storage.UseSSL,
	}

	var fileStorage storage.StorageInterface
	switch cfg.Storage.Type {
	case "minio":
		fileStorage, err = storage.NewMinIOStorage(storageConfig)
		if err != nil {
			log.Fatalf("Ошибка создания MinIO клиента: %v", err)
		}
	case "s3":
		fileStorage, err = storage.NewS3Storage(storageConfig)
		if err != nil {
			log.Fatalf("Ошибка создания S3 клиента: %v", err)
		}
	default:
		log.Fatalf("Неподдерживаемый тип хранилища: %s", cfg.Storage.Type)
	}

	// Создание bucket если не существует
	ctx := context.Background()
	if err := fileStorage.CreateBucket(ctx, cfg.Storage.Bucket); err != nil {
		log.Printf("Ошибка создания bucket: %v", err)
	}

	// Создание handler
	fileHandler := handlers.NewFileHandler(db, fileStorage)

	// Запуск gRPC сервера
	go startGRPCServer(fileHandler, cfg)

	// Запуск HTTP сервера
	startHTTPServer(fileHandler, cfg)
}

func connectToDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Настройка пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func startGRPCServer(handler *handlers.FileHandler, cfg *config.Config) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Ошибка создания gRPC listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterFileServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("gRPC сервер запущен на порту %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}

func startHTTPServer(handler *handlers.FileHandler, cfg *config.Config) {
	// Настройка Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middleware для CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "file-service"})
	})

	// Регистрация маршрутов
	handler.RegisterHTTPRoutes(r)

	// Запуск HTTP сервера
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	log.Printf("HTTP сервер запущен на порту %s", httpPort)
	log.Fatal(r.Run(":" + httpPort))
}
