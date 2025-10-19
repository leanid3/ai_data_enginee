package tests

import (
	"ai-data-engineer-backend/internal/api/handlers"
	"ai-data-engineer-backend/internal/service"
	"ai-data-engineer-backend/pkg/client"
	"ai-data-engineer-backend/pkg/logger"
	"bytes"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// createTestServices создает реальные сервисы для тестирования
func createTestServices(t *testing.T) (*service.FileService, logger.Logger) {
	// Создаем логгер
	testLogger := logger.NewLogger("info", "json", "stdout")

	// Создаем реальный MinIO клиент для тестирования
	minioClient, err := client.NewMinIOClient(
		"localhost:9000", // MinIO endpoint
		"minioadmin",     // Access key
		"minioadmin",     // Secret key
		false,            // Use SSL
		testLogger,
	)
	if err != nil {
		t.Fatalf("Не удалось создать MinIO клиент: %v", err)
	}

	// Создаем реальный FileService с реальным MinIO клиентом
	fileService := service.NewFileService(minioClient, testLogger)

	return fileService, testLogger
}

// isMinIOAvailable проверяет, доступен ли MinIO сервер
func isMinIOAvailable() bool {
	conn, err := net.DialTimeout("tcp", "localhost:9000", 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func executeRequest(t *testing.T, router *gin.Engine, buf *bytes.Buffer, writer *multipart.Writer) {
	// Создаем новый HTTP   
	req := httptest.NewRequest("POST", "/api/v1/files/upload", bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
    rr := &httptest.ResponseRecorder{}
	
    router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %d, получили %d", http.StatusOK, status)
		return
	}

	// Проверяем, что ответ содержит ожидаемые поля
	responseBody := rr.Body.String()
	if !bytes.Contains([]byte(responseBody), []byte("file_id")) {
		t.Errorf("Ответ не содержит file_id: %s", responseBody)
		return
	}

	if !bytes.Contains([]byte(responseBody), []byte("uploaded")) {
		t.Errorf("Ответ не содержит статус uploaded: %s", responseBody)
		return
	}

	// Проверяем, что ответ содержит JSON структуру
	if !bytes.Contains([]byte(responseBody), []byte("{")) {
		t.Errorf("Ответ не является JSON: %s", responseBody)
		return
	}

}
func TestUploadFile(t *testing.T) {
	// Проверяем, что MinIO сервер доступен
	if !isMinIOAvailable() {
		t.Skip("MinIO сервер недоступен. Запустите: docker-compose up minio")
	}

	// Создаем буфер для multipart данных
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Создаем файл в форме
	part, err := writer.CreateFormFile("file", "test_data.csv")
	if err != nil {
		t.Fatalf("Не удалось создать форму файла %v", err)
	}

	fileContent := []byte("name,age\nJohn,30\nJane,25")
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatalf("Не удалось записать содержимое файла %v", err)
	}

	writer.Close()

}
