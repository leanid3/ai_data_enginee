package config

import (
	"os"
)

type Config struct {
	LLMBaseURL string
	LLMAPIKey  string
}

func LoadConfig() *Config {
	return &Config{
		LLMBaseURL: getEnv("LLM_BASE_URL", "http://ollama:11434/api/chat"),
		LLMAPIKey:  getEnv("LLM_API_KEY", ""), // Ollama не требует API ключ
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
