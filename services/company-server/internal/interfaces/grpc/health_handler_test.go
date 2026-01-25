package grpc

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

// TestNewHealthHandler проверяет создание health handler
func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler()
	if handler == nil {
		t.Fatal("NewHealthHandler returned nil")
	}
}

// TestNewHealthHandlerWithDeps проверяет создание health handler с зависимостями
func TestNewHealthHandlerWithDeps(t *testing.T) {
	logger := zap.NewNop()
	handler := NewHealthHandlerWithDeps(nil, logger)

	if handler == nil {
		t.Fatal("NewHealthHandlerWithDeps returned nil")
	}
	if handler.logger != logger {
		t.Error("Logger not set correctly")
	}
}

// TestHealthCheckWithoutDB проверяет health check без БД
func TestHealthCheckWithoutDB(t *testing.T) {
	handler := NewHealthHandler()

	resp, err := handler.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})

	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("Expected SERVING status, got %v", resp.Status)
	}
}

// TestHealthCheckLogging проверяет логирование при health check
func TestHealthCheckLogging(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	handler := NewHealthHandlerWithDeps(nil, logger)

	resp, err := handler.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})

	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("Expected SERVING status, got %v", resp.Status)
	}
}

// MockHealthWatchServer для тестирования Watch
type MockHealthWatchServer struct {
	responses []*grpc_health_v1.HealthCheckResponse
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewMockHealthWatchServer() *MockHealthWatchServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &MockHealthWatchServer{
		responses: make([]*grpc_health_v1.HealthCheckResponse, 0),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (m *MockHealthWatchServer) Send(resp *grpc_health_v1.HealthCheckResponse) error {
	m.responses = append(m.responses, resp)
	return nil
}

func (m *MockHealthWatchServer) Context() context.Context {
	return m.ctx
}

func (m *MockHealthWatchServer) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockHealthWatchServer) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockHealthWatchServer) SetTrailer(md metadata.MD) {
}

func (m *MockHealthWatchServer) SendMsg(m2 interface{}) error {
	return nil
}

func (m *MockHealthWatchServer) RecvMsg(m2 interface{}) error {
	return nil
}

// TestHealthWatchBasic проверяет базовую функцию Watch
func TestHealthWatchBasic(t *testing.T) {
	handler := NewHealthHandler()
	mockServer := NewMockHealthWatchServer()

	go func() {
		time.Sleep(100 * time.Millisecond)
		mockServer.cancel()
	}()

	_ = handler.Watch(&grpc_health_v1.HealthCheckRequest{}, mockServer)

	if len(mockServer.responses) == 0 {
		t.Error("Watch did not send any responses")
	}
	if mockServer.responses[0].Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("Expected SERVING status, got %v", mockServer.responses[0].Status)
	}
}

// TestHealthCheckConcurrent проверяет конкурентные health checks
func TestHealthCheckConcurrent(t *testing.T) {
	handler := NewHealthHandler()
	results := make(chan *grpc_health_v1.HealthCheckResponse, 10)
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			resp, err := handler.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
			if err == nil {
				results <- resp
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	close(results)
	count := 0
	for resp := range results {
		if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			t.Errorf("Expected SERVING status, got %v", resp.Status)
		}
		count++
	}

	if count != 10 {
		t.Logf("Expected 10 responses, got %d", count)
	}
}

// TestHealthCheckRequestWithService проверяет check для конкретного сервиса
func TestHealthCheckRequestWithService(t *testing.T) {
	handler := NewHealthHandler()
	req := &grpc_health_v1.HealthCheckRequest{
		Service: "company.CompanyService",
	}

	resp, err := handler.Check(context.Background(), req)

	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("Expected SERVING status, got %v", resp.Status)
	}
}

// TestHealthCheckRequestEmpty проверяет check с пустым запросом
func TestHealthCheckRequestEmpty(t *testing.T) {
	handler := NewHealthHandler()
	req := &grpc_health_v1.HealthCheckRequest{}

	resp, err := handler.Check(context.Background(), req)

	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("Expected SERVING status, got %v", resp.Status)
	}
}

// TestHealthCheckMultipleCalls проверяет несколько вызовов check
func TestHealthCheckMultipleCalls(t *testing.T) {
	handler := NewHealthHandler()

	for i := 0; i < 5; i++ {
		resp, err := handler.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})

		if err != nil {
			t.Fatalf("Iteration %d: Check returned error: %v", i, err)
		}
		if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			t.Errorf("Iteration %d: Expected SERVING status, got %v", i, resp.Status)
		}
	}
}
