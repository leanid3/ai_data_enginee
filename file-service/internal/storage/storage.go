package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/minio/minio-go/v7"
	minioCredentials "github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageInterface определяет интерфейс для работы с хранилищем
type StorageInterface interface {
	UploadFile(ctx context.Context, bucket, key string, reader io.Reader, size int64, contentType string) error
	DownloadFile(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, bucket, key string) error
	GetFileInfo(ctx context.Context, bucket, key string) (*FileInfo, error)
	GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error)
	CreateBucket(ctx context.Context, bucket string) error
	BucketExists(ctx context.Context, bucket string) (bool, error)
}

// FileInfo содержит информацию о файле в хранилище
type FileInfo struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
}

// StorageConfig содержит конфигурацию хранилища
type StorageConfig struct {
	Type      string
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	Region    string
}

// MinIOStorage реализует StorageInterface для MinIO
type MinIOStorage struct {
	client *minio.Client
	bucket string
}

// NewMinIOStorage создает новый экземпляр MinIOStorage
func NewMinIOStorage(config StorageConfig) (*MinIOStorage, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  minioCredentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания MinIO клиента: %w", err)
	}

	return &MinIOStorage{
		client: client,
		bucket: config.Bucket,
	}, nil
}

// UploadFile загружает файл в MinIO
func (m *MinIOStorage) UploadFile(ctx context.Context, bucket, key string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// DownloadFile скачивает файл из MinIO
func (m *MinIOStorage) DownloadFile(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

// DeleteFile удаляет файл из MinIO
func (m *MinIOStorage) DeleteFile(ctx context.Context, bucket, key string) error {
	return m.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

// GetFileInfo получает информацию о файле в MinIO
func (m *MinIOStorage) GetFileInfo(ctx context.Context, bucket, key string) (*FileInfo, error) {
	stat, err := m.client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Key:          key,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		LastModified: stat.LastModified,
		ETag:         stat.ETag,
	}, nil
}

// GeneratePresignedURL генерирует предварительно подписанный URL для MinIO
func (m *MinIOStorage) GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, bucket, key, expiration, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// CreateBucket создает bucket в MinIO
func (m *MinIOStorage) CreateBucket(ctx context.Context, bucket string) error {
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		return m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	return nil
}

// BucketExists проверяет существование bucket в MinIO
func (m *MinIOStorage) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return m.client.BucketExists(ctx, bucket)
}

// S3Storage реализует StorageInterface для AWS S3
type S3Storage struct {
	client *s3.S3
	bucket string
	region string
}

// NewS3Storage создает новый экземпляр S3Storage
func NewS3Storage(config StorageConfig) (*S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKey,
			config.SecretKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания S3 сессии: %w", err)
	}

	return &S3Storage{
		client: s3.New(sess),
		bucket: config.Bucket,
		region: config.Region,
	}, nil
}

// UploadFile загружает файл в S3
func (s *S3Storage) UploadFile(ctx context.Context, bucket, key string, reader io.Reader, size int64, contentType string) error {
	uploader := s3manager.NewUploaderWithClient(s.client)
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	return err
}

// DownloadFile скачивает файл из S3
func (s *S3Storage) DownloadFile(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

// DeleteFile удаляет файл из S3
func (s *S3Storage) DeleteFile(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

// GetFileInfo получает информацию о файле в S3
func (s *S3Storage) GetFileInfo(ctx context.Context, bucket, key string) (*FileInfo, error) {
	result, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Key:          key,
		Size:         *result.ContentLength,
		ContentType:  *result.ContentType,
		LastModified: *result.LastModified,
		ETag:         *result.ETag,
	}, nil
}

// GeneratePresignedURL генерирует предварительно подписанный URL для S3
func (s *S3Storage) GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(expiration)
}

// CreateBucket создает bucket в S3
func (s *S3Storage) CreateBucket(ctx context.Context, bucket string) error {
	_, err := s.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}

// BucketExists проверяет существование bucket в S3
func (s *S3Storage) BucketExists(ctx context.Context, bucket string) (bool, error) {
	_, err := s.client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	return err == nil, nil
}

// NewStorage создает новый экземпляр хранилища на основе конфигурации
func NewStorage(config StorageConfig) (StorageInterface, error) {
	switch strings.ToLower(config.Type) {
	case "minio":
		return NewMinIOStorage(config)
	case "s3":
		return NewS3Storage(config)
	default:
		return nil, fmt.Errorf("неподдерживаемый тип хранилища: %s", config.Type)
	}
}

// GenerateFilePath генерирует путь к файлу в хранилище
func GenerateFilePath(userID, fileID, filename string) string {
	return filepath.Join("users", userID, fileID, filename)
}
