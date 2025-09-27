package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type StorageConfig struct {
	Type      string `mapstructure:"type"` // "minio", "s3", "local"
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
	LocalPath string `mapstructure:"local_path"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Установка значений по умолчанию
	viper.SetDefault("server.port", "50054")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("database.host", "postgres")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "fileuser")
	viper.SetDefault("database.password", "filepass")
	viper.SetDefault("database.dbname", "fileservice")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("storage.type", "minio")
	viper.SetDefault("storage.endpoint", "minio:9000")
	viper.SetDefault("storage.access_key", "minioadmin")
	viper.SetDefault("storage.secret_key", "minioadmin")
	viper.SetDefault("storage.bucket", "files")
	viper.SetDefault("storage.use_ssl", false)
	viper.SetDefault("storage.local_path", "/tmp/files")

	// Чтение переменных окружения
	viper.AutomaticEnv()

	// Чтение конфигурационного файла
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Ошибка чтения конфигурации: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Ошибка парсинга конфигурации: %v", err)
	}

	// Переопределение из переменных окружения
	if port := os.Getenv("FILE_SERVICE_PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("FILE_SERVICE_HOST"); host != "" {
		config.Server.Host = host
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		config.Database.Port = dbPort
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}

	return &config
}
