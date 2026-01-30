// Файл invoice-server/internal/interfaces/grpc/health_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

// HealthHandler обработчик health проверок
type HealthHandler struct {
	pb.UnimplementedHealthServer
}

// NewHealthHandler создает новый health обработчик
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check реализует health проверку
func (h *HealthHandler) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

// Watch реализует потоковую health проверку
func (h *HealthHandler) Watch(req *pb.HealthCheckRequest, stream pb.Health_WatchServer) error {
	return stream.Send(&pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING})
}
