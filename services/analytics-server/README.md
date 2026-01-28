# Analytics Server

gRPC сервис для аналитики счетов-фактур.

## Описание

Analytics Server предоставляет аналитические данные на основе счетов-фактур:

- Общая статистика (количество, сумма, средний чек)
- Данные для графиков продаж
- Распределение по статусам
- Распределение по типам операций
- Топ контрагентов
- Месячная выручка

## Технологии

- Go 1.24
- gRPC
- PostgreSQL (materialized views для производительности)
- Redis (кеширование с TTL 5 минут)

## Структура

```
analytics-server/
├── cmd/                          # Точка входа
│   └── main.go
├── internal/
│   ├── application/
│   │   └── services/            # Бизнес-логика
│   ├── domain/
│   │   └── repository/          # Интерфейсы репозиториев
│   ├── infrastructure/
│   │   ├── cache/               # Redis cache
│   │   ├── config/              # Конфигурация
│   │   └── repository/          # PostgreSQL реализация
│   │       └── migrations/      # SQL миграции
│   └── interfaces/
│       └── grpc/                # gRPC handlers
├── Dockerfile
├── go.mod
└── README.md
```

## Переменные окружения

### Сервер

- `GRPC_PORT` - gRPC порт (default: 50070)
- `METRICS_PORT` - Metrics порт (default: 9115)

### База данных

- `DB_HOST` - хост PostgreSQL (default: localhost)
- `DB_PORT` - порт PostgreSQL (default: 5432)
- `DB_USER` - пользователь БД (default: analytics_svc)
- `DB_PASSWORD` - пароль БД (default: analytics_pass_2026)
- `DB_NAME` - название БД (default: analytics_db)
- `DB_SSLMODE` - SSL режим (default: disable)
- `DB_MAX_OPEN_CONNS` - макс. открытых соединений (default: 25)
- `DB_MAX_IDLE_CONNS` - макс. idle соединений (default: 5)

### Redis

- `REDIS_ADDR` - адрес Redis (default: localhost:6379)
- `REDIS_PASSWORD` - пароль Redis
- `REDIS_DB` - номер БД Redis (default: 0)
- `REDIS_TTL` - TTL кеша в секундах (default: 300)

### Observability

- `SERVICE_NAME` - имя сервиса (default: analytics-server)
- `LOG_LEVEL` - уровень логирования (default: info)
- `ENABLE_METRICS` - включить метрики (default: true)
- `ENABLE_TRACING` - включить tracing (default: false)

## Запуск

### Локально

```bash
cd services/analytics-server
go run cmd/main.go
```

### Docker

```bash
docker build -t analytics-server:latest -f services/analytics-server/Dockerfile .
docker run -p 50070:50070 -p 9115:9115 \
  -e DB_HOST=postgres-analytics \
  -e REDIS_ADDR=redis:6379 \
  analytics-server:latest
```

### Docker Compose

```bash
docker-compose up -d analytics-server
```

## Миграции

Миграции применяются автоматически при создании контейнера postgres-analytics.

Для ручного применения:

```bash
psql -h localhost -U analytics_svc -d analytics_db < internal/infrastructure/repository/migrations/001_create_analytics_views.up.sql
```

## Materialized Views

Для производительности используются materialized views:

- `analytics_daily_stats` - дневная статистика
- `analytics_weekly_stats` - недельная статистика
- `analytics_monthly_stats` - месячная статистика
- `analytics_status_distribution` - распределение по статусам
- `analytics_top_contractors` - топ контрагентов
- `analytics_operation_type` - типы операций

Обновление views:

```sql
SELECT refresh_analytics_materialized_views();
```

Рекомендуется настроить автоматическое обновление через pg_cron каждые 5 минут.

## API

gRPC методы:

- `GetDashboardStats` - общая статистика
- `GetSalesChart` - данные для графика продаж
- `GetStatusDistribution` - распределение по статусам
- `GetOperationTypeDistribution` - распределение по типам операций
- `GetTopContractors` - топ контрагентов
- `GetMonthlyRevenue` - месячная выручка

## Производительность

- **PostgreSQL materialized views** - предагрегированные данные
- **Redis кеширование** - TTL 5 минут для всех запросов
- **Connection pooling** - оптимизированный пул соединений
- **Индексы** - все materialized views индексированы

## Мониторинг

- Health check: `grpc://localhost:50070/grpc.health.v1.Health/Check`
- Metrics: `http://localhost:9115/metrics` (Prometheus)

## Разработка

После изменения proto файлов:

```bash
cd proto
make generate-analytics
```

Запуск тестов:

```bash
go test ./...
```
