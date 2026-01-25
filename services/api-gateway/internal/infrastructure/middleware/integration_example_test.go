package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestRequestResponseLogger_IntegrationExample демонстрирует полный workflow
func TestRequestResponseLogger_IntegrationExample(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаём production-like logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Настраиваем router с полным middleware stack
	router := gin.New()

	// 1. Request ID для correlation
	router.Use(middleware.RequestIDMiddleware())

	// 2. Recovery для обработки паник
	router.Use(gin.Recovery())

	// 3. Наш Request/Response Logger
	config := middleware.DefaultRequestResponseLoggerConfig()
	config.LogRequestBody = true
	config.LogResponseBody = true
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	// Регистрируем тестовый endpoint
	router.POST("/api/v1/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_request",
				"message": "Invalid JSON",
			})
			return
		}

		// Имитация успешного логина
		c.JSON(http.StatusOK, gin.H{
			"token":    "jwt-token-example",
			"user_id":  "123",
			"username": req.Username,
		})
	})

	// Выполняем тестовый запрос
	loginReq := map[string]string{
		"username": "testuser",
		"password": "secret123",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Integration-Test/1.0")
	req.Header.Set("Authorization", "Bearer previous-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "testuser", response["username"])
	assert.Equal(t, "123", response["user_id"])

	// В реальном логе будет:
	// 1. Incoming request с замаскированным Authorization header
	// 2. Request body с замаскированным password полем
	// 3. Completed request с токеном в response body
	// 4. Все с одинаковым request_id для корреляции
}

// TestRequestResponseLogger_ErrorHandling тест обработки ошибок
func TestRequestResponseLogger_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())
	config := middleware.DefaultRequestResponseLoggerConfig()
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	// Endpoint с ошибкой
	router.GET("/api/v1/error", func(c *gin.Context) {
		c.Error(assert.AnError) // Добавляем ошибку в контекст
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Something went wrong",
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/error", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Лог будет с level=error и содержать ошибку из c.Errors
}

// TestRequestResponseLogger_HighTrafficSimulation симуляция high-traffic scenario
func TestRequestResponseLogger_HighTrafficSimulation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping high-traffic simulation in short mode")
	}

	gin.SetMode(gin.ReleaseMode)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())

	// Production config - без логирования тел для производительности
	config := middleware.DefaultRequestResponseLoggerConfig()
	config.LogRequestBody = false
	config.LogResponseBody = false
	config.SkipPaths = []string{"/health"}
	router.Use(middleware.RequestResponseLoggerMiddleware(logger, config))

	router.GET("/api/v1/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "response"})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Симулируем 100 запросов
	for i := 0; i < 100; i++ {
		// API запросы логируются
		req := httptest.NewRequest("GET", "/api/v1/data", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Health checks не логируются (skip path)
		req = httptest.NewRequest("GET", "/health", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Все 100 запросов обработаны с минимальным overhead
}
