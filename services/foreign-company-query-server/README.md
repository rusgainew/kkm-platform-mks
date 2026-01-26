# Foreign Company Query Server

Foreign Company Query Server - сервис для чтения (query) данных об иностранных компаниях в рамках CQRS архитектуры.

## Описание

Этот сервис реализует Query-сторону CQRS паттерна для работы с иностранными компаниями. Он предоставляет оптимизированный API для чтения данных с поддержкой:

- Пагинации
- Фильтрации по PIN, названию, стране
- Полнотекстового поиска
- Кэширования результатов (Redis)
- Трассировки запросов (Jaeger)
- Метрик производительности (Prometheus)

## Архитектура

### CQRS Pattern

- **Command Side**: `foreign-company-server` - обрабатывает создание, обновление, удаление
- **Query Side**: `foreign-company-query-server` (этот сервис) - оптимизирован для чтения

### Технологии

- **Go 1.24**
- **gRPC** - для межсервисного взаимодействия
- **PostgreSQL** - хранилище данных для чтения
- **Redis** - кэширование результатов запросов
- **OpenTelemetry/Jaeger** - distributed tracing
- **Prometheus** - метрики

### Структура

```
foreign-company-query-server/
├── cmd/                  # Точка входа приложения
│   └── main.go
├── internal/
│   ├── domain/          # Бизнес-логика
│   │   └── ports/       # Интерфейсы репозиториев
│   ├── infrastructure/  # Инфраструктурный слой
│   │   ├── cache/       # Redis кэш
│   │   ├── config/      # Конфигурация
│   │   ├── middleware/  # JWT аутентификация
│   │   ├── migration/   # Миграции БД
│   │   ├── observability/ # Метрики, трассировка
│   │   └── repository/  # Репозитории БД
│   └── interfaces/      # Интерфейсы
│       └── grpc/        # gRPC handlers
└── migrations/          # SQL миграции

```

## Запуск

### Локальный запуск

1. Установите зависимости:

```bash
cd services/foreign-company-query-server
go mod download
```

2. Настройте окружение (скопируйте `.env.example` в `.env`):

```bash
cp .env.example .env
```

3. Запустите PostgreSQL и Redis:

```bash
docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:15
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

4. Запустите сервер:

```bash
go run cmd/main.go
```

### Docker

```bash
docker build -t foreign-company-query-server .
docker run -p 50057:50057 -p 8057:8057 --env-file .env foreign-company-query-server
```

## API

### gRPC Endpoints

#### ListForeignCompanies

Получение списка иностранных компаний с пагинацией.

```protobuf
rpc ListForeignCompanies(PageInfo) returns (APIResponse);
```

#### ListForeignCompaniesWithFilter

Получение списка с фильтрацией и сортировкой.

```protobuf
rpc ListForeignCompaniesWithFilter(ForeignCompanyFilterRequest) returns (APIResponse);
```

Параметры фильтрации:

- `pin` - Идентификационный номер (точное совпадение)
- `full_name` - Название компании (ILIKE поиск)
- `country_code` - Код страны
- `search_text` - Полнотекстовый поиск по PIN и названию
- `sort_field` - Поле сортировки (pin, full_name, country_code)
- `sort_order` - Порядок сортировки (ASC, DESC)

#### SearchForeignCompanies

Полнотекстовый поиск по иностранным компаниям.

```protobuf
rpc SearchForeignCompanies(SearchRequest) returns (APIResponse);
```

### Health Check

```bash
curl http://localhost:8057/health
```

### Metrics

```bash
curl http://localhost:8057/metrics
```

## Конфигурация

Все настройки задаются через переменные окружения:

| Переменная                          | Описание          | По умолчанию                      |
| ----------------------------------- | ----------------- | --------------------------------- |
| `FOREIGN_COMPANY_QUERY_GRPC_PORT`   | gRPC порт         | 50057                             |
| `FOREIGN_COMPANY_QUERY_HEALTH_PORT` | Health check порт | 8057                              |
| `FOREIGN_COMPANY_QUERY_DB_URL`      | URL базы данных   | -                                 |
| `REDIS_URL`                         | URL Redis         | redis://localhost:6379/1          |
| `CACHE_TTL`                         | TTL кэша          | 10m                               |
| `JAEGER_ENDPOINT`                   | Jaeger endpoint   | http://localhost:14268/api/traces |
| `JWT_SECRET`                        | JWT секрет        | -                                 |

## Миграции

Миграции выполняются автоматически при запуске сервера. SQL файлы находятся в `internal/infrastructure/migration/sql/`.

### Создание новой миграции

```bash
# Формат: XXX_description.up.sql и XXX_description.down.sql
touch internal/infrastructure/migration/sql/002_add_index.up.sql
touch internal/infrastructure/migration/sql/002_add_index.down.sql
```

## Мониторинг

### Prometheus Metrics

- `foreign_company_query_service_requests_total` - Количество запросов
- `foreign_company_query_service_request_duration_ms` - Длительность запросов
- `foreign_company_query_service_errors_total` - Количество ошибок
- `foreign_company_query_service_query_execution_time_ms` - Время выполнения SQL запросов

### Jaeger Tracing

Распределенная трассировка запросов доступна в Jaeger UI (по умолчанию http://localhost:16686).

## Кэширование

Результаты запросов кэшируются в Redis с TTL по умолчанию 10 минут. Кэш инвалидируется автоматически при изменении данных через command-сервис.

## Аутентификация

Все gRPC методы (кроме health checks) требуют JWT токен в metadata:

```
authorization: Bearer <jwt_token>
```

## Тестирование

```bash
go test ./...
```

## Связанные сервисы

- **foreign-company-server** - Command-сервис для записи данных
- **api-gateway** - API Gateway для HTTP доступа
- **PostgreSQL** - База данных
- **Redis** - Кэш

## Лицензия

Proprietary
