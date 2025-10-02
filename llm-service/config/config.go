package config

import (
	"os"
)

type Config struct {
	ServerPort   string
	CustomLLMURL string
	LLMModel     string
	LLMAPIKey    string
	Database     DatabaseConfig
	MinIO        MinIOConfig
	Airflow      AirflowConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type MinIOConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

type AirflowConfig struct {
	BaseURL  string
	Username string
	Password string
}

func Load() *Config {
	return &Config{
		ServerPort:   getEnv("LLM_SERVER_PORT", "50056"),
		CustomLLMURL: getEnv("CUSTOM_LLM_URL", "http://localhost:8124/api/v1/process"),
		LLMModel:     getEnv("LLM_MODEL", "openrouter/auto"),
		LLMAPIKey:    getEnv("LLM_API_KEY", ""),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "aien_db"),
		},
		MinIO: MinIOConfig{
			Endpoint:   getEnv("MINIO_ENDPOINT", "minio:9000"),
			AccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin"),
			BucketName: getEnv("MINIO_BUCKET", "files"),
			UseSSL:     false,
		},
		Airflow: AirflowConfig{
			BaseURL:  getEnv("AIRFLOW_BASE_URL", "http://airflow-webserver:8080"),
			Username: getEnv("AIRFLOW_USERNAME", "admin"),
			Password: getEnv("AIRFLOW_PASSWORD", "admin"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
