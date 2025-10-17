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

// AnalyzeHandler обработчик для анализа данных
type AnalyzeHandler struct {
	analyzeService service.AnalyzeService
	logger         logger.Logger
}

// NewAnalyzeHandler создает новый AnalyzeHandler
func NewAnalyzeHandler(analyzeService service.AnalyzeService, logger logger.Logger) *AnalyzeHandler {
	return &AnalyzeHandler{
		analyzeService: analyzeService,
		logger:         logger,
	}
}

// StartAnalysis запускает анализ файла
func (h *AnalyzeHandler) StartAnalysis(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	requestLogger.Info("Starting analysis")

	var req models.AnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		requestLogger.WithField("error", err.Error()).Warn("Invalid request body")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "validation_error",
			Message:   "Неверный формат запроса",
			Timestamp: time.Now(),
		})
		return
	}

	analysis, err := h.analyzeService.StartAnalysis(c.Request.Context(), &service.StartAnalysisRequest{
		FileID:   req.FileID,
		UserID:   req.UserID,
		FilePath: req.FilePath,
	})
	if err != nil {
		requestLogger.WithField("error", err.Error()).Error("Failed to start analysis")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "analysis_failed",
			Message:   "Ошибка запуска анализа",
			Timestamp: time.Now(),
		})
		return
	}

	response := models.AnalysisResponse{
		AnalysisID: analysis.AnalysisID,
		Status:     string(analysis.Status),
		Message:    "Анализ запущен",
		Result:     map[string]interface{}{"analysis_id": analysis.AnalysisID},
		CreatedAt:  analysis.CreatedAt,
	}

	requestLogger.WithField("analysis_id", analysis.AnalysisID).Info("Analysis started")
	c.JSON(http.StatusOK, response)
}

// GetAnalysisStatus получает статус анализа
func (h *AnalyzeHandler) GetAnalysisStatus(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	analysisID := c.Param("id")

	if analysisID == "" {
		requestLogger.Warn("Missing analysis ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID анализа обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	analysis, err := h.analyzeService.GetAnalysisStatus(c.Request.Context(), analysisID)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("analysis_id", analysisID).Error("Failed to get analysis status")
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:     "analysis_not_found",
			Message:   "Анализ не найден",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("analysis_id", analysisID).Info("Analysis status retrieved")
	c.JSON(http.StatusOK, analysis)
}

// GetAnalysisResult получает результат анализа
func (h *AnalyzeHandler) GetAnalysisResult(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	analysisID := c.Param("id")

	if analysisID == "" {
		requestLogger.Warn("Missing analysis ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:     "missing_field",
			Message:   "ID анализа обязателен",
			Timestamp: time.Now(),
		})
		return
	}

	result, err := h.analyzeService.GetAnalysisResult(c.Request.Context(), analysisID)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("analysis_id", analysisID).Error("Failed to get analysis result")
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:     "analysis_not_found",
			Message:   "Результат анализа не найден",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("analysis_id", analysisID).Info("Analysis result retrieved")
	c.JSON(http.StatusOK, result)
}

// ListAnalyses получает список анализов пользователя
func (h *AnalyzeHandler) ListAnalyses(c *gin.Context) {
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

	analyses, err := h.analyzeService.ListAnalyses(c.Request.Context(), userID, limit, offset)
	if err != nil {
		requestLogger.WithField("error", err.Error()).WithField("user_id", userID).Error("Failed to list analyses")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "list_failed",
			Message:   "Ошибка получения списка анализов",
			Timestamp: time.Now(),
		})
		return
	}

	requestLogger.WithField("user_id", userID).WithField("count", len(analyses)).Info("Analyses listed")
	c.JSON(http.StatusOK, gin.H{
		"analyses": analyses,
		"limit":    limit,
		"offset":   offset,
		"count":    len(analyses),
	})
}
