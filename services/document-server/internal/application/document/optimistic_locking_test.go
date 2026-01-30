// Файл document-server/internal/application/document/optimistic_locking_test.go содержит реализацию пакета document.
package document

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockRepositoryWithVersion для тестирования optimistic locking
type MockRepositoryWithVersion struct {
	documents map[string]map[string]interface{}
}

func NewMockRepositoryWithVersion() *MockRepositoryWithVersion {
	return &MockRepositoryWithVersion{
		documents: make(map[string]map[string]interface{}),
	}
}

func (m *MockRepositoryWithVersion) Create(ctx context.Context, id, organizationID, title, content, createdBy string, createdAt int64) error {
	m.documents[id] = map[string]interface{}{
		"id":                id,
		"organization_id":   organizationID,
		"title":             title,
		"content":           content,
		"status":            "draft",
		"created_by":        createdBy,
		"created_at":        createdAt,
		"updated_at":        createdAt,
		"status_changed_at": createdAt,
		"version":           1,
	}
	return nil
}

func (m *MockRepositoryWithVersion) Get(ctx context.Context, documentID string) (map[string]interface{}, error) {
	if doc, ok := m.documents[documentID]; ok {
		return doc, nil
	}
	return nil, nil
}

func (m *MockRepositoryWithVersion) GetWithVersion(ctx context.Context, documentID string) (map[string]interface{}, error) {
	if doc, ok := m.documents[documentID]; ok {
		return doc, nil
	}
	return nil, nil
}

func (m *MockRepositoryWithVersion) Update(ctx context.Context, documentID, title, content string, updatedAt int64) error {
	if doc, ok := m.documents[documentID]; ok {
		doc["title"] = title
		doc["content"] = content
		doc["updated_at"] = updatedAt
		version := doc["version"].(int)
		doc["version"] = version + 1
		return nil
	}
	return domain.ErrDocumentNotFound
}

func (m *MockRepositoryWithVersion) UpdateWithVersion(ctx context.Context, documentID, title, content string, expectedVersion int) error {
	if doc, ok := m.documents[documentID]; ok {
		currentVersion := doc["version"].(int)
		if currentVersion != expectedVersion {
			return domain.ErrVersionConflict
		}
		doc["title"] = title
		doc["content"] = content
		doc["version"] = expectedVersion + 1
		return nil
	}
	return domain.ErrDocumentNotFound
}

func (m *MockRepositoryWithVersion) UpdateStatus(ctx context.Context, documentID, status string, statusChangedAt int64) error {
	if doc, ok := m.documents[documentID]; ok {
		doc["status"] = status
		doc["status_changed_at"] = statusChangedAt
		return nil
	}
	return nil
}

func (m *MockRepositoryWithVersion) List(ctx context.Context, organizationID string, status string, page, perPage int) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}

func (m *MockRepositoryWithVersion) Delete(ctx context.Context, documentID string) error {
	delete(m.documents, documentID)
	return nil
}

// TestOptimisticLocking тестирует optimistic locking функциональность
func TestOptimisticLocking(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	repo := NewMockRepositoryWithVersion()
	service := NewService(repo, nil, logger)

	ctx := context.Background()

	t.Run("GetDocumentWithVersion returns document with version", func(t *testing.T) {
		documentID := "test-doc-1"
		repo.Create(ctx, documentID, "org-1", "Test Doc", "Content", "user-1", 1000)

		doc, err := service.GetDocumentWithVersion(ctx, documentID)
		require.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 1, doc["version"])
	})

	t.Run("UpdateDocumentWithVersion succeeds with correct version", func(t *testing.T) {
		documentID := "test-doc-2"
		repo.Create(ctx, documentID, "org-1", "Test Doc", "Content", "user-1", 1000)

		err := service.UpdateDocumentWithVersion(ctx, documentID, "Updated Title", "Updated Content", 1)
		require.NoError(t, err)

		doc, _ := service.GetDocumentWithVersion(ctx, documentID)
		assert.Equal(t, "Updated Title", doc["title"])
		assert.Equal(t, 2, doc["version"]) // Version incremented
	})

	t.Run("UpdateDocumentWithVersion fails with wrong version", func(t *testing.T) {
		documentID := "test-doc-3"
		repo.Create(ctx, documentID, "org-1", "Test Doc", "Content", "user-1", 1000)

		err := service.UpdateDocumentWithVersion(ctx, documentID, "Updated", "Updated", 999)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrVersionConflict, err)
	})

	t.Run("Multiple updates increment version correctly", func(t *testing.T) {
		documentID := "test-doc-4"
		repo.Create(ctx, documentID, "org-1", "Test Doc", "Content", "user-1", 1000)

		// Update 1
		service.UpdateDocumentWithVersion(ctx, documentID, "Update 1", "Content 1", 1)
		doc1, _ := service.GetDocumentWithVersion(ctx, documentID)
		assert.Equal(t, 2, doc1["version"])

		// Update 2
		service.UpdateDocumentWithVersion(ctx, documentID, "Update 2", "Content 2", 2)
		doc2, _ := service.GetDocumentWithVersion(ctx, documentID)
		assert.Equal(t, 3, doc2["version"])

		// Update 3
		service.UpdateDocumentWithVersion(ctx, documentID, "Update 3", "Content 3", 3)
		doc3, _ := service.GetDocumentWithVersion(ctx, documentID)
		assert.Equal(t, 4, doc3["version"])
	})

	t.Run("Concurrent update attempts with different versions", func(t *testing.T) {
		documentID := "test-doc-5"
		repo.Create(ctx, documentID, "org-1", "Test Doc", "Content", "user-1", 1000)

		// User 1 updates successfully
		err1 := service.UpdateDocumentWithVersion(ctx, documentID, "User 1 Update", "Content 1", 1)
		require.NoError(t, err1)

		// User 2 tries to update with old version - should fail
		err2 := service.UpdateDocumentWithVersion(ctx, documentID, "User 2 Update", "Content 2", 1)
		assert.Error(t, err2)
		assert.Equal(t, domain.ErrVersionConflict, err2)

		// User 2 updates with correct version - should succeed
		err3 := service.UpdateDocumentWithVersion(ctx, documentID, "User 2 Update", "Content 2", 2)
		require.NoError(t, err3)
	})
}
