// Файл company-server/internal/interfaces/grpc/health_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthHandler реализует health check
type HealthHandler struct {
	grpc_health_v1.UnimplementedHealthServer
	db     *sqlx.DB
	logger *zap.Logger
}

// NewHealthHandler создает новый health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		db:     nil,
		logger: nil,
	}
}

// NewHealthHandlerWithDeps создает health handler с зависимостями
func NewHealthHandlerWithDeps(db *sqlx.DB, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// Check выполняет health check
func (h *HealthHandler) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	// Проверка БД
	if h.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := h.db.PingContext(ctx); err != nil {
			if h.logger != nil {
				h.logger.Error("Health check failed: database not accessible", zap.Error(err))
			}
			return &grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
			}, nil
		}
	}

	if h.logger != nil {
		h.logger.Debug("Health check passed")
	}
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch наблюдает за изменениями статуса здоровья
func (h *HealthHandler) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	// Отправить начальный статус
	if err := server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}); err != nil {
		if h.logger != nil {
			h.logger.Error("Failed to send health status", zap.Error(err))
		}
		return err
	}

	// Периодическая проверка здоровья каждые 10 секунд
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Проверка здоровья
			healthStatus := grpc_health_v1.HealthCheckResponse_SERVING

			if h.db != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				if err := h.db.PingContext(ctx); err != nil {
					healthStatus = grpc_health_v1.HealthCheckResponse_NOT_SERVING
					if h.logger != nil {
						h.logger.Warn("Health check watch: database not accessible", zap.Error(err))
					}
				}
				cancel()
			}

			// Отправить обновленный статус
			if err := server.Send(&grpc_health_v1.HealthCheckResponse{
				Status: healthStatus,
			}); err != nil {
				if h.logger != nil {
					h.logger.Error("Failed to send health status in watch", zap.Error(err))
				}
				return err
			}

		case <-server.Context().Done():
			return server.Context().Err()
		}
	}
}
