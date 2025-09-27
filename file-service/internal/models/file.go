package models

import (
	"time"

	"gorm.io/gorm"
)

// File представляет файл в системе
type File struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	Filename    string         `gorm:"not null" json:"filename"`
	ContentType string         `gorm:"not null" json:"content_type"`
	FileSize    int64          `gorm:"not null" json:"file_size"`
	Status      string         `gorm:"not null;default:'uploaded'" json:"status"` // uploaded, processing, failed, deleted
	StoragePath string         `gorm:"not null" json:"storage_path"`
	StorageType string         `gorm:"not null;default:'minio'" json:"storage_type"` // minio, s3, local
	Checksum    string         `json:"checksum"`
	Metadata    string         `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON с дополнительными метаданными
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// FileUploadSession представляет сессию загрузки файла
type FileUploadSession struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"not null;index" json:"user_id"`
	Filename     string    `gorm:"not null" json:"filename"`
	ContentType  string    `gorm:"not null" json:"content_type"`
	TotalSize    int64     `gorm:"not null" json:"total_size"`
	UploadedSize int64     `gorm:"not null;default:0" json:"uploaded_size"`
	Status       string    `gorm:"not null;default:'uploading'" json:"status"` // uploading, completed, failed
	ChunkSize    int64     `gorm:"not null;default:5242880" json:"chunk_size"` // 5MB по умолчанию
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// FileAccessLog представляет лог доступа к файлу
type FileAccessLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FileID    string    `gorm:"not null;index" json:"file_id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Action    string    `gorm:"not null" json:"action"` // download, view, delete
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName возвращает имя таблицы для модели File
func (File) TableName() string {
	return "files"
}

// TableName возвращает имя таблицы для модели FileUploadSession
func (FileUploadSession) TableName() string {
	return "file_upload_sessions"
}

// TableName возвращает имя таблицы для модели FileAccessLog
func (FileAccessLog) TableName() string {
	return "file_access_logs"
}
