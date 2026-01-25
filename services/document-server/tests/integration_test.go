package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/application/document"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/messaging"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/observability"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestDocumentServiceIntegration выполняет интеграционный тест полного цикла жизни документа
func TestDocumentServiceIntegration(t *testing.T) {
	// Пропускаем тест если нет БД
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Подключение к тестовой БД
	db, err := setupTestDatabase(t)
	require.NoError(t, err)
	defer db.Close()

	// Инициализация логгера
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Инициализация компонентов
	repo := repository.NewPostgresDocumentRepository(db)
	publisher := messaging.NewNoOpPublisher(logger)
	service := document.NewService(repo, publisher, logger)
	metrics := observability.InitMetrics()

	ctx := context.Background()

	t.Run("Metrics operations", func(t *testing.T) {
		// Проверяем что методы работают без паники
		start := time.Now()
		metrics.IncActiveRequests()
		metrics.DecActiveRequests()

		duration := time.Since(start)
		metrics.RecordRequestDuration(duration)

		metrics.RecordDocumentCreated()
		metrics.RecordDocumentUpdated()
		metrics.RecordDocumentApproved()
		metrics.RecordDocumentRejected()
		metrics.RecordDocumentArchived()
		metrics.RecordDocumentSent()

		metrics.RecordValidationError()
		metrics.RecordRepositoryError()
		metrics.RecordRabbitMQPublishError()

		metrics.RecordPaginationRequest(50)

		assert.True(t, true, "Metrics recording completed without panic")
	})

	t.Run("Document lifecycle and status transitions", func(t *testing.T) {
		organizationID := uuid.NewString()
		userID := uuid.NewString()
		cleanupDocuments(t, db, organizationID)

		docID, err := service.CreateDocument(ctx, organizationID, "Integration Lifecycle", "Lifecycle content", userID)
		require.NoError(t, err)
		require.NotEmpty(t, docID)

		doc, err := service.GetDocument(ctx, docID)
		require.NoError(t, err)
		require.Equal(t, "draft", doc["status"])
		require.Equal(t, 1, versionFromDoc(t, doc))

		require.NoError(t, service.UpdateDocument(ctx, docID, "Updated title", "Updated content"))
		doc, err = service.GetDocument(ctx, docID)
		require.NoError(t, err)
		require.Equal(t, "draft", doc["status"])
		require.Equal(t, 2, versionFromDoc(t, doc))

		require.NoError(t, service.SendDocument(ctx, docID, "recipient-1", "msg"))
		doc, err = service.GetDocument(ctx, docID)
		require.NoError(t, err)
		require.Equal(t, "sent", doc["status"])
		require.Equal(t, 3, versionFromDoc(t, doc))

		require.NoError(t, service.ApproveDocument(ctx, docID, "approver-1", "looks good"))
		doc, err = service.GetDocument(ctx, docID)
		require.NoError(t, err)
		require.Equal(t, "approved", doc["status"])
		require.Equal(t, 4, versionFromDoc(t, doc))

		require.NoError(t, service.ArchiveDocument(ctx, docID))
		doc, err = service.GetDocument(ctx, docID)
		require.NoError(t, err)
		require.Equal(t, "archived", doc["status"])
		require.Equal(t, 5, versionFromDoc(t, doc))

		docs, total, err := service.ListDocuments(ctx, organizationID, "", 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(docs), 1)
		require.True(t, containsDocumentID(docs, docID))
	})

	t.Run("Pagination returns total and pages", func(t *testing.T) {
		organizationID := uuid.NewString()
		userID := uuid.NewString()
		cleanupDocuments(t, db, organizationID)

		for i := 0; i < 3; i++ {
			_, err := service.CreateDocument(ctx, organizationID, fmt.Sprintf("Paginated %d", i), "content", userID)
			require.NoError(t, err)
		}

		page1, total, err := service.ListDocuments(ctx, organizationID, "", 1, 2)
		require.NoError(t, err)
		require.Equal(t, int64(3), total)
		require.Equal(t, 2, len(page1))

		page2, total2, err := service.ListDocuments(ctx, organizationID, "", 2, 2)
		require.NoError(t, err)
		require.Equal(t, int64(3), total2)
		require.GreaterOrEqual(t, len(page2), 1)
	})

	t.Run("Optimistic locking conflicts", func(t *testing.T) {
		organizationID := uuid.NewString()
		userID := uuid.NewString()
		cleanupDocuments(t, db, organizationID)

		docID, err := service.CreateDocument(ctx, organizationID, "Locking", "content", userID)
		require.NoError(t, err)

		docWithVersion, err := service.GetDocumentWithVersion(ctx, docID)
		require.NoError(t, err)
		initialVersion := versionFromDoc(t, docWithVersion)
		require.Equal(t, 1, initialVersion)

		require.NoError(t, service.UpdateDocumentWithVersion(ctx, docID, "Locking updated", "content", initialVersion))

		docAfterUpdate, err := service.GetDocumentWithVersion(ctx, docID)
		require.NoError(t, err)
		require.Equal(t, initialVersion+1, versionFromDoc(t, docAfterUpdate))

		err = service.UpdateDocumentWithVersion(ctx, docID, "Locking second", "content", initialVersion)
		require.Error(t, err)
		require.Equal(t, domain.ErrVersionConflict, err)
	})
}

// setupTestDatabase устанавливает тестовую БД
func setupTestDatabase(t *testing.T) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "document_user", "document_password", "document_db",
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Skipf("Cannot connect to test database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Cannot ping test database: %v", err)
		return nil, err
	}

	return db, nil
}

func cleanupDocuments(t *testing.T, db *sqlx.DB, organizationID string) {
	t.Helper()
	_, err := db.Exec("DELETE FROM documents WHERE organization_id = $1", organizationID)
	require.NoError(t, err)
}

func versionFromDoc(t *testing.T, doc map[string]interface{}) int {
	t.Helper()
	switch v := doc["version"].(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	default:
		t.Fatalf("unsupported version type %T", v)
		return 0
	}
}

func containsDocumentID(docs []map[string]interface{}, id string) bool {
	for _, d := range docs {
		if docID, ok := d["id"].(string); ok && docID == id {
			return true
		}
	}
	return false
}
