package handlers

import (
	"context"
	"fmt"
	"log"

	chatGen "backend/api-gateway/chat-gen"
	fileGen "backend/api-gateway/file-gen"
	"backend/api-gateway/gen"
	llmGen "backend/api-gateway/llm-gen"
	orchestrationGen "backend/api-gateway/orchestration-gen"
)

type GatewayHandler struct {
	gen.UnimplementedDataEngineeringAssistantServer
	FileClient          fileGen.FileServiceClient
	ChatClient          chatGen.ChatServiceClient
	LLMClient           llmGen.LLMServiceClient
	OrchestrationClient orchestrationGen.OrchestrationServiceClient
}

func NewGatewayHandler(
	fileClient fileGen.FileServiceClient,
	chatClient chatGen.ChatServiceClient,
	llmClient llmGen.LLMServiceClient,
	orchestrationClient orchestrationGen.OrchestrationServiceClient,
) *GatewayHandler {
	return &GatewayHandler{
		FileClient:          fileClient,
		ChatClient:          chatClient,
		LLMClient:           llmClient,
		OrchestrationClient: orchestrationClient,
	}
}

// UploadCSV загружает CSV файл
func (h *GatewayHandler) UploadCSV(ctx context.Context, req *gen.UploadCSVRequest) (*gen.UploadFileResponse, error) {
	// Перенаправляем запрос в File Service
	fileReq := &fileGen.UploadFileRequest{
		UserId:      req.UserId,
		Filename:    "uploaded.csv", // Будет переопределено из multipart
		ContentType: "text/csv",
		FileSize:    0, // Будет определено из multipart
	}

	resp, err := h.FileClient.UploadFile(ctx, fileReq)
	if err != nil {
		return &gen.UploadFileResponse{
			FileId:  "",
			Status:  "failed",
			Message: fmt.Sprintf("Ошибка загрузки файла: %v", err),
		}, nil
	}

	return &gen.UploadFileResponse{
		FileId:  resp.FileId,
		Status:  resp.Status,
		Message: resp.Message,
	}, nil
}

// UploadJSON загружает JSON файл
func (h *GatewayHandler) UploadJSON(ctx context.Context, req *gen.UploadJSONRequest) (*gen.UploadFileResponse, error) {
	fileReq := &fileGen.UploadFileRequest{
		UserId:      req.UserId,
		Filename:    "uploaded.json",
		ContentType: "application/json",
		FileSize:    0,
	}

	resp, err := h.FileClient.UploadFile(ctx, fileReq)
	if err != nil {
		return &gen.UploadFileResponse{
			FileId:  "",
			Status:  "failed",
			Message: fmt.Sprintf("Ошибка загрузки файла: %v", err),
		}, nil
	}

	return &gen.UploadFileResponse{
		FileId:  resp.FileId,
		Status:  resp.Status,
		Message: resp.Message,
	}, nil
}

// UploadXML загружает XML файл
func (h *GatewayHandler) UploadXML(ctx context.Context, req *gen.UploadXMLRequest) (*gen.UploadFileResponse, error) {
	fileReq := &fileGen.UploadFileRequest{
		UserId:      req.UserId,
		Filename:    "uploaded.xml",
		ContentType: "application/xml",
		FileSize:    0,
	}

	resp, err := h.FileClient.UploadFile(ctx, fileReq)
	if err != nil {
		return &gen.UploadFileResponse{
			FileId:  "",
			Status:  "failed",
			Message: fmt.Sprintf("Ошибка загрузки файла: %v", err),
		}, nil
	}

	return &gen.UploadFileResponse{
		FileId:  resp.FileId,
		Status:  resp.Status,
		Message: resp.Message,
	}, nil
}

// ConnectDatabase подключается к базе данных
func (h *GatewayHandler) ConnectDatabase(ctx context.Context, req *gen.DatabaseConnectionRequest) (*gen.DatabaseConnectionResponse, error) {
	// Пока что возвращаем заглушку
	return &gen.DatabaseConnectionResponse{
		ConnectionId: "conn_" + req.UserId,
		Status:       "connected",
		Message:      "Подключение к базе данных успешно",
	}, nil
}

// CreatePipeline создает пайплайн
func (h *GatewayHandler) CreatePipeline(ctx context.Context, req *gen.CreatePipelineRequest) (*gen.CreatePipelineResponse, error) {
	// Создаем диалог в Chat Service
	chatReq := &chatGen.CreateDialogRequest{
		UserId:         req.UserId,
		Title:          "Анализ данных",
		InitialMessage: "Создание пайплайна для анализа данных",
	}

	chatResp, err := h.ChatClient.CreateDialog(ctx, chatReq)
	if err != nil {
		log.Printf("Ошибка создания диалога: %v", err)
	}

	// Анализируем данные через LLM Service
	llmReq := &llmGen.AnalyzeDataStructureRequest{
		UserId:       req.UserId,
		FilePath:     req.Source.GetFile().Url,
		FileFormat:   req.Source.GetFile().Format,
		SampleSize:   1000,
		AnalysisType: "detailed",
	}

	llmResp, err := h.LLMClient.AnalyzeDataStructure(ctx, llmReq)
	if err != nil {
		return &gen.CreatePipelineResponse{
			PipelineId:  "",
			Status:      "failed",
			Explanation: fmt.Sprintf("Ошибка анализа данных: %v", err),
		}, nil
	}

	// Создаем DAG в Orchestration Service
	orchestrationReq := &orchestrationGen.CreateAndRunDAGRequest{
		UserId:           req.UserId,
		DagId:            "data_analysis_" + req.UserId,
		DagYaml:          generateDAGYAML(req),
		ScheduleInterval: "manual",
		StartImmediately: true,
		Description:      "Анализ данных пользователя",
	}

	orchestrationResp, err := h.OrchestrationClient.CreateAndRunDAG(ctx, orchestrationReq)
	if err != nil {
		log.Printf("Ошибка создания DAG: %v", err)
	}

	return &gen.CreatePipelineResponse{
		PipelineId:  orchestrationResp.DagId,
		Status:      "created",
		Explanation: "Пайплайн успешно создан и запущен",
		DataProfile: convertToUnifiedProfile(llmResp.DataProfile),
	}, nil
}

// GetPipeline получает информацию о пайплайне
func (h *GatewayHandler) GetPipeline(ctx context.Context, req *gen.GetPipelineRequest) (*gen.PipelineDetails, error) {
	// Получаем статус DAG из Orchestration Service
	statusReq := &orchestrationGen.GetDAGStatusRequest{
		DagId:  req.PipelineId,
		UserId: "system", // TODO: получить из контекста
	}

	statusResp, err := h.OrchestrationClient.GetDAGStatus(ctx, statusReq)
	if err != nil {
		return &gen.PipelineDetails{
			PipelineId:    req.PipelineId,
			Status:        "unknown",
			LastRunStatus: "unknown",
		}, nil
	}

	return &gen.PipelineDetails{
		PipelineId:    req.PipelineId,
		Status:        statusResp.Status,
		LastRunStatus: statusResp.State,
	}, nil
}

// StartPostgresTransfer запускает перенос в PostgreSQL
func (h *GatewayHandler) StartPostgresTransfer(ctx context.Context, req *gen.StartTransferRequest) (*gen.StartTransferResponse, error) {
	// Пока что возвращаем заглушку
	return &gen.StartTransferResponse{
		TransferId:        "transfer_" + req.UserId,
		Status:            "started",
		Message:           "Перенос в PostgreSQL запущен",
		EstimatedDuration: 300,         // 5 минут
		EstimatedSize:     1024 * 1024, // 1MB
	}, nil
}

// StartClickHouseTransfer запускает перенос в ClickHouse
func (h *GatewayHandler) StartClickHouseTransfer(ctx context.Context, req *gen.StartTransferRequest) (*gen.StartTransferResponse, error) {
	// Пока что возвращаем заглушку
	return &gen.StartTransferResponse{
		TransferId:        "transfer_" + req.UserId,
		Status:            "started",
		Message:           "Перенос в ClickHouse запущен",
		EstimatedDuration: 300,
		EstimatedSize:     1024 * 1024,
	}, nil
}

// StartHDFSTransfer запускает перенос в HDFS
func (h *GatewayHandler) StartHDFSTransfer(ctx context.Context, req *gen.StartTransferRequest) (*gen.StartTransferResponse, error) {
	// Пока что возвращаем заглушку
	return &gen.StartTransferResponse{
		TransferId:        "transfer_" + req.UserId,
		Status:            "started",
		Message:           "Перенос в HDFS запущен",
		EstimatedDuration: 300,
		EstimatedSize:     1024 * 1024,
	}, nil
}

// GetTransferStatus получает статус переноса
func (h *GatewayHandler) GetTransferStatus(ctx context.Context, req *gen.GetTransferStatusRequest) (*gen.GetTransferStatusResponse, error) {
	// Пока что возвращаем заглушку
	return &gen.GetTransferStatusResponse{
		TransferId:      req.TransferId,
		Status:          "completed",
		Message:         "Перенос завершен",
		ProgressPercent: 100,
		ProcessedRows:   1000,
		TotalRows:       1000,
		ProcessedBytes:  1024 * 1024,
		TotalBytes:      1024 * 1024,
		StartedAt:       "2024-01-01T00:00:00Z",
		CompletedAt:     "2024-01-01T00:05:00Z",
	}, nil
}

// CancelTransfer отменяет перенос
func (h *GatewayHandler) CancelTransfer(ctx context.Context, req *gen.CancelTransferRequest) (*gen.CancelTransferResponse, error) {
	// Пока что возвращаем заглушку
	return &gen.CancelTransferResponse{
		TransferId: req.TransferId,
		Status:     "cancelled",
		Message:    "Перенос отменен",
	}, nil
}

// Вспомогательные функции

func generateDAGYAML(req *gen.CreatePipelineRequest) string {
	// Простой YAML для DAG
	return fmt.Sprintf(`
dag_id: data_analysis_%s
description: Анализ данных пользователя
schedule_interval: manual
tasks:
  - task_id: analyze_data
    operator: PythonOperator
    python_callable: analyze_data_function
    params:
      user_id: %s
      file_path: %s
`, req.UserId, req.UserId, req.Source.GetFile().Url)
}

func convertToUnifiedProfile(profile *llmGen.DataProfile) *gen.UnifiedDataProfile {
	if profile == nil {
		return &gen.UnifiedDataProfile{}
	}

	var fields []*gen.UnifiedField
	for _, field := range profile.Fields {
		fields = append(fields, &gen.UnifiedField{
			Name:        field.Name,
			Type:        field.Type,
			Nullable:    field.Nullable,
			NullCount:   int32(field.NullCount),
			SampleValue: field.SampleValue,
			MinValue:    field.MinValue,
			MaxValue:    field.MaxValue,
			Description: field.Description,
		})
	}

	return &gen.UnifiedDataProfile{
		DataType:         profile.DataType,
		TotalRows:        int32(profile.TotalRows),
		SampledRows:      int32(profile.SampledRows),
		Fields:           fields,
		SampleData:       profile.SampleData,
		DataQualityScore: fmt.Sprintf("%.2f", profile.DataQualityScore),
	}
}
