// Файл api-gateway/internal/application/services/document_query_service_test.go содержит реализацию пакета services.
package services

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

func TestDocumentQueryService_GetDocument(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)
	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewDocumentQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	_, err := service.GetDocument(ctx, "test-doc-id")

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestDocumentQueryService_ListDocuments(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)
	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewDocumentQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	_, err := service.ListDocuments(ctx, 0, 20, "", "", "", "")

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}

func TestDocumentQueryService_SearchDocuments(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	tracer := observability.NewTracer("test", logger)
	connMgr := client.NewConnectionManager(5, logger)
	defer connMgr.CloseAll()

	service := NewDocumentQueryService(connMgr, "invalid:99999", metrics, tracer, logger)

	ctx := context.Background()
	_, err := service.SearchDocuments(ctx, "test query", 0, 20, "", "", "")

	if err == nil {
		t.Error("Expected connection error, got nil")
	}
}
