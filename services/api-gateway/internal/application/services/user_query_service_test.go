// Файл api-gateway/internal/application/services/user_query_service_test.go содержит реализацию пакета services.
package services

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

func TestUserQueryService_GetUser(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)
	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewUserQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	_, err := service.GetUser(ctx, "test-user-id")

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestUserQueryService_ListUsers(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)
	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewUserQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	_, err := service.ListUsers(ctx, 0, 20, "", "")

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestUserQueryService_SearchUsers(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)
	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewUserQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	_, err := service.SearchUsers(ctx, "test query", 0, 20, "")

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}
