// Файл api-gateway/internal/infrastructure/client/circuit_breaker.go содержит реализацию пакета client.
package client

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CircuitBreakerConfig конфигурация circuit breaker
type CircuitBreakerConfig struct {
	MaxFailures   int // максимум ошибок перед открытием (default 5)
	ResetTimeout  int // таймаут в секундах перед попыткой восстановления (default 30)
	ConsecutiveOK int // кол-во успешных запросов для закрытия (default 2)
	IgnoredCodes  map[codes.Code]bool
}

// DefaultCircuitBreakerConfig дефолтная конфигурация
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeout:  30,
		ConsecutiveOK: 2,
		IgnoredCodes: map[codes.Code]bool{
			codes.InvalidArgument:  true, // не считаем client ошибки
			codes.NotFound:         true,
			codes.PermissionDenied: true,
		},
	}
}

// CircuitBreaker состояние circuit breaker
type CircuitBreaker struct {
	state           string // "closed", "open", "half-open"
	failureCount    int
	successCount    int
	lastFailureTime int64
	config          *CircuitBreakerConfig
	logger          *zap.Logger
}

// NewCircuitBreaker создает новый circuit breaker
func NewCircuitBreaker(cfg *CircuitBreakerConfig, logger *zap.Logger) *CircuitBreaker {
	if cfg == nil {
		cfg = DefaultCircuitBreakerConfig()
	}
	return &CircuitBreaker{
		state:  "closed",
		config: cfg,
		logger: logger,
	}
}

// IsOpen проверяет открыт ли circuit
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.state == "open"
}

// RecordSuccess записывает успешный вызов
func (cb *CircuitBreaker) RecordSuccess() {
	if cb.state == "half-open" {
		cb.successCount++
		if cb.successCount >= cb.config.ConsecutiveOK {
			cb.state = "closed"
			cb.failureCount = 0
			cb.successCount = 0
			cb.logger.Info("Circuit breaker closed")
		}
	}
}

// RecordFailure записывает ошибку
func (cb *CircuitBreaker) RecordFailure(err error) {
	// Не считаем client ошибки
	if st, ok := status.FromError(err); ok {
		if cb.config.IgnoredCodes[st.Code()] {
			return
		}
	}

	cb.failureCount++
	cb.successCount = 0

	if cb.state == "closed" && cb.failureCount >= cb.config.MaxFailures {
		cb.state = "open"
		cb.logger.Warn("Circuit breaker opened", zap.Int("failures", cb.failureCount))
	}
}

// CircuitBreakerUnaryInterceptor unary interceptor для circuit breaker
func CircuitBreakerUnaryInterceptor(cb *CircuitBreaker) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if cb.IsOpen() {
			return status.Error(codes.Unavailable, "circuit breaker is open")
		}

		err := invoker(ctx, method, req, reply, cc, opts...)

		if err == nil {
			cb.RecordSuccess()
		} else {
			cb.RecordFailure(err)
		}

		return err
	}
}

// CircuitBreakerStreamInterceptor stream interceptor для circuit breaker
func CircuitBreakerStreamInterceptor(cb *CircuitBreaker) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		if cb.IsOpen() {
			return nil, status.Error(codes.Unavailable, "circuit breaker is open")
		}

		stream, err := streamer(ctx, desc, cc, method, opts...)

		if err == nil {
			cb.RecordSuccess()
		} else {
			cb.RecordFailure(err)
		}

		return stream, err
	}
}

// GetDialOptions возвращает options для grpc.Dial с circuit breaker и retry
func GetDialOptions(logger *zap.Logger) []grpc.DialOption {
	cb := NewCircuitBreaker(nil, logger)

	return []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10 * 1024 * 1024), // 10MB
		),
		grpc.WithUnaryInterceptor(CircuitBreakerUnaryInterceptor(cb)),
		grpc.WithStreamInterceptor(CircuitBreakerStreamInterceptor(cb)),
	}
}
