package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит конфигурацию сервиса
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	RabbitMQ      RabbitMQConfig
	Auth          AuthConfig
	Observability ObservabilityConfig
}

// ServerConfig настройки сервера
type ServerConfig struct {
	Host        string
	Port        string
	MetricsPort string
}

// DatabaseConfig настройки базы данных
type DatabaseConfig struct {
	Driver string
	URL    string
}

// RedisConfig настройки Redis
type RedisConfig struct {
	URL string
}

// RabbitMQConfig настройки RabbitMQ
type RabbitMQConfig struct {
	URL string
}

// AuthConfig настройки аутентификации
type AuthConfig struct {
	JWTSecret string
}

// ObservabilityConfig настройки observability
type ObservabilityConfig struct {
	JaegerURL   string
	LogLevel    string
	Environment string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	// Загрузка .env файла (игнорируем ошибку, если файл не существует)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host:        getEnv("CATALOG_SERVER_HOST", "0.0.0.0"),
			Port:        getEnv("CATALOG_SERVER_PORT", "50053"),
			MetricsPort: getEnv("CATALOG_SERVER_METRICS_PORT", "9093"),
		},
		Database: DatabaseConfig{
			Driver: getEnv("CATALOG_SERVER_DB_DRIVER", "postgres"),
			URL:    getEnv("CATALOG_SERVER_DB_URL", "postgres://postgres:postgres@localhost:5432/catalogdb?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("CATALOG_SERVER_REDIS_URL", "redis://localhost:6379/0"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("CATALOG_SERVER_RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("CATALOG_SERVER_JWT_SECRET", getEnv("JWT_SECRET", "your-secret-key-change-in-production")),
		},
		Observability: ObservabilityConfig{
			JaegerURL:   getEnv("CATALOG_SERVER_JAEGER_URL", "http://localhost:14268/api/traces"),
			LogLevel:    getEnv("CATALOG_SERVER_LOG_LEVEL", "info"),
			Environment: getEnv("CATALOG_SERVER_ENVIRONMENT", "development"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
