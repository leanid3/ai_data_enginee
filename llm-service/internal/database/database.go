package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

type FileMetadata struct {
	ID         string `gorm:"primaryKey"`
	UserID     string `gorm:"not null"`
	FileName   string `gorm:"not null"`
	FilePath   string `gorm:"not null"`
	FileSize   int64  `gorm:"not null"`
	FileFormat string `gorm:"not null"`
	Status     string `gorm:"not null"`
	MinIOPath  string `gorm:"not null"`
	CreatedAt  string `gorm:"not null"`
	UpdatedAt  string `gorm:"not null"`
}

type AnalysisResult struct {
	ID              string  `gorm:"primaryKey"`
	FileID          string  `gorm:"not null"`
	UserID          string  `gorm:"not null"`
	AnalysisType    string  `gorm:"not null"`
	DataProfile     string  `gorm:"type:text"`
	QualityScore    float64 `gorm:"not null"`
	Recommendations string  `gorm:"type:text"`
	DDLScript       string  `gorm:"type:text"`
	Status          string  `gorm:"not null"`
	CreatedAt       string  `gorm:"not null"`
	UpdatedAt       string  `gorm:"not null"`
}

func NewDatabase(host, port, user, password, dbname string) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Автомиграция таблиц
	err = db.AutoMigrate(&FileMetadata{}, &AnalysisResult{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) SaveFileMetadata(metadata *FileMetadata) error {
	return d.DB.Create(metadata).Error
}

func (d *Database) GetFileMetadata(fileID string) (*FileMetadata, error) {
	var metadata FileMetadata
	err := d.DB.Where("id = ?", fileID).First(&metadata).Error
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

func (d *Database) SaveAnalysisResult(result *AnalysisResult) error {
	return d.DB.Create(result).Error
}

func (d *Database) GetAnalysisResult(fileID string) (*AnalysisResult, error) {
	var result AnalysisResult
	err := d.DB.Where("file_id = ?", fileID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *Database) UpdateAnalysisResult(fileID string, updates map[string]interface{}) error {
	return d.DB.Model(&AnalysisResult{}).Where("file_id = ?", fileID).Updates(updates).Error
}
