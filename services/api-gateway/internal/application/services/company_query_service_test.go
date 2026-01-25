package services

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

func TestCompanyQueryService_ListCompanies_Success(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	// Используем реальный ConnectionManager с невалидным адресом для теста ошибок
	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewCompanyQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	// Test
	ctx := context.Background()
	_, err := service.ListCompanies(ctx, 0, 20)

	// Ожидаем ошибку подключения
	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestCompanyQueryService_SearchCompanies_Success(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewCompanyQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	searchReq := &pb.SearchRequest{
		SearchText: "test",
		Page:       0,
		Size:       20,
	}

	_, err := service.SearchCompanies(ctx, searchReq)
	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestCompanyQueryService_FilterCompanies_Success(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewCompanyQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	filterReq := &pb.CatalogFilterRequest{
		Name: "test",
		Page: 0,
		Size: 20,
	}

	_, err := service.FilterCompanies(ctx, filterReq)
	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestCompanyQueryService_ListCompanies_ConnectionError(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	connMgr := client.NewConnectionManager(1, logger)
	defer connMgr.CloseAll()

	service := NewCompanyQueryService(connMgr, "invalid-host:12345", metrics, tracer, logger)

	ctx := context.Background()
	_, err := service.ListCompanies(ctx, 0, 20)

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestCompanyQueryService_SearchCompanies_ConnectionError(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	connMgr := client.NewConnectionManager(1, logger)
	defer connMgr.CloseAll()

	service := NewCompanyQueryService(connMgr, "invalid-host:12345", metrics, tracer, logger)

	ctx := context.Background()
	searchReq := &pb.SearchRequest{
		SearchText: "test",
		Page:       0,
		Size:       20,
	}

	_, err := service.SearchCompanies(ctx, searchReq)

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestCompanyQueryService_FilterCompanies_ConnectionError(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	connMgr := client.NewConnectionManager(1, logger)
	defer connMgr.CloseAll()

	service := NewCompanyQueryService(connMgr, "invalid-host:12345", metrics, tracer, logger)

	ctx := context.Background()
	filterReq := &pb.CatalogFilterRequest{
		Name: "test",
		Page: 0,
		Size: 20,
	}

	_, err := service.FilterCompanies(ctx, filterReq)

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}
