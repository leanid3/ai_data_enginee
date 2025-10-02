package config

import (
	"os"
)

type Config struct {
	LLMBaseURL string
	LLMAPIKey  string
	LLMModel   string
}

func LoadConfig() *Config {
	return &Config{
		LLMBaseURL: getEnv("LLM_BASE_URL", "http://localhost:8124/api/v1/process"),
		LLMAPIKey:  getEnv("LLM_API_KEY", ""), // API ключ для кастомной LLM
		LLMModel:   getEnv("LLM_MODEL", "openrouter/auto"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
