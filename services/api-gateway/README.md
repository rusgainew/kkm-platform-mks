# API Gateway

API Gateway для микросервисной архитектуры KKM Project. Реализован с использованием чистой архитектуры (Clean Architecture).

## Архитектура

Проект следует принципам чистой архитектуры с четким разделением слоев:

### Структура проекта

```
api-gateway/
├── cmd/
│   └── api/
│       └── main.go                    # Точка входа приложения
├── internal/
│   ├── domain/                        # Доменный слой
│   │   ├── models/                    # Доменные модели
│   │   └── ports/                     # Интерфейсы (порты)
│   ├── application/                   # Слой бизнес-логики
│   │   └── services/                  # Use cases и сервисы
│   ├── infrastructure/                # Слой инфраструктуры
│   │   ├── config/                    # Конфигурация
│   │   ├── auth/                      # JWT аутентификация
│   │   ├── client/                    # gRPC клиенты
│   │   ├── middleware/                # HTTP middleware
│   │   └── observability/             # Метрики и трейсинг
│   └── interfaces/                    # Слой интерфейсов
│       └── http/                      # HTTP handlers
└── pkg/                               # Публичные пакеты
```

## Возможности

- ✅ **HTTP REST API** - Преобразование REST запросов в gRPC
- ✅ **JWT Аутентификация** - Проверка и генерация JWT токенов
- ✅ **Rate Limiting** - Ограничение частоты запросов
- ✅ **CORS** - Настраиваемая поддержка CORS
- ✅ **Метрики** - Prometheus метрики
- ✅ **Трейсинг** - OpenTelemetry трейсинг с Jaeger
- ✅ **Graceful Shutdown** - Корректное завершение работы
- ✅ **Health Checks** - Проверка здоровья сервиса
- ✅ **Middleware Pipeline** - Логирование, восстановление после паники, таймауты

## Конфигурация

Создайте файл `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

### Основные параметры

- `HTTP_PORT` - порт HTTP сервера (по умолчанию: 8080)
- `METRICS_PORT` - порт для метрик (по умолчанию: 9090)
- `JWT_SECRET` - секрет для JWT токенов
- `JWT_EXPIRATION` - срок действия JWT токенов (по умолчанию: 24h)

### Адреса gRPC сервисов

- `USER_SERVICE_URL` - адрес user service
- `COMPANY_SERVICE_URL` - адрес company service
- `INVOICE_SERVICE_URL` - адрес invoice service
- `CATALOG_SERVICE_URL` - адрес catalog service
- `BANK_ACCOUNT_SERVICE_URL` - адрес bank account service
- `FOREIGN_COMPANY_SERVICE_URL` - адрес foreign company service

### Query сервисы

- `INVOICE_QUERY_SERVICE_URL` - адрес invoice query service
- `CATALOG_QUERY_SERVICE_URL` - адрес catalog query service
- `BANK_ACCOUNT_QUERY_SERVICE_URL` - адрес bank account query service

## Запуск

### Локальный запуск

```bash
# Установка зависимостей
go mod download

# Запуск
go run cmd/api/main.go
```

### Docker

```bash
# Сборка образа
docker build -t api-gateway:latest .

# Запуск контейнера
docker run -p 8080:8080 -p 9090:9090 --env-file .env api-gateway:latest
```

### Docker Compose

```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f api-gateway
```

## API Endpoints

### Health Checks

```
GET /api/v1/health  - Проверка здоровья сервиса
GET /api/v1/ready   - Проверка готовности сервиса
```

### Companies (требуется аутентификация)

```
POST   /api/v1/companies       - Создать компанию
GET    /api/v1/companies/:id   - Получить компанию по ID
PUT    /api/v1/companies/:id   - Обновить компанию
DELETE /api/v1/companies/:id   - Удалить компанию
GET    /api/v1/companies       - Получить список компаний (с пагинацией)
```

### Invoices (требуется аутентификация)

```
POST   /api/v1/invoices        - Создать счет
GET    /api/v1/invoices/:id    - Получить счет по ID
PUT    /api/v1/invoices/:id    - Обновить счет
DELETE /api/v1/invoices/:id    - Удалить счет
GET    /api/v1/invoices        - Получить список счетов (с пагинацией)
```

## Аутентификация

API Gateway использует JWT токены для аутентификации. Токен должен быть передан в заголовке `Authorization`:

```
Authorization: Bearer <token>
```

### Формат токена

JWT токен содержит следующие claims:

- `user_id` - ID пользователя
- `username` - имя пользователя
- `email` - email пользователя
- `roles` - роли пользователя
- `exp` - время истечения токена
- `iat` - время создания токена

## Метрики

Метрики доступны на порту `METRICS_PORT` (по умолчанию 9090):

```
http://localhost:9090/metrics
```

### Доступные метрики

- `api_gateway_http_requests_total` - общее количество HTTP запросов
- `api_gateway_http_request_duration_seconds` - длительность HTTP запросов
- `api_gateway_grpc_calls_total` - общее количество gRPC вызовов
- `api_gateway_grpc_call_duration_seconds` - длительность gRPC вызовов
- `api_gateway_errors_total` - общее количество ошибок
- `api_gateway_active_connections` - количество активных подключений

## Трейсинг

Трейсы отправляются в Jaeger на адрес `JAEGER_ENDPOINT`.

Просмотр трейсов: `http://localhost:16686` (если Jaeger запущен локально)

## Middleware

API Gateway использует следующие middleware (в порядке выполнения):

1. **RecoveryMiddleware** - восстановление после паники
2. **LoggingMiddleware** - логирование запросов
3. **MetricsMiddleware** - сбор метрик
4. **TracingMiddleware** - трейсинг запросов
5. **CORSMiddleware** - обработка CORS
6. **RateLimitMiddleware** - ограничение частоты запросов
7. **TimeoutMiddleware** - установка таймаута запроса
8. **AuthMiddleware** - проверка JWT токена (только для защищенных роутов)

## Разработка

### Добавление нового сервиса

1. Создайте новый сервис в `internal/application/services/`
2. Создайте handler в `internal/interfaces/http/`
3. Зарегистрируйте роуты в `cmd/api/main.go`

### Добавление новой модели

Модели добавляются в `internal/domain/models/`

### Добавление нового middleware

Middleware добавляются в `internal/infrastructure/middleware/`

## Чистая архитектура

Проект следует принципам чистой архитектуры:

### Domain Layer (Доменный слой)

- Содержит бизнес-логику и доменные модели
- Не зависит от других слоев
- Определяет интерфейсы (порты) для внешних зависимостей

### Application Layer (Слой приложения)

- Содержит use cases и application services
- Оркестрирует выполнение бизнес-логики
- Зависит только от domain layer

### Infrastructure Layer (Слой инфраструктуры)

- Реализует интерфейсы из domain layer
- Содержит технические детали (БД, API клиенты, конфигурация)
- Зависит от domain и application layers

### Interfaces Layer (Слой интерфейсов)

- HTTP handlers, gRPC handlers
- Преобразует внешние запросы в вызовы application layer
- Зависит от всех внутренних слоев

## Зависимости

Основные зависимости:

- `gin-gonic/gin` - HTTP фреймворк
- `golang-jwt/jwt` - JWT аутентификация
- `grpc` - gRPC клиенты
- `prometheus` - метрики
- `opentelemetry` - трейсинг
- `zap` - логирование

## Безопасность

- JWT токены с configurable expiration
- Rate limiting для предотвращения DDoS
- CORS protection
- Request timeouts
- Secure headers (можно расширить)

## Мониторинг

### Prometheus

Метрики собираются через Prometheus на порту `METRICS_PORT`

### Grafana

Можно создать дашборды для визуализации метрик

### Jaeger

Трейсы доступны в Jaeger UI для анализа производительности

## Производительность

- Connection pooling для gRPC клиентов
- Graceful shutdown для корректного завершения запросов
- Настраиваемые таймауты
- Rate limiting для защиты от перегрузки

## TODO

- [ ] Реализовать полные gRPC вызовы в сервисах
- [ ] Добавить handlers для остальных сервисов (user, catalog, bank-account, foreign-company)
- [ ] Добавить интеграционные тесты
- [ ] Добавить Swagger/OpenAPI документацию
- [ ] Реализовать circuit breaker для gRPC вызовов
- [ ] Добавить кэширование
- [ ] Реализовать более сложный rate limiting (по пользователю, по IP)
- [ ] Добавить WebSocket поддержку

## Лицензия

Copyright © 2026 KKM Project
