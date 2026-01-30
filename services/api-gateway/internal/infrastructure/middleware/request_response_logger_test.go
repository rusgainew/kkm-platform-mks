// Файл api-gateway/internal/infrastructure/middleware/request_response_logger_test.go содержит реализацию пакета middleware_test.
package middleware_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestResponseLoggerMiddleware_BasicFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаём observer logger для проверки логов
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	config := middleware.DefaultRequestResponseLoggerConfig()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Выполняем запрос
	body := map[string]string{"test": "data"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем что логи записаны
	assert.Equal(t, http.StatusOK, w.Code)
	assert.GreaterOrEqual(t, logs.Len(), 2) // Минимум 2 лога: incoming и completed

	// Проверяем incoming request log
	incomingLog := logs.All()[0]
	assert.Equal(t, "Incoming HTTP request", incomingLog.Message)
	assert.Contains(t, incomingLog.ContextMap(), "request_id")
	assert.Equal(t, "POST", incomingLog.ContextMap()["method"])
	assert.Equal(t, "/test", incomingLog.ContextMap()["path"])

	// Проверяем completed request log
	completedLog := logs.All()[logs.Len()-1]
	assert.Equal(t, "HTTP request completed", completedLog.Message)
	assert.Equal(t, int64(200), completedLog.ContextMap()["status"])
	assert.Contains(t, completedLog.ContextMap(), "duration")
	assert.Contains(t, completedLog.ContextMap(), "duration_ms")
}

func TestRequestResponseLoggerMiddleware_SensitiveHeadersMasking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	config := middleware.DefaultRequestResponseLoggerConfig()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Запрос с чувствительными заголовками
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	req.Header.Set("Cookie", "session=12345")
	req.Header.Set("X-Api-Key", "my-api-key")
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем что чувствительные заголовки замаскированы
	incomingLog := logs.All()[0]
	requestHeaders := incomingLog.ContextMap()["request_headers"].(string)

	var headers map[string]interface{}
	json.Unmarshal([]byte(requestHeaders), &headers)

	assert.Equal(t, "***MASKED***", headers["Authorization"])
	assert.Equal(t, "***MASKED***", headers["Cookie"])
	assert.Equal(t, "***MASKED***", headers["X-Api-Key"])
	assert.Equal(t, "test-agent", headers["User-Agent"]) // Не чувствительный
}

func TestRequestResponseLoggerMiddleware_BodyLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	config := middleware.DefaultRequestResponseLoggerConfig()
	config.LogRequestBody = true
	config.LogResponseBody = true
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.POST("/test", func(c *gin.Context) {
		var body map[string]interface{}
		c.BindJSON(&body)
		c.JSON(http.StatusOK, gin.H{"received": body["message"]})
	})

	// Запрос с телом
	requestBody := map[string]string{"message": "hello world"}
	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем что тело запроса залогировано
	incomingLog := logs.All()[0]
	assert.Contains(t, incomingLog.ContextMap(), "request_body")
	requestBodyLogged := incomingLog.ContextMap()["request_body"].(string)
	assert.Contains(t, requestBodyLogged, "hello world")

	// Проверяем что тело ответа залогировано
	completedLog := logs.All()[logs.Len()-1]
	assert.Contains(t, completedLog.ContextMap(), "response_body")
	responseBodyLogged := completedLog.ContextMap()["response_body"].(string)
	assert.Contains(t, responseBodyLogged, "hello world")
}

func TestRequestResponseLoggerMiddleware_SensitiveDataRedaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	config := middleware.DefaultRequestResponseLoggerConfig()
	config.LogRequestBody = true
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Запрос с чувствительными данными
	requestBody := map[string]string{
		"username": "testuser",
		"password": "secret123",
		"api_key":  "my-secret-key",
	}
	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем что пароль и api_key замаскированы
	incomingLog := logs.All()[0]
	requestBodyLogged := incomingLog.ContextMap()["request_body"].(string)

	var loggedBody map[string]interface{}
	json.Unmarshal([]byte(requestBodyLogged), &loggedBody)

	assert.Equal(t, "testuser", loggedBody["username"])
	assert.Equal(t, "***REDACTED***", loggedBody["password"])
	assert.Equal(t, "***REDACTED***", loggedBody["api_key"])
}

func TestRequestResponseLoggerMiddleware_SkipPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	config := middleware.DefaultRequestResponseLoggerConfig()
	config.SkipPaths = []string{"/health", "/metrics"}
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "test"})
	})

	// Запрос к /health не должен логироваться
	req1 := httptest.NewRequest("GET", "/health", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	initialLogCount := logs.Len()

	// Запрос к /api/test должен логироваться
	req2 := httptest.NewRequest("GET", "/api/test", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Проверяем что для /health логов не добавилось, а для /api/test добавились
	assert.Equal(t, initialLogCount, 0)     // /health не логировался
	assert.GreaterOrEqual(t, logs.Len(), 2) // /api/test залогирован
}

func TestRequestResponseLoggerMiddleware_ErrorStatusCodes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		statusCode     int
		expectedLevel  string
		expectedLogMsg string
	}{
		{
			name:           "4xx_warning",
			statusCode:     http.StatusBadRequest,
			expectedLevel:  "warn",
			expectedLogMsg: "HTTP request completed with warning",
		},
		{
			name:           "5xx_error",
			statusCode:     http.StatusInternalServerError,
			expectedLevel:  "error",
			expectedLogMsg: "HTTP request completed with error",
		},
		{
			name:           "2xx_info",
			statusCode:     http.StatusOK,
			expectedLevel:  "info",
			expectedLogMsg: "HTTP request completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			router := gin.New()
			config := middleware.DefaultRequestResponseLoggerConfig()
			router.Use(middleware.RequestIDMiddleware())
			router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

			router.GET("/test", func(c *gin.Context) {
				c.JSON(tt.statusCode, gin.H{"message": "test"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Проверяем уровень логирования и сообщение
			completedLog := logs.All()[logs.Len()-1]
			assert.Equal(t, tt.expectedLogMsg, completedLog.Message)
			assert.Equal(t, int64(tt.statusCode), completedLog.ContextMap()["status"])
		})
	}
}

func TestRequestResponseLoggerMiddleware_LargeBodyTruncation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	config := middleware.DefaultRequestResponseLoggerConfig()
	config.LogRequestBody = true
	config.MaxBodyLogSize = 100 // Только 100 байт
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.POST("/test", func(c *gin.Context) {
		io.Copy(io.Discard, c.Request.Body)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Создаём большое тело (больше 100 байт)
	largeBody := make([]byte, 500)
	for i := range largeBody {
		largeBody[i] = 'A'
	}

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(largeBody))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем что тело обрезано
	incomingLog := logs.All()[0]
	if requestBody, ok := incomingLog.ContextMap()["request_body"]; ok {
		bodyStr := requestBody.(string)
		assert.LessOrEqual(t, len(bodyStr), 110) // 100 + "(truncated)"
	}
}

func TestRequestResponseLoggerMiddleware_NonJSONBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	config := middleware.DefaultRequestResponseLoggerConfig()
	config.LogRequestBody = true
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "plain text response")
	})

	// Запрос с plain text
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("plain text body"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем что plain text тело залогировано как есть
	incomingLog := logs.All()[0]
	assert.Contains(t, incomingLog.ContextMap(), "request_body")
	requestBodyLogged := incomingLog.ContextMap()["request_body"].(string)
	assert.Equal(t, "plain text body", requestBodyLogged)
}

func TestRequestResponseLoggerMiddleware_CorrelationID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := gin.New()
	config := middleware.DefaultRequestResponseLoggerConfig()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Запрос с существующим X-Request-ID
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "custom-request-id-123")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем что request_id присутствует в логах
	incomingLog := logs.All()[0]
	completedLog := logs.All()[logs.Len()-1]

	assert.Equal(t, "custom-request-id-123", incomingLog.ContextMap()["request_id"])
	assert.Equal(t, "custom-request-id-123", completedLog.ContextMap()["request_id"])

	// Проверяем что X-Request-ID добавлен в ответ
	assert.Equal(t, "custom-request-id-123", w.Header().Get("X-Request-ID"))
}
