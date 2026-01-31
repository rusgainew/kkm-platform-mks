# 🔧 Analytics Backend Integration Guide

## ✅ Текущий статус

### Выполнено:

- ✅ Analytics Repository (interface + PostgreSQL implementation)
- ✅ Analytics Service с Redis caching
- ✅ Analytics HTTP Handler с 6 endpoints
- ✅ Materialized Views SQL миграции
- ✅ Prometheus metrics wrapper
- ✅ DI Container интеграция (Container + Initializer)
- ✅ Router configurator обновлен

### ⚠️ **ВАЖНО: Текущее состояние**

Analytics Service создается **БЕЗ database repository** (nil).
Endpoints будут возвращать ошибку до тех пор, пока не будет подключена база данных.

---

## 🚀 Шаги для полной интеграции

### **Шаг 1: Применить SQL миграции**

Миграции находятся в:

```
services/api-gateway/internal/infrastructure/repository/migrations/
├── 001_create_analytics_views.up.sql   (141 строка)
└── 001_create_analytics_views.down.sql (24 строки)
```

**Применить миграцию:**

```bash
# Через psql
psql -h localhost -U postgres -d invoice_db -f services/api-gateway/internal/infrastructure/repository/migrations/001_create_analytics_views.up.sql

# Или через Docker
docker exec -i postgres psql -U postgres -d invoice_db < services/api-gateway/internal/infrastructure/repository/migrations/001_create_analytics_views.up.sql
```

**Что создается:**

- 6 materialized views для быстрых агрегаций
- Индексы для оптимизации
- Функция `refresh_analytics_materialized_views()` для обновления

---

### **Шаг 2: Настроить автообновление Materialized Views**

**Вариант A: pg_cron (рекомендуется для production)**

```sql
-- Установить pg_cron расширение
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Настроить обновление каждые 5 минут
SELECT cron.schedule(
    'refresh-analytics',           -- job name
    '*/5 * * * *',                  -- every 5 minutes
    'SELECT refresh_analytics_materialized_views();'
);

-- Проверить задачи
SELECT * FROM cron.job;
```

**Вариант B: System Cron (для dev/staging)**

```bash
# Добавить в crontab
crontab -e

# Добавить строку:
*/5 * * * * psql -h localhost -U postgres -d invoice_db -c "SELECT refresh_analytics_materialized_views();"
```

**Вариант C: Ручное обновление (для тестирования)**

```sql
SELECT refresh_analytics_materialized_views();
```

---

### **Шаг 3: Подключить PostgreSQL в DI Container**

**Файл:** `services/api-gateway/pkg/di_container/initializer.go`

**Текущий код (строка ~270):**

```go
func createAnalyticsService(
	cfg *config.Config,
	redisCache *cache.RedisCache,
	metrics *observability.Metrics,
	logger *zap.Logger,
) *services.AnalyticsService {
	logger.Warn("AnalyticsService created WITHOUT repository - analytics endpoints will fail")
	return services.NewAnalyticsService(nil, redisCache, logger)
}
```

**Заменить на:**

```go
func createAnalyticsService(
	cfg *config.Config,
	redisCache *cache.RedisCache,
	metrics *observability.Metrics,
	logger *zap.Logger,
) *services.AnalyticsService {
	// Получить DB connection для invoice database
	db, err := getInvoiceDBConnection(cfg)
	if err != nil {
		logger.Error("Failed to connect to invoice DB for analytics", zap.Error(err))
		return services.NewAnalyticsService(nil, redisCache, logger)
	}

	// Создать repository
	repo := repository.NewPostgresAnalyticsRepository(db, logger)

	// Обернуть в metrics wrapper
	repoWithMetrics := repository.NewAnalyticsRepositoryWithMetrics(repo, metrics, logger)

	// Создать service с кешированием
	return services.NewAnalyticsService(repoWithMetrics, redisCache, logger)
}

// Вспомогательная функция для получения DB connection
func getInvoiceDBConnection(cfg *config.Config) (*sql.DB, error) {
	// TODO: Реализовать подключение к invoice_db
	// Можно переиспользовать существующую логику из других сервисов
	// Пример:
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	return sql.Open("postgres", dsn)
}
```

**Необходимые импорты:**

```go
import (
	"database/sql"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/repository"
	_ "github.com/lib/pq" // PostgreSQL driver
)
```

---

### **Шаг 4: Добавить конфигурацию БД (если отсутствует)**

**Файл:** `services/api-gateway/internal/infrastructure/config/config.go`

**Добавить в Config struct:**

```go
type Config struct {
	// ... existing fields ...

	Database DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}
```

**Environment variables:**

```bash
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=your_password
export DATABASE_DBNAME=invoice_db
```

---

### **Шаг 5: Проверить Redis конфигурацию**

Analytics Service использует Redis для кеширования. Убедитесь, что Redis настроен:

**Environment variables:**

```bash
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=
export REDIS_TTL=300  # 5 минут
```

**Проверить подключение:**

```bash
redis-cli ping
# Должен вернуть: PONG
```

---

## 🧪 Тестирование

### **1. Проверить миграции**

```sql
-- Проверить наличие materialized views
SELECT schemaname, matviewname, definition
FROM pg_matviews
WHERE matviewname LIKE 'analytics%';

-- Должно вернуть 6 views:
-- analytics_daily_stats
-- analytics_weekly_stats
-- analytics_monthly_stats
-- analytics_status_distribution
-- analytics_top_contractors
-- analytics_operation_type
```

### **2. Тестовые данные**

```sql
-- Добавить тестовые invoice для проверки
INSERT INTO invoices (id, document_uuid, invoice_number, number, invoice_date, created_date, total_amount, is_resident, status, legal_person_id, contractor_id, created_by)
VALUES
    ('test-1', 'uuid-1', 'INV-001', '001', NOW(), NOW(), 10000.00, true, 'accepted', 'lp-1', 'c-1', 'user-1'),
    ('test-2', 'uuid-2', 'INV-002', '002', NOW(), NOW(), 15000.00, true, 'signed', 'lp-1', 'c-2', 'user-1'),
    ('test-3', 'uuid-3', 'INV-003', '003', NOW(), NOW(), 20000.00, false, 'draft', 'lp-1', 'c-1', 'user-1');

-- Обновить materialized views
SELECT refresh_analytics_materialized_views();
```

### **3. Тестировать API endpoints**

```bash
# Получить токен (если требуется аутентификация)
TOKEN="your_jwt_token"

# 1. Dashboard Stats
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/stats?period=month"

# 2. Sales Chart
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/sales-chart?period=week"

# 3. Status Statistics
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/status-stats?period=month"

# 4. Operation Type Stats
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/operation-stats?period=month"

# 5. Top Contractors
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/top-contractors?limit=10"

# 6. Revenue by Month
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/revenue-by-month?period=year"
```

### **4. Проверить метрики Prometheus**

```bash
# Метрики доступны на /metrics
curl http://localhost:8080/metrics | grep analytics

# Ожидаемые метрики:
# - analytics_query_latency_seconds
# - analytics_error_count_total
# - cache_hit_total
# - cache_miss_total
```

### **5. Проверить кеш Redis**

```bash
# Посмотреть ключи кеша
redis-cli KEYS "analytics:*"

# Проверить TTL ключа
redis-cli TTL "analytics:stats:2026-01-01:2026-01-31"

# Должен вернуть ~300 (5 минут)
```

---

## 📊 Ожидаемая производительность

### **Без кеша (первый запрос):**

- Simple stats: **< 200ms**
- Sales chart: **< 300ms**
- Status distribution: **< 150ms**
- Top contractors: **< 250ms**

### **С кешем (последующие запросы):**

- Все endpoints: **< 50ms**
- Cache hit rate: **> 80%**

### **С Materialized Views:**

- Сложные агрегации: **< 100ms**
- Обновление views: **< 5 секунд**

---

## 🔍 Отладка

### **Проблема: Analytics endpoints возвращают 500**

**Проверить:**

1. DB connection подключен?

   ```go
   logger.Info("Analytics service initialized",
       zap.Bool("has_repository", analyticsService != nil))
   ```

2. Миграции применены?

   ```sql
   SELECT * FROM pg_matviews WHERE matviewname LIKE 'analytics%';
   ```

3. Логи ошибок:
   ```bash
   docker logs api-gateway | grep analytics
   ```

### **Проблема: Медленные запросы**

**Проверить:**

1. Materialized views обновлены?

   ```sql
   SELECT last_refresh FROM pg_stat_user_tables
   WHERE schemaname = 'public' AND relname LIKE 'analytics%';
   ```

2. Индексы созданы?

   ```sql
   SELECT indexname FROM pg_indexes
   WHERE tablename LIKE 'analytics%';
   ```

3. Redis работает?
   ```bash
   redis-cli PING
   ```

### **Проблема: Кеш не работает**

**Проверить:**

1. Redis подключен?

   ```bash
   redis-cli INFO | grep connected_clients
   ```

2. TTL настроен?

   ```bash
   echo $REDIS_TTL
   # Должен быть 300
   ```

3. Метрики кеша:
   ```bash
   curl localhost:8080/metrics | grep cache
   ```

---

## 📝 Checklist для production

- [ ] SQL миграции применены в production DB
- [ ] pg_cron настроен для автообновления views
- [ ] DB connection pool настроен (max 25 connections)
- [ ] Redis cluster настроен для HA
- [ ] Prometheus alerts настроены:
  - [ ] Slow query alert (> 3 seconds)
  - [ ] High error rate (> 5%)
  - [ ] Low cache hit rate (< 70%)
- [ ] Backup strategy для materialized views
- [ ] Monitoring dashboard создан в Grafana
- [ ] Load testing выполнен (1000 concurrent users)

---

## 🎯 Производительность (цели)

| Метрика           | Target    | Measured |
| ----------------- | --------- | -------- |
| Avg Response Time | < 200ms   | _TBD_    |
| P95 Response Time | < 500ms   | _TBD_    |
| P99 Response Time | < 1s      | _TBD_    |
| Cache Hit Rate    | > 80%     | _TBD_    |
| Error Rate        | < 1%      | _TBD_    |
| Throughput        | > 100 RPS | _TBD_    |

---

## 📚 Связанные документы

- [ANALYTICS_API.md](ANALYTICS_API.md) - API спецификация
- [ANALYTICS_BACKEND_INTEGRATION.md](ANALYTICS_BACKEND_INTEGRATION.md) - Детали реализации
- Database schema: `001_create_analytics_views.up.sql`

---

## 🆘 Поддержка

При возникновении проблем:

1. Проверить логи: `docker logs api-gateway`
2. Проверить метрики: `curl localhost:8080/metrics`
3. Проверить health: `curl localhost:8080/health`
4. Открыть issue в репозитории с логами и метриками
