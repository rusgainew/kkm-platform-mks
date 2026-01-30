// Файл api-gateway/internal/infrastructure/middleware/middleware.go содержит реализацию пакета middleware.
package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/auth"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// AuthMiddleware middleware для проверки JWT токена
func AuthMiddleware(authService *auth.Service, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из заголовка
		authHeader := c.GetHeader("Authorization")
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			logger.Warn("Failed to extract token from header",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(401, gin.H{"error": "unauthorized", "message": "missing or invalid authorization header"})
			c.Abort()
			return
		}

		// Валидируем токен
		claims, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			logger.Warn("Token validation failed",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(401, gin.H{"error": "unauthorized", "message": "invalid or expired token"})
			c.Abort()
			return
		}

		// Сохраняем claims в контексте
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("claims", claims)
		// Сохраняем сам token для передачи в gRPC запросы
		c.Set("user_token", token)

		// Добавляем token в контекст запроса для использования в сервисах
		ctx := context.WithValue(c.Request.Context(), "user_token", token)
		ctx = context.WithValue(ctx, "authorization", authHeader)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// MetricsMiddleware middleware для сбора метрик
func MetricsMiddleware(metrics *observability.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Увеличиваем счетчик активных подключений
		metrics.IncrementActiveConnections()
		defer metrics.DecrementActiveConnections()

		// Обрабатываем запрос
		c.Next()

		// Записываем метрики
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		metrics.RecordHTTPRequest(
			c.Request.Method,
			c.FullPath(),
			statusToString(status),
			duration,
		)
	}
}

// TracingMiddleware middleware для трассировки
func TracingMiddleware(tracer *observability.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, endSpan := tracer.StartSpan(c.Request.Context(), c.Request.Method+" "+c.FullPath())
		defer endSpan()

		// Добавляем атрибуты
		tracer.SetAttributes(ctx, map[string]interface{}{
			"http.method":      c.Request.Method,
			"http.url":         c.Request.URL.String(),
			"http.user_agent":  c.Request.UserAgent(),
			"http.remote_addr": c.ClientIP(),
		})

		// Обновляем контекст в запросе
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// Добавляем информацию о статусе
		tracer.SetAttributes(ctx, map[string]interface{}{
			"http.status_code": c.Writer.Status(),
		})

		// Записываем ошибки если они есть
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				tracer.RecordError(ctx, err)
			}
		}
	}
}

// LoggingMiddleware middleware для логирования запросов
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		logger.Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// RateLimitMiddleware middleware для ограничения частоты запросов
func RateLimitMiddleware(requests int, window time.Duration, logger *zap.Logger) gin.HandlerFunc {
	// Создаем rate limiter
	limiter := rate.NewLimiter(rate.Every(window/time.Duration(requests)), requests)

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if !limiter.Allow() {
			logger.Warn("Rate limit exceeded",
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(429, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		// Для более сложного rate limiting можно использовать контекст
		_ = ctx

		c.Next()
	}
}

// CORSMiddleware middleware для обработки CORS
func CORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders []string, maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Проверяем разрешенные origins
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", joinStrings(allowedMethods, ", "))
			c.Header("Access-Control-Allow-Headers", joinStrings(allowedHeaders, ", "))
			c.Header("Access-Control-Max-Age", intToString(maxAge))
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Обрабатываем preflight запросы
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware middleware для восстановления после паники
func RecoveryMiddleware(logger *zap.Logger, metrics *observability.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				metrics.IncrementErrorCount("panic", "api-gateway")

				c.JSON(500, gin.H{
					"error":   "internal_server_error",
					"message": "an unexpected error occurred",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

// TimeoutMiddleware middleware для установки таймаута запроса
func TimeoutMiddleware(timeout time.Duration, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		// Канал для отслеживания завершения
		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// Запрос завершен успешно
		case <-ctx.Done():
			// Таймаут
			logger.Warn("Request timeout",
				zap.String("path", c.Request.URL.Path),
				zap.Duration("timeout", timeout),
			)
			c.JSON(504, gin.H{
				"error":   "timeout",
				"message": "request timeout exceeded",
			})
			c.Abort()
		}
	}
}

// Helper functions

func statusToString(status int) string {
	if status >= 200 && status < 300 {
		return "2xx"
	} else if status >= 300 && status < 400 {
		return "3xx"
	} else if status >= 400 && status < 500 {
		return "4xx"
	} else if status >= 500 {
		return "5xx"
	}
	return "unknown"
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}
