package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type APIGateway struct {
	FileServiceURL     string
	AnalysisServiceURL string
	Client             *http.Client
}

type FileUploadRequest struct {
	File   []byte `json:"file"`
	UserID string `json:"user_id"`
	Format string `json:"format"`
}

type FileUploadResponse struct {
	FileID    string `json:"file_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	UploadURL string `json:"upload_url"`
	CreatedAt string `json:"created_at"`
}

type AnalysisRequest struct {
	FileID   string `json:"file_id"`
	UserID   string `json:"user_id"`
	FilePath string `json:"file_path"`
}

type AnalysisResponse struct {
	AnalysisID string                 `json:"analysis_id"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Result     map[string]interface{} `json:"result,omitempty"`
}

func NewAPIGateway() *APIGateway {
	return &APIGateway{
		FileServiceURL:     "http://file-service:50054",
		AnalysisServiceURL: "http://data-analysis-service:8080",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (gw *APIGateway) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (gw *APIGateway) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем файл
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Определяем формат файла
	format := "csv"
	if handler.Filename != "" {
		ext := handler.Filename[len(handler.Filename)-4:]
		if ext == ".csv" {
			format = "csv"
		} else if ext == ".json" {
			format = "json"
		}
	}

	// Создаем multipart form для отправки в File Service
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	
	fileWriter, err := writer.CreateFormFile("file", handler.Filename)
	if err != nil {
		http.Error(w, "Failed to create form file", http.StatusInternalServerError)
		return
	}
	
	_, err = fileWriter.Write(fileBytes)
	if err != nil {
		http.Error(w, "Failed to write file to form", http.StatusInternalServerError)
		return
	}
	
	writer.Close()

	// Отправляем запрос в File Service
	resp, err := gw.Client.Post(
		fmt.Sprintf("%s/v1/files/upload/%s", gw.FileServiceURL, format),
		writer.FormDataContentType(),
		&b,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("File service error: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (gw *APIGateway) StartAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в Data Analysis Service
	jsonData, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	resp, err := gw.Client.Post(
		fmt.Sprintf("%s/api/v1/analysis/start", gw.AnalysisServiceURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Analysis service error: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (gw *APIGateway) GetAnalysisStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем analysis_id из URL
	analysisID := r.URL.Path[len("/api/v1/analysis/status/"):]
	if analysisID == "" {
		http.Error(w, "Analysis ID is required", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в Data Analysis Service
	resp, err := gw.Client.Get(
		fmt.Sprintf("%s/api/v1/analysis/status/%s", gw.AnalysisServiceURL, analysisID),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Analysis service error: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func main() {
	gw := NewAPIGateway()

	// Настраиваем маршруты
	http.HandleFunc("/health", gw.Health)
	http.HandleFunc("/api/v1/files/upload", gw.UploadFile)
	http.HandleFunc("/api/v1/analysis/start", gw.StartAnalysis)
	http.HandleFunc("/api/v1/analysis/status/", gw.GetAnalysisStatus)

	log.Println("API Gateway запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
