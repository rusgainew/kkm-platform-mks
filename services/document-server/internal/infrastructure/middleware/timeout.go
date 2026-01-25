package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TimeoutInterceptor определяет timeout для разных методов
type TimeoutInterceptor struct {
	logger *zap.Logger
	// timeoutMap хранит timeout для каждого метода
	timeoutMap map[string]time.Duration
}

// NewTimeoutInterceptor создает новый timeout interceptor
func NewTimeoutInterceptor(logger *zap.Logger) *TimeoutInterceptor {
	return &TimeoutInterceptor{
		logger: logger,
		timeoutMap: map[string]time.Duration{
			// Read operations - 5 seconds
			"/api.DocumentService/GetDocument":   5 * time.Second,
			"/api.DocumentService/ListDocuments": 5 * time.Second,

			// Write operations - 10 seconds
			"/api.DocumentService/CreateDocument":  10 * time.Second,
			"/api.DocumentService/UpdateDocument":  10 * time.Second,
			"/api.DocumentService/SendDocument":    10 * time.Second,
			"/api.DocumentService/ApproveDocument": 10 * time.Second,
			"/api.DocumentService/RejectDocument":  10 * time.Second,
			"/api.DocumentService/ArchiveDocument": 10 * time.Second,
		},
	}
}

// UnaryServerInterceptor возвращает unary interceptor для timeout
func (ti *TimeoutInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Получаем timeout для метода
		timeout, ok := ti.timeoutMap[info.FullMethod]
		if !ok {
			// Default timeout 30 seconds если метод не найден
			timeout = 30 * time.Second
		}

		// Проверяем есть ли уже deadline в контексте
		if _, hasDeadline := ctx.Deadline(); hasDeadline {
			// Если уже есть deadline, используем его
			return handler(ctx, req)
		}

		// Создаем новый контекст с timeout
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Логируем timeout для debugging
		ti.logger.Debug("Setting timeout for method",
			zap.String("method", info.FullMethod),
			zap.Duration("timeout", timeout),
		)

		// Обрабатываем запрос с timeout
		resultChan := make(chan interface{}, 1)
		errChan := make(chan error, 1)

		go func() {
			result, err := handler(ctx, req)
			if err != nil {
				errChan <- err
			} else {
				resultChan <- result
			}
		}()

		// Ждем результат или timeout
		select {
		case result := <-resultChan:
			return result, nil
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			ti.logger.Warn("Request timeout exceeded",
				zap.String("method", info.FullMethod),
				zap.Duration("timeout", timeout),
				zap.Error(ctx.Err()),
			)
			return nil, status.Error(codes.DeadlineExceeded, "request timeout exceeded")
		}
	}
}
