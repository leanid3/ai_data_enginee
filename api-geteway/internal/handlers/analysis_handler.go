package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AnalysisHandler struct {
	airflowURL string
	ollamaURL  string
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

type AirflowDAGRunRequest struct {
	Conf map[string]interface{} `json:"conf"`
}

type AirflowDAGRunResponse struct {
	DagRunID string `json:"dag_run_id"`
	Message  string `json:"message"`
}

func NewAnalysisHandler(airflowURL, ollamaURL string) *AnalysisHandler {
	return &AnalysisHandler{
		airflowURL: airflowURL,
		ollamaURL:  ollamaURL,
	}
}

// RegisterAnalysisRoutes регистрирует маршруты для анализа
func (h *AnalysisHandler) RegisterAnalysisRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/analysis/start", h.StartAnalysisHTTP)
	mux.HandleFunc("/api/v1/analysis/status/", h.GetAnalysisStatusHTTP)
}

// StartAnalysisHTTP HTTP обработчик для запуска анализа
func (h *AnalysisHandler) StartAnalysisHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if req.FileID == "" || req.UserID == "" || req.FilePath == "" {
		http.Error(w, "file_id, user_id и file_path обязательны", http.StatusBadRequest)
		return
	}

	// Подготовка данных для Airflow DAG
	dagRunData := AirflowDAGRunRequest{
		Conf: map[string]interface{}{
			"file_id":   req.FileID,
			"user_id":   req.UserID,
			"file_path": req.FilePath,
		},
	}

	// Отправка запроса в Airflow
	dagRunResponse, err := h.triggerAirflowDAG(dagRunData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Ошибка запуска анализа",
			"details": err.Error(),
		})
		return
	}

	// Генерация уникального ID анализа
	analysisID := fmt.Sprintf("analysis_%s_%d", req.FileID, time.Now().Unix())

	response := AnalysisResponse{
		AnalysisID: analysisID,
		Status:     "started",
		Message:    "Анализ данных запущен",
		Result: map[string]interface{}{
			"dag_run_id":  dagRunResponse.DagRunID,
			"airflow_url": fmt.Sprintf("%s/dags/data_analysis_pipeline/grid", h.airflowURL),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAnalysisStatusHTTP HTTP обработчик для получения статуса анализа
func (h *AnalysisHandler) GetAnalysisStatusHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	// Здесь должна быть логика получения статуса из Airflow
	// Для демонстрации возвращаем фиктивный статус

	status := map[string]interface{}{
		"analysis_id": analysisID,
		"status":      "completed",
		"progress":    100,
		"message":     "Анализ завершен",
		"result": map[string]interface{}{
			"data_quality_score": 0.85,
			"recommendations": []string{
				"Обнаружены пропущенные значения в полях name, age, city",
				"Рекомендуется очистка данных перед загрузкой в целевую систему",
				"Поле email содержит некорректные форматы",
				"Данные подходят для аналитической обработки",
			},
			"storage_recommendation": "PostgreSQL",
			"analysis_timestamp":     time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// triggerAirflowDAG запускает DAG в Airflow
func (h *AnalysisHandler) triggerAirflowDAG(dagRunData AirflowDAGRunRequest) (*AirflowDAGRunResponse, error) {
	// Подготовка URL для Airflow API
	url := fmt.Sprintf("%s/api/v1/dags/data_analysis_pipeline/dagRuns", h.airflowURL)

	// Подготовка данных для отправки
	jsonData, err := json.Marshal(dagRunData)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	// Создание HTTP запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Отправка запроса
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка Airflow API: %s", string(body))
	}

	// Парсинг ответа
	var dagRunResponse AirflowDAGRunResponse
	if err := json.Unmarshal(body, &dagRunResponse); err != nil {
		return nil, fmt.Errorf("ошибка парсинг ответа: %w", err)
	}

	return &dagRunResponse, nil
}
