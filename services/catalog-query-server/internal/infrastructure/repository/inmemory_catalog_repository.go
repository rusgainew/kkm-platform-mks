// Файл catalog-query-server/internal/infrastructure/repository/inmemory_catalog_repository.go содержит реализацию пакета repository.
package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rusgainew/kkm-project-mks/catalog-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/lib/conversion"
	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	defaultMaxEntries = 100000 // 100k каталогов
	evictionRatio     = 0.3    // Удалять 30% при превышении
)

var (
	catalogCacheSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "catalog_cache_size",
		Help: "Current number of catalogs in memory",
	})
	catalogOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "catalog_operation_duration_seconds",
		Help:    "Duration of catalog operations",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
	catalogEvictions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "catalog_memory_evictions_total",
		Help: "Total number of memory evictions",
	})
	catalogIndexSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "catalog_index_size",
		Help: "Size of secondary indexes",
	}, []string{"index_type"})
)

type catalogEntry struct {
	data       *dictionaries.Catalog
	lastAccess time.Time
}

// InMemoryCatalogRepository implements in-memory catalog repository
// Data is populated via RabbitMQ events
type InMemoryCatalogRepository struct {
	mu         sync.RWMutex
	catalogs   map[string]*catalogEntry // key: code
	indexTnved map[string][]string      // index: tnved_code -> []catalog_code
	indexGked  map[string][]string      // index: gked_code -> []catalog_code
	maxEntries int
	logger     *zap.Logger
	lastUpdate time.Time // Track last update time for staleness checks
}

// NewInMemoryCatalogRepository creates new in-memory repository
func NewInMemoryCatalogRepository(logger *zap.Logger) ports.CatalogQueryRepository {
	return &InMemoryCatalogRepository{
		catalogs:   make(map[string]*catalogEntry),
		indexTnved: make(map[string][]string),
		indexGked:  make(map[string][]string),
		maxEntries: defaultMaxEntries,
		logger:     logger,
	}
}

// GetLastUpdate returns the time of the last cache update
func (r *InMemoryCatalogRepository) GetLastUpdate() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdate
}

// GetCacheSize returns the current number of items in cache
func (r *InMemoryCatalogRepository) GetCacheSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.catalogs)
}

// ListCatalogs returns paginated list of catalogs
func (r *InMemoryCatalogRepository) ListCatalogs(ctx context.Context, page, size int32) ([]*dictionaries.Catalog, int32, error) {
	start := time.Now()
	defer func() {
		catalogOperationDuration.WithLabelValues("list").Observe(time.Since(start).Seconds())
	}()

	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}
	if page < 1 {
		page = 1
	}
	if size > 100 {
		size = 10
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	allCatalogs := make([]*dictionaries.Catalog, 0, len(r.catalogs))
	for _, entry := range r.catalogs {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		allCatalogs = append(allCatalogs, proto.Clone(entry.data).(*dictionaries.Catalog))
	}

	totalCount := conversion.SafeIntToInt32WithDefault(len(allCatalogs), 0)
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.Catalog{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return allCatalogs[offset:end], totalCount, nil
}

// ListCatalogsWithFilter returns filtered and sorted list of catalogs
func (r *InMemoryCatalogRepository) ListCatalogsWithFilter(ctx context.Context, filter *ports.CatalogFilter, sort *ports.CatalogSort, page, size int32) ([]*dictionaries.Catalog, int32, error) {
	start := time.Now()
	defer func() {
		catalogOperationDuration.WithLabelValues("list_with_filter").Observe(time.Since(start).Seconds())
	}()

	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}
	if page < 1 {
		page = 1
	}
	if size > 100 {
		size = 10
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Filter catalogs
	filtered := make([]*dictionaries.Catalog, 0)
	for _, entry := range r.catalogs {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		if !matchesFilter(entry.data, filter) {
			continue
		}
		filtered = append(filtered, proto.Clone(entry.data).(*dictionaries.Catalog))
	}

	// Apply sorting
	if sort != nil && sort.Field != "" {
		sortCatalogs(filtered, sort)
	}

	totalCount := conversion.SafeIntToInt32WithDefault(len(filtered), 0)
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.Catalog{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// SearchCatalogs performs full-text search
func (r *InMemoryCatalogRepository) SearchCatalogs(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.Catalog, int32, error) {
	if searchText == "" {
		return []*dictionaries.Catalog{}, 0, nil
	}
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	searchLower := strings.ToLower(searchText)
	filtered := make([]*dictionaries.Catalog, 0)

	for _, entry := range r.catalogs {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		catalog := entry.data
		if strings.Contains(strings.ToLower(catalog.Name), searchLower) ||
			strings.Contains(strings.ToLower(catalog.Number), searchLower) {
			filtered = append(filtered, proto.Clone(catalog).(*dictionaries.Catalog))
		}
	}

	totalCount := conversion.SafeIntToInt32WithDefault(len(filtered), 0)
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.Catalog{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// GetCatalogByNumber returns catalog by number
func (r *InMemoryCatalogRepository) GetCatalogByNumber(ctx context.Context, number string) (*dictionaries.Catalog, error) {
	start := time.Now()
	defer func() {
		catalogOperationDuration.WithLabelValues("get_by_number").Observe(time.Since(start).Seconds())
	}()

	if number == "" {
		return nil, nil
	}
	r.mu.RLock()
	entry, exists := r.catalogs[number]
	r.mu.RUnlock()

	if !exists {
		return nil, nil
	}

	// Обновляем lastAccess (требует Write lock)
	r.mu.Lock()
	entry.lastAccess = time.Now()
	r.mu.Unlock()

	return proto.Clone(entry.data).(*dictionaries.Catalog), nil
}

// GetCatalogsByTnvedCode returns catalogs by TNVED code
func (r *InMemoryCatalogRepository) GetCatalogsByTnvedCode(ctx context.Context, tnvedCode string, page, size int32) ([]*dictionaries.Catalog, int32, error) {
	start := time.Now()
	defer func() {
		catalogOperationDuration.WithLabelValues("get_by_tnved").Observe(time.Since(start).Seconds())
	}()

	if tnvedCode == "" {
		return []*dictionaries.Catalog{}, 0, nil
	}
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Используем secondary index для O(1) поиска
	catalogCodes, exists := r.indexTnved[tnvedCode]
	if !exists || len(catalogCodes) == 0 {
		return []*dictionaries.Catalog{}, 0, nil
	}

	// Собираем каталоги из индекса
	filtered := make([]*dictionaries.Catalog, 0, len(catalogCodes))
	for _, code := range catalogCodes {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}

		if entry, exists := r.catalogs[code]; exists {
			filtered = append(filtered, proto.Clone(entry.data).(*dictionaries.Catalog))
		}
	}

	totalCount := conversion.SafeIntToInt32WithDefault(len(filtered), 0)
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.Catalog{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// UpsertCatalog adds or updates catalog (used by event handlers)
func (r *InMemoryCatalogRepository) UpsertCatalog(catalog *dictionaries.Catalog) {
	start := time.Now()
	defer func() {
		catalogOperationDuration.WithLabelValues("upsert").Observe(time.Since(start).Seconds())
	}()

	if catalog == nil || catalog.Number == "" {
		r.logger.Warn("Attempted to upsert catalog with empty number")
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	action := "update"
	if _, exists := r.catalogs[catalog.Number]; !exists {
		action = "insert"
	}

	// Проверяем лимит памяти
	if len(r.catalogs) >= r.maxEntries {
		r.evictOldEntries()
	}

	// Удаляем старые индексы если запись существует
	if oldEntry, exists := r.catalogs[catalog.Number]; exists {
		r.removeFromIndex(catalog.Number, oldEntry.data)
	}

	// Сохраняем каталог
	r.catalogs[catalog.Number] = &catalogEntry{
		data:       proto.Clone(catalog).(*dictionaries.Catalog),
		lastAccess: time.Now(),
	}

	// Обновляем время последнего обновления кэша
	r.lastUpdate = time.Now()

	// Добавляем в индексы
	r.addToIndex(catalog.Number, catalog)

	// Обновляем метрики
	catalogCacheSize.Set(float64(len(r.catalogs)))
	catalogIndexSize.WithLabelValues("tnved").Set(float64(len(r.indexTnved)))
	catalogIndexSize.WithLabelValues("gked").Set(float64(len(r.indexGked)))

	r.logger.Debug("Catalog upserted",
		zap.String("action", action),
		zap.String("number", catalog.Number),
		zap.String("tnved_code", catalog.TnvedCode),
		zap.Int("cache_size", len(r.catalogs)))
}

// DeleteCatalog removes catalog (used by event handlers)
func (r *InMemoryCatalogRepository) DeleteCatalog(code string) {
	start := time.Now()
	defer func() {
		catalogOperationDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	if code == "" {
		r.logger.Warn("Attempted to delete catalog with empty code")
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	// Удаляем из индексов
	if entry, exists := r.catalogs[code]; exists {
		r.removeFromIndex(code, entry.data)
	}

	delete(r.catalogs, code)

	// Обновляем метрики
	catalogCacheSize.Set(float64(len(r.catalogs)))
	catalogIndexSize.WithLabelValues("tnved").Set(float64(len(r.indexTnved)))
	catalogIndexSize.WithLabelValues("gked").Set(float64(len(r.indexGked)))

	r.logger.Debug("Catalog deleted", zap.String("code", code), zap.Int("cache_size", len(r.catalogs)))
}

// matchesFilter checks if catalog matches filter criteria
func matchesFilter(catalog *dictionaries.Catalog, filter *ports.CatalogFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Name != "" && !strings.Contains(strings.ToLower(catalog.Name), strings.ToLower(filter.Name)) {
		return false
	}

	if filter.Number != "" && catalog.Number != filter.Number {
		return false
	}

	if filter.TnvedCode != "" && catalog.TnvedCode != filter.TnvedCode {
		return false
	}

	if filter.GkedCode != "" && catalog.GkedCode != filter.GkedCode {
		return false
	}

	if filter.SearchText != "" {
		searchLower := strings.ToLower(filter.SearchText)
		if !strings.Contains(strings.ToLower(catalog.Name), searchLower) &&
			!strings.Contains(strings.ToLower(catalog.Number), searchLower) {
			return false
		}
	}

	return true
}

// evictOldEntries удаляет старые записи при превышении лимита (только для внутреннего использования, вызывается с активным Lock)
func (r *InMemoryCatalogRepository) evictOldEntries() {
	toEvict := int(float64(len(r.catalogs)) * evictionRatio)
	if toEvict == 0 {
		toEvict = 1
	}

	// Собираем все ключи с lastAccess
	type entryInfo struct {
		key        string
		lastAccess time.Time
	}
	entries := make([]entryInfo, 0, len(r.catalogs))
	for key, entry := range r.catalogs {
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
		// Удаляем из индексов
		if catalogEntry, exists := r.catalogs[entry.key]; exists {
			r.removeFromIndex(entry.key, catalogEntry.data)
		}
		delete(r.catalogs, entry.key)
		evicted++
	}

	// Обновляем метрики
	catalogEvictions.Add(float64(evicted))
	catalogCacheSize.Set(float64(len(r.catalogs)))

	r.logger.Info("Memory eviction performed",
		zap.Int("evicted_count", evicted),
		zap.Int("remaining_count", len(r.catalogs)),
		zap.Int("max_entries", r.maxEntries))
}

// addToIndex добавляет каталог в secondary indexes (вызывается с активным Lock)
func (r *InMemoryCatalogRepository) addToIndex(code string, catalog *dictionaries.Catalog) {
	// Индекс по TNVED
	if catalog.TnvedCode != "" {
		r.indexTnved[catalog.TnvedCode] = append(r.indexTnved[catalog.TnvedCode], code)
	}
	// Индекс по GKED
	if catalog.GkedCode != "" {
		r.indexGked[catalog.GkedCode] = append(r.indexGked[catalog.GkedCode], code)
	}
}

// removeFromIndex удаляет каталог из secondary indexes (вызывается с активным Lock)
func (r *InMemoryCatalogRepository) removeFromIndex(code string, catalog *dictionaries.Catalog) {
	// Удаляем из индекса TNVED
	if catalog.TnvedCode != "" {
		codes := r.indexTnved[catalog.TnvedCode]
		for i, c := range codes {
			if c == code {
				r.indexTnved[catalog.TnvedCode] = append(codes[:i], codes[i+1:]...)
				break
			}
		}
		// Удаляем пустой ключ
		if len(r.indexTnved[catalog.TnvedCode]) == 0 {
			delete(r.indexTnved, catalog.TnvedCode)
		}
	}
	// Удаляем из индекса GKED
	if catalog.GkedCode != "" {
		codes := r.indexGked[catalog.GkedCode]
		for i, c := range codes {
			if c == code {
				r.indexGked[catalog.GkedCode] = append(codes[:i], codes[i+1:]...)
				break
			}
		}
		if len(r.indexGked[catalog.GkedCode]) == 0 {
			delete(r.indexGked, catalog.GkedCode)
		}
	}
}

// sortCatalogs sorts catalog slice according to sort parameters
func sortCatalogs(catalogs []*dictionaries.Catalog, sortParams *ports.CatalogSort) {
	if sortParams == nil || sortParams.Field == "" {
		return
	}

	less := func(i, j int) bool {
		var result bool
		switch sortParams.Field {
		case "name":
			result = strings.ToLower(catalogs[i].Name) < strings.ToLower(catalogs[j].Name)
		case "number":
			result = catalogs[i].Number < catalogs[j].Number
		case "tnved_code":
			result = catalogs[i].TnvedCode < catalogs[j].TnvedCode
		case "gked_code":
			result = catalogs[i].GkedCode < catalogs[j].GkedCode
		default:
			// Default sort by number
			result = catalogs[i].Number < catalogs[j].Number
		}

		// Reverse for DESC order
		if sortParams.Order == ports.SortOrderDesc {
			return !result
		}
		return result
	}

	sort.Slice(catalogs, less)
}
