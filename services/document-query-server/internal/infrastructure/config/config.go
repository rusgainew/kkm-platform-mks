package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCPort         int
	MetricsPort      int
	DatabaseURL      string
	DatabaseDriver   string
	LogLevel         string
	ShutdownTimeout  time.Duration
	JaegerEndpoint   string
	RedisURL         string
	CacheTTL         time.Duration
	RabbitMQURL      string
	RabbitMQEnabled  bool
	RabbitMQExchange string
}

func Load() *Config {
	cfg := &Config{
		GRPCPort:         loadPortEnv("DOCUMENT_QUERY_GRPC_PORT", 50062),
		MetricsPort:      loadPortEnv("DOCUMENT_QUERY_METRICS_PORT", 9103),
		DatabaseDriver:   getEnv("DOCUMENT_QUERY_DB_DRIVER", "postgres"),
		DatabaseURL:      getEnv("DOCUMENT_QUERY_DB_URL", ""),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		ShutdownTimeout:  getEnvDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		JaegerEndpoint:   getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		RedisURL:         getEnv("DOCUMENT_QUERY_REDIS_URL", "redis://localhost:6379/3"),
		CacheTTL:         getEnvDuration("DOCUMENT_QUERY_CACHE_TTL", 10*time.Minute),
		RabbitMQURL:      getEnv("DOCUMENT_QUERY_RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		RabbitMQEnabled:  getEnvBool("DOCUMENT_QUERY_RABBITMQ_ENABLED", true),
		RabbitMQExchange: getEnv("DOCUMENT_QUERY_RABBITMQ_EXCHANGE", "documents"),
	}

	return cfg
}

func LoadValidated() (*Config, error) {
	cfg := Load()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.GRPCPort < 1 || c.GRPCPort > 65535 {
		return fmt.Errorf("invalid GRPC port: %d (must be 1-65535)", c.GRPCPort)
	}
	if c.MetricsPort < 1 || c.MetricsPort > 65535 {
		return fmt.Errorf("invalid Metrics port: %d (must be 1-65535)", c.MetricsPort)
	}
	if c.GRPCPort == c.MetricsPort {
		return fmt.Errorf("GRPC port (%d) cannot be same as Metrics port", c.GRPCPort)
	}
	if c.DatabaseURL == "" {
		return errors.New("DOCUMENT_QUERY_DB_URL is required")
	}
	if c.DatabaseDriver == "" {
		return errors.New("DOCUMENT_QUERY_DB_DRIVER is required")
	}
	return nil
}

func loadPortEnv(key string, defaultValue int) int {
	portStr := os.Getenv(key)
	if portStr == "" {
		return defaultValue
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return defaultValue
	}
	return port
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
