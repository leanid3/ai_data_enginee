package handlers

import (
	"net/http"
	"time"

	"ai-data-engineer-backend/internal/models"
	"ai-data-engineer-backend/internal/service"
	"ai-data-engineer-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// HealthHandler обработчик для health checks
type HealthHandler struct {
	healthService service.HealthService
	logger        logger.Logger
}

// NewHealthHandler создает новый HealthHandler
func NewHealthHandler(healthService service.HealthService, logger logger.Logger) *HealthHandler {
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

	result, err := h.healthService.TestDatabaseConnection(c.Request.Context(), &service.DatabaseTestRequest{
		Type:     req.Type,
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		DBName:   req.DBName,
		SSLMode:  req.SSLMode,
		Secure:   req.Secure,
	})
	if err != nil {
		requestLogger.WithField("error", err.Error()).Error("Database test failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "test_failed",
			Message:   "Ошибка тестирования подключения к БД",
			Timestamp: time.Now(),
		})
		return
	}

	response := models.DatabaseTestResponse{
		Status:    result.Status,
		Message:   result.Message,
		Connected: result.Connected,
		TestedAt:  time.Now(),
		Details:   result.Details,
	}

	requestLogger.WithField("type", req.Type).WithField("connected", result.Connected).Info("Database test completed")
	c.JSON(http.StatusOK, response)
}
