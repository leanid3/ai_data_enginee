package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaClient struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

type OllamaError struct {
	Error string `json:"error"`
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		BaseURL: baseURL,
		Model:   model,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *OllamaClient) GenerateResponse(prompt string) (string, error) {
	request := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.Client.Post(
		fmt.Sprintf("%s/api/generate", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var ollamaErr OllamaError
		if err := json.Unmarshal(body, &ollamaErr); err == nil {
			return "", fmt.Errorf("ollama error: %s", ollamaErr.Error)
		}
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Response, nil
}

func (c *OllamaClient) CheckHealth() error {
	resp, err := c.Client.Get(fmt.Sprintf("%s/api/tags", c.BaseURL))
	if err != nil {
		return fmt.Errorf("failed to check Ollama health: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama is not healthy, status: %d", resp.StatusCode)
	}

	return nil
}
