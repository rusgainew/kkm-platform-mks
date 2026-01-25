package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// EndpointRateLimiter управляет rate limiting для конкретных endpoint'ов
type EndpointRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	logger   *zap.Logger
}

// NewEndpointRateLimiter создает новый endpoint rate limiter
func NewEndpointRateLimiter(logger *zap.Logger) *EndpointRateLimiter {
	return &EndpointRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		logger:   logger,
	}
}

// AddLimit добавляет лимит для конкретного endpoint
func (e *EndpointRateLimiter) AddLimit(endpoint string, requestsPerSecond int, burst int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
	e.limiters[endpoint] = limiter

	e.logger.Info("Added rate limit for endpoint",
		zap.String("endpoint", endpoint),
		zap.Int("requests_per_second", requestsPerSecond),
		zap.Int("burst", burst),
	)
}

// Middleware возвращает middleware для rate limiting
func (e *EndpointRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()

		e.mu.RLock()
		limiter, exists := e.limiters[path]
		e.mu.RUnlock()

		// Если лимит не установлен для этого endpoint, пропускаем
		if !exists {
			c.Next()
			return
		}

		// Проверяем лимит
		if !limiter.Allow() {
			e.logger.Warn("Rate limit exceeded",
				zap.String("path", path),
				zap.String("client_ip", c.ClientIP()),
			)

			c.Header("X-RateLimit-Limit", "")
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", "1")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPBasedRateLimiter управляет rate limiting на основе IP
type IPBasedRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	limit    rate.Limit
	burst    int
	logger   *zap.Logger
	stop     chan struct{}
}

// NewIPBasedRateLimiter создает новый IP-based rate limiter
func NewIPBasedRateLimiter(requestsPerSecond int, burst int, logger *zap.Logger) *IPBasedRateLimiter {
	return &IPBasedRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		limit:    rate.Limit(requestsPerSecond),
		burst:    burst,
		logger:   logger,
		stop:     make(chan struct{}),
	}
}

// getLimiter получает или создает limiter для IP
func (i *IPBasedRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(i.limit, i.burst)
		i.limiters[ip] = limiter
	}

	return limiter
}

// Middleware возвращает middleware для IP-based rate limiting
func (i *IPBasedRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := i.getLimiter(ip)

		if !limiter.Allow() {
			i.logger.Warn("IP rate limit exceeded",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
			)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "too many requests from your IP, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Cleanup периодически очищает неактивные limiters
func (i *IPBasedRateLimiter) Cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				i.mu.Lock()
				// Простая очистка - удаляем все limiters
				// В production лучше отслеживать last access time
				i.limiters = make(map[string]*rate.Limiter)
				i.mu.Unlock()
			case <-i.stop:
				return
			}
		}
	}()
}

// Close останавливает фоновую goroutine cleanup
func (i *IPBasedRateLimiter) Close() {
	close(i.stop)
}
