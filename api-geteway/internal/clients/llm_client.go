package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LLMClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type LLMRequest struct {
	UserQuery     string                 `json:"user_query"`
	SourceConfig  map[string]interface{} `json:"source_config"`
	TargetConfig  map[string]interface{} `json:"target_config"`
	OperationType string                 `json:"operation_type"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	PipelineID string    `json:"pipeline_id"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	Error      *LLMError `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type LLMError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type GenerateDDLRequest struct {
	DataProfile DataProfile `json:"data_profile"`
	Target      DataTarget  `json:"target"`
}

type DataProfile struct {
	DataType         string      `json:"data_type"`
	TotalRows        int32       `json:"total_rows"`
	SampledRows      int32       `json:"sampled_rows"`
	Fields           []DataField `json:"fields"`
	SampleData       string      `json:"sample_data"`
	DataQualityScore string      `json:"data_quality_score"`
}

type DataField struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Nullable    bool    `json:"nullable"`
	NullCount   int32   `json:"null_count"`
	SampleValue string  `json:"sample_value"`
	MinValue    float64 `json:"min_value"`
	MaxValue    float64 `json:"max_value"`
	Description string  `json:"description"`
}

type DataTarget struct {
	Type             string `json:"type"`
	ConnectionString string `json:"connection_string"`
	TableName        string `json:"table_name"`
	CredentialsRef   string `json:"credentials_ref"`
}

type GenerateDDLResponse struct {
	DDLScript   string `json:"ddl_script"`
	Status      string `json:"status"`
	Explanation string `json:"explanation"`
}

func NewLLMClient(baseURL, apiKey string) *LLMClient {
	return &LLMClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *LLMClient) GenerateDDL(dataProfile DataProfile, target DataTarget) (string, error) {
	// Формируем запрос для кастомной LLM
	userQuery := c.buildPrompt(dataProfile, target)
	
	request := LLMRequest{
		UserQuery: userQuery,
		SourceConfig: map[string]interface{}{
			"type":      dataProfile.DataType,
			"file_path": "data.csv",
		},
		TargetConfig: map[string]interface{}{
			"type":             target.Type,
			"table_name":       target.TableName,
			"connection_string": target.ConnectionString,
		},
		OperationType: "ddl_generation",
	}

	// Отправляем запрос к кастомной LLM
	response, err := c.callCustomLLM(request)
	if err != nil {
		return "", fmt.Errorf("failed to call custom LLM: %v", err)
	}

	return response, nil
}

func (c *LLMClient) buildPrompt(dataProfile DataProfile, target DataTarget) string {
	// Преобразуем профиль данных в JSON для промпта
	profileJSON, _ := json.MarshalIndent(dataProfile, "", "  ")

	prompt := fmt.Sprintf(`
Ты - эксперт по проектированию баз данных. Проанализируй профиль данных и создай оптимальный DDL скрипт.

ПРОФИЛЬ ДАННЫХ:
%s

ЦЕЛЕВАЯ СИСТЕМА:
- Тип: %s
- Таблица: %s

ТРЕБОВАНИЯ:
1. Создай DDL скрипт для создания таблицы
2. Учти типы данных и ограничения
3. Добавь индексы для оптимизации
4. Включи комментарии к полям
5. Учти качество данных: %s

Верни только SQL скрипт без дополнительных объяснений.
`, profileJSON, target.Type, target.TableName, dataProfile.DataQualityScore)

	return prompt
}

func (c *LLMClient) callCustomLLM(request LLMRequest) (string, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Custom LLM API error: %s", string(body))
	}

	var llmResp LLMResponse
	if err := json.Unmarshal(body, &llmResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if llmResp.Error != nil {
		return "", fmt.Errorf("Custom LLM error: %s", llmResp.Error.Message)
	}

	if llmResp.Status != "success" {
		return "", fmt.Errorf("Custom LLM returned error status: %s", llmResp.Message)
	}

	return llmResp.Message, nil
}

// Мок-реализация для тестирования без реального LLM
func (c *LLMClient) GenerateDDLMock(dataProfile DataProfile, target DataTarget) (string, error) {
	ddl := fmt.Sprintf(`-- DDL скрипт для таблицы %s
-- Сгенерировано на основе анализа данных
-- Качество данных: %s
-- Тип данных: %s

CREATE TABLE IF NOT EXISTS %s (
`, target.TableName, dataProfile.DataQualityScore, dataProfile.DataType, target.TableName)

	// Добавляем поля на основе профиля
	for i, field := range dataProfile.Fields {
		if i > 0 {
			ddl += ",\n"
		}

		// Определяем SQL тип
		sqlType := c.mapToSQLType(field.Type)
		if field.Nullable {
			sqlType += " NULL"
		} else {
			sqlType += " NOT NULL"
		}

		ddl += fmt.Sprintf("    %s %s -- %s", field.Name, sqlType, field.Description)
	}

	ddl += "\n);\n\n"

	// Добавляем индексы
	ddl += "-- Индексы для оптимизации\n"
	for _, field := range dataProfile.Fields {
		if field.Type == "int" || field.Type == "datetime" {
			ddl += fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_%s ON %s (%s);\n",
				target.TableName, field.Name, target.TableName, field.Name)
		}
	}

	return ddl, nil
}

func (c *LLMClient) mapToSQLType(goType string) string {
	switch goType {
	case "int":
		return "INTEGER"
	case "float":
		return "DECIMAL(10,2)"
	case "string":
		return "VARCHAR(255)"
	case "datetime":
		return "TIMESTAMP"
	case "bool":
		return "BOOLEAN"
	case "json":
		return "JSON"
	default:
		return "TEXT"
	}
}
