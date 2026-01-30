// Файл document-server/internal/infrastructure/config/config.go содержит реализацию пакета config.
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port        string
	MetricsPort string
}

type ObservabilityConfig struct {
	ServiceName    string
	LogLevel       string
	EnableTracing  bool
	JaegerEndpoint string
}

type RabbitMQConfig struct {
	Enabled      bool
	URL          string
	Exchange     string
	ExchangeType string
}

type Config struct {
	Database      DatabaseConfig
	Server        ServerConfig
	Observability ObservabilityConfig
	RabbitMQ      RabbitMQConfig
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "document_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "50054"),
			MetricsPort: getEnv("METRICS_PORT", "9094"),
		},
		Observability: ObservabilityConfig{
			ServiceName:    getEnv("SERVICE_NAME", "document-server"),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			EnableTracing:  parseBool(getEnv("ENABLE_TRACING", "false")),
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://jaeger:14268/v1/traces"),
		},
		RabbitMQ: RabbitMQConfig{
			Enabled:      parseBool(getEnv("RABBITMQ_ENABLED", "true")),
			URL:          getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
			Exchange:     getEnv("RABBITMQ_EXCHANGE", "documents"),
			ExchangeType: getEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func parseBool(value string) bool {
	b, _ := strconv.ParseBool(value)
	return b
}
