# Analytics Server - Структура и Миграция

## Выполненная работа

### 1. Создана структура нового gRPC сервиса `analytics-server`

```
services/analytics-server/
├── cmd/
│   └── main.go                                      # Точка входа приложения
├── internal/
│   ├── application/
│   │   └── services/
│   │       └── analytics_service.go                 # Скопирован из api-gateway
│   ├── domain/
│   │   └── repository/
│   │       └── analytics_repository.go              # Скопирован из api-gateway
│   ├── infrastructure/
│   │   ├── cache/
│   │   │   └── redis_cache.go                       # Новый Redis клиент
│   │   ├── config/
│   │   │   └── config.go                            # Конфигурация сервиса
│   │   └── repository/
│   │       ├── analytics_repository_metrics.go      # Скопирован из api-gateway
│   │       ├── postgres_analytics_repository.go     # Скопирован из api-gateway
│   │       └── migrations/
│   │           ├── 001_create_analytics_views.up.sql
│   │           └── 001_create_analytics_views.down.sql
│   └── interfaces/
│       └── grpc/
│           └── analytics_handler.go                 # Заглушка для gRPC handler
├── Dockerfile
├── go.mod
└── README.md
```

### 2. Добавлена новая БД `postgres-analytics` в docker-compose.yml

```yaml
postgres-analytics:
  image: postgres:15-alpine
  container_name: postgres-analytics
  environment:
    POSTGRES_USER: analytics_svc
    POSTGRES_PASSWORD: analytics_pass_2026
    POSTGRES_DB: analytics_db
  ports:
    - "5439:5432"
  volumes:
    - postgres_analytics_data:/var/lib/postgresql/data
```

### 3. Добавлен сервис `analytics-server` в docker-compose.yml

```yaml
analytics-server:
  build:
    context: .
    dockerfile: services/analytics-server/Dockerfile
  container_name: analytics-server
  environment:
    GRPC_PORT: 50070
    METRICS_PORT: 9115
    DB_HOST: postgres-analytics
    DB_USER: analytics_svc
    DB_PASSWORD: analytics_pass_2026
    DB_NAME: analytics_db
    REDIS_ADDR: redis:6379
  ports:
    - "50070:50070" # gRPC
    - "9115:9115" # Metrics
  depends_on:
    - postgres-analytics
    - redis
```

### 4. Исправлен docker-compose.yml

- ✅ Исправлена структура volumes (вынесена из networks)
- ✅ Добавлен `postgres_analytics_data` volume
- ✅ Исправлена опечатка `rabbitmqres` → `rabbitmq`
- ✅ Валидация успешна

## Следующие шаги

### Шаг 1: Генерация proto кода

Сейчас у нас есть структура сервиса и вся логика. Нужно:

1. Создать proto файл `proto/analytics/analytics.proto` на основе существующих интерфейсов
2. Сгенерировать Go код: `protoc --go_out=... --go-grpc_out=...`
3. Реализовать gRPC handler в `internal/interfaces/grpc/analytics_handler.go`

### Шаг 2: Применить миграции к postgres-analytics

```bash
# Запустить postgres-analytics
docker compose up -d postgres-analytics

# Применить миграции
docker exec -i postgres-analytics psql -U analytics_svc -d analytics_db < \
  services/analytics-server/internal/infrastructure/repository/migrations/001_create_analytics_views.up.sql
```

### Шаг 3: Запустить analytics-server

```bash
docker compose up -d --build analytics-server
```

### Шаг 4: Интегрировать в api-gateway

После того как analytics-server работает, обновить api-gateway:

1. Добавить gRPC клиент для analytics-server
2. Обновить analytics_handler.go для вызова gRPC вместо прямого доступа к БД
3. Удалить старые файлы аналитики из api-gateway

## Технические детали

### Порты

- **50070** - gRPC порт analytics-server
- **5439** - PostgreSQL порт postgres-analytics (хост)
- **9115** - Metrics порт analytics-server

### База данных

- **Host**: postgres-analytics
- **Port**: 5432 (internal), 5439 (host)
- **User**: analytics_svc
- **Password**: analytics_pass_2026
- **Database**: analytics_db

### Redis

- **DB**: 6 (выделенная для analytics)
- **TTL**: 300 секунд (5 минут)

### Materialized Views

Все 6 views из оригинальной миграции:

1. `analytics_daily_stats`
2. `analytics_weekly_stats`
3. `analytics_monthly_stats`
4. `analytics_status_distribution`
5. `analytics_top_contractors`
6. `analytics_operation_type`

## Преимущества новой архитектуры

1. **Изоляция**: Аналитика в отдельном сервисе
2. **Масштабируемость**: Можно масштабировать независимо
3. **Отдельная БД**: Не нагружает основную БД счетов
4. **gRPC**: Эффективный протокол для межсервисного взаимодействия
5. **Кеширование**: Redis с TTL 5 минут
6. **Мониторинг**: Отдельные метрики на порту 9115

## Что осталось сделать

- [ ] Создать proto файл на основе интерфейсов
- [ ] Сгенерировать Go код из proto
- [ ] Реализовать gRPC handler
- [ ] Применить миграции к postgres-analytics
- [ ] Собрать и запустить analytics-server
- [ ] Обновить api-gateway для использования gRPC клиента
- [ ] Удалить старую логику из api-gateway
