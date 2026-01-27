# 🔍 Полный анализ кода Go-сервисов

**Дата:** 27 января 2026 г.  
**Объём кода:** 336 файлов Go (~47,224 строк), 49 тестовых файлов  
**Архитектура:** Микросервисная, CQRS, Event-driven, gRPC

---

## 📊 Общая статистика

### Сервисы (15 микросервисов)

**Command Services (Write-side):**

- ✅ `api-gateway` - HTTP/gRPC шлюз, JWT auth, кэширование, rate limiting
- ✅ `company-server` - Управление организациями и сотрудниками
- ✅ `bank-account-server` - Управление банковскими счетами
- ✅ `catalog-server` - Управление каталогами товаров/услуг
- ✅ `foreign-company-server` - Управление иностранными компаниями
- ✅ `invoice-server` - Управление счетами-фактурами
- ✅ `document-server` - Управление документами
- ⚠️ `user-server` - Управление пользователями

**Query Services (Read-side):**

- ✅ `catalog-query-server` - Redis + RabbitMQ + In-memory (мигрирован с PostgreSQL)
- ✅ `invoice-query-server` - Redis + RabbitMQ + In-memory (мигрирован с PostgreSQL)
- ✅ `bank-account-query-server` - Redis + RabbitMQ + In-memory (мигрирован с PostgreSQL)
- ✅ `foreign-company-query-server` - Redis + RabbitMQ + In-memory (мигрирован с PostgreSQL)
- ✅ `user-query-server` - Redis + RabbitMQ + In-memory (мигрирован с PostgreSQL)
- ✅ `document-query-server` - Redis + RabbitMQ + In-memory (мигрирован с PostgreSQL)

**Инфраструктура:**

- ✅ `nginx-proxy` - Reverse proxy, SSL termination

---

## ⭐ Сильные стороны

### 1. **Архитектурные решения** ✅

#### ✅ CQRS Pattern реализован корректно

```go
// Разделение команд и запросов
services/
  catalog-server/         # Команды (write)
  catalog-query-server/   # Запросы (read)
```

#### ✅ Event-Driven Architecture с RabbitMQ

- Все write-сервисы публикуют события
- Query-сервисы потребляют события для синхронизации
- Dead Letter Queue для обработки ошибок
- Механизмы retry и идемпотентности

```go
// user-query-server: Обработка событий
func (h *UserEventHandler) HandleUserRegistered(ctx context.Context, event *UserRegisteredEvent) error {
    // Идемпотентная вставка в read-модель
}
```

#### ✅ Observability (Метрики, Логи, Трассировка)

**Prometheus Metrics:**

- Все сервисы экспортируют метрики на `/metrics`
- Стандартизированные метрики: `requests_total`, `request_duration_ms`, `errors_total`
- Бизнес-метрики: `documents_created`, `invoices_sent`, `bank_accounts_active`

**Structured Logging (zap):**

```go
logger.Info("Document created successfully",
    zap.String("document_id", docID),
    zap.String("organization_id", req.OrganizationId))
```

**OpenTelemetry Tracing:**

- Интеграция с Jaeger
- Распределенная трассировка через микросервисы

### 2. **Качество кода** ✅

#### ✅ Правильное использование Context

```go
// Timeout в main.go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err = db.PingContext(ctx); err != nil { ... }

// Graceful shutdown
shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
defer cancel()
grpcServer.GracefulStop()
```

#### ✅ Resource Cleanup с defer

```go
defer logger.Sync()
defer db.Close()
defer publisher.Close()
defer redisCache.Close()
```

#### ✅ Proper Error Wrapping

```go
return nil, fmt.Errorf("failed to query catalogs: %w", err)
```

#### ✅ Sync primitives правильно используются

```go
type InMemoryCatalogRepository struct {
    mu       sync.RWMutex
    catalogs map[string]*dictionaries.Catalog
}

func (r *InMemoryCatalogRepository) ListCatalogs(ctx context.Context, page, size int32) ([]*dictionaries.Catalog, int32, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    // ...
}
```

### 3. **Тестирование** ✅

#### ✅ Unit тесты с высоким покрытием (примеры)

- `company-server`: Metrics, Validation, Service - 9+ тестов
- `user-query-server`: Handler тесты - 8+ тестов
- `invoice-server`: Validation тесты

#### ✅ Integration тесты для query-servers

```go
// catalog-query-server/tests/integration/filtering_test.go
func TestListCatalogsWithFilter(t *testing.T)
func TestCatalogCacheHit(t *testing.T)
func TestCatalogPaginationEdgeCases(t *testing.T)
```

### 4. **Security** ✅

#### ✅ JWT Authentication

```go
// api-gateway: JWT middleware
func JWTAuth() gin.HandlerFunc {
    token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, ...)
    // Валидация подписи
}
```

#### ✅ Authorization Middleware

```go
// company-server: Role-based access
func (s *Service) CheckAccess(ctx context.Context, userID, orgID string, action string) error {
    // Проверка прав доступа
}
```

#### ✅ Input Validation

```go
// company-server: Валидация запросов
func ValidateCreateOrganizationRequest(req *pb.CreateOrganizationRequest) []string {
    if req.Name == "" { errors = append(errors, "name is required") }
    if req.Tin == "" { errors = append(errors, "tin is required") }
}
```

---

## ⚠️ Критические проблемы

### 1. **Pagination Inconsistency** 🔴 КРИТИЧНО

**Проблема:** Разные формулы пагинации в разных сервисах.

**Примеры:**

```go
// ❌ invoice-query-server (НЕПРАВИЛЬНО)
offset := page * size  // Page 0 -> offset 0, Page 1 -> offset 10

// ✅ catalog-query-server (ПРАВИЛЬНО)
offset := (page - 1) * size  // Page 1 -> offset 0, Page 2 -> offset 10
```

**Решение:**

```go
// Стандартизировать на всех сервисах
func normalizePagination(page, size int32) (int32, int32) {
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 10
    }
    offset := (page - 1) * size
    return offset, size
}
```

**Затронутые файлы:**

- `services/invoice-query-server/internal/infrastructure/repository/inmemory_invoice_repository.go:97`
- `services/catalog-query-server/internal/infrastructure/repository/inmemory_catalog_repository.go:44`
- `services/bank-account-query-server/internal/infrastructure/repository/inmemory_bank_account_repository.go:44`

**Влияние:** Клиенты получают неожиданные результаты, API некорректный.

---

### 2. **Race Conditions в In-Memory Repositories** 🔴 КРИТИЧНО

**Проблема:** Прямое присваивание указателей без копирования.

```go
// ❌ ОПАСНО: race condition
func (r *InMemoryCatalogRepository) UpsertCatalog(ctx context.Context, catalog *dictionaries.Catalog) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Если внешний код изменит catalog после вызова, данные повредятся
    r.catalogs[catalog.Code] = catalog
    return nil
}
```

**Что может сломаться:**

- Если RabbitMQ consumer переиспользует буфер, данные перезапишутся
- Одновременные чтения вернут некорректные данные
- Data corruption при concurrent access

**Решение:**

```go
// ✅ ПРАВИЛЬНО: deep copy
func (r *InMemoryCatalogRepository) UpsertCatalog(ctx context.Context, catalog *dictionaries.Catalog) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Deep copy proto message
    catalogCopy := proto.Clone(catalog).(*dictionaries.Catalog)
    r.catalogs[catalog.Code] = catalogCopy
    return nil
}
```

**Затронутые файлы:**

- Все 4 in-memory репозитория (catalog, invoice, bank-account, foreign-company)

---

### 3. **Отсутствие Context Cancellation Checks** 🔴 КРИТИЧНО

**Проблема:** Долгие операции не проверяют `ctx.Done()`.

```go
// ❌ ПРОБЛЕМА: O(n) фильтрация без проверки контекста
func (r *InMemoryCatalogRepository) ListCatalogsWithFilter(ctx context.Context, ...) ([]*dictionaries.Catalog, int32, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    // Если массив огромный и клиент отменил запрос, мы продолжаем работать
    for _, catalog := range r.catalogs {
        if matchesFilter(catalog, filter) {
            filtered = append(filtered, catalog)
        }
    }
}
```

**Последствия:**

- Waste CPU на отмененные запросы
- Не соблюдаются timeouts
- Потенциальные goroutine leaks

**Решение:**

```go
// ✅ С проверкой контекста
for _, catalog := range r.catalogs {
    select {
    case <-ctx.Done():
        return nil, 0, ctx.Err()
    default:
    }

    if matchesFilter(catalog, filter) {
        filtered = append(filtered, catalog)
    }
}
```

---

### 4. **Отсутствие Input Validation** 🔴 КРИТИЧНО

**Проблема:** In-memory репозитории не валидируют входные данные.

```go
// ❌ НЕТ ВАЛИДАЦИИ
func (r *InMemoryCatalogRepository) UpsertCatalog(ctx context.Context, catalog *dictionaries.Catalog) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.catalogs[catalog.Code] = catalog  // Что если catalog == nil? Что если Code == ""?
    return nil
}
```

**Что может сломаться:**

- Nil pointer dereference → panic
- Пустые ключи → data corruption
- Invalid data в storage

**Решение:**

```go
// ✅ С валидацией
func (r *InMemoryCatalogRepository) UpsertCatalog(ctx context.Context, catalog *dictionaries.Catalog) error {
    if catalog == nil {
        return errors.New("catalog cannot be nil")
    }
    if catalog.Code == "" {
        return errors.New("catalog code cannot be empty")
    }

    r.mu.Lock()
    defer r.mu.Unlock()
    catalogCopy := proto.Clone(catalog).(*dictionaries.Catalog)
    r.catalogs[catalog.Code] = catalogCopy
    return nil
}
```

---

### 5. **Nil Check проблемы в InvoiceDate** 🟡 ВАЖНО

**Проблема:** Потенциальные nil pointer dereference.

```go
// ⚠️ invoice-query-server/internal/infrastructure/repository/inmemory_invoice_repository.go:102
if invoice.InvoiceDate != nil && invoice.InvoiceDate.Value != nil {
    // Date handling
}
```

**Проблема:** Проверка на `Value != nil` не защищает от panic, если `InvoiceDate` станет nil ПОСЛЕ проверки (race condition).

**Решение:**

```go
invoiceDate := invoice.InvoiceDate
if invoiceDate == nil || invoiceDate.Value == nil {
    continue
}
// Теперь безопасно работать с invoiceDate.Value
```

---

### 6. **Memory Leaks - Unbounded Maps** 🔴 КРИТИЧНО

**Проблема:** In-memory maps растут бесконечно, нет механизмов очистки.

```go
type InMemoryCatalogRepository struct {
    mu       sync.RWMutex
    catalogs map[string]*dictionaries.Catalog  // ❌ Растет бесконечно
}
```

**Последствия:**

- OOM (Out of Memory) в production
- Memory leaks при удалении записей (DeleteCatalog не всегда вызывается)

**Решение:**

```go
// Вариант 1: TTL с автоочисткой
type InMemoryCatalogRepository struct {
    mu         sync.RWMutex
    catalogs   map[string]*catalogEntry
    maxSize    int
}

type catalogEntry struct {
    data      *dictionaries.Catalog
    timestamp time.Time
}

// Вариант 2: LRU Cache
import "github.com/hashicorp/golang-lru/v2"
lruCache, _ := lru.New[string, *dictionaries.Catalog](10000)

// Вариант 3: Memory limit
const maxMemoryMB = 512
if r.estimateMemoryUsage() > maxMemoryMB * 1024 * 1024 {
    r.evictOldestEntries()
}
```

---

## 🟡 Важные проблемы

### 7. **O(n) Filtering без индексов** 🟡

**Проблема:** Все фильтрации и поиски делают full scan.

```go
// ❌ O(n) на каждый запрос
func (r *InMemoryCatalogRepository) GetCatalogsByTnvedCode(ctx context.Context, tnvedCode string) ([]*dictionaries.Catalog, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    var result []*dictionaries.Catalog
    for _, catalog := range r.catalogs {  // Full scan
        if catalog.TnvedCode == tnvedCode {
            result = append(result, catalog)
        }
    }
}
```

**Решение:**

```go
type InMemoryCatalogRepository struct {
    mu              sync.RWMutex
    catalogs        map[string]*dictionaries.Catalog
    indexByTnved    map[string][]string  // tnvedCode -> []catalogCode
    indexByGked     map[string][]string  // gkedCode -> []catalogCode
    indexByNumber   map[string]string    // number -> catalogCode
}

func (r *InMemoryCatalogRepository) GetCatalogsByTnvedCode(ctx context.Context, tnvedCode string) ([]*dictionaries.Catalog, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    codes := r.indexByTnved[tnvedCode]  // O(1)
    result := make([]*dictionaries.Catalog, 0, len(codes))
    for _, code := range codes {
        if catalog, ok := r.catalogs[code]; ok {
            result = append(result, catalog)
        }
    }
    return result, nil
}
```

---

### 8. **Отсутствие Logging в In-Memory Repos** 🟡

**Проблема:** Невозможно отладить проблемы в production.

```go
// ❌ Нет логов
func (r *InMemoryCatalogRepository) UpsertCatalog(ctx context.Context, catalog *dictionaries.Catalog) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.catalogs[catalog.Code] = catalog
    return nil
}
```

**Решение:**

```go
type InMemoryCatalogRepository struct {
    mu       sync.RWMutex
    catalogs map[string]*dictionaries.Catalog
    logger   *zap.Logger  // ✅ Добавить logger
}

func (r *InMemoryCatalogRepository) UpsertCatalog(ctx context.Context, catalog *dictionaries.Catalog) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    action := "update"
    if _, exists := r.catalogs[catalog.Code]; !exists {
        action = "insert"
    }

    r.catalogs[catalog.Code] = catalog

    r.logger.Debug("Catalog upserted",
        zap.String("code", catalog.Code),
        zap.String("action", action),
        zap.Int("total_count", len(r.catalogs)))

    return nil
}
```

---

### 9. **Отсутствие Metrics в In-Memory Repos** 🟡

**Проблема:** Нет метрик для мониторинга состояния кэша.

**Решение:**

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    cacheSize = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "inmemory_cache_size",
            Help: "Number of entries in in-memory cache",
        },
        []string{"service", "entity"},
    )

    cacheHits = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "inmemory_cache_hits_total",
            Help: "Total cache hits",
        },
        []string{"service", "entity"},
    )
)

func (r *InMemoryCatalogRepository) ListCatalogs(...) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    cacheSize.WithLabelValues("catalog-query", "catalog").Set(float64(len(r.catalogs)))
    // ...
}
```

---

### 10. **No Memory Limits** 🟡

**Проблема:** Нет защиты от OOM.

**Решение:**

```go
type Config struct {
    MaxCatalogEntries int    `env:"MAX_CATALOG_ENTRIES" envDefault:"100000"`
    MaxMemoryMB       int    `env:"MAX_MEMORY_MB" envDefault:"512"`
}

func (r *InMemoryCatalogRepository) UpsertCatalog(ctx context.Context, catalog *dictionaries.Catalog) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if len(r.catalogs) >= r.maxEntries {
        return errors.New("cache capacity exceeded")
    }

    // Или: evict oldest entries
    if len(r.catalogs) >= r.maxEntries {
        r.evictOldest()
    }

    r.catalogs[catalog.Code] = catalog
    return nil
}
```

---

### 11. **Sorting Not Implemented** 🟡

**Проблема:** Sort параметры принимаются но игнорируются.

```go
// ❌ catalog-query-server
func (r *InMemoryCatalogRepository) ListCatalogsWithFilter(ctx context.Context, filter *ports.CatalogFilter, sort *ports.CatalogSort, ...) {
    // sort parameter полностью игнорируется!
}
```

**Решение:**

```go
import "sort"

func (r *InMemoryCatalogRepository) ListCatalogsWithFilter(ctx context.Context, filter *ports.CatalogFilter, sort *ports.CatalogSort, ...) {
    // ... filtering ...

    // Sort results
    if sort != nil {
        sort.Slice(filtered, func(i, j int) bool {
            switch sort.Field {
            case "name":
                return compareName(filtered[i], filtered[j], sort.Order)
            case "number":
                return compareNumber(filtered[i], filtered[j], sort.Order)
            case "created_at":
                return compareCreatedAt(filtered[i], filtered[j], sort.Order)
            default:
                return false
            }
        })
    }
}
```

---

### 12. **Unsafe String Search** 🟡

**Проблема:** Case-sensitive поиск, нет нормализации.

```go
// ❌ Не найдет "Товар" если ищем "товар"
if strings.Contains(catalog.Name, filter.SearchText) {
    // ...
}
```

**Решение:**

```go
import (
    "strings"
    "unicode"
    "golang.org/x/text/unicode/norm"
)

func normalizeString(s string) string {
    s = strings.TrimSpace(s)
    s = strings.ToLower(s)
    s = norm.NFC.String(s)  // Unicode normalization
    return s
}

if strings.Contains(normalizeString(catalog.Name), normalizeString(filter.SearchText)) {
    // ...
}
```

---

### 13. **Duplicate Code между Repositories** 🟡

**Проблема:** Один и тот же код повторяется в 4 репозиториях.

**Решение:**

```go
// Создать базовый репозиторий
package repository

type BaseInMemoryRepository[K comparable, V any] struct {
    mu      sync.RWMutex
    data    map[K]V
    logger  *zap.Logger
    metrics MetricsCollector
}

func (r *BaseInMemoryRepository[K, V]) Get(key K) (V, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    val, ok := r.data[key]
    return val, ok
}

func (r *BaseInMemoryRepository[K, V]) Set(key K, val V) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.data[key] = val
    r.metrics.IncrementCacheSize()
}

// Использование
type InMemoryCatalogRepository struct {
    *BaseInMemoryRepository[string, *dictionaries.Catalog]
    indexByTnved map[string][]string
}
```

---

## 💡 Рекомендации по улучшению

### 14. **Отсутствие Unit Tests для In-Memory Repos** 📝

**Текущее состояние:**

- `catalog-query-server/tests/` - только integration тесты
- 0% покрытие для in-memory репозиториев

**Рекомендация:**

```go
// catalog-query-server/internal/infrastructure/repository/inmemory_catalog_repository_test.go
func TestInMemoryCatalogRepository_UpsertCatalog(t *testing.T) {
    repo := NewInMemoryCatalogRepository()
    catalog := &dictionaries.Catalog{Code: "CAT001", Name: "Test"}

    err := repo.UpsertCatalog(context.Background(), catalog)
    require.NoError(t, err)

    // Verify
    result, err := repo.GetCatalogByNumber(context.Background(), "CAT001")
    require.NoError(t, err)
    assert.Equal(t, "Test", result.Name)
}

func TestInMemoryCatalogRepository_ConcurrentAccess(t *testing.T) {
    repo := NewInMemoryCatalogRepository()

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            catalog := &dictionaries.Catalog{
                Code: fmt.Sprintf("CAT%03d", index),
                Name: fmt.Sprintf("Catalog %d", index),
            }
            repo.UpsertCatalog(context.Background(), catalog)
        }(i)
    }
    wg.Wait()

    catalogs, total, err := repo.ListCatalogs(context.Background(), 1, 100)
    require.NoError(t, err)
    assert.Equal(t, int32(100), total)
}
```

---

### 15. **Нет Benchmark Tests** 📝

**Рекомендация:**

```go
func BenchmarkInMemoryCatalogRepository_ListCatalogs(b *testing.B) {
    repo := setupRepoWithData(10000) // 10k entries

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _, _ = repo.ListCatalogs(context.Background(), 1, 10)
    }
}

func BenchmarkInMemoryCatalogRepository_FilterByTnved(b *testing.B) {
    repo := setupRepoWithData(10000)
    filter := &ports.CatalogFilter{TnvedCode: "123456"}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _, _ = repo.ListCatalogsWithFilter(context.Background(), filter, nil, 1, 10)
    }
}
```

---

### 16. **TODO комментарии не отслеживаются** 📝

**Найдено 20+ TODO:**

```go
// TODO: получить настоящий organization_id из контекста/параметра
// TODO: Implement default account logic
// TODO: добавить маппинг специфичных ошибок
// TODO: Implement when Country field is added to ForeignCompany proto
```

**Рекомендация:**

1. Создать GitHub Issues для каждого TODO
2. Добавить ссылки на issues в комментариях:

```go
// TODO(#123): получить organization_id из JWT claims
```

---

### 17. **Нет Circuit Breaker для внешних зависимостей** 📝

**Проблема:** Сервисы падают при недоступности PostgreSQL/RabbitMQ.

**Рекомендация:**

```go
import "github.com/sony/gobreaker"

var dbCircuitBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "database",
    MaxRequests: 3,
    Interval:    time.Minute,
    Timeout:     10 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        return counts.Requests >= 3 && failureRatio >= 0.6
    },
})

func (r *PostgresCatalogRepository) ListCatalogs(...) {
    result, err := dbCircuitBreaker.Execute(func() (interface{}, error) {
        return r.queryDatabase(...)
    })
    if err != nil {
        // Fallback: return from cache or return empty
        return r.fallbackStrategy()
    }
}
```

---

### 18. **Отсутствие API Documentation (Swagger)** 📝

**Текущее состояние:**

- `api-gateway` имеет Swagger (`services/api-gateway/SWAGGER_DOCUMENTATION.md`)
- Остальные сервисы - нет

**Рекомендация:**

```go
// Добавить swagger annotations
// @Summary List catalogs
// @Description Get paginated list of catalogs
// @Tags catalog
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} CatalogListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/catalogs [get]
func (h *CatalogHandler) ListCatalogs(ctx context.Context, req *pb.ListCatalogsRequest) (*pb.APIResponse, error)
```

---

### 19. **Hardcoded Constants** 📝

**Проблема:**

```go
// ❌ Hardcoded
if size < 1 || size > 100 {
    size = 10
}
```

**Решение:**

```go
// ✅ Configurable
type PaginationConfig struct {
    DefaultPageSize int `env:"DEFAULT_PAGE_SIZE" envDefault:"10"`
    MaxPageSize     int `env:"MAX_PAGE_SIZE" envDefault:"100"`
}

func normalizePagination(page, size int32, cfg PaginationConfig) (int32, int32) {
    if page < 1 {
        page = 1
    }
    if size < 1 {
        size = int32(cfg.DefaultPageSize)
    }
    if size > int32(cfg.MaxPageSize) {
        size = int32(cfg.MaxPageSize)
    }
    return page, size
}
```

---

### 20. **Panic в Middleware** 📝

**Проблема:**

```go
// ❌ bank-account-server/internal/infrastructure/middleware/auth.go:149
func GetUserIDFromContext(ctx context.Context) string {
    userID, ok := ctx.Value(UserIDKey).(string)
    if !ok {
        panic("user_id not found in context")  // ❌ ОПАСНО
    }
    return userID
}
```

**Решение:**

```go
// ✅ Return error instead of panic
func GetUserIDFromContext(ctx context.Context) (string, error) {
    userID, ok := ctx.Value(UserIDKey).(string)
    if !ok {
        return "", errors.New("user_id not found in context")
    }
    return userID, nil
}

// В handler
userID, err := GetUserIDFromContext(ctx)
if err != nil {
    return nil, status.Error(codes.Unauthenticated, "authentication required")
}
```

---

### 21. **Graceful Shutdown для RabbitMQ Consumer** 📝

**Текущий код:**

```go
// user-query-server: Consumer запущен в горутине без graceful shutdown
go func() {
    if err := consumer.Start(ctx); err != nil {
        logger.Error("Consumer error", zap.Error(err))
    }
}()
```

**Улучшение:**

```go
// В main.go
consumerDone := make(chan struct{})
go func() {
    if err := consumer.Start(ctx); err != nil {
        logger.Error("Consumer error", zap.Error(err))
    }
    close(consumerDone)
}()

// Graceful shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

<-sigChan
logger.Info("Shutting down...")

// Stop consumer
consumer.Close()
<-consumerDone  // Wait for consumer to finish

// Stop gRPC
grpcServer.GracefulStop()
```

---

### 22. **Health Checks для In-Memory Repos** 📝

**Проблема:** Health checks не проверяют состояние in-memory кэша.

**Решение:**

```go
func (h *HealthHandler) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
    // Check cache size
    cacheSize := h.repo.GetSize()
    if cacheSize == 0 {
        return &grpc_health_v1.HealthCheckResponse{
            Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
        }, nil
    }

    // Check last update time
    lastUpdate := h.repo.GetLastUpdateTime()
    if time.Since(lastUpdate) > 5*time.Minute {
        logger.Warn("Cache is stale", zap.Duration("age", time.Since(lastUpdate)))
    }

    return &grpc_health_v1.HealthCheckResponse{
        Status: grpc_health_v1.HealthCheckResponse_SERVING,
    }, nil
}
```

---

### 23. **No Request ID Propagation** 📝

**Рекомендация:**

```go
// Middleware для генерации request ID
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}

// Логирование с request_id
logger.Info("Processing request",
    zap.String("request_id", requestID),
    zap.String("method", req.Method))
```

---

## 📋 План исправления

### Фаза 1: Критические баги (Неделя 1)

**Priority 1: Pagination** (1 день)

- [ ] Стандартизировать формулу пагинации
- [ ] Написать тесты для pagination edge cases
- [ ] Обновить 4 in-memory репозитория

**Priority 2: Race Conditions** (1 день)

- [ ] Добавить deep copy в UpsertCatalog
- [ ] Добавить deep copy в UpsertInvoice
- [ ] Добавить deep copy в UpsertBankAccount
- [ ] Добавить deep copy в UpsertForeignCompany
- [ ] Написать race detector тесты

**Priority 3: Input Validation** (1 день)

- [ ] Добавить nil checks во все Upsert методы
- [ ] Добавить empty key checks
- [ ] Написать unit тесты для validation

**Priority 4: Context Cancellation** (1 день)

- [ ] Добавить ctx.Done() checks в фильтрации
- [ ] Добавить ctx.Done() checks в search
- [ ] Написать timeout тесты

### Фаза 2: Performance & Reliability (Неделя 2)

**Priority 5: Indexes** (2 дня)

- [ ] Реализовать secondary indexes для TnvedCode
- [ ] Реализовать secondary indexes для GkedCode
- [ ] Реализовать secondary indexes для Number
- [ ] Benchmark до/после

**Priority 6: Memory Limits** (1 день)

- [ ] Добавить MaxEntries config
- [ ] Реализовать LRU eviction
- [ ] Добавить memory usage metrics

**Priority 7: Logging & Metrics** (1 день)

- [ ] Добавить structured logging в repos
- [ ] Добавить cache size metrics
- [ ] Добавить operation duration metrics

### Фаза 3: Features (Неделя 3)

**Priority 8: Sorting** (1 день)

- [ ] Реализовать sorting для всех полей
- [ ] Написать тесты

**Priority 9: Health Checks** (1 день)

- [ ] Улучшить health checks
- [ ] Добавить cache staleness check

**Priority 10: Documentation** (1 день)

- [ ] Написать README для каждого сервиса
- [ ] Добавить архитектурные диаграммы
- [ ] Документировать API

### Фаза 4: Testing (Неделя 4)

- [ ] Unit tests для всех repos (80%+ coverage)
- [ ] Integration tests для happy path
- [ ] Load tests (1000 RPS)
- [ ] Chaos engineering (kill pods)

---

## 🎯 Метрики качества

### Текущее состояние

| Метрика            | Значение                         | Оценка     |
| ------------------ | -------------------------------- | ---------- |
| **Архитектура**    | CQRS + Event-driven              | ✅ Отлично |
| **Observability**  | Prometheus + zap + OpenTelemetry | ✅ Отлично |
| **Security**       | JWT + Validation                 | ✅ Хорошо  |
| **Error Handling** | fmt.Errorf + wrapping            | ✅ Хорошо  |
| **Concurrency**    | sync.RWMutex + channels          | ✅ Хорошо  |
| **Testing**        | 49 тестов, но 0% для новых repos | ⚠️ Удовл.  |
| **Performance**    | O(n) фильтры, no indexes         | ⚠️ Удовл.  |
| **Memory Safety**  | Race conditions, no limits       | 🔴 Плохо   |
| **Documentation**  | Partial, no Swagger              | ⚠️ Удовл.  |

### Целевые метрики (после исправлений)

| Метрика             | Целевое значение    |
| ------------------- | ------------------- |
| Test Coverage       | > 80%               |
| Response Time (p95) | < 100ms             |
| Error Rate          | < 0.1%              |
| Memory Usage        | < 512MB per service |
| Availability        | > 99.9%             |

---

## 🔧 Инструменты для мониторинга

### Рекомендуемые инструменты

1. **Race Detector**

```bash
go test -race ./...
```

2. **Memory Profiler**

```bash
go test -memprofile=mem.out
go tool pprof mem.out
```

3. **Static Analysis**

```bash
go vet ./...
staticcheck ./...
golangci-lint run
```

4. **Coverage Report**

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📚 Best Practices

### ✅ Что делается правильно

1. **Clean Architecture** - четкое разделение слоев
2. **Dependency Injection** - repositories, services передаются через конструкторы
3. **Interface segregation** - маленькие интерфейсы (CatalogQueryRepository)
4. **Error wrapping** - везде используется `fmt.Errorf(..., %w, err)`
5. **Context propagation** - ctx передается везде
6. **Structured logging** - zap с полями
7. **Metrics** - Prometheus везде
8. **Health checks** - grpc_health_v1 реализован

### ⚠️ Что нужно улучшить

1. **Memory management** - нет лимитов, нет eviction
2. **Performance** - O(n) фильтры вместо O(1) индексов
3. **Testing** - недостаточно unit tests
4. **Documentation** - нет API docs
5. **Validation** - недостаточно проверок входных данных

---

## 📞 Контакты для вопросов

- **Архитектура CQRS:** [services/catalog-query-server/](services/catalog-query-server/)
- **Event-driven:** [services/user-query-server/internal/infrastructure/messaging/](services/user-query-server/internal/infrastructure/messaging/)
- **Best Example:** `company-server` (полное тестирование, tracing, metrics)

---

## 🚀 Приоритет действий

**СЕЙЧАС (критично):**

1. ✅ Исправить pagination inconsistency
2. ✅ Добавить deep copy в Upsert методы
3. ✅ Добавить input validation
4. ✅ Добавить context cancellation checks

**СКОРО (важно):** 5. Добавить secondary indexes 6. Добавить memory limits 7. Написать unit tests

**ПОЗЖЕ (желательно):** 8. Реализовать sorting 9. Улучшить health checks 10. Написать документацию

---

**Итог:** Архитектура отличная, но есть критические баги в новых in-memory репозиториях. После исправления 6 критических проблем система будет production-ready.
