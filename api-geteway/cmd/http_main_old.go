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


func NewAPIGateway() *APIGateway {
	return &APIGateway{
		FileServiceURL:     "http://file-service:50054",
		AnalysisServiceURL: "http://data-analysis-service:8080",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (gw *APIGateway) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := r.FormValue("user_id")
	format := r.FormValue("format")
	if format == "" {
		format = "csv"
	}

	// Читаем файл
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Создаем multipart form для File Service
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Добавляем файл
	fileWriter, err := writer.CreateFormFile("file", "uploaded_file")
	if err != nil {
		http.Error(w, "Failed to create form file", http.StatusInternalServerError)
		return
	}
	fileWriter.Write(fileBytes)

	// Добавляем user_id
	writer.WriteField("user_id", userID)

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

func (gw *APIGateway) CreateDialog(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DialogCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в Chat Service
	jsonData, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	resp, err := gw.Client.Post(
		fmt.Sprintf("%s/v1/dialogs", gw.ChatServiceURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Chat service error: %v", err), http.StatusInternalServerError)
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

func (gw *APIGateway) AnalyzeData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в LLM Service
	jsonData, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	resp, err := gw.Client.Post(
		fmt.Sprintf("%s/v1/analyze", gw.LLMServiceURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("LLM service error: %v", err), http.StatusInternalServerError)
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

func (gw *APIGateway) GenerateDDL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateDDLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в LLM Service
	jsonData, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	resp, err := gw.Client.Post(
		fmt.Sprintf("%s/v1/generate-ddl", gw.LLMServiceURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("LLM service error: %v", err), http.StatusInternalServerError)
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

	var req struct {
		FileID   string `json:"file_id"`
		UserID   string `json:"user_id"`
		FilePath string `json:"file_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if req.FileID == "" || req.UserID == "" || req.FilePath == "" {
		http.Error(w, "file_id, user_id и file_path обязательны", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в сервис анализа данных
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
	path := r.URL.Path
	analysisID := path[len("/api/v1/analysis/status/"):]

	if analysisID == "" {
		http.Error(w, "analysis_id обязателен", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в сервис анализа данных
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

func (gw *APIGateway) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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
