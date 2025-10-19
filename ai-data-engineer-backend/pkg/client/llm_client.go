package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ai-data-engineer-backend/domain/models"
	"ai-data-engineer-backend/pkg/logger"
)

type LLMClient interface {
	ProcessRequest(ctx context.Context, req *models.LLMRequest) (*models.LLMResponse, error)
	AnalyzeSchema(ctx context.Context, schema *models.DataSchema) (*models.AnalysisResult, error)
	GenerateDDL(ctx context.Context, req *models.GenerateDDLRequest) (*models.GenerateDDLResponse, error)
	AnalyzeFile(ctx context.Context, userID, objectName string) (string, error)
}

// llmClient реализация LLMClient
type llmClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     logger.Logger
}

// NewLLMClient создает новый LLM клиент
func NewLLMClient(baseURL, apiKey string, logger logger.Logger) LLMClient {
	return &llmClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// ProcessRequest отправляет запрос к LLM сервису
func (c *llmClient) ProcessRequest(ctx context.Context, req *models.LLMRequest) (*models.LLMResponse, error) {
	c.logger.WithField("operation_type", req.OperationType).Info("Processing LLM request")

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Используем правильный endpoint для LLM сервиса
	endpoint := c.baseURL + "/api/v1/process"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var llmResp models.LLMResponse
	if err := json.Unmarshal(body, &llmResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if llmResp.Error != nil {
		return nil, fmt.Errorf("LLM error: %s", *llmResp.Error)
	}

	c.logger.WithField("pipeline_id", llmResp.PipelineID).Info("LLM request processed successfully")
	return &llmResp, nil
}

// AnalyzeSchema анализирует схему данных
func (c *llmClient) AnalyzeSchema(ctx context.Context, schema *models.DataSchema) (*models.AnalysisResult, error) {
	c.logger.Info("Analyzing data schema with LLM")

	// Преобразуем схему в формат для LLM
	fields := make([]models.DataField, len(schema.Fields))
	for i, field := range schema.Fields {
		fields[i] = models.DataField{
			Name:        field.Name,
			Type:        field.Type,
			Nullable:    field.Nullable,
			SampleValue: field.SampleValue,
			Description: field.Description,
		}
	}

	dataProfile := models.DataProfile{
		DataType:         "csv", // По умолчанию
		TotalRows:        1000,  // Заглушка
		SampledRows:      100,   // Заглушка
		Fields:           fields,
		SampleData:       "",   // Будет заполнено из schema.Sample
		DataQualityScore: 0.85, // Заглушка
	}
	_ = dataProfile // Используем переменную

	req := &models.LLMRequest{
		UserQuery: "Проанализируй структуру данных и дай рекомендации по оптимизации",
		SourceConfig: map[string]interface{}{
			"type": "csv",
		},
		TargetConfig: map[string]interface{}{
			"type": "analysis",
		},
		OperationType: "data_analysis",
	}

	resp, err := c.ProcessRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze schema: %w", err)
	}

	// Преобразуем ответ в AnalysisResult
	return &models.AnalysisResult{
		AnalysisID:  resp.PipelineID,
		Status:      models.AnalysisStatusCompleted,
		LLMAnalysis: resp.Message,
	}, nil
}

// GenerateDDL генерирует DDL скрипт
func (c *llmClient) GenerateDDL(ctx context.Context, req *models.GenerateDDLRequest) (*models.GenerateDDLResponse, error) {
	c.logger.Info("Generating DDL with LLM")

	// Преобразуем запрос в формат для LLM
	llmReq := &models.LLMRequest{
		UserQuery: fmt.Sprintf("Создай DDL скрипт для таблицы %s на основе профиля данных", req.Target.TableName),
		SourceConfig: map[string]interface{}{
			"type": req.DataProfile.DataType,
		},
		TargetConfig: map[string]interface{}{
			"type":       req.Target.Type,
			"table_name": req.Target.TableName,
		},
		OperationType: "ddl_generation",
	}

	resp, err := c.ProcessRequest(ctx, llmReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DDL: %w", err)
	}

	return &models.GenerateDDLResponse{
		DDL:      resp.Message,
		Database: req.Database,
		Metadata: map[string]interface{}{
			"generated_by": "LLM",
			"status":       resp.Status,
		},
	}, nil
}

// AnalyzeFile отправляет запрос на анализ файла в LLM
func (c *llmClient) AnalyzeFile(ctx context.Context, userID, objectName string) (string, error) {
	c.logger.WithField("user_id", userID).WithField("object_name", objectName).Info("Analyzing file with LLM")

	// Создаем запрос для анализа файла
	req := &models.LLMRequest{
		UserQuery:     fmt.Sprintf("Проанализируй файл пользователя %s: %s", userID, objectName),
		SourceConfig:  map[string]interface{}{"type": "file_analysis"},
		TargetConfig:  map[string]interface{}{"type": "analysis"},
		OperationType: "file_analysis",
		DataProfile:   nil, // Не передаем DataProfile, только userID и objectName
	}

	// Отправляем запрос через ProcessRequest
	resp, err := c.ProcessRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to analyze file: %w", err)
	}

	// Возвращаем результат анализа
	if resp.UserReport != "" {
		return resp.UserReport, nil
	}
	return resp.Message, nil
}
