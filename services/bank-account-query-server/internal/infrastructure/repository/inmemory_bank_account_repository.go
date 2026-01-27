package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rusgainew/kkm-project-mks/bank-account-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
	"github.com/rusgainew/kkm-project-mks/services/pkg/conversion"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	defaultMaxEntries = 100000 // 100k банковских счетов
	evictionRatio     = 0.3    // Удалять 30% при превышении
)

var (
	bankAccountCacheSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bank_account_cache_size",
		Help: "Current number of bank accounts in memory",
	})
	bankAccountOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "bank_account_operation_duration_seconds",
		Help:    "Duration of bank account operations",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
	bankAccountEvictions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bank_account_memory_evictions_total",
		Help: "Total number of memory evictions",
	})
)

type bankAccountEntry struct {
	data       *dictionaries.BankAccount
	lastAccess time.Time
}

// InMemoryBankAccountRepository implements in-memory catalog repository
// Data is populated via RabbitMQ events
type InMemoryBankAccountRepository struct {
	mu           sync.RWMutex
	bankAccounts map[string]*bankAccountEntry // key: code
	maxEntries   int
	logger       *zap.Logger
	lastUpdate   time.Time // Track last update time for staleness checks
}

// NewInMemoryBankAccountRepository creates new in-memory repository
func NewInMemoryBankAccountRepository(logger *zap.Logger) ports.BankAccountQueryRepository {
	return &InMemoryBankAccountRepository{
		bankAccounts: make(map[string]*bankAccountEntry),
		maxEntries:   defaultMaxEntries,
		logger:       logger,
	}
}

// GetLastUpdate returns the time of the last cache update
func (r *InMemoryBankAccountRepository) GetLastUpdate() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdate
}

// GetCacheSize returns the current number of items in cache
func (r *InMemoryBankAccountRepository) GetCacheSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.bankAccounts)
}

// ListBankAccounts returns paginated list of bankAccounts
func (r *InMemoryBankAccountRepository) ListBankAccounts(ctx context.Context, page, size int32) ([]*dictionaries.BankAccount, int32, error) {
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
	allBankAccounts := make([]*dictionaries.BankAccount, 0, len(r.bankAccounts))
	for _, entry := range r.bankAccounts {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		allBankAccounts = append(allBankAccounts, proto.Clone(entry.data).(*dictionaries.BankAccount))
	}

	totalCount := conversion.SafeIntToInt32WithDefault(len(allBankAccounts), 0)
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.BankAccount{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return allBankAccounts[offset:end], totalCount, nil
}

// ListBankAccountsWithFilter returns filtered and sorted bankAccounts
func (r *InMemoryBankAccountRepository) ListBankAccountsWithFilter(ctx context.Context, filter *ports.BankAccountFilter, sort *ports.BankAccountSort, page, size int32) ([]*dictionaries.BankAccount, int32, error) {
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

	// Filter bankAccounts
	filtered := make([]*dictionaries.BankAccount, 0)
	for _, entry := range r.bankAccounts {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		if !matchesFilter(entry.data, filter) {
			continue
		}
		filtered = append(filtered, proto.Clone(entry.data).(*dictionaries.BankAccount))
	}

	totalCount := conversion.SafeIntToInt32WithDefault(len(filtered), 0)
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.BankAccount{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// SearchBankAccounts performs full-text search
func (r *InMemoryBankAccountRepository) SearchBankAccounts(ctx context.Context, searchText string, page, size int32) ([]*dictionaries.BankAccount, int32, error) {
	if searchText == "" {
		return []*dictionaries.BankAccount{}, 0, nil
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
	filtered := make([]*dictionaries.BankAccount, 0)

	for _, entry := range r.bankAccounts {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		bankAccount := entry.data
		if strings.Contains(strings.ToLower(bankAccount.AccountName), searchLower) ||
			strings.Contains(strings.ToLower(bankAccount.BankAccount), searchLower) {
			filtered = append(filtered, proto.Clone(bankAccount).(*dictionaries.BankAccount))
		}
	}

	totalCount := conversion.SafeIntToInt32WithDefault(len(filtered), 0)
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.BankAccount{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// GetBankAccountByNumber returns bank account by number
func (r *InMemoryBankAccountRepository) GetBankAccountByNumber(ctx context.Context, number string) (*dictionaries.BankAccount, error) {
	if number == "" {
		return nil, nil
	}
	r.mu.RLock()
	entry, exists := r.bankAccounts[number]
	r.mu.RUnlock()

	if !exists {
		return nil, nil
	}

	// Обновляем lastAccess (требует Write lock)
	r.mu.Lock()
	entry.lastAccess = time.Now()
	r.mu.Unlock()

	return proto.Clone(entry.data).(*dictionaries.BankAccount), nil
}

// GetActiveBankAccounts returns only active bank accounts
func (r *InMemoryBankAccountRepository) GetActiveBankAccounts(ctx context.Context, page, size int32) ([]*dictionaries.BankAccount, int32, error) {
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

	filtered := make([]*dictionaries.BankAccount, 0)
	for _, entry := range r.bankAccounts {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}
		if entry.data.IsActive { // Only active accounts
			filtered = append(filtered, proto.Clone(entry.data).(*dictionaries.BankAccount))
		}
	}

	totalCount := conversion.SafeIntToInt32WithDefault(len(filtered), 0)
	offset := (page - 1) * size

	if offset >= totalCount {
		return []*dictionaries.BankAccount{}, totalCount, nil
	}

	end := offset + size
	if end > totalCount {
		end = totalCount
	}

	return filtered[offset:end], totalCount, nil
}

// UpsertBankAccount adds or updates bank account (used by event handlers)
func (r *InMemoryBankAccountRepository) UpsertBankAccount(bankAccount *dictionaries.BankAccount) {
	start := time.Now()
	defer func() {
		bankAccountOperationDuration.WithLabelValues("upsert").Observe(time.Since(start).Seconds())
	}()

	if bankAccount == nil || bankAccount.BankAccount == "" {
		r.logger.Warn("Attempted to upsert bank account with empty bank_account")
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	action := "update"
	if _, exists := r.bankAccounts[bankAccount.BankAccount]; !exists {
		action = "insert"
	}

	// Проверяем лимит памяти
	if len(r.bankAccounts) >= r.maxEntries {
		r.evictOldEntries()
	}

	r.bankAccounts[bankAccount.BankAccount] = &bankAccountEntry{
		data:       proto.Clone(bankAccount).(*dictionaries.BankAccount),
		lastAccess: time.Now(),
	}

	r.lastUpdate = time.Now()

	// Обновляем метрики
	bankAccountCacheSize.Set(float64(len(r.bankAccounts)))

	r.logger.Debug("Bank account upserted",
		zap.String("action", action),
		zap.String("bank_account", bankAccount.BankAccount),
		zap.Int("cache_size", len(r.bankAccounts)))
}

// DeleteBankAccount removes bank account (used by event handlers)
func (r *InMemoryBankAccountRepository) DeleteBankAccount(bankAccountNumber string) {
	start := time.Now()
	defer func() {
		bankAccountOperationDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	if bankAccountNumber == "" {
		r.logger.Warn("Attempted to delete bank account with empty bank_account_number")
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.bankAccounts, bankAccountNumber)

	// Обновляем метрики
	bankAccountCacheSize.Set(float64(len(r.bankAccounts)))

	r.logger.Debug("Bank account deleted",
		zap.String("bank_account_number", bankAccountNumber),
		zap.Int("cache_size", len(r.bankAccounts)))
}

// matchesFilter checks if bank account matches filter criteria
func matchesFilter(bankAccount *dictionaries.BankAccount, filter *ports.BankAccountFilter) bool {
	if filter == nil {
		return true
	}

	if filter.AccountName != "" && !strings.Contains(strings.ToLower(bankAccount.AccountName), strings.ToLower(filter.AccountName)) {
		return false
	}

	if filter.BankAccount != "" && bankAccount.BankAccount != filter.BankAccount {
		return false
	}

	if filter.ContractorTin != "" && bankAccount.ContractorTin != filter.ContractorTin {
		return false
	}

	if filter.IsActive != nil && bankAccount.IsActive != *filter.IsActive {
		return false
	}

	if filter.SearchText != "" {
		searchLower := strings.ToLower(filter.SearchText)
		if !strings.Contains(strings.ToLower(bankAccount.AccountName), searchLower) &&
			!strings.Contains(strings.ToLower(bankAccount.BankAccount), searchLower) {
			return false
		}
	}

	return true
}

// sortBankAccounts сортирует банковские счета по указанному полю и порядку
func sortBankAccounts(accounts []*dictionaries.BankAccount, field string, order string) {
	if field == "" {
		return
	}

	less := func(i, j int) bool {
		var result bool
		switch field {
		case "account_number":
			result = strings.ToLower(accounts[i].BankAccount) < strings.ToLower(accounts[j].BankAccount)
		case "bank_name":
			result = strings.ToLower(accounts[i].AccountName) < strings.ToLower(accounts[j].AccountName)
		default:
			result = strings.ToLower(accounts[i].BankAccount) < strings.ToLower(accounts[j].BankAccount)
		}

		if order == "DESC" {
			return !result
		}
		return result
	}

	sort.Slice(accounts, less)
}

// evictOldEntries удаляет старые записи при превышении лимита (только для внутреннего использования, вызывается с активным Lock)
func (r *InMemoryBankAccountRepository) evictOldEntries() {
	toEvict := int(float64(len(r.bankAccounts)) * evictionRatio)
	if toEvict == 0 {
		toEvict = 1
	}

	// Собираем все ключи с lastAccess
	type entryInfo struct {
		key        string
		lastAccess time.Time
	}
	entries := make([]entryInfo, 0, len(r.bankAccounts))
	for key, entry := range r.bankAccounts {
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
		delete(r.bankAccounts, entry.key)
		evicted++
	}

	// Обновляем метрики
	bankAccountEvictions.Add(float64(evicted))
	bankAccountCacheSize.Set(float64(len(r.bankAccounts)))

	r.logger.Info("Memory eviction performed",
		zap.Int("evicted_count", evicted),
		zap.Int("remaining_count", len(r.bankAccounts)),
		zap.Int("max_entries", r.maxEntries))
}
