package handlers

import (
	"net/http"
	"strconv"
	"time"

	"ai-data-engineer-backend/internal/models"
	"ai-data-engineer-backend/internal/service"
	"ai-data-engineer-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// PipelineHandler обработчик для работы с пайплайнами
type PipelineHandler struct {
	pipelineService service.PipelineService
	logger          logger.Logger
}

// NewPipelineHandler создает новый PipelineHandler
func NewPipelineHandler(pipelineService service.PipelineService, logger logger.Logger) *PipelineHandler {
	return &PipelineHandler{
		pipelineService: pipelineService,
		logger:          logger,
	}
}

// CreatePipeline создает новый пайплайн
func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	requestLogger.Info("Creating pipeline")

	var req models.PipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		requestLogger.WithField("error", err.Error()).Warn("Invalid request body")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "validation_error",
			Message:   "Неверный формат запроса",
			Timestamp: time.Now(),
		})
		return
	}

	pipeline, err := h.pipelineService.CreatePipeline(c.Request.Context(), &service.CreatePipelineRequest{
		AnalysisID: req.AnalysisID,
		UserID:     req.UserID,
		Config:     req.Config,
	})
	if err != nil {
		requestLogger.WithField("error", err.Error()).Error("Failed to create pipeline")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "creation_failed",
			Message:   "Ошибка создания пайплайна",
			Timestamp: time.Now(),
		})
		return
	}

	response := models.PipelineResponse{
		PipelineID: pipeline.ID,
		Status:     string(pipeline.Status),
		Message:    "Пайплайн успешно создан",
		Config:     pipeline.Config,
		CreatedAt:  pipeline.CreatedAt,
	}

	requestLogger.WithField("pipeline_id", pipeline.ID).Info("Pipeline created successfully")
	c.JSON(http.StatusCreated, response)
}

// GetPipeline получает пайплайн по ID
func (h *PipelineHandler) GetPipeline(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	pipelineID := c.Param("id")

	if pipelineID == "" {
		requestLogger.Warn("Missing pipeline ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID пайплайна обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	pipeline, err := h.pipelineService.GetPipeline(c.Request.Context(), pipelineID)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("pipeline_id", pipelineID).Error("Failed to get pipeline")
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:     "pipeline_not_found",
			Message:   "Пайплайн не найден",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("pipeline_id", pipelineID).Info("Pipeline retrieved")
	c.JSON(http.StatusOK, pipeline)
}

// ExecutePipeline выполняет пайплайн
func (h *PipelineHandler) ExecutePipeline(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	pipelineID := c.Param("id")

	if pipelineID == "" {
		requestLogger.Warn("Missing pipeline ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID пайплайна обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	var req models.ExecutePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		requestLogger.WithField("error", err.Error()).Warn("Invalid request body")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "validation_error",
			Message:   "Неверный формат запроса",
			Timestamp: time.Now(),
		})
		return
	}

	execution, err := h.pipelineService.ExecutePipeline(c.Request.Context(), &service.ExecutePipelineRequest{
		PipelineID: pipelineID,
		UserID:     req.UserID,
		Parameters: req.Parameters,
	})
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("pipeline_id", pipelineID).Error("Failed to execute pipeline")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "execution_failed",
			Message:   "Ошибка выполнения пайплайна",
			Timestamp: time.Now(),
		})
		return
	}

	response := models.ExecutePipelineResponse{
		ExecutionID: execution.ID,
		Status:      string(execution.Status),
		Message:     "Пайплайн запущен на выполнение",
		Parameters:  execution.Parameters,
		StartedAt:   execution.StartedAt,
	}

	requestLogger.WithField("pipeline_id", pipelineID).WithField("execution_id", execution.ID).Info("Pipeline execution started")
	c.JSON(http.StatusOK, response)
}

// DeletePipeline удаляет пайплайн
func (h *PipelineHandler) DeletePipeline(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	pipelineID := c.Param("id")

	if pipelineID == "" {
		requestLogger.Warn("Missing pipeline ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID пайплайна обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	err := h.pipelineService.DeletePipeline(c.Request.Context(), pipelineID)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("pipeline_id", pipelineID).Error("Failed to delete pipeline")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "delete_failed",
			Message:   "Ошибка удаления пайплайна",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("pipeline_id", pipelineID).Info("Pipeline deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":     "Пайплайн успешно удален",
		"pipeline_id": pipelineID,
	})
}

// ListPipelines получает список пайплайнов пользователя
func (h *PipelineHandler) ListPipelines(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	userID := c.Query("user_id")

	if userID == "" {
		requestLogger.Warn("Missing user ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID пользователя обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	// Параметры пагинации
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	pipelines, err := h.pipelineService.ListPipelines(c.Request.Context(), userID, limit, offset)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("user_id", userID).Error("Failed to list pipelines")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "list_failed",
			Message:   "Ошибка получения списка пайплайнов",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("user_id", userID).WithField("count", len(pipelines)).Info("Pipelines listed")
	c.JSON(http.StatusOK, gin.H{
		"pipelines": pipelines,
		"limit":     limit,
		"offset":    offset,
		"count":     len(pipelines),
	})
}
