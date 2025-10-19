package tests

import (
	"ai-data-engineer-backend/internal/api/handlers"
	"ai-data-engineer-backend/internal/service"
	"ai-data-engineer-backend/pkg/client"
	"ai-data-engineer-backend/pkg/logger"
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
	llmClient := client.NewLLMClient("http://localhost:8124", "", logger.NewLogger("info", "json", "stdout"), map[string]string{"analyze_file": "/api/v1/analyze-file"})
	analyzeService := service.NewDataAnalyzer(logger.NewLogger("info", "json", "stdout"), llmClient)
	analyzeHandler := handlers.NewAnalyzeHandler(analyzeService, logger.NewLogger("info", "json", "stdout"))
	router.POST("/api/v1/analyze-file", analyzeHandler.AnalyzeFile)

	req := httptest.NewRequest("POST", "/api/v1/analyze-file?user_id=default_user", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %d, получили %d", http.StatusOK, status)
		return
	}
}
