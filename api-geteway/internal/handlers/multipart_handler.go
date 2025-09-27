package handlers

import (
	"backend/api-gateway/config"
	"backend/api-gateway/internal/clients"
	profilerGen "backend/api-gateway/profiler-gen"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MultipartHandler struct {
	ProfilerClient profilerGen.DataProfilerClient
	LLMClient      *clients.LLMClient
}

func NewMultipartHandler(profilerClient profilerGen.DataProfilerClient) *MultipartHandler {
	cfg := config.LoadConfig()
	return &MultipartHandler{
		ProfilerClient: profilerClient,
		LLMClient:      clients.NewLLMClient(cfg.LLMBaseURL, cfg.LLMAPIKey),
	}
}

// UploadCSV обрабатывает загрузку CSV файлов
func (h *MultipartHandler) UploadCSV(w http.ResponseWriter, r *http.Request) {
	h.handleFileUpload(w, r, "csv")
}

// UploadJSON обрабатывает загрузку JSON файлов
func (h *MultipartHandler) UploadJSON(w http.ResponseWriter, r *http.Request) {
	h.handleFileUpload(w, r, "json")
}

// UploadXML обрабатывает загрузку XML файлов
func (h *MultipartHandler) UploadXML(w http.ResponseWriter, r *http.Request) {
	h.handleFileUpload(w, r, "xml")
}

// handleFileUpload общий метод для обработки загрузки файлов
func (h *MultipartHandler) handleFileUpload(w http.ResponseWriter, r *http.Request, fileType string) {
	// Устанавливаем максимальный размер файла (500MB)
	const maxMemory = 500 << 20 // 500MB
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse multipart form: %v", err), http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file from form: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем user_id из формы
	userID := r.FormValue("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	// Читаем содержимое файла
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
		return
	}

	// Определяем content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Создаем запрос к Data Profiler
	profilerReq := &profilerGen.UploadFileRequest{
		Filename:    header.Filename,
		ContentType: contentType,
		FileData:    fileData,
		UserId:      userID,
	}

	// Вызываем Data Profiler
	profilerResp, err := h.ProfilerClient.UploadFile(context.Background(), profilerReq)
	if err != nil {
		response := map[string]interface{}{
			"fileId":  "",
			"status":  "failed",
			"message": fmt.Sprintf("Upload failed: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]interface{}{
		"fileId":  profilerResp.FileId,
		"status":  profilerResp.Status,
		"message": profilerResp.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ProcessFileMultipart обрабатывает файл после загрузки
func (h *MultipartHandler) ProcessFileMultipart(w http.ResponseWriter, r *http.Request) {
	// Извлекаем file_id из URL или body
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		// Пытаемся получить из body
		var req struct {
			FileID string `json:"file_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			fileID = req.FileID
		}
	}

	if fileID == "" {
		http.Error(w, "file_id is required", http.StatusBadRequest)
		return
	}

	// Создаем запрос к Data Profiler
	profilerReq := &profilerGen.ProcessFileRequest{
		FileId: fileID,
	}

	// Вызываем Data Profiler
	profilerResp, err := h.ProfilerClient.ProcessFile(context.Background(), profilerReq)
	if err != nil {
		response := map[string]interface{}{
			"status":  "failed",
			"message": fmt.Sprintf("Processing failed: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profilerResp)
}
