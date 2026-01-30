// Файл foreign-company-query-server/internal/infrastructure/config/config.go содержит реализацию пакета config.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCPort        int
	HealthPort      int
	LogLevel        string
	ShutdownTimeout time.Duration
	JaegerEndpoint  string
	RedisURL        string
	CacheTTL        time.Duration
	RabbitMQURL     string
}

func Load() *Config {
	cfg := &Config{
		GRPCPort:        loadPortEnv("FOREIGN_COMPANY_QUERY_GRPC_PORT", 50057),
		HealthPort:      loadPortEnv("FOREIGN_COMPANY_QUERY_HEALTH_PORT", 8057),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		JaegerEndpoint:  getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379/1"),
		CacheTTL:        getEnvDuration("CACHE_TTL", 10*time.Minute),
		RabbitMQURL:     getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
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
	if c.HealthPort < 1 || c.HealthPort > 65535 {
		return fmt.Errorf("invalid Health port: %d (must be 1-65535)", c.HealthPort)
	}
	if c.GRPCPort == c.HealthPort {
		return fmt.Errorf("GRPC port (%d) cannot be same as Health port", c.GRPCPort)
	}
	if c.RedisURL == "" {
		return errors.New("REDIS_URL is required")
	}
	if c.RabbitMQURL == "" {
		return errors.New("RABBITMQ_URL is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadPortEnv(key string, defaultPort int) int {
	if port := os.Getenv(key); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			return p
		}
	}
	return defaultPort
}

func getEnvDuration(key string, defaultDuration time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultDuration
}
