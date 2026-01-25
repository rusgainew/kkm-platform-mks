package services

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

func TestForeignCompanyQueryService_ListForeignCompanies(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewForeignCompanyQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	_, err := service.ListForeignCompanies(ctx, 0, 20)

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestForeignCompanyQueryService_SearchForeignCompanies(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewForeignCompanyQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	searchReq := &pb.SearchRequest{
		SearchText: "test",
		Page:       0,
		Size:       20,
	}

	_, err := service.SearchForeignCompanies(ctx, searchReq)
	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestForeignCompanyQueryService_FilterForeignCompanies(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)

	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewForeignCompanyQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	filterReq := &pb.CatalogFilterRequest{
		Name: "test",
		Page: 0,
		Size: 20,
	}

	_, err := service.FilterForeignCompanies(ctx, filterReq)
	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}
