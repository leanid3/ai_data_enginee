package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"ai-data-engineer-backend/pkg/logger"
)

type StorageClient interface {
	UploadFile(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error
	DownloadFile(ctx context.Context, bucket, objectName string) (io.ReadCloser, error)
	DownloadFileAsBytes(ctx context.Context, bucket, objectName string) ([]byte, error)
	DeleteFile(ctx context.Context, bucket, objectName string) error
	ListFiles(ctx context.Context, bucket, prefix string) ([]string, error)
	FileExists(ctx context.Context, bucket, objectName string) (bool, error)
}

// FileService реализация FileService
type FileService struct {
	storageClient StorageClient
	logger        logger.Logger
}

// NewFileService создает новый FileService
func NewFileService(
	storageClient StorageClient,
	logger logger.Logger,
) *FileService {
	return &FileService{
		storageClient: storageClient,
		logger:        logger,
	}
}

// * UploadFile загружает файл в MinIO
func (s *FileService) UploadFile(ctx context.Context, userID, filename string, file io.Reader) (string, error) {
	s.logger.WithField("user_id", userID).WithField("filename", filename).Info("Starting analyze file")

	//! Генерируем уникальное имя файла для MinIO
	s.logger.Info("Generating unique object name for MinIO")
	objectName := generateFileName(userID, filename)

	// Читаем содержимое файла для получения размера
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	userBucket := "ai-data-engineer"
	userfilePath := fmt.Sprintf("users/%s/files/%s", userID, objectName)
	uploadErr := s.storageClient.UploadFile(ctx, userBucket, userfilePath, strings.NewReader(string(content)), int64(len(content)), "application/octet-stream")
	if uploadErr != nil {
		s.logger.WithField("error", uploadErr.Error()).Error("Failed to save file to MinIO")
		return "", fmt.Errorf("failed to save file to MinIO: %w", uploadErr)
	}

	return objectName, nil
}

// GetFileInfo получает информацию о файле
func (s *FileService) GetFileInfo(ctx context.Context, fileID string) (interface{}, error) {
	// TODO: Implement file info retrieval
	return map[string]interface{}{"file_id": fileID}, nil
}

// DeleteFile удаляет файл
func (s *FileService) DeleteFile(ctx context.Context, fileID string) error {
	// TODO: Implement file deletion
	return nil
}

// ListFiles получает список файлов
func (s *FileService) ListFiles(ctx context.Context, userID string, limit, offset int) ([]interface{}, error) {
	// TODO: Implement file listing
	return []interface{}{}, nil
}

func generateFileName(userID, filename string) string {
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	cleanName := strings.ReplaceAll(nameWithoutExt, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "/", "_")
	cleanName = strings.ReplaceAll(cleanName, "\\", "_")
	return fmt.Sprintf("%s_%s%s", timestamp, cleanName, ext)
}
