// Файл document-server/internal/application/document/caching_test.go содержит реализацию пакета document.
package document

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockCacheForTesting для тестирования кэширования
type MockCacheForTesting struct {
	Data               map[string][]byte
	GetCalls           int
	SetCalls           int
	DeletePatternCalls int
	GetErr             error
	SetErr             error
	DeleteErr          error
}

func NewMockCacheForTesting() *MockCacheForTesting {
	return &MockCacheForTesting{
		Data: make(map[string][]byte),
	}
}

func (mc *MockCacheForTesting) Get(ctx context.Context, key string) ([]byte, error) {
	mc.GetCalls++
	if mc.GetErr != nil {
		return nil, mc.GetErr
	}
	return mc.Data[key], nil
}

func (mc *MockCacheForTesting) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	mc.SetCalls++
	if mc.SetErr != nil {
		return mc.SetErr
	}
	mc.Data[key] = value
	return nil
}

func (mc *MockCacheForTesting) Delete(ctx context.Context, key string) error {
	if mc.DeleteErr != nil {
		return mc.DeleteErr
	}
	delete(mc.Data, key)
	return nil
}

func (mc *MockCacheForTesting) DeletePattern(ctx context.Context, pattern string) error {
	mc.DeletePatternCalls++
	if mc.DeleteErr != nil {
		return mc.DeleteErr
	}
	prefix := pattern[:len(pattern)-1]
	for key := range mc.Data {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(mc.Data, key)
		}
	}
	return nil
}

func (mc *MockCacheForTesting) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := mc.Data[key]
	return exists, nil
}

func (mc *MockCacheForTesting) TTL(ctx context.Context, key string) (int64, error) {
	if _, exists := mc.Data[key]; !exists {
		return -2, nil
	}
	return 300, nil
}

func (mc *MockCacheForTesting) Close() error {
	return nil
}

// TestGetDocumentCached тестирует получение документа с кэша
func TestGetDocumentCached(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockRepository()
	mockPublisher := &MockPublisher{}
	mockCache := NewMockCacheForTesting()

	svc := NewService(mockRepo, mockPublisher, logger)
	cachingSvc := NewCachingService(svc, mockCache, 5*time.Minute)

	ctx := context.Background()
	documentID := "test-doc-1"

	t.Run("First access populates cache", func(t *testing.T) {
		doc := map[string]interface{}{
			"id":                documentID,
			"organization_id":   "org-123",
			"title":             "Test Document",
			"content":           "Test Content",
			"status":            "draft",
			"created_by":        "user-456",
			"created_at":        int64(1000000),
			"updated_at":        int64(1000000),
			"status_changed_at": int64(1000000),
			"version":           int32(1),
		}
		mockRepo.GetResult = doc

		result, err := cachingSvc.GetDocumentCached(ctx, documentID)
		require.NoError(t, err)
		assert.Equal(t, documentID, result["id"])

		initialGets := mockRepo.GetCalls
		initialCacheGets := mockCache.GetCalls

		// Второй доступ должен идти из кэша
		result2, err := cachingSvc.GetDocumentCached(ctx, documentID)
		require.NoError(t, err)
		assert.Equal(t, documentID, result2["id"])
		assert.Equal(t, initialGets, mockRepo.GetCalls) // Не обращались к репо!
		assert.Greater(t, mockCache.GetCalls, initialCacheGets)
	})

	t.Run("Cache hit returns cached data", func(t *testing.T) {
		mockCache.Data = make(map[string][]byte)
		mockCache.GetCalls = 0
		mockRepo.GetCalls = 0

		doc := map[string]interface{}{
			"id":    "test-doc-2",
			"title": "Cached Document",
		}
		docJSON, _ := cache.SerializeDocument(doc)
		mockCache.Data[cache.GetDocumentCacheKey("test-doc-2")] = docJSON

		result, err := cachingSvc.GetDocumentCached(ctx, "test-doc-2")
		require.NoError(t, err)
		assert.Equal(t, "test-doc-2", result["id"])
		assert.Equal(t, "Cached Document", result["title"])
		assert.Equal(t, 0, mockRepo.GetCalls) // Не обращались к репо!
		assert.Greater(t, mockCache.GetCalls, 0)
	})

	t.Run("Cache error falls back to database", func(t *testing.T) {
		mockCache.GetErr = errors.New("redis down")
		mockRepo.GetCalls = 0
		mockRepo.GetResult = map[string]interface{}{
			"id":    "test-doc-3",
			"title": "Fallback Document",
		}

		result, err := cachingSvc.GetDocumentCached(ctx, "test-doc-3")
		require.NoError(t, err)
		assert.Equal(t, "test-doc-3", result["id"])
		assert.Greater(t, mockRepo.GetCalls, 0) // Обращились к репо из-за ошибки кэша
		mockCache.GetErr = nil                  // Очищаем ошибку для следующих тестов
	})
}

// TestListDocumentsCached тестирует получение списка документов с кэша
func TestListDocumentsCached(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockRepository()
	mockPublisher := &MockPublisher{}
	mockCache := NewMockCacheForTesting()

	svc := NewService(mockRepo, mockPublisher, logger)
	cachingSvc := NewCachingService(svc, mockCache, 5*time.Minute)

	ctx := context.Background()

	t.Run("First list call populates cache", func(t *testing.T) {
		mockRepo.ListResult = []map[string]interface{}{
			{"id": "doc-1", "title": "Doc 1"},
			{"id": "doc-2", "title": "Doc 2"},
		}
		mockRepo.ListTotal = 2

		docs, total, err := cachingSvc.ListDocumentsCached(ctx, "org-123", "draft", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, 2, len(docs))
		assert.Equal(t, int64(2), total)
		assert.Greater(t, mockCache.SetCalls, 0)
	})
}

// TestInvalidateCache тестирует инвалидацию кэша при обновлении
func TestInvalidateCache(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockRepository()
	mockPublisher := &MockPublisher{}
	mockCache := NewMockCacheForTesting()

	svc := NewService(mockRepo, mockPublisher, logger)
	cachingSvc := NewCachingService(svc, mockCache, 5*time.Minute)

	ctx := context.Background()
	documentID := "test-doc-4"

	t.Run("Update invalidates document cache", func(t *testing.T) {
		mockCache.Data = make(map[string][]byte)
		mockCache.DeletePatternCalls = 0

		doc := map[string]interface{}{"id": documentID, "title": "Old Title"}
		docJSON, _ := cache.SerializeDocument(doc)
		mockCache.Data[cache.GetDocumentCacheKey(documentID)] = docJSON

		err := cachingSvc.UpdateDocumentCached(ctx, documentID, "New Title", "New Content")
		require.NoError(t, err)

		// DeletePattern должен быть вызван
		assert.Greater(t, mockCache.DeletePatternCalls, 0)
	})

	t.Run("Create invalidates list cache", func(t *testing.T) {
		mockCache.DeletePatternCalls = 0

		_, err := cachingSvc.CreateDocumentCached(ctx, "org-123", "Title", "Content", "user-1")
		require.NoError(t, err)

		assert.Greater(t, mockCache.DeletePatternCalls, 0)
	})
}

// TestCacheTTLConfiguration тестирует конфигурацию TTL
func TestCacheTTLConfiguration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockRepository()
	mockPublisher := &MockPublisher{}
	mockCache := NewMockCacheForTesting()

	svc := NewService(mockRepo, mockPublisher, logger)
	cachingSvc := NewCachingService(svc, mockCache, 5*time.Minute)

	t.Run("TTL can be configured", func(t *testing.T) {
		assert.Equal(t, 5*time.Minute, cachingSvc.GetCacheTTL())

		cachingSvc.SetCacheTTL(10 * time.Minute)
		assert.Equal(t, 10*time.Minute, cachingSvc.GetCacheTTL())
	})

	t.Run("Zero TTL is ignored", func(t *testing.T) {
		currentTTL := cachingSvc.GetCacheTTL()
		cachingSvc.SetCacheTTL(0)
		assert.Equal(t, currentTTL, cachingSvc.GetCacheTTL())
	})
}

// TestClearAllCache тестирует очистку всего кэша
func TestClearAllCache(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := NewMockRepository()
	mockPublisher := &MockPublisher{}
	mockCache := NewMockCacheForTesting()

	svc := NewService(mockRepo, mockPublisher, logger)
	cachingSvc := NewCachingService(svc, mockCache, 5*time.Minute)

	ctx := context.Background()

	t.Run("ClearAllCache removes all document cache entries", func(t *testing.T) {
		mockCache.Data = make(map[string][]byte)
		mockCache.DeletePatternCalls = 0

		mockCache.Data["document:doc-1"] = []byte("data1")
		mockCache.Data["document:doc-2"] = []byte("data2")

		err := cachingSvc.ClearAllCache(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, mockCache.DeletePatternCalls)
	})
}
