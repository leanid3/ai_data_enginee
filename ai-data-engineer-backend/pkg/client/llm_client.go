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
	SendRequest(ctx context.Context, req *models.LLMRequest, endpoint string) (*models.LLMResponse, error)
	GenerateDDL(ctx context.Context, req *models.GenerateDDLRequest) (*models.GenerateDDLResponse, error)
	AnalyzeFile(ctx context.Context, userID string) (string, error)
}

// llmClient реализация LLMClient
type llmClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     logger.Logger
	endpoints  map[string]string
}

// NewLLMClient создает новый LLM клиент
func NewLLMClient(baseURL, apiKey string, logger logger.Logger, endpoints map[string]string) LLMClient {
	return &llmClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:    logger,
		endpoints: endpoints,
	}
}

// SendRequest отправляет запрос к LLM сервису
func (c *llmClient) SendRequest(ctx context.Context, req *models.LLMRequest, endpoint string) (*models.LLMResponse, error) {
	c.logger.Info("llmClient.SendRequest: Starting")

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + endpoint
	c.logger.WithField("url", url).Info("llmClient.SendRequest: Sending request")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var llmResp models.LLMResponse
	if err := json.Unmarshal(body, &llmResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.Info("llmClient.SendRequest: Ending")
	return &llmResp, nil
}

// GenerateDDL генерирует DDL скрипт
func (c *llmClient) GenerateDDL(ctx context.Context, req *models.GenerateDDLRequest) (*models.GenerateDDLResponse, error) {
	//TODO: send request to LLM to generate DDL
	return nil, nil
}

// AnalyzeFile отправляет запрос на анализ файла в LLM
func (c *llmClient) AnalyzeFile(ctx context.Context, userID string) (string, error) {
	c.logger.WithField("user_id", userID).Info("LLMClient.AnalyzeFile: Starting")

	// Создаем запрос для анализа файла
	req := &models.LLMRequest{
		UserID: userID,
	}

	// Получаем endpoint для анализа файла
	endpoint := c.endpoints["analyze_file"]
	if endpoint == "" {
		return "", fmt.Errorf("analyze_file endpoint is not configured")
	}

	// Отправляем запрос через sendRequest
	resp, err := c.SendRequest(ctx, req, endpoint)
	if err != nil {
		c.logger.WithField("error", err.Error()).Error("LLMClient.AnalyzeFile: Failed to send request")
		return "", fmt.Errorf("failed to analyze file: %w", err)
	}

	// Преобразуем ответ в строку
	//TODO заменить на Response Model
	var result string
	if content, ok := resp.Content.(string); ok {
		result = content
	} else {
		// Если content не строка, преобразуем в JSON
		jsonBytes, err := json.Marshal(resp.Content)
		if err != nil {
			c.logger.WithField("error", err.Error()).Error("LLMClient.AnalyzeFile: Failed to marshal response content")
			return "", fmt.Errorf("failed to marshal response content: %w", err)
		}
		result = string(jsonBytes)
	}

	c.logger.WithField("result", result).Info("LLMClient.AnalyzeFile: Ending")
	return result, nil
}
