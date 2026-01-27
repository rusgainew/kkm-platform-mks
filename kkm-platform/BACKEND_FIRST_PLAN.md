# 📋 План доработки: Backend-First Development

**Дата:** 27 января 2026 г.  
**Версия:** 2.0  
**Подход:** Реализация бизнес-логики в api-gateway → Интеграция в kkm-platform

---

## 🎯 Принцип Backend-First

**Для каждой функции:**

1. ✅ **Backend** - Реализация в `services/api-gateway` (Go)
2. ✅ **Testing** - Unit + Integration тесты
3. ✅ **API Contract** - Swagger/OpenAPI документация
4. ✅ **Frontend** - Интеграция в `kkm-platform` (Next.js/TypeScript)
5. ✅ **E2E** - Полное end-to-end тестирование

---

## 🚀 ФАЗА 1: Rate Limiting и защита API

**Приоритет:** HIGH 🔥  
**Backend:** 8-10 часов  
**Frontend:** 4-6 часов  
**Итого:** 12-16 часов

---

### 📍 1.1. Backend Implementation (Go)

#### Архитектура

```
services/api-gateway/
├── internal/
│   ├── domain/
│   │   └── ports/
│   │       └── rate_limiter.go           # Интерфейс
│   ├── infrastructure/
│   │   └── ratelimit/
│   │       ├── token_bucket.go           # In-memory алгоритм
│   │       ├── redis_limiter.go          # Distributed Redis limiter
│   │       └── limiter_factory.go        # Фабрика лимитеров
│   └── interfaces/
│       └── http/
│           └── middleware/
│               └── rate_limit.go         # Gin middleware
└── cmd/api/main.go                       # Интеграция
```

---

#### Шаг 1.1.1: Domain Layer - Интерфейс

**Файл:** `internal/domain/ports/rate_limiter.go`

```go
package ports

import (
	"context"
	"time"
)

// RateLimiter интерфейс для контроля частоты запросов
type RateLimiter interface {
	// Allow проверяет, разрешен ли один запрос
	Allow(ctx context.Context, key string) (allowed bool, err error)

	// GetLimit возвращает информацию о лимите
	GetLimit(ctx context.Context, key string) (limit, remaining int, resetAt time.Time, err error)

	// Reset сбрасывает счетчик для ключа
	Reset(ctx context.Context, key string) error
}

// RateLimiterConfig конфигурация rate limiter
type RateLimiterConfig struct {
	RequestsPerSecond int           // Запросов в секунду
	BurstSize         int           // Размер burst
	WindowSize        time.Duration // Размер окна (для sliding window)
	Strategy          string        // "token_bucket" или "redis"
}
```

**Тест:** `internal/domain/ports/rate_limiter_test.go`

```go
package ports_test

import (
	"context"
	"testing"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"github.com/stretchr/testify/assert"
)

// MockRateLimiter для тестирования
type MockRateLimiter struct {
	AllowFunc    func(ctx context.Context, key string) (bool, error)
	GetLimitFunc func(ctx context.Context, key string) (int, int, time.Time, error)
	ResetFunc    func(ctx context.Context, key string) error
}

func (m *MockRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	if m.AllowFunc != nil {
		return m.AllowFunc(ctx, key)
	}
	return true, nil
}

func (m *MockRateLimiter) GetLimit(ctx context.Context, key string) (int, int, time.Time, error) {
	if m.GetLimitFunc != nil {
		return m.GetLimitFunc(ctx, key)
	}
	return 100, 50, time.Now().Add(time.Minute), nil
}

func (m *MockRateLimiter) Reset(ctx context.Context, key string) error {
	if m.ResetFunc != nil {
		return m.ResetFunc(ctx, key)
	}
	return nil
}

func TestMockRateLimiter(t *testing.T) {
	mock := &MockRateLimiter{}

	ctx := context.Background()
	allowed, err := mock.Allow(ctx, "test-key")

	assert.NoError(t, err)
	assert.True(t, allowed)
}
```

---

#### Шаг 1.1.2: Infrastructure Layer - Token Bucket

**Файл:** `internal/infrastructure/ratelimit/token_bucket.go`

```go
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// TokenBucket реализует алгоритм Token Bucket
type TokenBucket struct {
	rate       float64                 // Скорость пополнения (токенов/сек)
	capacity   int                     // Максимальная емкость
	tokens     map[string]float64      // Текущие токены
	lastUpdate map[string]time.Time    // Последнее обновление
	mu         sync.RWMutex
}

// NewTokenBucket создает новый Token Bucket limiter
func NewTokenBucket(requestsPerSecond, burstSize int) *TokenBucket {
	return &TokenBucket{
		rate:       float64(requestsPerSecond),
		capacity:   burstSize,
		tokens:     make(map[string]float64),
		lastUpdate: make(map[string]time.Time),
	}
}

// Allow проверяет доступность запроса
func (tb *TokenBucket) Allow(ctx context.Context, key string) (bool, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()

	// Инициализация для нового ключа
	if _, exists := tb.tokens[key]; !exists {
		tb.tokens[key] = float64(tb.capacity)
		tb.lastUpdate[key] = now
		return true, nil
	}

	// Пополняем токены
	elapsed := now.Sub(tb.lastUpdate[key]).Seconds()
	tb.tokens[key] = min(tb.tokens[key]+elapsed*tb.rate, float64(tb.capacity))
	tb.lastUpdate[key] = now

	// Проверяем доступность
	if tb.tokens[key] >= 1.0 {
		tb.tokens[key] -= 1.0
		return true, nil
	}

	return false, nil
}

// GetLimit возвращает информацию о лимите
func (tb *TokenBucket) GetLimit(ctx context.Context, key string) (limit, remaining int, resetAt time.Time, err error) {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	tokens, exists := tb.tokens[key]
	if !exists {
		tokens = float64(tb.capacity)
	}

	// Время до полного восстановления
	timeToFull := (float64(tb.capacity) - tokens) / tb.rate
	resetAt = time.Now().Add(time.Duration(timeToFull * float64(time.Second)))

	return tb.capacity, int(tokens), resetAt, nil
}

// Reset сбрасывает счетчик
func (tb *TokenBucket) Reset(ctx context.Context, key string) error {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.tokens[key] = float64(tb.capacity)
	tb.lastUpdate[key] = time.Now()
	return nil
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
```

**Тест:** `internal/infrastructure/ratelimit/token_bucket_test.go`

```go
package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/ratelimit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenBucket_Allow(t *testing.T) {
	tb := ratelimit.NewTokenBucket(10, 10) // 10 req/sec, burst 10
	ctx := context.Background()

	t.Run("First requests should be allowed", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			allowed, err := tb.Allow(ctx, "test-key")
			require.NoError(t, err)
			assert.True(t, allowed, "Request %d should be allowed", i+1)
		}
	})

	t.Run("Burst exceeded - should be denied", func(t *testing.T) {
		allowed, err := tb.Allow(ctx, "test-key")
		require.NoError(t, err)
		assert.False(t, allowed, "Request after burst should be denied")
	})

	t.Run("After wait - should be allowed", func(t *testing.T) {
		time.Sleep(200 * time.Millisecond) // Ждем пополнения токенов
		allowed, err := tb.Allow(ctx, "test-key")
		require.NoError(t, err)
		assert.True(t, allowed, "Request after refill should be allowed")
	})
}

func TestTokenBucket_GetLimit(t *testing.T) {
	tb := ratelimit.NewTokenBucket(10, 10)
	ctx := context.Background()

	limit, remaining, resetAt, err := tb.GetLimit(ctx, "test-key")

	require.NoError(t, err)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 10, remaining)
	assert.True(t, resetAt.After(time.Now()))
}

func TestTokenBucket_Reset(t *testing.T) {
	tb := ratelimit.NewTokenBucket(10, 10)
	ctx := context.Background()

	// Исчерпываем токены
	for i := 0; i < 10; i++ {
		tb.Allow(ctx, "test-key")
	}

	// Проверяем, что больше нельзя
	allowed, _ := tb.Allow(ctx, "test-key")
	assert.False(t, allowed)

	// Сбрасываем
	err := tb.Reset(ctx, "test-key")
	require.NoError(t, err)

	// Проверяем, что снова можно
	allowed, _ = tb.Allow(ctx, "test-key")
	assert.True(t, allowed)
}
```

---

#### Шаг 1.1.3: Infrastructure Layer - Redis Limiter

**Файл:** `internal/infrastructure/ratelimit/redis_limiter.go`

```go
package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter distributed rate limiter на основе Redis
type RedisRateLimiter struct {
	client    *redis.Client
	rate      int           // Максимум запросов
	window    time.Duration // За период
	keyPrefix string
}

// NewRedisRateLimiter создает Redis limiter
func NewRedisRateLimiter(client *redis.Client, rate int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:    client,
		rate:      rate,
		window:    window,
		keyPrefix: "ratelimit:",
	}
}

// Allow проверяет доступность запроса (Sliding Window алгоритм)
func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := rl.keyPrefix + key
	now := time.Now()
	windowStart := now.Add(-rl.window)

	pipe := rl.client.Pipeline()

	// Удаляем старые записи
	pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", windowStart.UnixNano()))

	// Добавляем текущий запрос
	pipe.ZAdd(ctx, redisKey, redis.Z{
		Score:  float64(now.UnixNano()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})

	// Подсчитываем количество
	pipe.ZCard(ctx, redisKey)

	// TTL для автоочистки
	pipe.Expire(ctx, redisKey, rl.window*2)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// Результат ZCard (3-я команда)
	countCmd := cmds[2].(*redis.IntCmd)
	count, err := countCmd.Result()
	if err != nil {
		return false, err
	}

	return count <= int64(rl.rate), nil
}

// GetLimit возвращает информацию о лимите
func (rl *RedisRateLimiter) GetLimit(ctx context.Context, key string) (limit, remaining int, resetAt time.Time, err error) {
	redisKey := rl.keyPrefix + key
	now := time.Now()
	windowStart := now.Add(-rl.window)

	count, err := rl.client.ZCount(ctx, redisKey,
		fmt.Sprintf("%d", windowStart.UnixNano()),
		fmt.Sprintf("%d", now.UnixNano())).Result()
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	remaining = rl.rate - int(count)
	if remaining < 0 {
		remaining = 0
	}

	resetAt = now.Add(rl.window)

	return rl.rate, remaining, resetAt, nil
}

// Reset сбрасывает счетчик
func (rl *RedisRateLimiter) Reset(ctx context.Context, key string) error {
	redisKey := rl.keyPrefix + key
	return rl.client.Del(ctx, redisKey).Err()
}
```

**Integration Test:** `internal/infrastructure/ratelimit/redis_limiter_integration_test.go`

```go
// +build integration

package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/ratelimit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisRateLimiter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Подключаемся к тестовому Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Тестовая БД
	})
	defer client.Close()

	ctx := context.Background()

	// Проверяем подключение
	err := client.Ping(ctx).Err()
	require.NoError(t, err, "Redis должен быть доступен для integration тестов")

	limiter := ratelimit.NewRedisRateLimiter(client, 5, 1*time.Second)

	t.Run("Allow up to rate limit", func(t *testing.T) {
		testKey := "test-key-1"
		client.Del(ctx, "ratelimit:"+testKey) // Очистка

		// Первые 5 запросов должны пройти
		for i := 0; i < 5; i++ {
			allowed, err := limiter.Allow(ctx, testKey)
			require.NoError(t, err)
			assert.True(t, allowed, "Request %d should be allowed", i+1)
		}

		// 6-й запрос должен быть отклонен
		allowed, err := limiter.Allow(ctx, testKey)
		require.NoError(t, err)
		assert.False(t, allowed, "Request 6 should be denied")
	})

	t.Run("Reset after window", func(t *testing.T) {
		testKey := "test-key-2"
		client.Del(ctx, "ratelimit:"+testKey)

		// Исчерпываем лимит
		for i := 0; i < 5; i++ {
			limiter.Allow(ctx, testKey)
		}

		// Ждем окончания окна
		time.Sleep(1100 * time.Millisecond)

		// Должен быть разрешен
		allowed, err := limiter.Allow(ctx, testKey)
		require.NoError(t, err)
		assert.True(t, allowed, "After window reset, should be allowed")
	})
}
```

---

#### Шаг 1.1.4: Interface Layer - Middleware

**Файл:** `internal/interfaces/http/middleware/rate_limit.go`

```go
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"go.uber.org/zap"
)

// RateLimitMiddleware создает middleware для контроля частоты запросов
func RateLimitMiddleware(
	limiter ports.RateLimiter,
	logger *zap.Logger,
	keyFunc func(*gin.Context) string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)

		// Проверяем лимит
		allowed, err := limiter.Allow(c.Request.Context(), key)
		if err != nil {
			logger.Error("Rate limiter error",
				zap.Error(err),
				zap.String("key", key),
				zap.String("ip", c.ClientIP()),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Rate limiter error",
			})
			c.Abort()
			return
		}

		// Получаем информацию для заголовков
		limit, remaining, resetAt, _ := limiter.GetLimit(c.Request.Context(), key)

		// Добавляем стандартные X-RateLimit заголовки
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt.Unix()))

		if !allowed {
			retryAfter := int(time.Until(resetAt).Seconds())
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))

			logger.Warn("Rate limit exceeded",
				zap.String("key", key),
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate_limit_exceeded",
				"message":     "Too many requests, please try again later",
				"retry_after": retryAfter,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Key generation functions

// IPBasedKey использует IP адрес клиента
func IPBasedKey(c *gin.Context) string {
	return c.ClientIP()
}

// UserBasedKey использует user ID из JWT
func UserBasedKey(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return c.ClientIP() // Fallback на IP
	}
	return fmt.Sprintf("user:%v", userID)
}

// EndpointBasedKey комбинирует IP + endpoint
func EndpointBasedKey(c *gin.Context) string {
	return fmt.Sprintf("%s:%s:%s", c.ClientIP(), c.Request.Method, c.Request.URL.Path)
}
```

**Тест:** `internal/interfaces/http/middleware/rate_limit_test.go`

```go
package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockLimiter struct {
	allowFunc    func(ctx context.Context, key string) (bool, error)
	getLimitFunc func(ctx context.Context, key string) (int, int, time.Time, error)
}

func (m *mockLimiter) Allow(ctx context.Context, key string) (bool, error) {
	if m.allowFunc != nil {
		return m.allowFunc(ctx, key)
	}
	return true, nil
}

func (m *mockLimiter) GetLimit(ctx context.Context, key string) (int, int, time.Time, error) {
	if m.getLimitFunc != nil {
		return m.getLimitFunc(ctx, key)
	}
	return 100, 50, time.Now().Add(time.Minute), nil
}

func (m *mockLimiter) Reset(ctx context.Context, key string) error {
	return nil
}

func TestRateLimitMiddleware_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := &mockLimiter{
		allowFunc: func(ctx context.Context, key string) (bool, error) {
			return true, nil
		},
	}

	logger := zap.NewNop()

	router := gin.New()
	router.Use(middleware.RateLimitMiddleware(limiter, logger, middleware.IPBasedKey))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header(), "X-RateLimit-Limit")
	assert.Contains(t, w.Header(), "X-RateLimit-Remaining")
}

func TestRateLimitMiddleware_Denied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := &mockLimiter{
		allowFunc: func(ctx context.Context, key string) (bool, error) {
			return false, nil
		},
		getLimitFunc: func(ctx context.Context, key string) (int, int, time.Time, error) {
			return 100, 0, time.Now().Add(time.Minute), nil
		},
	}

	logger := zap.NewNop()

	router := gin.New()
	router.Use(middleware.RateLimitMiddleware(limiter, logger, middleware.IPBasedKey))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Header(), "Retry-After")
	assert.Contains(t, w.Body.String(), "rate_limit_exceeded")
}
```

---

#### Шаг 1.1.5: Configuration

**Файл:** `internal/infrastructure/config/config.go` (добавить)

```go
type RateLimitConfig struct {
	Enabled           bool          `env:"RATE_LIMIT_ENABLED" envDefault:"true"`
	Strategy          string        `env:"RATE_LIMIT_STRATEGY" envDefault:"memory"` // "memory" или "redis"
	RequestsPerSecond int           `env:"RATE_LIMIT_RPS" envDefault:"100"`
	BurstSize         int           `env:"RATE_LIMIT_BURST" envDefault:"200"`
	WindowSeconds     int           `env:"RATE_LIMIT_WINDOW_SEC" envDefault:"60"`

	// Специфичные лимиты
	AuthRPS       int `env:"RATE_LIMIT_AUTH_RPS" envDefault:"5"`
	AuthBurst     int `env:"RATE_LIMIT_AUTH_BURST" envDefault:"10"`
}

type Config struct {
	// ... существующие поля ...
	RateLimit RateLimitConfig
}
```

**Файл:** `.env` (добавить)

```bash
# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_STRATEGY=redis
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200
RATE_LIMIT_WINDOW_SEC=60

# Auth endpoints
RATE_LIMIT_AUTH_RPS=5
RATE_LIMIT_AUTH_BURST=10
```

---

#### Шаг 1.1.6: Integration в Router

**Файл:** `internal/interfaces/http/router_configurator.go` (обновить)

```go
func (rc *RouteConfigurator) Configure(router *gin.Engine) error {
	// ... существующий код ...

	if !rc.cfg.RateLimit.Enabled {
		rc.logger.Info("Rate limiting disabled")
		return rc.configureRoutesWithoutRateLimit(router)
	}

	// Создаем лимитер на основе стратегии
	var apiLimiter ports.RateLimiter
	var authLimiter ports.RateLimiter

	if rc.cfg.RateLimit.Strategy == "redis" && rc.container.RedisClient() != nil {
		rc.logger.Info("Using Redis rate limiter")

		apiLimiter = ratelimit.NewRedisRateLimiter(
			rc.container.RedisClient(),
			rc.cfg.RateLimit.RequestsPerSecond,
			time.Duration(rc.cfg.RateLimit.WindowSeconds)*time.Second,
		)

		authLimiter = ratelimit.NewRedisRateLimiter(
			rc.container.RedisClient(),
			rc.cfg.RateLimit.AuthRPS,
			time.Duration(rc.cfg.RateLimit.WindowSeconds)*time.Second,
		)
	} else {
		rc.logger.Info("Using in-memory rate limiter")

		apiLimiter = ratelimit.NewTokenBucket(
			rc.cfg.RateLimit.RequestsPerSecond,
			rc.cfg.RateLimit.BurstSize,
		)

		authLimiter = ratelimit.NewTokenBucket(
			rc.cfg.RateLimit.AuthRPS,
			rc.cfg.RateLimit.AuthBurst,
		)
	}

	api := router.Group("/api/v1")

	// Auth endpoints с строгим rate limiting (по IP)
	authGroup := api.Group("/users")
	authGroup.Use(middleware.RateLimitMiddleware(authLimiter, rc.logger, middleware.IPBasedKey))
	{
		userHandler := NewUserHandler(rc.container.UserService(), rc.logger)
		authGroup.POST("/register", userHandler.Register)
		authGroup.POST("/login", userHandler.Login)
		authGroup.POST("/refresh", userHandler.RefreshToken)
	}

	// Protected routes с rate limiting по user ID
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(rc.container.AuthService(), rc.logger))
	protected.Use(middleware.RateLimitMiddleware(apiLimiter, rc.logger, middleware.UserBasedKey))
	{
		rc.configureProtectedRoutes(protected)
	}

	return nil
}
```

---

#### Шаг 1.1.7: Metrics

**Файл:** `internal/infrastructure/observability/metrics.go` (добавить)

```go
func (m *Metrics) InitRateLimitMetrics() {
	m.rateLimitExceeded = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Name:      "rate_limit_exceeded_total",
			Help:      "Total number of rate limit exceeded events",
		},
		[]string{"limiter", "key_type"},
	)

	m.rateLimitAllowed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Name:      "rate_limit_allowed_total",
			Help:      "Total number of allowed requests",
		},
		[]string{"limiter", "key_type"},
	)

	prometheus.MustRegister(m.rateLimitExceeded, m.rateLimitAllowed)
}

func (m *Metrics) RecordRateLimitExceeded(limiter, keyType string) {
	m.rateLimitExceeded.WithLabelValues(limiter, keyType).Inc()
}

func (m *Metrics) RecordRateLimitAllowed(limiter, keyType string) {
	m.rateLimitAllowed.WithLabelValues(limiter, keyType).Inc()
}
```

---

### 📍 1.2. API Contract (Swagger)

**Файл:** `docs/swagger_extensions.yaml`

```yaml
paths:
  /api/v1/users/login:
    post:
      summary: User login
      tags:
        - Authentication
      parameters:
        - name: X-RateLimit-Limit
          in: header
          schema:
            type: integer
          description: Maximum number of requests allowed
        - name: X-RateLimit-Remaining
          in: header
          schema:
            type: integer
          description: Number of requests remaining
        - name: X-RateLimit-Reset
          in: header
          schema:
            type: integer
          description: Unix timestamp when the rate limit resets
      responses:
        "200":
          description: Successful login
        "429":
          description: Rate limit exceeded
          headers:
            Retry-After:
              schema:
                type: integer
              description: Seconds to wait before retrying
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/RateLimitError"

components:
  schemas:
    RateLimitError:
      type: object
      properties:
        error:
          type: string
          example: rate_limit_exceeded
        message:
          type: string
          example: Too many requests, please try again later
        retry_after:
          type: integer
          example: 60
```

---

### 📍 1.3. Testing Plan

#### Unit Tests

- [x] `ports/rate_limiter_test.go` - Mock интерфейса
- [x] `ratelimit/token_bucket_test.go` - Token Bucket алгоритм
- [x] `middleware/rate_limit_test.go` - Middleware логика

#### Integration Tests

- [x] `ratelimit/redis_limiter_integration_test.go` - Redis limiter с реальным Redis
- [ ] `router_integration_test.go` - Полная интеграция в API Gateway

#### Manual Testing

```bash
# Тест Token Bucket (memory)
for i in {1..15}; do
  curl -i http://localhost:8080/api/v1/users/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"password"}'
  echo "\nRequest $i completed"
  sleep 0.1
done

# Ожидаемый результат:
# - Первые 10 запросов: 200 OK (burst)
# - Запросы 11-15: 429 Too Many Requests
# - Headers: X-RateLimit-Limit, X-RateLimit-Remaining, Retry-After

# Тест Redis (distributed)
# Запустить из разных терминалов одновременно
curl -i http://localhost:8080/api/v1/users/login ...
```

---

### 📍 1.4. Frontend Integration (TypeScript/React)

#### Шаг 1.4.1: API Client - Rate Limit Handler

**Файл:** `kkm-platform/lib/api/rateLimitHandler.ts`

```typescript
import { AxiosError, AxiosInstance } from "axios";

export interface RateLimitInfo {
  limit: number;
  remaining: number;
  reset: Date;
  retryAfter?: number;
}

export class RateLimitHandler {
  private rateLimitInfo: Map<string, RateLimitInfo> = new Map();

  /**
   * Извлекает информацию о rate limit из заголовков ответа
   */
  extractRateLimitInfo(headers: Record<string, string>): RateLimitInfo | null {
    const limit = headers["x-ratelimit-limit"];
    const remaining = headers["x-ratelimit-remaining"];
    const reset = headers["x-ratelimit-reset"];

    if (!limit || !remaining || !reset) {
      return null;
    }

    return {
      limit: parseInt(limit, 10),
      remaining: parseInt(remaining, 10),
      reset: new Date(parseInt(reset, 10) * 1000),
    };
  }

  /**
   * Проверяет 429 ошибку и извлекает retry-after
   */
  handleRateLimitError(error: AxiosError): RateLimitInfo | null {
    if (error.response?.status !== 429) {
      return null;
    }

    const retryAfter = error.response.headers["retry-after"];
    const info = this.extractRateLimitInfo(error.response.headers);

    if (info && retryAfter) {
      info.retryAfter = parseInt(retryAfter, 10);
    }

    return info;
  }

  /**
   * Сохраняет rate limit info для endpoint
   */
  storeRateLimitInfo(endpoint: string, info: RateLimitInfo) {
    this.rateLimitInfo.set(endpoint, info);
  }

  /**
   * Получает сохраненную rate limit info
   */
  getRateLimitInfo(endpoint: string): RateLimitInfo | undefined {
    return this.rateLimitInfo.get(endpoint);
  }

  /**
   * Проверяет, нужно ли ждать перед следующим запросом
   */
  shouldWait(endpoint: string): number {
    const info = this.rateLimitInfo.get(endpoint);
    if (!info || info.remaining > 0) {
      return 0;
    }

    const now = new Date();
    const waitMs = info.reset.getTime() - now.getTime();

    return Math.max(0, waitMs);
  }
}

/**
 * Axios interceptor для автоматической обработки rate limits
 */
export function setupRateLimitInterceptor(axios: AxiosInstance) {
  const handler = new RateLimitHandler();

  // Response interceptor
  axios.interceptors.response.use(
    (response) => {
      // Сохраняем rate limit info из заголовков
      const info = handler.extractRateLimitInfo(response.headers);
      if (info) {
        const endpoint = response.config.url || "";
        handler.storeRateLimitInfo(endpoint, info);
      }
      return response;
    },
    async (error: AxiosError) => {
      // Обрабатываем 429 ошибку
      const info = handler.handleRateLimitError(error);

      if (info && info.retryAfter) {
        console.warn(
          `Rate limit exceeded. Retry after ${info.retryAfter} seconds`,
        );

        // Можно автоматически повторить запрос после задержки
        if (error.config && info.retryAfter < 60) {
          // Только если < 1 минуты
          await new Promise((resolve) =>
            setTimeout(resolve, info.retryAfter * 1000),
          );
          return axios.request(error.config);
        }
      }

      return Promise.reject(error);
    },
  );

  return handler;
}
```

#### Шаг 1.4.2: React Hook для мониторинга Rate Limits

**Файл:** `kkm-platform/hooks/useRateLimit.ts`

```typescript
import { useState, useEffect } from "react";
import { rateLimitHandler } from "@/lib/api/client";

export interface UseRateLimitResult {
  limit: number | null;
  remaining: number | null;
  reset: Date | null;
  isLimited: boolean;
  secondsUntilReset: number;
}

export function useRateLimit(endpoint?: string): UseRateLimitResult {
  const [state, setState] = useState<UseRateLimitResult>({
    limit: null,
    remaining: null,
    reset: null,
    isLimited: false,
    secondsUntilReset: 0,
  });

  useEffect(() => {
    if (!endpoint) return;

    const updateState = () => {
      const info = rateLimitHandler.getRateLimitInfo(endpoint);

      if (!info) {
        setState({
          limit: null,
          remaining: null,
          reset: null,
          isLimited: false,
          secondsUntilReset: 0,
        });
        return;
      }

      const now = new Date();
      const secondsUntilReset = Math.max(
        0,
        Math.floor((info.reset.getTime() - now.getTime()) / 1000),
      );

      setState({
        limit: info.limit,
        remaining: info.remaining,
        reset: info.reset,
        isLimited: info.remaining === 0,
        secondsUntilReset,
      });
    };

    // Обновляем сразу
    updateState();

    // И каждую секунду
    const interval = setInterval(updateState, 1000);

    return () => clearInterval(interval);
  }, [endpoint]);

  return state;
}
```

#### Шаг 1.4.3: UI Component - Rate Limit Warning

**Файл:** `kkm-platform/components/ui/RateLimitWarning.tsx`

```typescript
'use client';

import { useRateLimit } from '@/hooks/useRateLimit';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { AlertTriangle, Clock } from 'lucide-react';

interface RateLimitWarningProps {
  endpoint: string;
  threshold?: number; // Показывать предупреждение когда remaining < threshold
}

export function RateLimitWarning({
  endpoint,
  threshold = 10
}: RateLimitWarningProps) {
  const { remaining, isLimited, secondsUntilReset, limit } = useRateLimit(endpoint);

  if (remaining === null || remaining > threshold) {
    return null;
  }

  if (isLimited) {
    return (
      <Alert variant="destructive" className="mb-4">
        <AlertTriangle className="h-4 w-4" />
        <AlertTitle>Rate Limit Exceeded</AlertTitle>
        <AlertDescription>
          You've reached the maximum number of requests ({limit}).
          Please wait {secondsUntilReset} seconds before trying again.
        </AlertDescription>
      </Alert>
    );
  }

  return (
    <Alert variant="warning" className="mb-4">
      <Clock className="h-4 w-4" />
      <AlertTitle>Approaching Rate Limit</AlertTitle>
      <AlertDescription>
        You have {remaining} out of {limit} requests remaining.
        Limit resets in {secondsUntilReset} seconds.
      </AlertDescription>
    </Alert>
  );
}
```

#### Шаг 1.4.4: Integration в API Client

**Файл:** `kkm-platform/lib/api/client.ts` (обновить)

```typescript
import axios from "axios";
import { setupRateLimitInterceptor } from "./rateLimitHandler";

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1",
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Setup rate limit handling
export const rateLimitHandler = setupRateLimitInterceptor(apiClient);

export default apiClient;
```

---

### 📍 1.5. E2E Testing

**Файл:** `kkm-platform/e2e/rateLimit.spec.ts`

```typescript
import { test, expect } from "@playwright/test";

test.describe("Rate Limiting", () => {
  test("should show rate limit warning when approaching limit", async ({
    page,
  }) => {
    await page.goto("/login");

    // Делаем несколько быстрых запросов
    for (let i = 0; i < 8; i++) {
      await page.fill('[name="email"]', `test${i}@test.com`);
      await page.fill('[name="password"]', "wrongpassword");
      await page.click('button[type="submit"]');
      await page.waitForTimeout(100);
    }

    // Должно появиться предупреждение
    await expect(
      page.locator('[data-testid="rate-limit-warning"]'),
    ).toBeVisible();
    await expect(
      page.locator('[data-testid="rate-limit-warning"]'),
    ).toContainText("Approaching Rate Limit");
  });

  test("should block requests when rate limit exceeded", async ({ page }) => {
    await page.goto("/login");

    // Исчерпываем лимит
    for (let i = 0; i < 12; i++) {
      await page.fill('[name="email"]', `test${i}@test.com`);
      await page.fill('[name="password"]', "wrongpassword");
      await page.click('button[type="submit"]');
      await page.waitForTimeout(100);
    }

    // Должна быть блокировка
    await expect(
      page.locator('[data-testid="rate-limit-error"]'),
    ).toBeVisible();
    await expect(
      page.locator('[data-testid="rate-limit-error"]'),
    ).toContainText("Rate Limit Exceeded");
  });
});
```

---

## ✅ Checklist для Фазы 1

### Backend

- [ ] Реализован интерфейс RateLimiter
- [ ] Реализован TokenBucket алгоритм
- [ ] Реализован RedisRateLimiter
- [ ] Создан Gin middleware
- [ ] Добавлена конфигурация
- [ ] Интегрировано в router
- [ ] Добавлены метрики Prometheus
- [ ] Написаны unit тесты (coverage > 80%)
- [ ] Написаны integration тесты
- [ ] Manual testing пройден

### API Contract

- [ ] Обновлена Swagger документация
- [ ] Добавлены примеры 429 ответов
- [ ] Документированы заголовки X-RateLimit-\*
- [ ] Проверено через Swagger UI

### Frontend

- [ ] Реализован RateLimitHandler
- [ ] Создан useRateLimit hook
- [ ] Создан RateLimitWarning компонент
- [ ] Интегрировано в API client
- [ ] Добавлены E2E тесты
- [ ] UI/UX проверен

### Documentation

- [ ] Обновлен README.md api-gateway
- [ ] Обновлен README.md kkm-platform
- [ ] Создан RATE_LIMITING_GUIDE.md
- [ ] Примеры использования в документации

---

## 🚀 Следующие фазы

### ФАЗА 2: Request/Response Logging & Audit (14-18 часов)

### ФАЗА 3: Health Checks & Monitoring Dashboard (12-16 часов)

### ФАЗА 4: Distributed Tracing (OpenTelemetry) (16-20 часов)

### ФАЗА 5: Circuit Breaker Pattern (10-14 часов)

### ФАЗА 6: Redis Cache Strategy (12-16 часов)

---

**Итого Фаза 1:** Backend (8-10ч) + Frontend (4-6ч) + Testing (2ч) = **14-18 часов**

При необходимости можем детализировать следующие фазы.
