package custom_llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CustomLLMClient struct {
	BaseURL string
	APIKey  string
	Model   string
	Client  *http.Client
}

type CustomLLMRequest struct {
	UserQuery     string                 `json:"user_query"`
	SourceConfig  map[string]interface{} `json:"source_config"`
	TargetConfig  map[string]interface{} `json:"target_config"`
	OperationType string                 `json:"operation_type"`
}

type CustomLLMResponse struct {
	PipelineID string    `json:"pipeline_id"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	Error      *LLMError `json:"error,omitempty"`
}

type LLMError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewCustomLLMClient(baseURL, apiKey, model string) *CustomLLMClient {
	return &CustomLLMClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
		Client: &http.Client{
			Timeout: 10 * time.Minute, // Увеличенный таймаут для больших файлов
		},
	}
}

func (c *CustomLLMClient) GenerateResponse(prompt string) (string, error) {
	request := CustomLLMRequest{
		UserQuery: prompt,
		SourceConfig: map[string]interface{}{
			"type": "text",
		},
		TargetConfig: map[string]interface{}{
			"type": "response",
		},
		OperationType: "text_generation",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Custom LLM: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("custom LLM returned status %d: %s", resp.StatusCode, string(body))
	}

	var response CustomLLMResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("custom LLM error: %s", response.Error.Message)
	}

	if response.Status != "success" {
		return "", fmt.Errorf("custom LLM returned error status: %s", response.Message)
	}

	return response.Message, nil
}

func (c *CustomLLMClient) CheckHealth() error {
	// Простая проверка доступности кастомной LLM
	resp, err := c.Client.Get(c.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to check Custom LLM health: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("custom LLM is not healthy, status: %d", resp.StatusCode)
	}

	return nil
}
