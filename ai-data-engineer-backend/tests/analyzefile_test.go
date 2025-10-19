package tests

import (
	"ai-data-engineer-backend/internal/api/handlers"
	"ai-data-engineer-backend/internal/service"
	"ai-data-engineer-backend/pkg/client"
	"ai-data-engineer-backend/pkg/logger"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAnalyzeFile(t *testing.T) {
	executeAnalyzeFileRequest(t)
}

func executeAnalyzeFileRequest(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.New()
	llmClient := client.NewLLMClient("http://llm_service:8124", "", logger.NewLogger("info", "json", "stdout"), map[string]string{"analyze_file": "/api/v1/analyze-file"})
	analyzeHandler := handlers.NewAnalyzeHandler(service.NewDataAnalyzer(logger.NewLogger("info", "json", "stdout"), llmClient), logger.NewLogger("info", "json", "stdout"))
	router.POST("/api/v1/analyze-file", analyzeHandler.AnalyzeFile)

	req := httptest.NewRequest("POST", "/api/v1/analyze-file", bytes.NewBuffer([]byte(`{"user_id": "test_user"}`)))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %d, получили %d", http.StatusOK, status)
		return
	}
}
