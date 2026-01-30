// Файл document-server/internal/interfaces/grpc/health_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthHandler struct {
	grpc_health_v1.UnimplementedHealthServer
}

// NewHealthHandler создает новый health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check возвращает статус сервиса
func (h *HealthHandler) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch отправляет статусы сервиса
func (h *HealthHandler) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
