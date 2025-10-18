package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"ai-data-engineer-backend/internal/models"
	"ai-data-engineer-backend/pkg/logger"
)

// fileService реализация FileService
type fileService struct {
	fileProcessor FileProcessor
	analyzer      DataAnalyzer
	llmClient     LLMClient
	minioClient   MinIOClient
	logger        logger.Logger
}

// NewFileService создает новый FileService
func NewFileService(
	fileProcessor FileProcessor,
	analyzer DataAnalyzer,
	llmClient LLMClient,
	minioClient MinIOClient,
	logger logger.Logger,
) FileService {
	return &fileService{
		fileProcessor: fileProcessor,
		analyzer:      analyzer,
		llmClient:     llmClient,
		minioClient:   minioClient,
		logger:        logger,
	}
}

// AnalyzeFile анализирует файл
func (s *fileService) AnalyzeFile(ctx context.Context, req *AnalyzeFileRequest) (*AnalyzeFileResponse, error) {
	s.logger.WithField("filename", req.Filename).WithField("user_id", req.UserID).Info("🔍 STEP 1: Starting file analysis")

	// Читаем содержимое файла
	s.logger.Info("📖 STEP 2: Reading file content")
	content, err := io.ReadAll(req.File)
	if err != nil {
		s.logger.WithField("error", err.Error()).Error("❌ Failed to read file")
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	s.logger.WithField("content_size", len(content)).Info("✅ File content read successfully")

	// Генерируем уникальное имя файла для MinIO
	s.logger.Info("🏷️ STEP 3: Generating unique object name for MinIO")
	objectName := GenerateObjectName(req.UserID, req.Filename)
	contentType := GetContentType(req.Filename)
	s.logger.WithField("object_name", objectName).WithField("content_type", contentType).Info("✅ Object name generated")

	// Сохраняем файл в MinIO
	s.logger.Info("💾 STEP 4: Saving file to MinIO storage")
	s.logger.WithField("bucket", "ai-data-engineer").WithField("object", objectName).Info("📤 Uploading to MinIO...")

	err = s.minioClient.UploadFile(ctx, "ai-data-engineer", objectName, strings.NewReader(string(content)), int64(len(content)), contentType)
	if err != nil {
		s.logger.WithField("error", err.Error()).Warn("⚠️ Failed to save file to MinIO, continuing with analysis")
		// Продолжаем анализ даже если не удалось сохранить в MinIO
	} else {
		s.logger.WithField("object_name", objectName).Info("✅ File saved to MinIO successfully")
	}

	// Определяем тип файла
	s.logger.Info("🔍 STEP 5: Detecting file type")
	fileType := s.fileProcessor.DetectFileType(req.Filename)
	if fileType == "unknown" {
		fileType = req.FileType // Используем переданный тип
	}
	s.logger.WithField("detected_type", fileType).Info("✅ File type detected")

	// Парсим файл в зависимости от типа
	s.logger.Info("📊 STEP 6: Parsing file content")
	var profile *models.DataProfile
	switch fileType {
	case "csv":
		s.logger.Info("📄 Parsing CSV file...")
		profile, err = s.fileProcessor.ParseCSV(ctx, content)
	case "json":
		s.logger.Info("📄 Parsing JSON file...")
		profile, err = s.fileProcessor.ParseJSON(ctx, content)
	case "xml":
		s.logger.Info("📄 Parsing XML file...")
		profile, err = s.fileProcessor.ParseXML(ctx, content)
	default:
		s.logger.WithField("file_type", fileType).Error("❌ Unsupported file type")
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}

	if err != nil {
		s.logger.WithField("error", err.Error()).Error("❌ Failed to parse file")
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	s.logger.WithField("fields_count", len(profile.Fields)).Info("✅ File parsed successfully")

	// Анализируем структуру данных
	s.logger.Info("🧠 STEP 7: Analyzing data structure")
	analysis, err := s.analyzer.AnalyzeDataStructure(ctx, profile)
	if err != nil {
		s.logger.WithField("error", err.Error()).Error("❌ Failed to analyze data structure")
		return nil, fmt.Errorf("failed to analyze data structure: %w", err)
	}
	s.logger.Info("✅ Data structure analysis completed")

	// Отправляем на анализ в LLM
	s.logger.Info("🤖 STEP 8: Sending to LLM for advanced analysis")
	llmAnalysis, err := s.llmClient.AnalyzeSchema(ctx, &models.DataSchema{
		Fields: profile.Fields,
		Sample: []map[string]interface{}{}, // TODO: преобразовать sample данные
	})
	if err != nil {
		s.logger.WithField("error", err.Error()).Warn("⚠️ LLM analysis failed, continuing without it")
		// Продолжаем без LLM анализа
	} else {
		s.logger.Info("✅ LLM analysis completed successfully")
		analysis.LLMAnalysis = llmAnalysis.LLMAnalysis
	}

	// Генерируем уникальный ID файла
	s.logger.Info("🆔 STEP 9: Generating unique IDs")
	fileID := fmt.Sprintf("file_%d", time.Now().Unix())
	analysisID := fmt.Sprintf("analysis_%d", time.Now().Unix())

	s.logger.WithField("file_id", fileID).WithField("analysis_id", analysisID).Info("🎉 File analysis completed successfully")

	return &AnalyzeFileResponse{
		FileID:     fileID,
		AnalysisID: analysisID,
		Status:     "analyzed",
		Message:    "File analyzed successfully",
	}, nil
}

// GetFileInfo получает информацию о файле
func (s *fileService) GetFileInfo(ctx context.Context, fileID string) (*models.FileMetadata, error) {
	s.logger.WithField("file_id", fileID).Info("Getting file info (stub implementation)")

	// Заглушка - будет реализована в следующих итерациях
	return &models.FileMetadata{
		ID:          fileID,
		UserID:      "stub-user",
		Filename:    "stub-file.csv",
		ContentType: "text/csv",
		Size:        1024,
		Path:        "/stub/path",
		Bucket:      "files",
		Status:      models.FileStatusUploaded,
	}, nil
}

// DeleteFile удаляет файл
func (s *fileService) DeleteFile(ctx context.Context, fileID string) error {
	s.logger.WithField("file_id", fileID).Info("Deleting file (stub implementation)")

	// Заглушка - будет реализована в следующих итерациях
	return nil
}

// ListFiles получает список файлов
func (s *fileService) ListFiles(ctx context.Context, userID string, limit, offset int) ([]*models.FileMetadata, error) {
	s.logger.WithField("user_id", userID).Info("Listing files (stub implementation)")

	// Заглушка - будет реализована в следующих итерациях
	return []*models.FileMetadata{}, nil
}
