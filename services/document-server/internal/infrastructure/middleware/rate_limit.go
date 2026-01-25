package middleware

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimiter ограничивает количество запросов
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      float64 // requests per second
	logger   *zap.Logger
}

// NewRateLimiter создает новый rate limiter
func NewRateLimiter(requestsPerSecond float64, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      requestsPerSecond,
		logger:   logger,
	}
}

// UnaryServerInterceptor возвращает unary interceptor для rate limiting
func (rl *RateLimiter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Используем RPC метод как ключ для limiting
		limiter := rl.getOrCreateLimiter(info.FullMethod)

		if !limiter.Allow() {
			rl.logger.Warn("Rate limit exceeded",
				zap.String("method", info.FullMethod))
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}

func (rl *RateLimiter) getOrCreateLimiter(method string) *rate.Limiter {
	rl.mu.RLock()
	if limiter, exists := rl.limiters[method]; exists {
		rl.mu.RUnlock()
		return limiter
	}
	rl.mu.RUnlock()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check pattern
	if limiter, exists := rl.limiters[method]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(rl.rps), 1)
	rl.limiters[method] = limiter
	return limiter
}
