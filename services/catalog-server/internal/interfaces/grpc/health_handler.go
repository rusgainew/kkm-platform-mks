// Файл catalog-server/internal/interfaces/grpc/health_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// HealthHandler реализует gRPC Health Check
type HealthHandler struct {
	grpc_health_v1.UnimplementedHealthServer
	logger *zap.Logger
}

// NewHealthHandler создает новый обработчик health check
func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// Check проверяет здоровье сервиса
func (h *HealthHandler) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	h.logger.Debug("health check called", zap.String("service", req.GetService()))

	// Здесь можно добавить проверки баз данных, очередей и т.д.
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch поддерживает streaming health check
func (h *HealthHandler) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	// Отправка начального статуса
	if err := stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// В реальном сценарии здесь можно мониторить состояние и отправлять обновления
	<-stream.Context().Done()
	return nil
}
