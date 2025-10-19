package service

import (
	"ai-data-engineer-backend/pkg/logger"
)

// PipelineService сервис для работы с пайплайнами
type PipelineService struct {
	logger logger.Logger
}

// NewPipelineService создает новый PipelineService
func NewPipelineService(logger logger.Logger) *PipelineService {
	return &PipelineService{
		logger: logger,
	}
}

// CreatePipeline создает новый пайплайн
func (p *PipelineService) CreatePipeline(name string) error {
	p.logger.WithField("pipeline_name", name).Info("Creating pipeline")
	// TODO: Implement pipeline creation
	return nil
}
