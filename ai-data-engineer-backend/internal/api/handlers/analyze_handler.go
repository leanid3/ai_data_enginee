package handlers

import (
	"ai-data-engineer-backend/domain/models"
	"ai-data-engineer-backend/pkg/logger"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DataAnalyzerService interface {
	AnalyzeFile(ctx context.Context, userID string) (models.AnalysisResult, error)
}

type AnalyzeHandler struct {
	dataAnalyzer DataAnalyzerService
	logger       logger.Logger
}

func NewAnalyzeHandler(dataAnalyzer DataAnalyzerService, logger logger.Logger) *AnalyzeHandler {
	return &AnalyzeHandler{
		dataAnalyzer: dataAnalyzer,
		logger:       logger,
	}
}

func (h *AnalyzeHandler) AnalyzeFile(c *gin.Context) {
	requestLogger := logger.GetLoggerFromContext(c.Request.Context())
	requestLogger.Info("Starting: Handler.AnalyzeHandler.AnalyzeFile")

	userID := c.Param("user_id")
	if userID == "" {
		userID = "default_user" // По умолчанию
	}

	analysisResult, err := h.dataAnalyzer.AnalyzeFile(c.Request.Context(), userID)
	if err != nil {
		requestLogger.WithField("error", err.Error()).Error("Failed to analyze file")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "analyze_failed",
			Message:   "Ошибка анализа файла",
			Timestamp: time.Now(),
		})
		return
	}

	// Парсим JSON строку в map
	var resultMap map[string]interface{}
	if err := json.Unmarshal([]byte(analysisResult.AnalysisResult.Content.(string)), &resultMap); err != nil {
		requestLogger.WithField("error", err.Error()).Error("Failed to parse analysis result")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:     "parsing_error",
			Message:   "Ошибка при обработке результата анализа",
			Timestamp: time.Now(),
		})
		return
	}

	response := models.AnalysisResponse{
		Status:    analysisResult.Status,
		Message:   "Файл успешно проанализирован",
		Result:    resultMap,
		CreatedAt: time.Now(),
	}
	requestLogger.WithField("analysis_id", analysisResult).Info("File analyzed successfully")
	c.JSON(http.StatusOK, response)
	requestLogger.Info("End: Handler.AnalyzeHandler.AnalyzeFile")
}
