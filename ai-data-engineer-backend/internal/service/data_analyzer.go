package service

import (
	"context"

	"ai-data-engineer-backend/pkg/client"
	"ai-data-engineer-backend/pkg/logger"
)

// DataAnalyzer реализация DataAnalyzer
type DataAnalyzer struct {
	logger    logger.Logger
	llmClient client.LLMClient
}

// NewDataAnalyzer создает новый анализатор данных
func NewDataAnalyzer(logger logger.Logger, llmClient client.LLMClient) *DataAnalyzer {
	return &DataAnalyzer{
		logger:    logger,
		llmClient: llmClient,
	}
}

// отправка запроса на анализ файла в LLM
func (d *DataAnalyzer) AnalyzeFile(ctx context.Context, userID string) (string, error) {
	d.logger.WithField("user_id", userID).Info("Starting analyze file")

	resp, err := d.llmClient.AnalyzeFile(ctx, userID)
	if err != nil {
		d.logger.WithField("error", err.Error()).Error("Failed to analyze file")
		return "", err
	}
	return resp, nil

}
