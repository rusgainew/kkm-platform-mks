// Файл document-server/internal/application/document/caching.go содержит реализацию пакета document.
package document

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/document-server/internal/infrastructure/cache"
	"go.uber.org/zap"
)

// CachingService обёртка для Service с кэшированием
type CachingService struct {
	*Service
	cacheLayer ports.CacheLayer
	cacheTTL   time.Duration
}

// NewCachingService создаёт сервис с кэшированием
func NewCachingService(svc *Service, cacheLayer ports.CacheLayer, cacheTTL time.Duration) *CachingService {
	if cacheTTL == 0 {
		cacheTTL = 5 * time.Minute // Значение по умолчанию
	}

	return &CachingService{
		Service:    svc,
		cacheLayer: cacheLayer,
		cacheTTL:   cacheTTL,
	}
}

// GetDocumentCached получает документ с использованием кэша
func (cs *CachingService) GetDocumentCached(ctx context.Context, documentID string) (map[string]interface{}, error) {
	cacheKey := cache.GetDocumentCacheKey(documentID)

	// Проверяем кэш
	cachedData, err := cs.cacheLayer.Get(ctx, cacheKey)
	if err != nil {
		cs.logger.Warn("Cache get error (will fallback to database)",
			zap.String("document_id", documentID),
			zap.Error(err),
		)
		// Не прерываем выполнение, идём в БД
	}

	if cachedData != nil {
		// Распаковываем из кэша
		doc, err := cache.DeserializeDocument(cachedData)
		if err != nil {
			cs.logger.Warn("Failed to deserialize document from cache",
				zap.String("document_id", documentID),
				zap.Error(err),
			)
			// Не прерываем, идём в БД
		} else {
			cs.logger.Debug("Document retrieved from cache", zap.String("document_id", documentID))
			return doc, nil
		}
	}

	// Получаем из БД
	doc, err := cs.repo.Get(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	docJSON, err := cache.SerializeDocument(doc)
	if err != nil {
		cs.logger.Warn("Failed to serialize document for cache",
			zap.String("document_id", documentID),
			zap.Error(err),
		)
		// Не прерываем, возвращаем документ
		return doc, nil
	}

	if err := cs.cacheLayer.Set(ctx, cacheKey, docJSON, cs.cacheTTL); err != nil {
		cs.logger.Warn("Failed to cache document",
			zap.String("document_id", documentID),
			zap.Error(err),
		)
		// Не прерываем, возвращаем документ
	}

	return doc, nil
}

// UpdateDocumentCached обновляет документ и инвалидирует кэш
func (cs *CachingService) UpdateDocumentCached(ctx context.Context, documentID, title, content string) error {
	// Выполняем обновление
	if err := cs.Service.UpdateDocument(ctx, documentID, title, content); err != nil {
		return err
	}

	// Инвалидируем кэш
	if err := cache.InvalidateDocumentCache(ctx, cs.cacheLayer, documentID); err != nil {
		cs.logger.Warn("Failed to invalidate cache after update",
			zap.String("document_id", documentID),
			zap.Error(err),
		)
		// Не прерываем, обновление уже выполнено
	}

	return nil
}

// CreateDocumentCached создаёт документ и инвалидирует кэш списков
func (cs *CachingService) CreateDocumentCached(ctx context.Context, organizationID, title, content, createdBy string) (string, error) {
	// Выполняем создание
	documentID, err := cs.Service.CreateDocument(ctx, organizationID, title, content, createdBy)
	if err != nil {
		return "", err
	}

	// Инвалидируем кэш списков организации
	if err := cache.InvalidateListCache(ctx, cs.cacheLayer, organizationID); err != nil {
		cs.logger.Warn("Failed to invalidate list cache after create",
			zap.String("organization_id", organizationID),
			zap.Error(err),
		)
		// Не прерываем, создание уже выполнено
	}

	return documentID, nil
}

// ListDocumentsCached получает список документов с использованием кэша
func (cs *CachingService) ListDocumentsCached(ctx context.Context, organizationID, status string, page, perPage int) ([]map[string]interface{}, int64, error) {
	cacheKey := cache.GetDocumentListCacheKey(organizationID, status, page, perPage)

	// Проверяем кэш
	cachedData, err := cs.cacheLayer.Get(ctx, cacheKey)
	if err != nil {
		cs.logger.Warn("Cache get error (will fallback to database)",
			zap.String("organization_id", organizationID),
			zap.Error(err),
		)
		// Не прерываем
	}

	if cachedData != nil {
		// Распаковываем из кэша
		var cacheResult struct {
			Documents []map[string]interface{} `json:"documents"`
			Total     int64                    `json:"total"`
		}

		if err := json.Unmarshal(cachedData, &cacheResult); err != nil {
			cs.logger.Warn("Failed to deserialize list from cache",
				zap.String("cache_key", cacheKey),
				zap.Error(err),
			)
			// Не прерываем
		} else {
			cs.logger.Debug("Document list retrieved from cache",
				zap.String("organization_id", organizationID),
				zap.Int("count", len(cacheResult.Documents)),
			)
			return cacheResult.Documents, cacheResult.Total, nil
		}
	}

	// Получаем из БД
	docs, total, err := cs.Service.ListDocuments(ctx, organizationID, status, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	// Сохраняем в кэш
	cacheResult := struct {
		Documents []map[string]interface{} `json:"documents"`
		Total     int64                    `json:"total"`
	}{
		Documents: docs,
		Total:     total,
	}

	listJSON, err := json.Marshal(cacheResult)
	if err != nil {
		cs.logger.Warn("Failed to serialize list for cache",
			zap.String("organization_id", organizationID),
			zap.Error(err),
		)
		// Не прерываем
		return docs, total, nil
	}

	if err := cs.cacheLayer.Set(ctx, cacheKey, listJSON, cs.cacheTTL); err != nil {
		cs.logger.Warn("Failed to cache list",
			zap.String("cache_key", cacheKey),
			zap.Error(err),
		)
		// Не прерываем
	}

	return docs, total, nil
}

// DeleteDocumentCached удаляет документ и инвалидирует кэш
func (cs *CachingService) DeleteDocumentCached(ctx context.Context, documentID string) error {
	// Сначала получаем организацию документа для инвалидации кэша
	var organizationID string
	doc, err := cs.repo.Get(ctx, documentID)
	if err == nil && doc != nil {
		if org, ok := doc["organization_id"].(string); ok {
			organizationID = org
		}
	}

	// Выполняем удаление
	if err := cs.repo.Delete(ctx, documentID); err != nil {
		return err
	}

	// Инвалидируем кэш
	if err := cache.InvalidateDocumentCache(ctx, cs.cacheLayer, documentID); err != nil {
		cs.logger.Warn("Failed to invalidate cache after delete",
			zap.String("document_id", documentID),
			zap.Error(err),
		)
		// Не прерываем
	}

	// Инвалидируем кэш списков организации
	if organizationID != "" {
		if err := cache.InvalidateListCache(ctx, cs.cacheLayer, organizationID); err != nil {
			cs.logger.Warn("Failed to invalidate list cache after delete",
				zap.String("organization_id", organizationID),
				zap.Error(err),
			)
			// Не прерываем
		}
	}

	return nil
}

// ClearAllCache очищает весь кэш документов
func (cs *CachingService) ClearAllCache(ctx context.Context) error {
	if err := cs.cacheLayer.DeletePattern(ctx, "document:*"); err != nil {
		cs.logger.Error("Failed to clear all document cache", zap.Error(err))
		return err
	}

	cs.logger.Info("All document cache cleared")
	return nil
}

// GetCacheTTL возвращает текущее значение TTL кэша
func (cs *CachingService) GetCacheTTL() time.Duration {
	return cs.cacheTTL
}

// SetCacheTTL устанавливает новое значение TTL кэша
func (cs *CachingService) SetCacheTTL(ttl time.Duration) {
	if ttl > 0 {
		cs.cacheTTL = ttl
		cs.logger.Info("Cache TTL updated", zap.Duration("ttl", ttl))
	}
}
