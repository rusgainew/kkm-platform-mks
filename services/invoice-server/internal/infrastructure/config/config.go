// Файл invoice-server/internal/infrastructure/config/config.go содержит реализацию пакета config.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит конфигурацию приложения
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	RabbitMQ      RabbitMQConfig
	Auth          AuthConfig
	Observability ObservabilityConfig
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Port        string
	MetricsPort string
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RabbitMQConfig конфигурация RabbitMQ
type RabbitMQConfig struct {
	URL          string
	Exchange     string
	ExchangeType string
	Enabled      bool
}

// AuthConfig настройки аутентификации
type AuthConfig struct {
	JWTSecret string
}

// ObservabilityConfig конфигурация observability
type ObservabilityConfig struct {
	JaegerEndpoint string
	ServiceName    string
	LogLevel       string
	EnableTracing  bool
	EnableMetrics  bool
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("GRPC_PORT", "50053"),
			MetricsPort: getEnv("METRICS_PORT", "9093"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "invoice_db"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		RabbitMQ: RabbitMQConfig{
			URL:          getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Exchange:     getEnv("RABBITMQ_EXCHANGE", "invoice_events"),
			ExchangeType: getEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
			Enabled:      getEnvAsBool("ENABLE_EVENTS", true),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("INVOICE_SERVER_JWT_SECRET", getEnv("JWT_SECRET", "")),
		},
		Observability: ObservabilityConfig{
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName:    getEnv("SERVICE_NAME", "invoice-server"),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			EnableTracing:  getEnvAsBool("ENABLE_TRACING", true),
			EnableMetrics:  getEnvAsBool("ENABLE_METRICS", true),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	return nil
}

// GetDSN возвращает строку подключения к базе данных
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
