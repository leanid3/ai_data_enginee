package handlers

import (
	"context"
	"time"

	"orchestration-service/gen"

	"github.com/google/uuid"
)

type OrchestrationHandler struct {
	gen.UnimplementedOrchestrationServiceServer
}

func NewOrchestrationHandler() *OrchestrationHandler {
	return &OrchestrationHandler{}
}

func (h *OrchestrationHandler) CreateAndRunDAG(ctx context.Context, req *gen.CreateAndRunDAGRequest) (*gen.CreateAndRunDAGResponse, error) {
	runID := uuid.New().String()

	return &gen.CreateAndRunDAGResponse{
		DagId:      req.DagId,
		RunId:      runID,
		Status:     "started",
		Message:    "DAG создан и запущен",
		CreatedAt:  time.Now().Format(time.RFC3339),
		AirflowUrl: "http://localhost:8082",
	}, nil
}

func (h *OrchestrationHandler) GetDAGStatus(ctx context.Context, req *gen.GetDAGStatusRequest) (*gen.GetDAGStatusResponse, error) {
	return &gen.GetDAGStatusResponse{
		DagId:           req.DagId,
		RunId:           req.RunId,
		Status:          "running",
		State:           "queued",
		ProgressPercent: 50,
		StartedAt:       time.Now().Format(time.RFC3339),
		Tasks:           []*gen.TaskStatus{},
	}, nil
}

func (h *OrchestrationHandler) GetDAGLogs(ctx context.Context, req *gen.GetDAGLogsRequest) (*gen.GetDAGLogsResponse, error) {
	logs := []*gen.LogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "INFO",
			Message:   "DAG запущен успешно",
			TaskId:    "start_task",
		},
	}

	return &gen.GetDAGLogsResponse{
		DagId:      req.DagId,
		RunId:      req.RunId,
		Logs:       logs,
		TotalLines: int32(len(logs)),
		HasMore:    false,
	}, nil
}

func (h *OrchestrationHandler) StopDAG(ctx context.Context, req *gen.StopDAGRequest) (*gen.StopDAGResponse, error) {
	return &gen.StopDAGResponse{
		DagId:   req.DagId,
		RunId:   req.RunId,
		Status:  "stopped",
		Message: "DAG остановлен",
	}, nil
}

func (h *OrchestrationHandler) ListDAGs(ctx context.Context, req *gen.ListDAGsRequest) (*gen.ListDAGsResponse, error) {
	dags := []*gen.DAGInfo{
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
	}

	return &gen.ListDAGsResponse{
		Dags:       dags,
		TotalCount: int32(len(dags)),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (h *OrchestrationHandler) GetDAGDetails(ctx context.Context, req *gen.GetDAGDetailsRequest) (*gen.GetDAGDetailsResponse, error) {
	tasks := []*gen.TaskInfo{
		{
			TaskId:       "extract_data",
			TaskType:     "PythonOperator",
			Description:  "Извлечение данных",
			Dependencies: []string{},
			RetryCount:   "3",
			Timeout:      "300",
		},
		{
			TaskId:       "transform_data",
			TaskType:     "PythonOperator",
			Description:  "Трансформация данных",
			Dependencies: []string{"extract_data"},
			RetryCount:   "3",
			Timeout:      "600",
		},
	}

	return &gen.GetDAGDetailsResponse{
		DagId:            req.DagId,
		Description:      "Пайплайн анализа данных",
		ScheduleInterval: "manual",
		LastRun:          time.Now().Format(time.RFC3339),
		NextRun:          time.Now().Add(time.Hour).Format(time.RFC3339),
		Status:           "active",
		Tags:             []string{"data", "analysis"},
		Tasks:            tasks,
		DagYaml:          "dag_id: data_analysis_pipeline",
		CreatedAt:        time.Now().Format(time.RFC3339),
	}, nil
}

func (h *OrchestrationHandler) RegisterCallback(ctx context.Context, req *gen.RegisterCallbackRequest) (*gen.RegisterCallbackResponse, error) {
	callbackID := uuid.New().String()

	return &gen.RegisterCallbackResponse{
		CallbackId: callbackID,
		Status:     "registered",
		Message:    "Callback зарегистрирован",
	}, nil
}

func (h *OrchestrationHandler) GetExecutionMetrics(ctx context.Context, req *gen.GetExecutionMetricsRequest) (*gen.GetExecutionMetricsResponse, error) {
	metrics := []*gen.ExecutionMetric{
		{
			MetricName: "dag_runs",
			Value:      10,
			Unit:       "count",
			Timestamp:  time.Now().Format(time.RFC3339),
			Labels:     map[string]string{"status": "success"},
		},
		{
			MetricName: "task_duration",
			Value:      120,
			Unit:       "seconds",
			Timestamp:  time.Now().Format(time.RFC3339),
			Labels:     map[string]string{"task": "extract_data"},
		},
	}

	return &gen.GetExecutionMetricsResponse{
		Metrics:     metrics,
		TimeRange:   req.TimeRange,
		GeneratedAt: time.Now().Format(time.RFC3339),
	}, nil
}
