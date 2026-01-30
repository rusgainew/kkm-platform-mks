// Файл foreign-company-server/internal/infrastructure/config/config.go содержит реализацию пакета config.
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
	// Загрузка .env файла (игнорируем ошибку если файл не существует)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host:        getEnv("FOREIGN_COMPANY_SERVER_HOST", "0.0.0.0"),
			Port:        getEnv("FOREIGN_COMPANY_SERVER_PORT", "50056"),
			MetricsPort: getEnv("FOREIGN_COMPANY_SERVER_METRICS_PORT", "9096"),
		},
		Database: DatabaseConfig{
			Driver: getEnv("FOREIGN_COMPANY_SERVER_DB_DRIVER", "postgres"),
			URL:    getEnv("FOREIGN_COMPANY_SERVER_DB_URL", ""),
		},
		Redis: RedisConfig{
			URL: getEnv("FOREIGN_COMPANY_SERVER_REDIS_URL", "redis://localhost:6379"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("FOREIGN_COMPANY_SERVER_RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("FOREIGN_COMPANY_SERVER_JWT_SECRET", getEnv("JWT_SECRET", "")),
		},
		Observability: ObservabilityConfig{
			JaegerURL:   getEnv("FOREIGN_COMPANY_SERVER_JAEGER_URL", "http://localhost:14268/api/traces"),
			LogLevel:    getEnv("FOREIGN_COMPANY_SERVER_LOG_LEVEL", "info"),
			Environment: getEnv("FOREIGN_COMPANY_SERVER_ENVIRONMENT", "development"),
		},
	}

	// Валидация обязательных параметров
	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("FOREIGN_COMPANY_SERVER_DB_URL is required")
	}
	if cfg.Auth.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
