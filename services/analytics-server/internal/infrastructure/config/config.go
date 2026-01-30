// Файл analytics-server/internal/infrastructure/config/config.go содержит реализацию пакета config.
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
	Redis         RedisConfig
	Observability ObservabilityConfig
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	GRPCPort    string
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

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	TTL      int // TTL в секундах для кеша
}

// ObservabilityConfig конфигурация observability
type ObservabilityConfig struct {
	ServiceName    string
	LogLevel       string
	EnableMetrics  bool
	EnableTracing  bool
	JaegerEndpoint string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			GRPCPort:    getEnv("GRPC_PORT", "50070"),
			MetricsPort: getEnv("METRICS_PORT", "9115"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "analytics_svc"),
			Password:        getEnv("DB_PASSWORD", "analytics_pass_2026"),
			DBName:          getEnv("DB_NAME", "analytics_db"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			TTL:      getEnvAsInt("REDIS_TTL", 300), // 5 минут
		},
		Observability: ObservabilityConfig{
			ServiceName:    getEnv("SERVICE_NAME", "analytics-server"),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			EnableMetrics:  getEnvAsBool("ENABLE_METRICS", true),
			EnableTracing:  getEnvAsBool("ENABLE_TRACING", false),
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
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
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
