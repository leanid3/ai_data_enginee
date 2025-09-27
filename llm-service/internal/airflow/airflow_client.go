package airflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AirflowClient struct {
	BaseURL  string
	Username string
	Password string
	Client   *http.Client
}

type DAGRunRequest struct {
	DagID         string                 `json:"dag_id"`
	Conf          map[string]interface{} `json:"conf"`
	ExecutionDate string                 `json:"execution_date"`
}

type DAGRunResponse struct {
	DagID     string `json:"dag_id"`
	DagRunID  string `json:"dag_run_id"`
	State     string `json:"state"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type TaskInstance struct {
	TaskID    string `json:"task_id"`
	State     string `json:"state"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type DAGRunStatus struct {
	DagID     string         `json:"dag_id"`
	DagRunID  string         `json:"dag_run_id"`
	State     string         `json:"state"`
	StartDate string         `json:"start_date"`
	EndDate   string         `json:"end_date"`
	Tasks     []TaskInstance `json:"tasks"`
}

func NewAirflowClient(baseURL, username, password string) *AirflowClient {
	return &AirflowClient{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (a *AirflowClient) TriggerDAG(dagID string, conf map[string]interface{}) (*DAGRunResponse, error) {
	request := DAGRunRequest{
		DagID:         dagID,
		Conf:          conf,
		ExecutionDate: time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/dags/%s/dagRuns", a.BaseURL, dagID),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(a.Username, a.Password)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airflow returned status %d: %s", resp.StatusCode, string(body))
	}

	var response DAGRunResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (a *AirflowClient) GetDAGRunStatus(dagID, dagRunID string) (*DAGRunStatus, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/dags/%s/dagRuns/%s", a.BaseURL, dagID, dagRunID),
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(a.Username, a.Password)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airflow returned status %d: %s", resp.StatusCode, string(body))
	}

	var response DAGRunStatus
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (a *AirflowClient) CheckHealth() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/health", a.BaseURL), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(a.Username, a.Password)

	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("airflow is not healthy, status: %d", resp.StatusCode)
	}

	return nil
}
