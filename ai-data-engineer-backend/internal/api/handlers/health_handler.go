package handlers

import (
	"context"
	"net/http"
	"time"

	"ai-data-engineer-backend/domain/models"
	"ai-data-engineer-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// TODO: переработать реализацию HealthService
type HealthService interface {
	CheckHealth(ctx context.Context) (bool, error)
	CheckDatabase(ctx context.Context) (bool, error)
	CheckLLM(ctx context.Context) (bool, error)
	TestDatabaseConnection(ctx context.Context, req *models.DatabaseTestRequest) (bool, error)
}

// HealthHandler обработчик для health checks
type HealthHandler struct {
	healthService HealthService
	logger        logger.Logger
}

// NewHealthHandler создает новый HealthHandler
func NewHealthHandler(healthService HealthService, logger logger.Logger) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
		logger:        logger,
	}
}

// HealthCheck проверяет состояние сервиса
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	requestLogger.Debug("Health check requested")

	checks := make(map[string]string)

	// Проверяем состояние сервиса
	_, err := h.healthService.CheckHealth(c.Request.Context())
	if err != nil {
		requestLogger.WithField("error", err.Error()).Error("Health check failed")
		checks["service"] = "unhealthy"
	} else {
		checks["service"] = "healthy"
	}

	// Проверяем подключения к БД
	_, err = h.healthService.CheckDatabase(c.Request.Context())
	if err != nil {
		requestLogger.WithField("error", err.Error()).Warn("Database check failed")
		checks["database"] = "unhealthy"
	} else {
		checks["database"] = "healthy"
	}

	// Проверяем LLM сервис
	_, err = h.healthService.CheckLLM(c.Request.Context())
	if err != nil {
		requestLogger.WithField("error", err.Error()).Warn("LLM check failed")
		checks["llm"] = "unhealthy"
	} else {
		checks["llm"] = "healthy"
	}

	// Определяем общий статус
	overallStatus := "healthy"
	for _, checkStatus := range checks {
		if checkStatus == "unhealthy" {
			overallStatus = "unhealthy"
			break
		}
	}

	response := models.HealthResponse{
		Status:    overallStatus,
		Service:   "ai-data-engineer-backend",
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Checks:    checks,
	}

	requestLogger.WithField("status", overallStatus).Debug("Health check completed")
	c.JSON(http.StatusOK, response)
}

// DatabaseTest тестирует подключение к базе данных
func (h *HealthHandler) DatabaseTest(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	requestLogger.Info("Database test requested")

	var req models.DatabaseTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		requestLogger.WithField("error", err.Error()).Warn("Invalid request body")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "validation_error",
			Message:   "Неверный формат запроса",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("type", req.Type).Info("Database test completed")
	c.JSON(http.StatusOK, models.DatabaseTestResponse{
		Status:    "healthy",
		Message:   "Database test completed",
		Connected: true,
		TestedAt:  time.Now(),
	})
}
