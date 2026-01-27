package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rusgainew/kkm-project-mks/foreign-company-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	defaultMaxEntries = 100000 // 100k иностранных компаний
	evictionRatio     = 0.3    // Удалять 30% при превышении
)

var (
	foreignCompanyCacheSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "foreign_company_cache_size",
		Help: "Current number of foreign companies in memory",
	})
	foreignCompanyOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "foreign_company_operation_duration_seconds",
		Help:    "Duration of foreign company operations",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
	foreignCompanyEvictions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "foreign_company_memory_evictions_total",
		Help: "Total number of memory evictions",
	})
)

type foreignCompanyEntry struct {
	data       *dictionaries.ForeignCompany
	lastAccess time.Time
}

// InMemoryForeignCompanyRepository implements in-memory catalog repository
// Data is populated via RabbitMQ events
type InMemoryForeignCompanyRepository struct {
	mu               sync.RWMutex
	foreignCompanies map[string]*foreignCompanyEntry // key: code
	maxEntries       int
	logger           *zap.Logger
	lastUpdate       time.Time // Track last update time for staleness checks
}

// NewInMemoryForeignCompanyRepository creates new in-memory repository
func NewInMemoryForeignCompanyRepository(logger *zap.Logger) ports.ForeignCompanyQueryRepository {
	return &InMemoryForeignCompanyRepository{
		foreignCompanies: make(map[string]*foreignCompanyEntry),
		maxEntries:       defaultMaxEntries,
		logger:           logger,
	}
}

// GetLastUpdate returns the time of the last cache update
func (r *InMemoryForeignCompanyRepository) GetLastUpdate() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdate
}

// GetCacheSize returns the current number of items in cache
func (r *InMemoryForeignCompanyRepository) GetCacheSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.foreignCompanies)
}

// ListForeignCompanies returns paginated list of foreignCompanies
func (r *InMemoryForeignCompanyRepository) ListForeignCompanies(ctx context.Context, page, size int32) ([]*dictionaries.ForeignCompany, int32, error) {
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
	allForeignCompanies := make([]*dictionaries.ForeignCompany, 0, len(r.foreignCompanies))
	for _, entry := range r.foreignCompanies {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		allForeignCompanies = append(allForeignCompanies, proto.Clone(entry.data).(*dictionaries.ForeignCompany))
	}

	totalCount := int32(len(allForeignCompanies))
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.ForeignCompany{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return allForeignCompanies[offset:end], totalCount, nil
}

// ListForeignCompaniesWithFilter returns filtered and sorted foreignCompanies
func (r *InMemoryForeignCompanyRepository) ListForeignCompaniesWithFilter(ctx context.Context, filter *ports.ForeignCompanyFilter, sort *ports.ForeignCompanySort, page, size int32) ([]*dictionaries.ForeignCompany, int32, error) {
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

	// Filter foreignCompanies
	filtered := make([]*dictionaries.ForeignCompany, 0)
	for _, entry := range r.foreignCompanies {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		if !matchesFilter(entry.data, filter) {
			continue
		}
		filtered = append(filtered, proto.Clone(entry.data).(*dictionaries.ForeignCompany))
	}

	totalCount := int32(len(filtered))
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.ForeignCompany{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// SearchForeignCompanies performs full-text search
func (r *InMemoryForeignCompanyRepository) SearchForeignCompanies(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.ForeignCompany, int32, error) {
	if searchText == "" {
		return []*dictionaries.ForeignCompany{}, 0, nil
	}
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

	searchLower := strings.ToLower(searchText)
	filtered := make([]*dictionaries.ForeignCompany, 0)

	for _, entry := range r.foreignCompanies {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		foreignCompany := entry.data
		if strings.Contains(strings.ToLower(foreignCompany.FullName), searchLower) ||
			strings.Contains(strings.ToLower(foreignCompany.Pin), searchLower) {
			filtered = append(filtered, proto.Clone(foreignCompany).(*dictionaries.ForeignCompany))
		}
	}

	totalCount := int32(len(filtered))
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.ForeignCompany{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// GetForeignCompanyByPIN returns foreign company by PIN
func (r *InMemoryForeignCompanyRepository) GetForeignCompanyByPIN(ctx context.Context, pin string) (*dictionaries.ForeignCompany, error) {
	if pin == "" {
		return nil, nil
	}
	r.mu.RLock()
	entry, exists := r.foreignCompanies[pin]
	r.mu.RUnlock()

	if !exists {
		return nil, nil
	}

	// Обновляем lastAccess (требует Write lock)
	r.mu.Lock()
	entry.lastAccess = time.Now()
	r.mu.Unlock()

	return proto.Clone(entry.data).(*dictionaries.ForeignCompany), nil
}

// GetForeignCompaniesByCountry returns foreign companies by country code
// Note: ForeignCompany proto doesn't have country field yet, returning empty for now
func (r *InMemoryForeignCompanyRepository) GetForeignCompaniesByCountry(ctx context.Context, countryCode string, page, size int32) ([]*dictionaries.ForeignCompany, int32, error) {
	// TODO: Implement when Country field is added to ForeignCompany proto
	return []*dictionaries.ForeignCompany{}, 0, nil
}

// UpsertForeignCompany adds or updates foreign company (used by event handlers)
func (r *InMemoryForeignCompanyRepository) UpsertForeignCompany(foreignCompany *dictionaries.ForeignCompany) {
	start := time.Now()
	defer func() {
		foreignCompanyOperationDuration.WithLabelValues("upsert").Observe(time.Since(start).Seconds())
	}()

	if foreignCompany == nil || foreignCompany.Pin == "" {
		r.logger.Warn("Attempted to upsert foreign company with empty pin")
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	action := "update"
	if _, exists := r.foreignCompanies[foreignCompany.Pin]; !exists {
		action = "insert"
	}

	// Проверяем лимит памяти
	if len(r.foreignCompanies) >= r.maxEntries {
		r.evictOldEntries()
	}

	r.foreignCompanies[foreignCompany.Pin] = &foreignCompanyEntry{
		data:       proto.Clone(foreignCompany).(*dictionaries.ForeignCompany),
		lastAccess: time.Now(),
	}

	// Обновляем время последнего обновления кэша
	r.lastUpdate = time.Now()

	// Обновляем метрики
	foreignCompanyCacheSize.Set(float64(len(r.foreignCompanies)))

	r.logger.Debug("Foreign company upserted",
		zap.String("action", action),
		zap.String("pin", foreignCompany.Pin),
		zap.Int("cache_size", len(r.foreignCompanies)))
}

// DeleteForeignCompany removes foreign company (used by event handlers)
func (r *InMemoryForeignCompanyRepository) DeleteForeignCompany(pin string) {
	start := time.Now()
	defer func() {
		foreignCompanyOperationDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	if pin == "" {
		r.logger.Warn("Attempted to delete foreign company with empty pin")
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.foreignCompanies, pin)

	// Обновляем метрики
	foreignCompanyCacheSize.Set(float64(len(r.foreignCompanies)))

	r.logger.Debug("Foreign company deleted",
		zap.String("pin", pin),
		zap.Int("cache_size", len(r.foreignCompanies)))
}

// matchesFilter checks if foreign company matches filter criteria
func matchesFilter(foreignCompany *dictionaries.ForeignCompany, filter *ports.ForeignCompanyFilter) bool {
	if filter == nil {
		return true
	}

	if filter.FullName != "" && !strings.Contains(strings.ToLower(foreignCompany.FullName), strings.ToLower(filter.FullName)) {
		return false
	}

	if filter.PIN != "" && foreignCompany.Pin != filter.PIN {
		return false
	}

	// TODO: Implement CountryCode filter when Country field is added to proto
	// if filter.CountryCode != "" { ... }

	if filter.SearchText != "" {
		searchLower := strings.ToLower(filter.SearchText)
		if !strings.Contains(strings.ToLower(foreignCompany.FullName), searchLower) &&
			!strings.Contains(strings.ToLower(foreignCompany.Pin), searchLower) {
			return false
		}
	}

	return true
}

// sortForeignCompanies сортирует иностранные компании по указанному полю и порядку
func sortForeignCompanies(companies []*dictionaries.ForeignCompany, field string, order string) {
	if field == "" {
		return
	}

	less := func(i, j int) bool {
		var result bool
		switch field {
		case "name":
			result = strings.ToLower(companies[i].FullName) < strings.ToLower(companies[j].FullName)
		case "code":
			result = strings.ToLower(companies[i].Pin) < strings.ToLower(companies[j].Pin)
		default:
			result = strings.ToLower(companies[i].FullName) < strings.ToLower(companies[j].FullName)
		}

		if order == "DESC" {
			return !result
		}
		return result
	}

	sort.Slice(companies, less)
}

// evictOldEntries удаляет старые записи при превышении лимита (только для внутреннего использования, вызывается с активным Lock)
func (r *InMemoryForeignCompanyRepository) evictOldEntries() {
	toEvict := int(float64(len(r.foreignCompanies)) * evictionRatio)
	if toEvict == 0 {
		toEvict = 1
	}

	// Собираем все ключи с lastAccess
	type entryInfo struct {
		key        string
		lastAccess time.Time
	}
	entries := make([]entryInfo, 0, len(r.foreignCompanies))
	for key, entry := range r.foreignCompanies {
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
		delete(r.foreignCompanies, entry.key)
		evicted++
	}

	// Обновляем метрики
	foreignCompanyEvictions.Add(float64(evicted))
	foreignCompanyCacheSize.Set(float64(len(r.foreignCompanies)))

	r.logger.Info("Memory eviction performed",
		zap.Int("evicted_count", evicted),
		zap.Int("remaining_count", len(r.foreignCompanies)),
		zap.Int("max_entries", r.maxEntries))
}
