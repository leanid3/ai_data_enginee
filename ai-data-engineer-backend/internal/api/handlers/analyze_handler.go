package handlers

import (
	"ai-data-engineer-backend/domain/models"
	"ai-data-engineer-backend/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DataAnalyzerService interface {
	AnalyzeFile(ctx context.Context, userID string) (string, error)
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

	userID := c.Query("user_id")
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
	response := models.AnalysisResponse{
		Status:    "completed",
		Message:   "Файл успешно проанализирован",
		Result:    map[string]interface{}{"analysis_result": analysisResult},
		CreatedAt: time.Now(),
	}
	requestLogger.WithField("analysis_id", analysisResult).Info("File analyzed successfully")
	c.JSON(http.StatusOK, response)
	requestLogger.Info("End: Handler.AnalyzeHandler.AnalyzeFile")
}
