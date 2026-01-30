// Файл document-server/internal/application/document/tracing_service_test.go содержит реализацию пакета document.
package document

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestServiceWithTracing тестирует сервис с трассированием
func TestServiceWithTracing(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockRepository()
	mockPublisher := &MockPublisher{}
	tracing := observability.InitTracingProvider("document-server-test")

	svc := NewServiceWithTracing(mockRepo, mockPublisher, logger, tracing)
	ctx := context.Background()

	t.Run("CreateDocumentWithTracing succeeds", func(t *testing.T) {
		documentID, err := svc.CreateDocumentWithTracing(ctx, "org-123", "Test Doc", "Content", "user-1")
		require.NoError(t, err)
		assert.NotEmpty(t, documentID)
		assert.Equal(t, 1, mockRepo.CreateCalls)
		assert.Greater(t, len(mockPublisher.Events), 0)
	})

	t.Run("CreateDocumentWithTracing fails with invalid input", func(t *testing.T) {
		documentID, err := svc.CreateDocumentWithTracing(ctx, "", "Test", "Content", "user-1")
		assert.Error(t, err)
		assert.Empty(t, documentID)
	})

	t.Run("GetDocumentWithTracing succeeds", func(t *testing.T) {
		mockRepo.GetResult = map[string]interface{}{
			"id":    "doc-1",
			"title": "Test Document",
		}
		mockRepo.GetCalls = 0

		doc, err := svc.GetDocumentWithTracing(ctx, "doc-1")
		require.NoError(t, err)
		assert.Equal(t, "doc-1", doc["id"])
		assert.Equal(t, 1, mockRepo.GetCalls)
	})

	t.Run("UpdateDocumentWithTracing succeeds", func(t *testing.T) {
		mockRepo.UpdateCalls = 0

		err := svc.UpdateDocumentWithTracing(ctx, "doc-1", "Updated Title", "Updated Content")
		require.NoError(t, err)
		assert.Equal(t, 1, mockRepo.UpdateCalls)
		assert.Greater(t, len(mockPublisher.Events), 0)
	})

	t.Run("ListDocumentsWithTracing succeeds", func(t *testing.T) {
		mockRepo.ListResult = []map[string]interface{}{
			{"id": "doc-1"},
			{"id": "doc-2"},
		}
		mockRepo.ListTotal = 2
		mockRepo.ListCalls = 0

		docs, total, err := svc.ListDocumentsWithTracing(ctx, "org-123", "draft", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, 2, len(docs))
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 1, mockRepo.ListCalls)
	})

	t.Run("SendDocumentWithTracing succeeds", func(t *testing.T) {
		err := svc.SendDocumentWithTracing(ctx, "doc-1", "user-2", "Please approve")
		require.NoError(t, err)
		assert.Greater(t, len(mockPublisher.Events), 0)
	})

	t.Run("ApproveDocumentWithTracing succeeds", func(t *testing.T) {
		err := svc.ApproveDocumentWithTracing(ctx, "doc-1", "user-2", "Approved")
		require.NoError(t, err)
		assert.Greater(t, len(mockPublisher.Events), 0)
	})

	t.Run("RejectDocumentWithTracing succeeds", func(t *testing.T) {
		err := svc.RejectDocumentWithTracing(ctx, "doc-1", "user-2", "Not suitable")
		require.NoError(t, err)
		assert.Greater(t, len(mockPublisher.Events), 0)
	})

	t.Run("ArchiveDocumentWithTracing succeeds", func(t *testing.T) {
		err := svc.ArchiveDocumentWithTracing(ctx, "doc-1")
		require.NoError(t, err)
		assert.Greater(t, len(mockPublisher.Events), 0)
	})
}

// TestTracingProvider тестирует провайдер трассирования
func TestTracingProvider(t *testing.T) {
	tracing := observability.NewTracingProvider("test-service")

	ctx := context.Background()

	t.Run("StartSpan creates span", func(t *testing.T) {
		newCtx, span := tracing.StartSpan(ctx, "test.operation")
		require.NotNil(t, span)
		assert.NotNil(t, newCtx)
		span.End()
	})

	t.Run("OperationSpan works correctly", func(t *testing.T) {
		opSpan := tracing.NewOperationSpan(ctx, "test_op", "test-id")
		require.NotNil(t, opSpan)
		assert.NotNil(t, opSpan.Context())

		opSpan.AddEvent("test_event")
		opSpan.End(nil)
	})

	t.Run("RecordError captures error", func(t *testing.T) {
		opSpan := tracing.NewOperationSpan(ctx, "error_op", "test-id")

		err := assert.AnError
		opSpan.RecordError(err)
		opSpan.End(nil)
	})

	t.Run("GetTracingProvider returns singleton", func(t *testing.T) {
		tp1 := observability.GetTracingProvider()
		tp2 := observability.GetTracingProvider()
		assert.NotNil(t, tp1)
		assert.NotNil(t, tp2)
	})
}
