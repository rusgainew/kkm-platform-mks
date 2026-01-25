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
	JWT           JWTConfig
	Services      ServicesConfig
	Observability ObservabilityConfig
	RateLimit     RateLimitConfig
	CORS          CORSConfig
	Timeouts      TimeoutsConfig
	Redis         *RedisConfig
	TLS           TLSConfig
}

// ServerConfig конфигурация HTTP сервера
type ServerConfig struct {
	HTTPPort    string
	MetricsPort string
}

// JWTConfig конфигурация JWT
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// ServicesConfig конфигурация адресов gRPC сервисов
type ServicesConfig struct {
	UserService             string
	UserQueryService        string
	CompanyService          string
	InvoiceService          string
	CatalogService          string
	BankAccountService      string
	ForeignCompanyService   string
	InvoiceQueryService     string
	CatalogQueryService     string
	BankAccountQueryService string
	DocumentService         string
	DocumentQueryService    string
}

// ObservabilityConfig конфигурация observability
type ObservabilityConfig struct {
	JaegerEndpoint string
	ServiceName    string
	LogLevel       string
	EnableTracing  bool
	EnableMetrics  bool
}

// RateLimitConfig конфигурация rate limiting
type RateLimitConfig struct {
	Enabled  bool
	Requests int
	Window   time.Duration
}

// CORSConfig конфигурация CORS
type CORSConfig struct {
	Enabled        bool
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

// TimeoutsConfig конфигурация таймаутов
type TimeoutsConfig struct {
	Request  time.Duration
	GRPCDial time.Duration
}

// TLSConfig конфигурация TLS для gRPC соединений
type TLSConfig struct {
	Enabled            bool
	CertFile           string
	KeyFile            string
	CAFile             string
	InsecureSkipVerify bool
}

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Addr     string
	Password string
	TTL      int // в секундах
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	// Загрузка .env файла (игнорируем ошибку если файл не найден)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			HTTPPort:    getEnv("HTTP_PORT", "8080"),
			MetricsPort: getEnv("METRICS_PORT", "9090"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		Services: ServicesConfig{
			UserService:             getEnv("USER_SERVICE_URL", "localhost:50051"),
			UserQueryService:        getEnv("USER_QUERY_SERVICE_URL", "localhost:50061"),
			CompanyService:          getEnv("COMPANY_SERVICE_URL", "localhost:50052"),
			InvoiceService:          getEnv("INVOICE_SERVICE_URL", "localhost:50055"),
			CatalogService:          getEnv("CATALOG_SERVICE_URL", "localhost:50053"),
			BankAccountService:      getEnv("BANK_ACCOUNT_SERVICE_URL", "localhost:50054"),
			ForeignCompanyService:   getEnv("FOREIGN_COMPANY_SERVICE_URL", "localhost:50056"),
			InvoiceQueryService:     getEnv("INVOICE_QUERY_SERVICE_URL", "localhost:50057"),
			CatalogQueryService:     getEnv("CATALOG_QUERY_SERVICE_URL", "localhost:50058"),
			BankAccountQueryService: getEnv("BANK_ACCOUNT_QUERY_SERVICE_URL", "localhost:50059"),
			DocumentService:         getEnv("DOCUMENT_SERVICE_URL", "localhost:50060"),
			DocumentQueryService:    getEnv("DOCUMENT_QUERY_SERVICE_URL", "localhost:50062"),
		},
		Observability: ObservabilityConfig{
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName:    getEnv("SERVICE_NAME", "api-gateway"),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			EnableTracing:  getEnvAsBool("ENABLE_TRACING", true),
			EnableMetrics:  getEnvAsBool("ENABLE_METRICS", true),
		},
		RateLimit: RateLimitConfig{
			Enabled:  getEnvAsBool("RATE_LIMIT_ENABLED", true),
			Requests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getEnvAsDuration("RATE_LIMIT_WINDOW", 1*time.Minute),
		},
		CORS: CORSConfig{
			Enabled:        getEnvAsBool("ENABLE_CORS", true),
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:8080"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
			MaxAge:         getEnvAsInt("CORS_MAX_AGE", 3600),
		},
		Timeouts: TimeoutsConfig{
			Request:  getEnvAsDuration("REQUEST_TIMEOUT", 30*time.Second),
			GRPCDial: getEnvAsDuration("GRPC_DIAL_TIMEOUT", 5*time.Second),
		},
		Redis: &RedisConfig{
			Addr:     getEnv("REDIS_ADDR", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			TTL:      getEnvAsInt("REDIS_TTL", 300), // 5 минут по умолчанию
		},
		TLS: TLSConfig{
			Enabled:            getEnvAsBool("GRPC_TLS_ENABLED", false),
			CertFile:           getEnv("TLS_CERT_FILE", ""),
			KeyFile:            getEnv("TLS_KEY_FILE", ""),
			CAFile:             getEnv("TLS_CA_FILE", ""),
			InsecureSkipVerify: getEnvAsBool("TLS_INSECURE_SKIP_VERIFY", false),
		},
	}

	// Если Redis адрес пуст, устанавливаем Redis в nil (кеш отключен)
	if cfg.Redis.Addr == "" {
		cfg.Redis = nil
	}

	// Валидация обязательных параметров
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.JWT.Secret == "" || c.JWT.Secret == "your-secret-key-change-in-production" {
		return fmt.Errorf("JWT_SECRET must be set to a secure value")
	}
	if c.Server.HTTPPort == "" {
		return fmt.Errorf("HTTP_PORT must be set")
	}
	return nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает значение переменной окружения как int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsBool получает значение переменной окружения как bool
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsDuration получает значение переменной окружения как time.Duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsSlice получает значение переменной окружения как слайс строк
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	// Разделяем по запятой
	var result []string
	for _, item := range splitAndTrim(valueStr, ",") {
		if item != "" {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}

// splitAndTrim разделяет строку и удаляет пробелы
func splitAndTrim(s, sep string) []string {
	var result []string
	for _, item := range splitString(s, sep) {
		trimmed := trimSpace(item)
		result = append(result, trimmed)
	}
	return result
}

// splitString разделяет строку по разделителю
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	var current string

	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)

	return result
}

// trimSpace удаляет пробелы в начале и конце строки
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// Удаляем пробелы в начале
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// Удаляем пробелы в конце
	for start < end && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
