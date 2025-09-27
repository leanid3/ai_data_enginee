package minio

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	Client     *minio.Client
	BucketName string
}

func NewMinIOClient(endpoint, accessKey, secretKey, bucketName string) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Проверяем существование bucket
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Created bucket: %s", bucketName)
	}

	return &MinIOClient{
		Client:     client,
		BucketName: bucketName,
	}, nil
}

func (m *MinIOClient) GetFile(objectName string) (io.ReadCloser, error) {
	ctx := context.Background()
	object, err := m.Client.GetObject(ctx, m.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return object, nil
}

func (m *MinIOClient) ReadFileContent(objectName string) ([]byte, error) {
	object, err := m.GetFile(objectName)
	if err != nil {
		return nil, err
	}
	defer object.Close()

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return content, nil
}

func (m *MinIOClient) GetFileInfo(objectName string) (minio.ObjectInfo, error) {
	ctx := context.Background()
	info, err := m.Client.StatObject(ctx, m.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("failed to get object info: %w", err)
	}
	return info, nil
}

func (m *MinIOClient) ListFiles(prefix string) ([]minio.ObjectInfo, error) {
	ctx := context.Background()
	var objects []minio.ObjectInfo

	objectCh := m.Client.ListObjects(ctx, m.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}
		objects = append(objects, object)
	}

	return objects, nil
}
