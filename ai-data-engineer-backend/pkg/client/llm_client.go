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
	AnalyzeFile(ctx context.Context, userID string) (*models.LLMResponse, error)
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
			Timeout: 2 * time.Minute,
		},
		logger:    logger,
		endpoints: endpoints,
	}
}

// SendRequest отправляет запрос к LLM сервису
func (c *llmClient) SendRequest(ctx context.Context, req *models.LLMRequest, endpoint string) (*models.LLMResponse, error) {
	c.logger.Info("llmClient.SendRequest: Starting")

	var jsonData []byte
	var err error

	if req == nil {
		jsonData = nil
	} else {
		jsonData, err = json.Marshal(req)
		if err != nil {
			c.logger.WithField("error", err.Error()).Error("llmClient.SendRequest: Failed to marshal request")
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	url := c.baseURL + endpoint
	c.logger.WithField("url", url).Info("llmClient.SendRequest: Sending request")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.WithField("error", err.Error()).Error("llmClient.SendRequest: Failed to send request")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.WithField("error", err.Error()).Error("llmClient.SendRequest: Failed to read response")
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	llmResp := models.LLMResponse{
		Content: string(body),
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
func (c *llmClient) AnalyzeFile(ctx context.Context, userID string) (*models.LLMResponse, error) {
	c.logger.WithField("user_id", userID).Info("LLMClient.AnalyzeFile: Starting")

	// Получаем endpoint для анализа файла
	endpoint := c.endpoints["analyze_file"]
	if endpoint == "" {
		return nil, fmt.Errorf("analyze_file endpoint is not configured")
	}

	// Добавляем параметр user_id
	if endpoint == "/api/v1/analyze-file" {
		endpoint = fmt.Sprintf("%s?user_id=%s", endpoint, userID)
	}

	// Отправляем запрос через sendRequest
	resp, err := c.SendRequest(ctx, nil, endpoint)
	if err != nil {
		c.logger.WithField("error", err.Error()).Error("LLMClient.AnalyzeFile: Failed to send request")
		return nil, fmt.Errorf("failed to analyze file: %w", err)
	}

	c.logger.WithField("result", resp.Content).Info("LLMClient.AnalyzeFile: Ending")
	return &models.LLMResponse{Content: resp.Content}, nil
}
