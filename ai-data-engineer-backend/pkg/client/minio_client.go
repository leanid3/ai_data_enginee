package client

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"ai-data-engineer-backend/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// FileInfo информация о файле в MinIO
type FileInfo struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
}

// minioClient реализация MinIOClient
type minioClient struct {
	client *minio.Client
	logger logger.Logger
}

// NewMinIOClient создает новый MinIO клиент
func NewMinIOClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool, logger logger.Logger) (*minioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &minioClient{
		client: client,
		logger: logger,
	}, nil
}

// * checkMinIOConnection проверяет подключение к MinIO
func (m *minioClient) checkMinIOConnection(ctx context.Context, bucket string) error {

	// Проверяем существование bucket
	m.logger.Info("Checking if bucket exists")
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		m.logger.WithField("error", err.Error()).Error("Failed to check bucket existence")
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		// Создаем bucket если не существует
		m.logger.WithField("bucket", bucket).Info("Bucket does not exist, creating new bucket")
		err = m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			m.logger.WithField("error", err.Error()).Error("Failed to create bucket")
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		m.logger.WithField("bucket", bucket).Info("Created new bucket successfully")
	} else {
		m.logger.WithField("bucket", bucket).Info("Bucket already exists")
	}
	return nil
}

// * UploadFile загружает файл в MinIO
func (m *minioClient) UploadFile(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error {
	m.logger.WithField("bucket", bucket).WithField("object", objectName).WithField("size", size).WithField("content_type", contentType).Info("🚀 Starting MinIO upload process")

	// Проверяем подключение к MinIO
	err := m.checkMinIOConnection(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check MinIO connection: %w", err)
	}

	// Загружаем файл
	m.logger.Info("Uploading file to MinIO")
	_, err = m.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		m.logger.WithField("error", err.Error()).Error("Failed to upload file to MinIO")
		return fmt.Errorf("failed to upload file: %w", err)
	}

	m.logger.WithField("bucket", bucket).WithField("object", objectName).Info("File uploaded to MinIO successfully")
	return nil
}

// DownloadFile скачивает файл из MinIO
func (m *minioClient) DownloadFile(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	m.logger.WithField("bucket", bucket).WithField("object", objectName).Info("Downloading file from MinIO")

	object, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return object, nil
}

// DownloadFileAsBytes скачивает файл из MinIO как []byte
func (m *minioClient) DownloadFileAsBytes(ctx context.Context, bucket, objectName string) ([]byte, error) {
	m.logger.WithField("bucket", bucket).WithField("object", objectName).Info("Downloading file as bytes from MinIO")

	object, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer object.Close()

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	m.logger.WithField("bucket", bucket).WithField("object", objectName).WithField("size", len(content)).Info("File downloaded successfully")
	return content, nil
}

// DeleteFile удаляет файл из MinIO
func (m *minioClient) DeleteFile(ctx context.Context, bucket, objectName string) error {
	m.logger.WithField("bucket", bucket).WithField("object", objectName).Info("Deleting file from MinIO")

	err := m.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	m.logger.WithField("bucket", bucket).WithField("object", objectName).Info("File deleted successfully")
	return nil
}

// ListFiles возвращает список файлов в bucket
func (m *minioClient) ListFiles(ctx context.Context, bucket, prefix string) ([]string, error) {
	m.logger.WithField("bucket", bucket).WithField("prefix", prefix).Info("Listing files in MinIO")

	var files []string
	objectCh := m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}
		files = append(files, object.Key)
	}

	return files, nil
}

// GetFileInfo получает информацию о файле
func (m *minioClient) GetFileInfo(ctx context.Context, bucket, objectName string) (*FileInfo, error) {
	m.logger.WithField("bucket", bucket).WithField("object", objectName).Info("Getting file info from MinIO")

	stat, err := m.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &FileInfo{
		Name:         stat.Key,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		LastModified: stat.LastModified,
		ETag:         stat.ETag,
	}, nil
}

// GenerateObjectName генерирует уникальное имя объекта для MinIO
func GenerateObjectName(userID, filename string) string {
	// Создаем путь: users/{userID}/files/{timestamp}_{filename}
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Очищаем имя файла от недопустимых символов
	cleanName := strings.ReplaceAll(nameWithoutExt, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "/", "_")
	cleanName = strings.ReplaceAll(cleanName, "\\", "_")

	return fmt.Sprintf("users/%s/files/%s_%s%s", userID, timestamp, cleanName, ext)
}

// GetContentType определяет Content-Type по расширению файла
func GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".csv":
		return "text/csv"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".txt":
		return "text/plain"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xls":
		return "application/vnd.ms-excel"
	default:
		return "application/octet-stream"
	}
}

// FileExists проверяет существование файла в MinIO
func (m *minioClient) FileExists(ctx context.Context, bucketName, objectName string) (bool, error) {
	m.logger.WithField("bucket", bucketName).WithField("object", objectName).Info("Checking file existence in MinIO")

	// Проверяем подключение к MinIO
	_, err := m.client.ListBuckets(ctx)
	if err != nil {
		m.logger.WithField("error", err.Error()).Error("MinIO connection failed")
		return false, fmt.Errorf("MinIO connection failed: %w", err)
	}

	// Проверяем существование объекта
	_, err = m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			m.logger.WithField("object", objectName).Info("File not found in MinIO")
			return false, nil
		}
		m.logger.WithField("error", err.Error()).Error("Failed to check file existence")
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	m.logger.WithField("object", objectName).Info("File exists in MinIO")
	return true, nil
}

// GetUserFromObjectName извлекает userID из имени объекта
func GetUserFromObjectName(objectName string) string {
	// Формат: users/{userID}/files/{filename}
	parts := strings.Split(objectName, "/")
	if len(parts) >= 2 && parts[0] == "users" {
		return parts[1]
	}
	return ""
}
