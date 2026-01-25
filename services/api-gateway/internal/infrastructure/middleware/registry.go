package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/config"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

// MiddlewareConfig содержит конфигурацию для middleware
type MiddlewareConfig struct {
	cfg     *config.Config
	logger  *zap.Logger
	metrics *observability.Metrics
	tracer  *observability.Tracer
}

// NewMiddlewareConfig создает новую конфигурацию middleware
func NewMiddlewareConfig(
	cfg *config.Config,
	logger *zap.Logger,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
) *MiddlewareConfig {
	return &MiddlewareConfig{
		cfg:     cfg,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

// MiddlewareRegistry отвечает за регистрацию всех middleware
type MiddlewareRegistry struct {
	cfg *MiddlewareConfig
}

// NewMiddlewareRegistry создает новый реестр middleware
func NewMiddlewareRegistry(cfg *MiddlewareConfig) *MiddlewareRegistry {
	return &MiddlewareRegistry{
		cfg: cfg,
	}
}

// RegisterAll регистрирует все global middleware в правильном порядке
func (mr *MiddlewareRegistry) RegisterAll(router *gin.Engine) error {
	mr.registerCoreMiddleware(router)

	if mr.cfg.cfg.Observability.EnableMetrics {
		mr.registerMetricsMiddleware(router)
	}

	if mr.cfg.cfg.Observability.EnableTracing {
		mr.registerTracingMiddleware(router)
	}

	if mr.cfg.cfg.CORS.Enabled {
		mr.registerCORSMiddleware(router)
	}

	if mr.cfg.cfg.RateLimit.Enabled {
		mr.registerRateLimitMiddleware(router)
	}

	mr.registerTimeoutMiddleware(router)

	return nil
}

// registerCoreMiddleware регистрирует основной middleware (recovery, request_id, logging, context)
func (mr *MiddlewareRegistry) registerCoreMiddleware(router *gin.Engine) {
	router.Use(RequestIDMiddleware())
	router.Use(RecoveryMiddleware(mr.cfg.logger, mr.cfg.metrics))

	// Детальное логирование запросов/ответов с correlation ID
	loggerConfig := DefaultRequestResponseLoggerConfig()
	// В продакшене можно отключить логирование тел запросов/ответов на основе log level
	if mr.cfg.cfg.Observability.LogLevel == "info" || mr.cfg.cfg.Observability.LogLevel == "warn" || mr.cfg.cfg.Observability.LogLevel == "error" {
		loggerConfig.LogRequestBody = false
		loggerConfig.LogResponseBody = false
	}
	router.Use(RequestResponseLoggerMiddleware(mr.cfg.logger, loggerConfig))

	router.Use(ContextPropagationMiddleware())
}

// registerMetricsMiddleware регистрирует metrics middleware
func (mr *MiddlewareRegistry) registerMetricsMiddleware(router *gin.Engine) {
	router.Use(MetricsMiddleware(mr.cfg.metrics))
}

// registerTracingMiddleware регистрирует tracing middleware
func (mr *MiddlewareRegistry) registerTracingMiddleware(router *gin.Engine) {
	router.Use(TracingMiddleware(mr.cfg.tracer))
}

// registerCORSMiddleware регистрирует CORS middleware
func (mr *MiddlewareRegistry) registerCORSMiddleware(router *gin.Engine) {
	router.Use(CORSMiddleware(
		mr.cfg.cfg.CORS.AllowedOrigins,
		mr.cfg.cfg.CORS.AllowedMethods,
		mr.cfg.cfg.CORS.AllowedHeaders,
		mr.cfg.cfg.CORS.MaxAge,
	))
}

// registerRateLimitMiddleware регистрирует rate limit middleware
func (mr *MiddlewareRegistry) registerRateLimitMiddleware(router *gin.Engine) {
	router.Use(RateLimitMiddleware(
		mr.cfg.cfg.RateLimit.Requests,
		mr.cfg.cfg.RateLimit.Window,
		mr.cfg.logger,
	))
}

// registerTimeoutMiddleware регистрирует timeout middleware
func (mr *MiddlewareRegistry) registerTimeoutMiddleware(router *gin.Engine) {
	router.Use(TimeoutMiddleware(mr.cfg.cfg.Timeouts.Request, mr.cfg.logger))
}
