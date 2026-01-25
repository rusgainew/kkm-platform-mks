package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	contextHelper "github.com/rusgainew/kkm-project-mks/api-gateway/pkg/context"
	"go.uber.org/zap"
)

// RequestIDMiddleware добавляет request_id к каждому запросу
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем есть ли уже X-Request-ID в заголовке
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = contextHelper.GenerateRequestID()
		}

		// Сохраняем в контексте
		contextHelper.SetRequestID(c, requestID)

		// Добавляем в response header
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// StructuredLoggingMiddleware логирует запросы с полным контекстом
func StructuredLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		// Извлекаем request_id
		requestID := contextHelper.GetRequestID(c)
		traceID := contextHelper.GetTraceID(c)

		c.Next()

		// После обработки запроса
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		bodySize := c.Writer.Size()

		// Извлекаем user_id если есть
		userID, _ := contextHelper.GetUserID(c)

		// Создаем structured log entry
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Int("body_size", bodySize),
			zap.String("request_id", requestID),
		}

		if userID != "" {
			fields = append(fields, zap.String("user_id", userID))
		}

		if traceID != "" {
			fields = append(fields, zap.String("trace_id", traceID))
		}

		// Логируем в зависимости от статуса
		if statusCode >= 500 {
			logger.Error("HTTP request failed", fields...)
		} else if statusCode >= 400 {
			logger.Warn("HTTP request client error", fields...)
		} else {
			logger.Info("HTTP request", fields...)
		}
	}
}

// ContextPropagationMiddleware пробрасывает контекст из Gin в Go context
func ContextPropagationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Создаем новый context с данными из Gin context
		ctx := contextHelper.PropagateToContext(c, c.Request.Context())

		// Обновляем request с новым context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
