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
	defaultMaxEntries = 100000 // 100k пользователей
	evictionRatio     = 0.3    // Удалять 30% при превышении
)

var (
	userCacheSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "user_cache_size",
		Help: "Current number of users in memory",
	})
	userOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "user_operation_duration_seconds",
		Help:    "Duration of user operations",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
	userEvictions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_memory_evictions_total",
		Help: "Total number of memory evictions",
	})
)

type userEntry struct {
	data       *pb.UserReadModel
	lastAccess time.Time
}

// InMemoryUserRepository хранит данные пользователей в памяти
type InMemoryUserRepository struct {
	mu         sync.RWMutex
	users      map[string]*userEntry // key: user_id
	logger     *zap.Logger
	maxEntries int
	lastUpdate time.Time // Track last update time for staleness checks
}

// NewInMemoryUserRepository создает новый in-memory репозиторий
func NewInMemoryUserRepository(logger *zap.Logger) *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:      make(map[string]*userEntry),
		logger:     logger,
		maxEntries: defaultMaxEntries,
	}
}

// GetLastUpdate returns the time of the last cache update
func (r *InMemoryUserRepository) GetLastUpdate() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdate
}

// GetCacheSize returns the current number of items in cache
func (r *InMemoryUserRepository) GetCacheSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users)
}

// GetUser получает пользователя по ID
func (r *InMemoryUserRepository) GetUser(ctx context.Context, userID string) (*pb.UserReadModel, error) {
	start := time.Now()
	defer func() {
		userOperationDuration.WithLabelValues("get").Observe(time.Since(start).Seconds())
	}()

	if userID == "" {
		r.logger.Warn("GetUser called with empty userID")
		return nil, nil
	}

	r.mu.RLock()
	entry, ok := r.users[userID]
	r.mu.RUnlock()

	if !ok {
		return nil, nil // Пользователь не найден
	}

	// Обновляем lastAccess (требует Write lock)
	r.mu.Lock()
	entry.lastAccess = time.Now()
	r.mu.Unlock()

	// Клонируем для предотвращения race conditions
	return proto.Clone(entry.data).(*pb.UserReadModel), nil
}

// ListUsers возвращает список пользователей с фильтрацией и пагинацией
func (r *InMemoryUserRepository) ListUsers(ctx context.Context, offset, limit int32, status, role string) ([]*pb.UserReadModel, int64, error) {
	start := time.Now()
	defer func() {
		userOperationDuration.WithLabelValues("list").Observe(time.Since(start).Seconds())
	}()

	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Собираем всех пользователей в slice
	allUsers := make([]*pb.UserReadModel, 0, len(r.users))
	for _, entry := range r.users {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}

		user := entry.data
		// Применяем фильтры
		if status != "" && user.Status != status {
			continue
		}
		if role != "" && user.Role != role {
			continue
		}
		// Клонируем для предотвращения race conditions
		allUsers = append(allUsers, proto.Clone(user).(*pb.UserReadModel))
	}

	totalCount := int64(len(allUsers))

	// Применяем пагинацию
	if offset >= int32(totalCount) {
		return []*pb.UserReadModel{}, totalCount, nil
	}

	end := offset + limit
	if end > int32(totalCount) {
		end = int32(totalCount)
	}

	return allUsers[offset:end], totalCount, nil
}

// SearchUsers ищет пользователей по запросу
func (r *InMemoryUserRepository) SearchUsers(ctx context.Context, query string, offset, limit int32, status string) ([]*pb.UserReadModel, int64, error) {
	start := time.Now()
	defer func() {
		userOperationDuration.WithLabelValues("search").Observe(time.Since(start).Seconds())
	}()

	if query == "" {
		r.logger.Warn("SearchUsers called with empty query")
		return []*pb.UserReadModel{}, 0, nil
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
	var filtered []*pb.UserReadModel

	for _, entry := range r.users {
		// Проверяем отмену контекста
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}

		user := entry.data
		// Применяем фильтр по статусу
		if status != "" && user.Status != status {
			continue
		}

		// Поиск по email, имени, фамилии, телефону
		if strings.Contains(strings.ToLower(user.Email), query) ||
			strings.Contains(strings.ToLower(user.FirstName), query) ||
			strings.Contains(strings.ToLower(user.LastName), query) ||
			strings.Contains(strings.ToLower(user.Phone), query) {
			// Клонируем для предотвращения race conditions
			filtered = append(filtered, proto.Clone(user).(*pb.UserReadModel))
		}
	}

	totalCount := int64(len(filtered))

	// Применяем пагинацию
	if offset >= int32(totalCount) {
		return []*pb.UserReadModel{}, totalCount, nil
	}

	end := offset + limit
	if end > int32(totalCount) {
		end = int32(totalCount)
	}

	return filtered[offset:end], totalCount, nil
}

// UpsertUser добавляет или обновляет пользователя (используется в RabbitMQ consumer)
func (r *InMemoryUserRepository) UpsertUser(ctx context.Context, user *pb.UserReadModel) error {
	start := time.Now()
	defer func() {
		userOperationDuration.WithLabelValues("upsert").Observe(time.Since(start).Seconds())
	}()

	if user == nil || user.Id == "" {
		r.logger.Warn("Attempted to upsert user with empty ID")
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	action := "created"
	if _, exists := r.users[user.Id]; exists {
		r.lastUpdate = time.Now()
		action = "updated"
	}

	// Проверяем лимит памяти
	if len(r.users) >= r.maxEntries {
		r.evictOldEntries()
	}

	// Клонируем для предотвращения изменения извне
	r.users[user.Id] = &userEntry{
		data:       proto.Clone(user).(*pb.UserReadModel),
		lastAccess: time.Now(),
	}

	userCacheSize.Set(float64(len(r.users)))

	r.logger.Debug("User upserted",
		zap.String("action", action),
		zap.String("user_id", user.Id),
		zap.Int("cache_size", len(r.users)))

	return nil
}

// DeleteUser удаляет пользователя из памяти (используется в RabbitMQ consumer)
func (r *InMemoryUserRepository) DeleteUser(ctx context.Context, userID string) error {
	start := time.Now()
	defer func() {
		userOperationDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	if userID == "" {
		r.logger.Warn("DeleteUser called with empty userID")
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, userID)
	userCacheSize.Set(float64(len(r.users)))

	r.logger.Debug("User deleted",
		zap.String("user_id", userID),
		zap.Int("cache_size", len(r.users)))

	return nil
}

// GetAllUsers возвращает всех пользователей (для отладки)
func (r *InMemoryUserRepository) GetAllUsers() []*pb.UserReadModel {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*pb.UserReadModel, 0, len(r.users))
	for _, entry := range r.users {
		// Клонируем для предотвращения race conditions
		users = append(users, proto.Clone(entry.data).(*pb.UserReadModel))
	}
	return users
}

// Count возвращает количество пользователей в памяти
func (r *InMemoryUserRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users)
}

// evictOldEntries удаляет старые записи при превышении лимита (только для внутреннего использования, вызывается с активным Lock)
func (r *InMemoryUserRepository) evictOldEntries() {
	toEvict := int(float64(len(r.users)) * evictionRatio)
	if toEvict == 0 {
		toEvict = 1
	}

	// Собираем все ключи с lastAccess
	type entryInfo struct {
		key        string
		lastAccess time.Time
	}
	entries := make([]entryInfo, 0, len(r.users))
	for key, entry := range r.users {
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
		delete(r.users, entry.key)
		evicted++
	}

	userEvictions.Add(float64(evicted))
	userCacheSize.Set(float64(len(r.users)))

	r.logger.Info("Memory eviction performed",
		zap.Int("evicted_count", evicted),
		zap.Int("remaining_count", len(r.users)),
		zap.Int("max_entries", r.maxEntries))
}

// sortUsers сортирует пользователей по указанному полю и порядку
func sortUsers(users []*pb.UserReadModel, field string, order string) {
	if field == "" {
		return
	}

	less := func(i, j int) bool {
		var result bool
		switch field {
		case "email":
			result = strings.ToLower(users[i].Email) < strings.ToLower(users[j].Email)
		case "first_name":
			result = strings.ToLower(users[i].FirstName) < strings.ToLower(users[j].FirstName)
		case "last_name":
			result = strings.ToLower(users[i].LastName) < strings.ToLower(users[j].LastName)
		case "created_at":
			result = users[i].CreatedAt < users[j].CreatedAt
		default:
			result = users[i].CreatedAt < users[j].CreatedAt
		}

		if order == "DESC" {
			return !result
		}
		return result
	}

	sort.Slice(users, less)
}
