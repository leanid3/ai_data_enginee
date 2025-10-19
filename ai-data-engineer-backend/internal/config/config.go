package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config представляет конфигурацию приложения
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Airflow  AirflowConfig  `mapstructure:"airflow"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig конфигурация HTTP сервера
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig конфигурация баз данных
type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig `mapstructure:"postgresql"`
	ClickHouse ClickHouseConfig `mapstructure:"clickhouse"`
}

// PostgreSQLConfig конфигурация PostgreSQL
type PostgreSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxOpen  int    `mapstructure:"max_open"`
	MaxIdle  int    `mapstructure:"max_idle"`
}

// ClickHouseConfig конфигурация ClickHouse
type ClickHouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Secure   bool   `mapstructure:"secure"`
}

// LLMConfig конфигурация LLM сервиса
type LLMConfig struct {
	BaseURL    string            `mapstructure:"base_url"`
	APIKey     string            `mapstructure:"api_key"`
	Model      string            `mapstructure:"model"`
	Timeout    time.Duration     `mapstructure:"timeout"`
	MaxRetries int               `mapstructure:"max_retries"`
	Endpoints  map[string]string `mapstructure:"endpoints"`
}

// StorageConfig конфигурация хранилища
type StorageConfig struct {
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

// AirflowConfig конфигурация Airflow
type AirflowConfig struct {
	DAGsPath string `mapstructure:"dags_path"`
	BaseURL  string `mapstructure:"base_url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// LoggingConfig конфигурация логирования
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// LoadConfig загружает конфигурацию из файла и переменных окружения
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Устанавливаем значения по умолчанию
	setDefaults()

	// Автоматически читаем переменные окружения
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Читаем конфигурационный файл
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("ошибка чтения конфигурации: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	return &config, nil
}

// setDefaults устанавливает значения по умолчанию
func setDefaults() {
	// Server
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	// PostgreSQL
	viper.SetDefault("database.postgresql.host", "localhost")
	viper.SetDefault("database.postgresql.port", "5432")
	viper.SetDefault("database.postgresql.user", "postgres")
	viper.SetDefault("database.postgresql.password", "postgres")
	viper.SetDefault("database.postgresql.dbname", "aien_db")
	viper.SetDefault("database.postgresql.sslmode", "disable")
	viper.SetDefault("database.postgresql.max_open", 25)
	viper.SetDefault("database.postgresql.max_idle", 5)

	// ClickHouse
	viper.SetDefault("database.clickhouse.host", "localhost")
	viper.SetDefault("database.clickhouse.port", "9000")
	viper.SetDefault("database.clickhouse.user", "default")
	viper.SetDefault("database.clickhouse.password", "")
	viper.SetDefault("database.clickhouse.dbname", "aien_db")
	viper.SetDefault("database.clickhouse.secure", false)

	// LLM
	viper.SetDefault("llm.base_url", "http://localhost:8124")
	viper.SetDefault("llm.api_key", "")
	viper.SetDefault("llm.model", "openrouter/auto")
	viper.SetDefault("llm.timeout", "30s")
	viper.SetDefault("llm.max_retries", 3)

	// Storage
	viper.SetDefault("storage.type", "minio")
	viper.SetDefault("storage.endpoint", "localhost:9000")
	viper.SetDefault("storage.access_key", "minioadmin")
	viper.SetDefault("storage.secret_key", "minioadmin")
	viper.SetDefault("storage.bucket", "files")
	viper.SetDefault("storage.use_ssl", false)

	// Airflow
	viper.SetDefault("airflow.dags_path", "/opt/airflow/dags")
	viper.SetDefault("airflow.base_url", "http://localhost:8081")
	viper.SetDefault("airflow.username", "admin")
	viper.SetDefault("airflow.password", "admin")

	// Logging
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
}

// GetPostgreSQLDSN возвращает DSN для PostgreSQL
func (c *Config) GetPostgreSQLDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.PostgreSQL.Host,
		c.Database.PostgreSQL.Port,
		c.Database.PostgreSQL.User,
		c.Database.PostgreSQL.Password,
		c.Database.PostgreSQL.DBName,
		c.Database.PostgreSQL.SSLMode,
	)
}

// GetClickHouseDSN возвращает DSN для ClickHouse
func (c *Config) GetClickHouseDSN() string {
	protocol := "http"
	if c.Database.ClickHouse.Secure {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		protocol,
		c.Database.ClickHouse.User,
		c.Database.ClickHouse.Password,
		c.Database.ClickHouse.Host,
		c.Database.ClickHouse.Port,
		c.Database.ClickHouse.DBName,
	)
}
