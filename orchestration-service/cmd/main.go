package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Простые структуры для тестирования
type CreateAndRunDAGRequest struct {
	UserId           string `json:"user_id"`
	DagId            string `json:"dag_id"`
	DagYaml          string `json:"dag_yaml"`
	ScheduleInterval string `json:"schedule_interval"`
	StartImmediately bool   `json:"start_immediately"`
	Description      string `json:"description"`
}

type CreateAndRunDAGResponse struct {
	DagId      string `json:"dag_id"`
	RunId      string `json:"run_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
	AirflowUrl string `json:"airflow_url"`
}

type GetDAGStatusRequest struct {
	DagId string `json:"dag_id"`
	RunId string `json:"run_id"`
}

type GetDAGStatusResponse struct {
	DagId           string        `json:"dag_id"`
	RunId           string        `json:"run_id"`
	Status          string        `json:"status"`
	State           string        `json:"state"`
	ProgressPercent int32         `json:"progress_percent"`
	StartedAt       string        `json:"started_at"`
	Tasks           []*TaskStatus `json:"tasks"`
}

type TaskStatus struct {
	TaskId    string `json:"task_id"`
	Status    string `json:"status"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}

type GetDAGLogsRequest struct {
	DagId string `json:"dag_id"`
	RunId string `json:"run_id"`
	Page  int32  `json:"page"`
	Size  int32  `json:"size"`
}

type GetDAGLogsResponse struct {
	DagId      string      `json:"dag_id"`
	RunId      string      `json:"run_id"`
	Logs       []*LogEntry `json:"logs"`
	TotalLines int32       `json:"total_lines"`
	HasMore    bool        `json:"has_more"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	TaskId    string `json:"task_id"`
}

type StopDAGRequest struct {
	DagId string `json:"dag_id"`
	RunId string `json:"run_id"`
}

type StopDAGResponse struct {
	DagId   string `json:"dag_id"`
	RunId   string `json:"run_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ListDAGsRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Status   string `json:"status"`
}

type ListDAGsResponse struct {
	Dags       []*DAGInfo `json:"dags"`
	TotalCount int32      `json:"total_count"`
	Page       int32      `json:"page"`
	PageSize   int32      `json:"page_size"`
}

type DAGInfo struct {
	DagId            string   `json:"dag_id"`
	Description      string   `json:"description"`
	Status           string   `json:"status"`
	LastRun          string   `json:"last_run"`
	NextRun          string   `json:"next_run"`
	ScheduleInterval string   `json:"schedule_interval"`
	Tags             []string `json:"tags"`
	TaskCount        int32    `json:"task_count"`
	CreatedAt        string   `json:"created_at"`
}

type GetDAGDetailsRequest struct {
	DagId string `json:"dag_id"`
}

type GetDAGDetailsResponse struct {
	DagId            string      `json:"dag_id"`
	Description      string      `json:"description"`
	ScheduleInterval string      `json:"schedule_interval"`
	LastRun          string      `json:"last_run"`
	NextRun          string      `json:"next_run"`
	Status           string      `json:"status"`
	Tags             []string    `json:"tags"`
	Tasks            []*TaskInfo `json:"tasks"`
	DagYaml          string      `json:"dag_yaml"`
	CreatedAt        string      `json:"created_at"`
}

type TaskInfo struct {
	TaskId       string   `json:"task_id"`
	TaskType     string   `json:"task_type"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"`
	RetryCount   string   `json:"retry_count"`
	Timeout      string   `json:"timeout"`
}

type RegisterCallbackRequest struct {
	Url        string   `json:"url"`
	EventTypes []string `json:"event_types"`
}

type RegisterCallbackResponse struct {
	CallbackId string `json:"callback_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type GetExecutionMetricsRequest struct {
	TimeRange string `json:"time_range"`
	DagId     string `json:"dag_id"`
}

type GetExecutionMetricsResponse struct {
	Metrics     []*ExecutionMetric `json:"metrics"`
	TimeRange   string             `json:"time_range"`
	GeneratedAt string             `json:"generated_at"`
}

type ExecutionMetric struct {
	MetricName string            `json:"metric_name"`
	Value      float64           `json:"value"`
	Unit       string            `json:"unit"`
	Timestamp  string            `json:"timestamp"`
	Labels     map[string]string `json:"labels"`
}

func main() {
	// Создаем HTTP сервер
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/v1/dags", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Создание DAG
			var req CreateAndRunDAGRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response := CreateAndRunDAGResponse{
				DagId:      req.DagId,
				RunId:      uuid.New().String(),
				Status:     "started",
				Message:    "DAG создан и запущен",
				CreatedAt:  time.Now().Format(time.RFC3339),
				AirflowUrl: "http://localhost:8082",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else if r.Method == "GET" {
			// Список DAG
			response := ListDAGsResponse{
				Dags: []*DAGInfo{
					{
						DagId:            "data_analysis_pipeline",
						Description:      "Пайплайн анализа данных",
						Status:           "active",
						LastRun:          time.Now().Format(time.RFC3339),
						NextRun:          time.Now().Add(time.Hour).Format(time.RFC3339),
						ScheduleInterval: "manual",
						Tags:             []string{"data", "analysis"},
						TaskCount:        5,
						CreatedAt:        time.Now().Format(time.RFC3339),
					},
				},
				TotalCount: 1,
				Page:       1,
				PageSize:   10,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/v1/dags/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Статус DAG
			response := GetDAGStatusResponse{
				DagId:           "data_analysis_pipeline",
				RunId:           uuid.New().String(),
				Status:          "running",
				State:           "queued",
				ProgressPercent: 50,
				StartedAt:       time.Now().Format(time.RFC3339),
				Tasks:           []*TaskStatus{},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Orchestration Service запущен на порту 50057")
	log.Fatal(http.ListenAndServe(":50057", nil))
}
