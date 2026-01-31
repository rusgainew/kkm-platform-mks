# Analytics Backend - Database & Caching Integration

## ✅ Completed Implementation

### 1. **Repository Layer** - PostgreSQL Integration

**File:** `internal/domain/repository/analytics_repository.go`

- Интерфейс для работы с аналитическими данными
- 6 методов для получения агрегированных данных:
  - `GetStats()` - общая статистика (выручка, кол-во накладных, средний чек, контрагенты)
  - `GetSalesData()` - данные продаж с гранулярностью (день/неделя/месяц)
  - `GetStatusDistribution()` - распределение по статусам
  - `GetOperationTypeDistribution()` - распределение по типам операций
  - `GetTopContractors()` - топ N контрагентов
  - `GetMonthlyRevenue()` - месячная выручка

**File:** `internal/infrastructure/repository/postgres_analytics_repository.go` (330 lines)

- Полная реализация всех методов с SQL запросами
- Использует агрегационные функции: `SUM()`, `COUNT()`, `AVG()`, `DATE_TRUNC()`
- Поддержка GROUP BY для группировки данных
- Обработка NULL значений через `COALESCE()`
- Логирование всех операций

### 2. **Service Layer** - Business Logic с Кешированием

**File:** `internal/application/services/analytics_service.go` (285 lines)

- Redis кеширование с TTL = 5 минут
- Умные ключи кеша: `analytics:{endpoint}:{startDate}:{endDate}:{params}`
- Автоматический fallback на БД при cache miss
- Метод `InvalidateCache()` для очистки при обновлении данных
- Проверка доступности кеша перед использованием

### 3. **HTTP Handler** - Обновленные endpoints

**File:** `internal/interfaces/http/analytics_handler.go` (595 lines)

- Заменен mock data на реальные запросы через `analyticsService`
- Умное определение гранулярности данных:
  - ≤7 дней → day
  - ≤90 дней → week
  - > 90 дней → month
- Локализация на русский язык (названия статусов, месяцев)
- Цветовая схема для визуализации (статусы, типы операций)
- Обработка ошибок с HTTP 500 и понятными сообщениями

### 4. **Materialized Views** - Performance Optimization

**Files:**

- `internal/infrastructure/repository/migrations/001_create_analytics_views.up.sql`
- `internal/infrastructure/repository/migrations/001_create_analytics_views.down.sql`

6 материализованных представлений:

1. `analytics_daily_stats` - дневная статистика
2. `analytics_weekly_stats` - еженедельная статистика
3. `analytics_monthly_stats` - месячная статистика
4. `analytics_status_distribution` - распределение по статусам
5. `analytics_top_contractors` - топ контрагенты
6. `analytics_operation_type` - распределение по типам операций

**Функция обновления:** `refresh_analytics_materialized_views()`

### 5. **Prometheus Metrics** - Monitoring

**File:** `internal/infrastructure/repository/analytics_repository_metrics.go` (140 lines)

- Wrapper для добавления метрик к repository
- Метрики для каждого метода:
  - `RecordQueryLatency()` - время выполнения запросов
  - `IncrementErrorCount()` - количество ошибок
- Позволяет отслеживать производительность в Grafana

---

## 📋 Integration Checklist

### Step 1: Database Migration

```bash
# Применить миграцию materialized views
cd services/api-gateway
psql -U your_user -d invoice_db -f internal/infrastructure/repository/migrations/001_create_analytics_views.up.sql

# Создать начальные данные в views
psql -U your_user -d invoice_db -c "SELECT refresh_analytics_materialized_views();"
```

### Step 2: Настройка автообновления Materialized Views

**Вариант A: pg_cron (рекомендуется)**

```sql
-- Установить pg_cron extension
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Настроить обновление каждые 5 минут
SELECT cron.schedule(
    'refresh-analytics',
    '*/5 * * * *',
    'SELECT refresh_analytics_materialized_views();'
);

-- Проверить расписание
SELECT * FROM cron.job;
```

**Вариант B: System Cron**

```bash
# Добавить в crontab
*/5 * * * * psql -U your_user -d invoice_db -c "SELECT refresh_analytics_materialized_views();"
```

**Вариант C: Из приложения (fallback)**

```go
// В main.go добавить горутину
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        _, err := db.Exec("SELECT refresh_analytics_materialized_views()")
        if err != nil {
            logger.Error("Failed to refresh analytics views", zap.Error(err))
        }
    }
}()
```

### Step 3: Dependency Injection

**File:** `cmd/api-gateway/main.go` (нужно добавить)

```go
import (
    "github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/repository"
    infraRepo "github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/repository"
    "github.com/rusgainew/kkm-project-mks/api-gateway/internal/application/services"
    "github.com/rusgainew/kkm-project-mks/api-gateway/internal/interfaces/http"
)

// 1. Создать Analytics Repository
analyticsRepo := infraRepo.NewPostgresAnalyticsRepository(db, logger)

// 2. Обернуть в metrics wrapper
analyticsRepoWithMetrics := infraRepo.NewAnalyticsRepositoryWithMetrics(
    analyticsRepo,
    metrics,
    logger,
)

// 3. Создать Analytics Service с Redis кешем
analyticsService := services.NewAnalyticsService(
    analyticsRepoWithMetrics,
    redisCache, // может быть nil если Redis отключен
    logger,
)

// 4. Создать Handler
analyticsHandler := http.NewAnalyticsHandler(analyticsService, logger)

// 5. Обновить router_configurator.go для использования нового handler
// Заменить старый NewAnalyticsHandler(invoiceQueryService) на новый
```

### Step 4: Router Configuration

**File:** `internal/interfaces/http/router_configurator.go`

```go
// Изменить метод NewRouterConfigurator signature:
func NewRouterConfigurator(
    // ... existing params ...
    analyticsHandler *AnalyticsHandler, // <-- Add this
) *RouterConfigurator {
    return &RouterConfigurator{
        // ... existing fields ...
        analyticsHandler: analyticsHandler, // <-- Add this
    }
}

// Обновить поле структуры:
type RouterConfigurator struct {
    // ... existing fields ...
    analyticsHandler *AnalyticsHandler // <-- Add this
}

// В методе configureAnalyticsRoutes() использовать поле вместо создания нового:
func (rc *RouterConfigurator) configureAnalyticsRoutes(protected *gin.RouterGroup) {
    analytics := protected.Group("/analytics")
    {
        analytics.GET("/stats", rc.analyticsHandler.GetDashboardStats)
        analytics.GET("/sales-chart", rc.analyticsHandler.GetSalesChart)
        analytics.GET("/status-stats", rc.analyticsHandler.GetStatusStats)
        analytics.GET("/operation-stats", rc.analyticsHandler.GetOperationTypeStats)
        analytics.GET("/top-contractors", rc.analyticsHandler.GetTopContractors)
        analytics.GET("/revenue-by-month", rc.analyticsHandler.GetRevenueByMonth)
    }
}
```

### Step 5: Redis Configuration (Optional)

Если Redis не настроен, сервис работает без кеша.

**File:** `.env`

```env
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_TTL=300  # 5 minutes
```

---

## 🧪 Testing

### 1. Unit Tests (TODO)

```bash
cd services/api-gateway
go test ./internal/infrastructure/repository -v
go test ./internal/application/services -v
go test ./internal/interfaces/http -run TestAnalytics -v
```

### 2. Integration Tests

```bash
# Запустить PostgreSQL и заполнить тестовыми данными
docker-compose up -d postgres

# Вставить тестовые накладные
psql -U your_user -d invoice_db < test_data/sample_invoices.sql

# Обновить materialized views
psql -U your_user -d invoice_db -c "SELECT refresh_analytics_materialized_views();"

# Запустить API Gateway
go run cmd/api-gateway/main.go

# Тестовые запросы
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/stats?period=month"

curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/sales-chart?period=week"

curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/top-contractors?limit=5"
```

### 3. Performance Tests

```bash
# Нагрузочное тестирование с Apache Bench
ab -n 1000 -c 10 \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/analytics/stats?period=month

# Ожидаемый результат:
# - С кешем: <50ms latency
# - Без кеша: <200ms latency
# - Cache hit rate: >80%
```

### 4. Cache Verification

```bash
# Подключиться к Redis
redis-cli

# Проверить ключи кеша
KEYS analytics:*

# Посмотреть значение
GET analytics:stats:2024-01-01:2024-01-31

# Проверить TTL
TTL analytics:stats:2024-01-01:2024-01-31
```

---

## 📊 Monitoring

### Prometheus Metrics

```
# Latency запросов к БД
api_gateway_query_latency{operation="analytics",method="get_stats"} 0.045

# Количество ошибок
api_gateway_errors_total{operation="analytics",method="get_stats"} 0

# Cache hit rate
api_gateway_cache_hits_total{name="default"} 850
api_gateway_cache_misses_total{name="default"} 150
```

### Grafana Dashboards

```
- Analytics Query Latency (95th percentile)
- Cache Hit Rate (должно быть >80%)
- Error Rate по endpoints
- Request Rate per endpoint
```

---

## 🚀 Performance Tips

1. **Materialized Views**
   - Обновляются каждые 5 минут
   - Значительно ускоряют агрегационные запросы
   - Для real-time данных используйте прямые запросы

2. **Redis Caching**
   - TTL = 5 минут соответствует обновлению views
   - Инвалидация кеша при создании/обновлении накладных
   - Graceful degradation если Redis недоступен

3. **Database Indexes**
   - `idx_invoices_created_date` - для фильтрации по датам
   - `idx_invoices_status` - для группировки по статусам
   - `idx_invoices_contractor_id` - для топ контрагентов

4. **Query Optimization**
   - Используйте `DATE_TRUNC()` вместо `DATE()`
   - `COALESCE()` для обработки NULL
   - LIMIT в запросах топ контрагентов

---

## 📝 Next Steps

### Phase 4: Advanced Features

- [ ] Сравнение с предыдущим периодом (growth %)
- [ ] Real-time WebSocket updates для live data
- [ ] Export to CSV/Excel
- [ ] Custom date range filters в UI
- [ ] Drill-down по контрагентам

### Phase 5: ML & Predictions

- [ ] Прогнозирование выручки
- [ ] Anomaly detection
- [ ] Рекомендации по оптимизации

---

## 🐛 Troubleshooting

### Problem: Slow queries

**Solution:**

- Проверить indexes: `EXPLAIN ANALYZE SELECT ...`
- Обновить статистику: `ANALYZE invoices;`
- Проверить materialized views: `SELECT * FROM analytics_daily_stats LIMIT 1;`

### Problem: Cache not working

**Solution:**

- Проверить подключение к Redis: `redis-cli PING`
- Проверить переменные окружения `REDIS_ADDR`
- Логи: `grep "cache" api-gateway.log`

### Problem: Incorrect data

**Solution:**

- Обновить views вручную: `SELECT refresh_analytics_materialized_views();`
- Проверить данные в invoices: `SELECT COUNT(*) FROM invoices WHERE created_date >= NOW() - INTERVAL '1 month';`
- Проверить кеш: `redis-cli KEYS analytics:*` и `FLUSHALL` если нужно

---

## 📚 Architecture Summary

```
┌─────────────┐
│   Frontend  │
│  Dashboard  │
└──────┬──────┘
       │ HTTP GET /api/v1/analytics/*
       ↓
┌──────────────────┐
│ AnalyticsHandler │ ← HTTP Layer (6 endpoints)
└────────┬─────────┘
         │
         ↓
┌──────────────────┐
│ AnalyticsService │ ← Business Logic + Caching
└────────┬─────────┘
         │
   ┌─────┴─────┐
   ↓           ↓
┌─────┐  ┌──────────────┐
│Redis│  │PostgresRepo  │ ← Data Layer
└─────┘  └──────┬───────┘
  (5min)        │
                ↓
         ┌──────────────┐
         │  PostgreSQL  │
         │              │
         │ • invoices   │
         │ • mat. views │
         └──────────────┘
```

**Слои:**

1. **HTTP Handler** - парсинг запросов, validation, response formatting
2. **Service** - бизнес-логика, кеширование, orchestration
3. **Repository** - SQL queries, data aggregation, error handling
4. **Database** - persistent storage, materialized views для performance

**Преимущества архитектуры:**

- ✅ Separation of Concerns
- ✅ Easy to test (mock каждый слой)
- ✅ Performance optimization на разных уровнях
- ✅ Graceful degradation (Redis optional)
- ✅ Observability (metrics на каждом слое)
