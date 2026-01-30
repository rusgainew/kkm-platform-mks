// Файл api-gateway/internal/infrastructure/middleware/request_response_logger.go содержит реализацию пакета middleware.
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	contextHelper "github.com/rusgainew/kkm-project-mks/api-gateway/pkg/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RequestResponseLoggerConfig конфигурация для детального логирования
type RequestResponseLoggerConfig struct {
	// LogRequestBody включает логирование тела запроса
	LogRequestBody bool
	// LogResponseBody включает логирование тела ответа
	LogResponseBody bool
	// MaxBodyLogSize максимальный размер тела для логирования (байты)
	MaxBodyLogSize int64
	// SkipPaths пути которые не нужно логировать детально
	SkipPaths []string
	// SensitiveHeaders заголовки которые нужно маскировать
	SensitiveHeaders []string
	// LogHeaders включает логирование заголовков
	LogHeaders bool
}

// DefaultRequestResponseLoggerConfig возвращает конфигурацию по умолчанию
func DefaultRequestResponseLoggerConfig() *RequestResponseLoggerConfig {
	return &RequestResponseLoggerConfig{
		LogRequestBody:  true,
		LogResponseBody: true,
		MaxBodyLogSize:  10 * 1024, // 10KB
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/swagger",
		},
		SensitiveHeaders: []string{
			"Authorization",
			"Cookie",
			"X-Api-Key",
			"X-Auth-Token",
		},
		LogHeaders: true,
	}
}

// bodyLogWriter оборачивает gin.ResponseWriter для захвата тела ответа
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// RequestResponseLoggerMiddleware детальное логирование HTTP запросов и ответов
func RequestResponseLoggerMiddleware(logger *zap.Logger, config *RequestResponseLoggerConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultRequestResponseLoggerConfig()
	}

	return func(c *gin.Context) {
		// Проверяем нужно ли скипать этот путь
		if shouldSkipPath(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		start := time.Now()
		requestID := contextHelper.GetRequestID(c)
		traceID := contextHelper.GetTraceID(c)

		// Логируем входящий запрос
		logFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("protocol", c.Request.Proto),
		}

		// Логируем заголовки запроса (маскируем чувствительные)
		if config.LogHeaders {
			headers := maskSensitiveHeaders(c.Request.Header, config.SensitiveHeaders)
			if headerJSON, err := json.Marshal(headers); err == nil {
				logFields = append(logFields, zap.String("request_headers", string(headerJSON)))
			}
		}

		// Логируем тело запроса (если есть и размер допустимый)
		var requestBody []byte
		if config.LogRequestBody && c.Request.Body != nil {
			requestBody = captureRequestBody(c, config.MaxBodyLogSize)
			if len(requestBody) > 0 && shouldLogBody(c.Request.Header.Get("Content-Type")) {
				logFields = append(logFields, zap.String("request_body", sanitizeBody(requestBody)))
			}
		}

		logger.Info("Incoming HTTP request", logFields...)

		// Оборачиваем ResponseWriter для захвата ответа
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// Обрабатываем запрос
		c.Next()

		// Вычисляем длительность
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Логируем исходящий ответ
		responseFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.Int("response_size", blw.body.Len()),
			zap.Float64("duration_ms", float64(duration.Milliseconds())),
		}

		// Логируем заголовки ответа
		if config.LogHeaders {
			headers := maskSensitiveHeaders(c.Writer.Header(), config.SensitiveHeaders)
			if headerJSON, err := json.Marshal(headers); err == nil {
				responseFields = append(responseFields, zap.String("response_headers", string(headerJSON)))
			}
		}

		// Логируем тело ответа (если размер допустимый)
		if config.LogResponseBody && blw.body.Len() > 0 && int64(blw.body.Len()) <= config.MaxBodyLogSize {
			responseBody := blw.body.String()
			if shouldLogBody(c.Writer.Header().Get("Content-Type")) {
				responseFields = append(responseFields, zap.String("response_body", sanitizeBody([]byte(responseBody))))
			}
		}

		// Добавляем информацию об ошибках если есть
		if len(c.Errors) > 0 {
			errorMessages := make([]string, len(c.Errors))
			for i, err := range c.Errors {
				errorMessages[i] = err.Error()
			}
			responseFields = append(responseFields, zap.Strings("errors", errorMessages))
		}

		// Уровень логирования зависит от статус кода
		logLevel := determineLogLevel(statusCode)
		switch logLevel {
		case zap.ErrorLevel:
			logger.Error("HTTP request completed with error", responseFields...)
		case zap.WarnLevel:
			logger.Warn("HTTP request completed with warning", responseFields...)
		default:
			logger.Info("HTTP request completed", responseFields...)
		}
	}
}

// captureRequestBody захватывает тело запроса (с лимитом размера)
func captureRequestBody(c *gin.Context, maxSize int64) []byte {
	if c.Request.Body == nil {
		return nil
	}

	// Читаем тело с ограничением размера
	bodyReader := io.LimitReader(c.Request.Body, maxSize)
	bodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil
	}

	// Восстанавливаем тело для дальнейшего использования
	c.Request.Body = io.NopCloser(io.MultiReader(
		bytes.NewReader(bodyBytes),
		c.Request.Body,
	))

	return bodyBytes
}

// maskSensitiveHeaders маскирует чувствительные заголовки
func maskSensitiveHeaders(headers map[string][]string, sensitiveHeaders []string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, values := range headers {
		isSensitive := false
		for _, sensitive := range sensitiveHeaders {
			if strings.EqualFold(key, sensitive) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			result[key] = "***MASKED***"
		} else {
			if len(values) == 1 {
				result[key] = values[0]
			} else {
				result[key] = values
			}
		}
	}

	return result
}

// shouldLogBody определяет нужно ли логировать тело по Content-Type
func shouldLogBody(contentType string) bool {
	// Логируем только текстовые форматы
	textTypes := []string{
		"application/json",
		"application/xml",
		"text/",
		"application/x-www-form-urlencoded",
	}

	contentType = strings.ToLower(contentType)
	for _, textType := range textTypes {
		if strings.Contains(contentType, textType) {
			return true
		}
	}

	return false
}

// sanitizeBody очищает тело от чувствительных данных (пароли, токены)
func sanitizeBody(body []byte) string {
	bodyStr := string(body)

	// Проверяем что это валидный JSON
	var data interface{}
	if err := json.Unmarshal(body, &data); err == nil {
		// Это JSON - маскируем чувствительные поля
		sanitized := sanitizeJSON(data)
		if sanitizedBytes, err := json.Marshal(sanitized); err == nil {
			return string(sanitizedBytes)
		}
	}

	// Если не JSON - возвращаем как есть (с ограничением длины)
	if len(bodyStr) > 1000 {
		return bodyStr[:1000] + "... (truncated)"
	}
	return bodyStr
}

// sanitizeJSON рекурсивно маскирует чувствительные поля в JSON
func sanitizeJSON(data interface{}) interface{} {
	sensitiveKeys := []string{
		"password",
		"secret",
		"token",
		"api_key",
		"apikey",
		"authorization",
		"access_token",
		"refresh_token",
		"private_key",
		"credit_card",
		"ssn",
	}

	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			lowerKey := strings.ToLower(key)
			isSensitive := false
			for _, sensitiveKey := range sensitiveKeys {
				if strings.Contains(lowerKey, sensitiveKey) {
					isSensitive = true
					break
				}
			}

			if isSensitive {
				result[key] = "***REDACTED***"
			} else {
				result[key] = sanitizeJSON(value)
			}
		}
		return result

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = sanitizeJSON(item)
		}
		return result

	default:
		return v
	}
}

// shouldSkipPath проверяет нужно ли скипать логирование для данного пути
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// determineLogLevel определяет уровень логирования по статус коду
func determineLogLevel(statusCode int) zapcore.Level {
	switch {
	case statusCode >= 500:
		return zapcore.ErrorLevel
	case statusCode >= 400:
		return zapcore.WarnLevel
	default:
		return zapcore.InfoLevel
	}
}
