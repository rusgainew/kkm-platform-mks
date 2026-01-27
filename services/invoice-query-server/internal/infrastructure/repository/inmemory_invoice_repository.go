package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/proto-lib/entities"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	defaultMaxEntries = 100000 // 100k инвойсов
	evictionRatio     = 0.3    // Удалять 30% при превышении
)

var (
	invoiceCacheSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "invoice_cache_size",
		Help: "Current number of invoices in memory",
	})
	invoiceOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "invoice_operation_duration_seconds",
		Help:    "Duration of invoice operations",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
	invoiceEvictions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "invoice_memory_evictions_total",
		Help: "Total number of memory evictions",
	})
)

type invoiceEntry struct {
	data       *entities.Invoice
	lastAccess time.Time
}

// InMemoryInvoiceRepository implements in-memory invoice repository
// Data is populated via RabbitMQ events
type InMemoryInvoiceRepository struct {
	mu         sync.RWMutex
	invoices   map[string]*invoiceEntry             // key: invoice_number
	details    map[string][]*entities.InvoiceDetail // key: invoice_number
	maxEntries int
	logger     *zap.Logger
	lastUpdate time.Time // Track last update time for staleness checks
}

// NewInMemoryInvoiceRepository creates new in-memory repository
func NewInMemoryInvoiceRepository(logger *zap.Logger) ports.InvoiceQueryRepository {
	return &InMemoryInvoiceRepository{
		invoices:   make(map[string]*invoiceEntry),
		details:    make(map[string][]*entities.InvoiceDetail),
		maxEntries: defaultMaxEntries,
		logger:     logger,
	}
}

// GetLastUpdate returns the time of the last cache update
func (r *InMemoryInvoiceRepository) GetLastUpdate() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdate
}

// GetCacheSize returns the current number of items in cache
func (r *InMemoryInvoiceRepository) GetCacheSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.invoices)
}

// ListInvoices returns list of invoices with pagination
func (r *InMemoryInvoiceRepository) ListInvoices(ctx context.Context, page, size int32) ([]*entities.Invoice, int32, error) {
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	allInvoices := make([]*entities.Invoice, 0, len(r.invoices))
	for _, entry := range r.invoices {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		allInvoices = append(allInvoices, proto.Clone(entry.data).(*entities.Invoice))
	}

	// Calculate pagination
	total := int32(len(allInvoices))
	start := (page - 1) * size
	end := start + size

	if start >= total {
		return []*entities.Invoice{}, total, nil
	}
	if end > total {
		end = total
	}

	return allInvoices[start:end], total, nil
}

// ListInvoiceDetails returns list of invoice details with pagination
func (r *InMemoryInvoiceRepository) ListInvoiceDetails(ctx context.Context, page, size int32) ([]*entities.InvoiceDetail, int32, error) {
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect all details
	allDetails := make([]*entities.InvoiceDetail, 0)
	for _, details := range r.details {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		for _, detail := range details {
			allDetails = append(allDetails, proto.Clone(detail).(*entities.InvoiceDetail))
		}
	}

	// Calculate pagination
	total := int32(len(allDetails))
	start := (page - 1) * size
	end := start + size

	if start >= total {
		return []*entities.InvoiceDetail{}, total, nil
	}
	if end > total {
		end = total
	}

	return allDetails[start:end], total, nil
}

// ListInvoicesWithFilter returns filtered, sorted and paginated invoices
func (r *InMemoryInvoiceRepository) ListInvoicesWithFilter(ctx context.Context, filter *ports.InvoiceFilter, sort *ports.InvoiceSort, page, size int32) ([]*entities.Invoice, int32, error) {
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Apply filters
	filtered := make([]*entities.Invoice, 0)
	for _, entry := range r.invoices {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		invoice := entry.data
		if filter != nil {
			// Filter by invoice number
			if filter.InvoiceNumber != "" && invoice.InvoiceNumber != filter.InvoiceNumber {
				continue
			}
			// Filter by search text (simplified - just check invoice number and note)
			if filter.SearchText != "" {
				searchLower := strings.ToLower(filter.SearchText)
				if !strings.Contains(strings.ToLower(invoice.InvoiceNumber), searchLower) &&
					!strings.Contains(strings.ToLower(invoice.Note), searchLower) {
					continue
				}
			}
			// TODO: Add date and amount filters when needed
		}
		filtered = append(filtered, proto.Clone(invoice).(*entities.Invoice))
	}

	// Применяем сортировку
	if sort != nil && sort.Field != "" {
		sortInvoices(filtered, sort.Field, string(sort.Order))
	}

	// Calculate pagination
	total := int32(len(filtered))
	start := (page - 1) * size
	end := start + size

	if start >= total {
		return []*entities.Invoice{}, total, nil
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

// SearchInvoices performs full-text search on invoices
func (r *InMemoryInvoiceRepository) SearchInvoices(ctx context.Context, searchText string, page, size int32) ([]*entities.Invoice, int32, error) {
	if searchText == "" {
		return []*entities.Invoice{}, 0, nil
	}
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	searchLower := strings.ToLower(searchText)
	results := make([]*entities.Invoice, 0)

	for _, entry := range r.invoices {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		invoice := entry.data
		// Search in invoice_number and note
		if strings.Contains(strings.ToLower(invoice.InvoiceNumber), searchLower) ||
			strings.Contains(strings.ToLower(invoice.Note), searchLower) {
			results = append(results, proto.Clone(invoice).(*entities.Invoice))
		}
	}

	// Calculate pagination
	total := int32(len(results))
	start := (page - 1) * size
	end := start + size

	if start >= total {
		return []*entities.Invoice{}, total, nil
	}
	if end > total {
		end = total
	}

	return results[start:end], total, nil
}

// GetInvoiceByNumber returns invoice by number
func (r *InMemoryInvoiceRepository) GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (*entities.Invoice, error) {
	if invoiceNumber == "" {
		return nil, nil
	}
	r.mu.RLock()
	entry, ok := r.invoices[invoiceNumber]
	r.mu.RUnlock()

	if !ok {
		return nil, nil // Invoice not found
	}

	// Обновляем lastAccess (требует Write lock)
	r.mu.Lock()
	entry.lastAccess = time.Now()
	r.mu.Unlock()

	return proto.Clone(entry.data).(*entities.Invoice), nil
}

// GetInvoicesByDateRange returns invoices within date range
func (r *InMemoryInvoiceRepository) GetInvoicesByDateRange(ctx context.Context, dateFrom, dateTo string, page, size int32) ([]*entities.Invoice, int32, error) {
	if dateFrom == "" || dateTo == "" {
		return []*entities.Invoice{}, 0, nil
	}
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Simple string comparison for dates (works with YYYY-MM-DD format)
	results := make([]*entities.Invoice, 0)
	for _, entry := range r.invoices {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		invoice := entry.data
		if invoice.InvoiceDate != nil && invoice.InvoiceDate.Value != "" {
			invoiceDate := invoice.InvoiceDate.Value
			if invoiceDate >= dateFrom && invoiceDate <= dateTo {
				results = append(results, proto.Clone(invoice).(*entities.Invoice))
			}
		}
	}

	// Calculate pagination
	total := int32(len(results))
	start := (page - 1) * size
	end := start + size

	if start >= total {
		return []*entities.Invoice{}, total, nil
	}
	if end > total {
		end = total
	}

	return results[start:end], total, nil
}

// evictOldEntries удаляет старые записи при превышении лимита (только для внутреннего использования, вызывается с активным Lock)
func (r *InMemoryInvoiceRepository) evictOldEntries() {
	toEvict := int(float64(len(r.invoices)) * evictionRatio)
	if toEvict == 0 {
		toEvict = 1
	}

	// Собираем все ключи с lastAccess
	type entryInfo struct {
		key        string
		lastAccess time.Time
	}
	entries := make([]entryInfo, 0, len(r.invoices))
	for key, entry := range r.invoices {
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
		delete(r.invoices, entry.key)
		evicted++
	}

	// Обновляем метрики
	invoiceEvictions.Add(float64(evicted))
	invoiceCacheSize.Set(float64(len(r.invoices)))

	r.logger.Info("Memory eviction performed",
		zap.Int("evicted_count", evicted),
		zap.Int("remaining_count", len(r.invoices)),
		zap.Int("max_entries", r.maxEntries))
}

// UpsertInvoice adds or updates invoice (used by event handlers)
func (r *InMemoryInvoiceRepository) UpsertInvoice(invoice *entities.Invoice) {
	start := time.Now()
	defer func() {
		invoiceOperationDuration.WithLabelValues("upsert").Observe(time.Since(start).Seconds())
	}()

	if invoice == nil || invoice.InvoiceNumber == "" {
		r.logger.Warn("Attempted to upsert invoice with empty invoice_number")
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	action := "update"
	if _, exists := r.invoices[invoice.InvoiceNumber]; !exists {
		action = "insert"
	}

	// Проверяем лимит памяти
	if len(r.invoices) >= r.maxEntries {
		r.evictOldEntries()
	}

	// Сохраняем инвойс
	r.invoices[invoice.InvoiceNumber] = &invoiceEntry{
		data:       proto.Clone(invoice).(*entities.Invoice),
		lastAccess: time.Now(),
	}

	r.lastUpdate = time.Now()

	// Обновляем метрики
	invoiceCacheSize.Set(float64(len(r.invoices)))

	r.logger.Debug("Invoice upserted",
		zap.String("action", action),
		zap.String("invoice_number", invoice.InvoiceNumber),
		zap.Int("cache_size", len(r.invoices)))
}

// DeleteInvoice removes invoice (used by event handlers)
func (r *InMemoryInvoiceRepository) DeleteInvoice(invoiceNumber string) {
	start := time.Now()
	defer func() {
		invoiceOperationDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	if invoiceNumber == "" {
		r.logger.Warn("Attempted to delete invoice with empty invoice_number")
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.invoices, invoiceNumber)
	delete(r.details, invoiceNumber)

	// Обновляем метрики
	invoiceCacheSize.Set(float64(len(r.invoices)))

	r.logger.Debug("Invoice deleted",
		zap.String("invoice_number", invoiceNumber),
		zap.Int("cache_size", len(r.invoices)))
}

// sortInvoices сортирует инвойсы по указанному полю и порядку
func sortInvoices(invoices []*entities.Invoice, field string, order string) {
	if field == "" {
		return
	}

	less := func(i, j int) bool {
		var result bool
		switch field {
		case "invoice_number":
			result = strings.ToLower(invoices[i].InvoiceNumber) < strings.ToLower(invoices[j].InvoiceNumber)
		case "invoice_date":
			// Сравниваем даты как строки (работает для формата YYYY-MM-DD)
			dateI := ""
			dateJ := ""
			if invoices[i].InvoiceDate != nil && invoices[i].InvoiceDate.Value != "" {
				dateI = invoices[i].InvoiceDate.Value
			}
			if invoices[j].InvoiceDate != nil && invoices[j].InvoiceDate.Value != "" {
				dateJ = invoices[j].InvoiceDate.Value
			}
			result = dateI < dateJ
		default:
			result = strings.ToLower(invoices[i].InvoiceNumber) < strings.ToLower(invoices[j].InvoiceNumber)
		}

		if order == "DESC" {
			return !result
		}
		return result
	}

	sort.Slice(invoices, less)
}
