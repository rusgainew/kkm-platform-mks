// Файл document-query-server/internal/infrastructure/repository/inmemory_document_repository.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

const (
	defaultMaxEntries = 100000 // 100k документов
	evictionRatio     = 0.3    // Удалять 30% при превышении
)

var (
	documentCacheSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "document_cache_size",
		Help: "Current number of documents in memory",
	})
	documentOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "document_operation_duration_seconds",
		Help:    "Duration of document operations",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
	documentEvictions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "document_memory_evictions_total",
		Help: "Total number of memory evictions",
	})
)

type documentEntry struct {
	data       *pb.DocumentReadModel
	lastAccess time.Time
}

// InMemoryDocumentRepository хранит данные документов в памяти
type InMemoryDocumentRepository struct {
	mu         sync.RWMutex
	documents  map[string]*documentEntry // key: document_id
	logger     *zap.Logger
	maxEntries int
	lastUpdate time.Time // Track last update time for staleness checks
}

// NewInMemoryDocumentRepository создает новый in-memory репозиторий
func NewInMemoryDocumentRepository(logger *zap.Logger) *InMemoryDocumentRepository {
	return &InMemoryDocumentRepository{
		documents:  make(map[string]*documentEntry),
		logger:     logger,
		maxEntries: defaultMaxEntries,
	}
}

// GetLastUpdate returns the time of the last cache update
func (r *InMemoryDocumentRepository) GetLastUpdate() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdate
}

// GetCacheSize returns the current number of items in cache
func (r *InMemoryDocumentRepository) GetCacheSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.documents)
}

// GetDocument получает документ по ID
func (r *InMemoryDocumentRepository) GetDocument(ctx context.Context, documentID string) (*pb.DocumentReadModel, error) {
	start := time.Now()
	defer func() {
		documentOperationDuration.WithLabelValues("get").Observe(time.Since(start).Seconds())
	}()

	if documentID == "" {
		r.logger.Warn("GetDocument called with empty documentID")
		return nil, nil
	}

	r.mu.RLock()
	entry, ok := r.documents[documentID]
	r.mu.RUnlock()

	if !ok {
		return nil, nil // Документ не найден
	}

	// Обновляем lastAccess (требует Write lock)
	r.mu.Lock()
	entry.lastAccess = time.Now()
	r.mu.Unlock()

	// Клонируем для предотвращения race conditions
	return proto.Clone(entry.data).(*pb.DocumentReadModel), nil
}

// ListDocuments возвращает список документов с фильтрацией и пагинацией
func (r *InMemoryDocumentRepository) ListDocuments(ctx context.Context, offset, limit int32, status, docType, companyID, approvalStatus string) ([]*pb.DocumentReadModel, int64, error) {
	start := time.Now()
	defer func() {
		documentOperationDuration.WithLabelValues("list").Observe(time.Since(start).Seconds())
	}()

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Собираем все документы в slice
	allDocs := make([]*pb.DocumentReadModel, 0, len(r.documents))
	for _, entry := range r.documents {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		doc := entry.data
		// Применяем фильтры
		if status != "" && doc.Status != status {
			continue
		}
		if docType != "" && doc.DocumentType != docType {
			continue
		}
		if companyID != "" && doc.CompanyId != companyID {
			continue
		}
		if approvalStatus != "" && doc.ApprovalStatus != approvalStatus {
			continue
		}
		// Клонируем для предотвращения race conditions
		allDocs = append(allDocs, proto.Clone(doc).(*pb.DocumentReadModel))
	}

	totalCount := int64(len(allDocs))

	// Применяем пагинацию
	if offset >= int32(totalCount) {
		return []*pb.DocumentReadModel{}, totalCount, nil
	}

	end := offset + limit
	if end > int32(totalCount) {
		end = int32(totalCount)
	}

	return allDocs[offset:end], totalCount, nil
}

// SearchDocuments ищет документы по запросу
func (r *InMemoryDocumentRepository) SearchDocuments(ctx context.Context, query string, offset, limit int32, docType, companyID, approvalStatus string) ([]*pb.DocumentReadModel, int64, error) {
	start := time.Now()
	defer func() {
		documentOperationDuration.WithLabelValues("search").Observe(time.Since(start).Seconds())
	}()

	if query == "" {
		r.logger.Warn("SearchDocuments called with empty query")
		return []*pb.DocumentReadModel{}, 0, nil
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(query)
	var filtered []*pb.DocumentReadModel

	for _, entry := range r.documents {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		doc := entry.data
		// Применяем фильтры
		if docType != "" && doc.DocumentType != docType {
			continue
		}
		if companyID != "" && doc.CompanyId != companyID {
			continue
		}
		if approvalStatus != "" && doc.ApprovalStatus != approvalStatus {
			continue
		}

		// Поиск по номеру, заголовку, описанию
		if strings.Contains(strings.ToLower(doc.DocumentNumber), query) ||
			strings.Contains(strings.ToLower(doc.Title), query) ||
			strings.Contains(strings.ToLower(doc.Description), query) ||
			strings.Contains(strings.ToLower(doc.CompanyName), query) ||
			strings.Contains(strings.ToLower(doc.CreatedByUserName), query) {
			// Клонируем для предотвращения race conditions
			filtered = append(filtered, proto.Clone(doc).(*pb.DocumentReadModel))
		}
	}

	totalCount := int64(len(filtered))

	// Применяем пагинацию
	if offset >= int32(totalCount) {
		return []*pb.DocumentReadModel{}, totalCount, nil
	}

	end := offset + limit
	if end > int32(totalCount) {
		end = int32(totalCount)
	}

	return filtered[offset:end], totalCount, nil
}

// GetPendingApprovalDocuments возвращает документы в ожидании утверждения
func (r *InMemoryDocumentRepository) GetPendingApprovalDocuments(ctx context.Context, companyID string, offset, limit int32) ([]*pb.DocumentReadModel, int64, error) {
	start := time.Now()
	defer func() {
		documentOperationDuration.WithLabelValues("get_pending").Observe(time.Since(start).Seconds())
	}()

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var pending []*pb.DocumentReadModel

	for _, entry := range r.documents {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		doc := entry.data
		// Фильтр по компании
		if companyID != "" && doc.CompanyId != companyID {
			continue
		}

		// Документы со статусом ожидания утверждения
		if doc.ApprovalStatus == "pending" || doc.ApprovalStatus == "awaiting_approval" {
			// Клонируем для предотвращения race conditions
			pending = append(pending, proto.Clone(doc).(*pb.DocumentReadModel))
		}
	}

	totalCount := int64(len(pending))

	// Применяем пагинацию
	if offset >= int32(totalCount) {
		return []*pb.DocumentReadModel{}, totalCount, nil
	}

	end := offset + limit
	if end > int32(totalCount) {
		end = int32(totalCount)
	}

	return pending[offset:end], totalCount, nil
}

// UpsertDocument добавляет или обновляет документ (используется в RabbitMQ consumer)
func (r *InMemoryDocumentRepository) UpsertDocument(ctx context.Context, doc *pb.DocumentReadModel) error {
	start := time.Now()
	defer func() {
		documentOperationDuration.WithLabelValues("upsert").Observe(time.Since(start).Seconds())
	}()

	if doc == nil || doc.Id == "" {
		r.logger.Warn("Attempted to upsert document with empty ID")
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	action := "created"
	if _, exists := r.documents[doc.Id]; exists {
		action = "updated"
	}

	// Проверяем лимит памяти
	if len(r.documents) >= r.maxEntries {
		r.evictOldEntries()
	}

	// Клонируем для предотвращения изменения извне
	r.documents[doc.Id] = &documentEntry{
		data:       proto.Clone(doc).(*pb.DocumentReadModel),
		lastAccess: time.Now(),
	}

	r.lastUpdate = time.Now()
	documentCacheSize.Set(float64(len(r.documents)))

	r.logger.Debug("Document upserted",
		zap.String("action", action),
		zap.String("document_id", doc.Id),
		zap.Int("cache_size", len(r.documents)))

	return nil
}

// DeleteDocument удаляет документ из памяти (используется в RabbitMQ consumer)
func (r *InMemoryDocumentRepository) DeleteDocument(ctx context.Context, documentID string) error {
	start := time.Now()
	defer func() {
		documentOperationDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	if documentID == "" {
		r.logger.Warn("DeleteDocument called with empty documentID")
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.documents, documentID)
	documentCacheSize.Set(float64(len(r.documents)))

	r.logger.Debug("Document deleted",
		zap.String("document_id", documentID),
		zap.Int("cache_size", len(r.documents)))

	return nil
}

// GetAllDocuments возвращает все документы (для отладки)
func (r *InMemoryDocumentRepository) GetAllDocuments() []*pb.DocumentReadModel {
	r.mu.RLock()
	defer r.mu.RUnlock()

	docs := make([]*pb.DocumentReadModel, 0, len(r.documents))
	for _, entry := range r.documents {
		// Клонируем для предотвращения race conditions
		docs = append(docs, proto.Clone(entry.data).(*pb.DocumentReadModel))
	}
	return docs
}

// Count возвращает количество документов в памяти
func (r *InMemoryDocumentRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.documents)
}

// evictOldEntries удаляет старые записи при превышении лимита (только для внутреннего использования, вызывается с активным Lock)
func (r *InMemoryDocumentRepository) evictOldEntries() {
	toEvict := int(float64(len(r.documents)) * evictionRatio)
	if toEvict == 0 {
		toEvict = 1
	}

	// Собираем все ключи с lastAccess
	type entryInfo struct {
		key        string
		lastAccess time.Time
	}
	entries := make([]entryInfo, 0, len(r.documents))
	for key, entry := range r.documents {
		entries = append(entries, entryInfo{key: key, lastAccess: entry.lastAccess})
	}

	// Сортируем по lastAccess (старые первыми)
	// Простой bubble sort для малого количества
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].lastAccess.After(entries[j].lastAccess) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Удаляем старые записи
	evicted := 0
	for _, entry := range entries {
		if evicted >= toEvict {
			break
		}
		delete(r.documents, entry.key)
		evicted++
	}

	documentEvictions.Add(float64(evicted))
	documentCacheSize.Set(float64(len(r.documents)))

	r.logger.Info("Memory eviction performed",
		zap.Int("evicted_count", evicted),
		zap.Int("remaining_count", len(r.documents)),
		zap.Int("max_entries", r.maxEntries))
}

// sortDocuments сортирует документы по указанному полю и порядку
func sortDocuments(documents []*pb.DocumentReadModel, field string, order string) {
	if field == "" {
		return
	}

	less := func(i, j int) bool {
		var result bool
		switch field {
		case "document_number":
			result = strings.ToLower(documents[i].DocumentNumber) < strings.ToLower(documents[j].DocumentNumber)
		case "title":
			result = strings.ToLower(documents[i].Title) < strings.ToLower(documents[j].Title)
		case "created_at":
			result = documents[i].CreatedAt < documents[j].CreatedAt
		case "status":
			result = strings.ToLower(documents[i].Status) < strings.ToLower(documents[j].Status)
		default:
			result = documents[i].CreatedAt < documents[j].CreatedAt
		}

		if order == "DESC" {
			return !result
		}
		return result
	}

	sort.Slice(documents, less)
}
