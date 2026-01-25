# Foreign Company Server

Микросервис для управления иностранными компаниями-контрагентами в системе KKM.

## 🎯 Назначение

Foreign Company Server предоставляет API для:

- Управления справочником иностранных компаний
- Регистрации зарубежных контрагентов
- Хранения реквизитов иностранных партнеров
- Поддержки экспортно-импортных операций

## 🏗️ Архитектура

Реализован с использованием **Clean Architecture**:

```
foreign-company-server/
├── cmd/
│   └── main.go                          # Точка входа
├── internal/
│   ├── domain/                          # Бизнес-логика
│   │   ├── entities.go                  # ForeignCompany сущность
│   │   ├── errors.go                    # Доменные ошибки
│   │   ├── ports/                       # Интерфейсы
│   │   │   ├── repository.go
│   │   │   └── publisher.go
│   │   └── events/                      # Доменные события
│   │       └── events.go
│   ├── application/                     # Use cases
│   │   └── foreigncompany/
│   │       └── service.go               # Бизнес-логика
│   ├── infrastructure/                  # Адаптеры
│   │   ├── config/                      # Конфигурация
│   │   ├── repository/                  # PostgreSQL
│   │   ├── messaging/                   # RabbitMQ
│   │   └── observability/               # Метрики + логирование
│   └── interfaces/                      # Входные точки
│       └── grpc/                        # gRPC handlers
├── migrations/                          # SQL миграции
├── Dockerfile                           # Multi-stage build
└── README.md                            # Эта документация
```

## 🚀 Быстрый старт

### Запуск с Docker Compose

```bash
# Запуск вместе со всеми сервисами
docker-compose up -d foreign-company-server

# Просмотр логов
docker-compose logs -f foreign-company-server

# Остановка
docker-compose stop foreign-company-server
```

### Локальная разработка

```bash
# 1. Установка зависимостей
cd foreign-company-server
go mod download

# 2. Настройка окружения
cp .env.example .env
# Отредактируйте .env файл

# 3. Запуск PostgreSQL (если еще не запущен)
docker-compose up -d postgres

# 4. Применение миграций
psql -U postgres -d foreigncompanydb -f migrations/001_create_foreign_companies_table.up.sql

# 5. Запуск сервера
go run cmd/main.go
```

## 📡 API Endpoints

### gRPC Service: `ForeignCompanyCommandService`

#### 1. CreateForeignCompany

Создание новой иностранной компании.

**Request:**

```protobuf
message CreateForeignCompanyRequest {
  string pin = 1;       // Идентификационный номер
  string full_name = 2; // Полное наименование
}
```

**Response:**

```protobuf
message ForeignCompany {
  int64 id = 1;
  string pin = 2;
  string full_name = 3;
}
```

**Пример с grpcurl:**

```bash
grpcurl -plaintext -d '{
  "pin": "DE123456789",
  "full_name": "Siemens AG"
}' localhost:50056 api.ForeignCompanyCommandService.CreateForeignCompany
```

#### 2. UpdateForeignCompany

Обновление данных компании.

**Request:**

```protobuf
message UpdateForeignCompanyRequest {
  int64 id = 1;
  string pin = 2;
  string full_name = 3;
}
```

**Пример с grpcurl:**

```bash
grpcurl -plaintext -d '{
  "id": 1,
  "pin": "DE987654321",
  "full_name": "Siemens AG (updated)"
}' localhost:50056 api.ForeignCompanyCommandService.UpdateForeignCompany
```

## 🗄️ Модель данных

### Таблица: `foreign_companies`

| Поле           | Тип          | Описание                                   |
| -------------- | ------------ | ------------------------------------------ |
| `id`           | BIGSERIAL    | Уникальный идентификатор                   |
| `pin`          | VARCHAR(50)  | Идентификационный номер (Tax ID, VAT, TIN) |
| `full_name`    | VARCHAR(500) | Полное наименование компании               |
| `country_code` | VARCHAR(2)   | Код страны ISO 3166-1 alpha-2              |
| `address`      | TEXT         | Юридический адрес (опционально)            |
| `is_active`    | BOOLEAN      | Активна ли компания                        |
| `created_by`   | UUID         | ID пользователя-создателя                  |
| `updated_by`   | UUID         | ID последнего редактора                    |
| `created_at`   | TIMESTAMP    | Дата создания                              |
| `updated_at`   | TIMESTAMP    | Дата обновления                            |

### Индексы

- `idx_foreign_companies_pin` - Поиск по PIN
- `idx_foreign_companies_full_name` - Full-text search по названию
- `idx_foreign_companies_country` - Фильтр по стране
- `idx_foreign_companies_pin_unique` - Уникальность PIN (только активные)

## 📊 События RabbitMQ

Сервис публикует события в exchange `foreign_company.events`:

### 1. `foreign_company.created`

```json
{
  "id": 1,
  "pin": "DE123456789",
  "full_name": "Siemens AG",
  "country_code": "DE",
  "created_by": "uuid-here",
  "created_at": "2026-01-08T12:00:00Z"
}
```

### 2. `foreign_company.updated`

```json
{
  "id": 1,
  "pin": "DE987654321",
  "full_name": "Siemens AG (updated)",
  "country_code": "DE",
  "address": "Werner-von-Siemens-Straße 1, Munich",
  "updated_by": "uuid-here",
  "updated_at": "2026-01-08T13:00:00Z"
}
```

### 3. `foreign_company.deleted`

```json
{
  "id": 1,
  "deleted_by": "uuid-here",
  "deleted_at": "2026-01-08T14:00:00Z"
}
```

### 4. `foreign_company.deactivated`

```json
{
  "id": 1,
  "deactivated_by": "uuid-here",
  "deactivated_at": "2026-01-08T15:00:00Z"
}
```

## 📈 Метрики Prometheus

Доступны на `http://localhost:9096/metrics`:

### Кастомные метрики

- `foreign_company_operations_total{operation, status}` - Счетчик операций
- `foreign_company_operation_duration_seconds{operation}` - Время выполнения операций
- `foreign_company_errors_total{operation, error_type}` - Счетчик ошибок
- `foreign_companies_active` - Количество активных компаний
- `foreign_companies_total` - Общее количество компаний

### Стандартные метрики

- `go_goroutines` - Количество горутин
- `go_memstats_*` - Статистика памяти
- `process_*` - Метрики процесса

## 🔍 Health Check

```bash
# gRPC health check
grpcurl -plaintext localhost:50056 grpc.health.v1.Health/Check

# Ожидаемый ответ
{
  "status": "SERVING"
}

# С grpc_health_probe (для Kubernetes)
grpc_health_probe -addr=localhost:50056
```

## 🛠️ Конфигурация

### Переменные окружения

```bash
# Server
FOREIGN_COMPANY_SERVER_HOST=0.0.0.0
FOREIGN_COMPANY_SERVER_PORT=50056
FOREIGN_COMPANY_SERVER_METRICS_PORT=9096

# Database
FOREIGN_COMPANY_SERVER_DB_DRIVER=postgres
FOREIGN_COMPANY_SERVER_DB_URL=postgres://user:pass@host:5432/foreigncompanydb?sslmode=disable

# Redis (опционально, для кэширования)
FOREIGN_COMPANY_SERVER_REDIS_URL=redis://localhost:6379

# RabbitMQ
FOREIGN_COMPANY_SERVER_RABBITMQ_URL=amqp://admin:admin123@localhost:5672/

# Observability
FOREIGN_COMPANY_SERVER_JAEGER_URL=http://localhost:14268/api/traces
FOREIGN_COMPANY_SERVER_LOG_LEVEL=info
FOREIGN_COMPANY_SERVER_ENVIRONMENT=production
```

## 🧪 Тестирование

### Unit тесты

```bash
go test ./internal/domain/... -v
go test ./internal/application/... -v
```

### Integration тесты

```bash
# Требует запущенный PostgreSQL
go test ./internal/infrastructure/repository/... -v
```

### Load тесты

```bash
# С ghz
ghz --insecure \
  --proto ../../proto/api/foreign_company_requests.proto \
  --call api.ForeignCompanyCommandService.CreateForeignCompany \
  --data '{"pin":"TEST123","full_name":"Test Company"}' \
  --total=1000 \
  --concurrency=50 \
  localhost:50056
```

## 📝 Примеры использования

### 1. Создание компании

```bash
grpcurl -plaintext -d '{
  "pin": "US987654321",
  "full_name": "Apple Inc."
}' localhost:50056 api.ForeignCompanyCommandService.CreateForeignCompany
```

### 2. Обновление компании

```bash
grpcurl -plaintext -d '{
  "id": 1,
  "pin": "US987654321",
  "full_name": "Apple Inc. (Cupertino)"
}' localhost:50056 api.ForeignCompanyCommandService.UpdateForeignCompany
```

### 3. Поиск компаний через БД

```sql
-- Все активные компании
SELECT * FROM foreign_companies WHERE is_active = true;

-- Поиск по названию
SELECT * FROM foreign_companies
WHERE to_tsvector('english', full_name) @@ plainto_tsquery('english', 'Apple');

-- Компании из США
SELECT * FROM foreign_companies WHERE country_code = 'US';
```

## 🐛 Troubleshooting

### Проблема: Порт 50056 занят

```bash
# Найти процесс
lsof -i :50056

# Убить процесс
kill -9 <PID>
```

### Проблема: База данных не создана

```bash
# Подключиться к PostgreSQL
docker exec -it postgres psql -U postgres

# Создать базу
CREATE DATABASE foreigncompanydb;

# Применить миграции
\c foreigncompanydb
\i /path/to/migrations/001_create_foreign_companies_table.up.sql
```

### Проблема: RabbitMQ недоступен

Сервис продолжит работу, но события не будут публиковаться. Проверьте:

```bash
# Статус RabbitMQ
docker-compose ps rabbitmq

# Логи RabbitMQ
docker-compose logs rabbitmq
```

## 📚 Связанные документы

- [Proto файлы](../proto/api/foreign_company_requests.proto)
- [Docker Compose Usage](../DOCKER_COMPOSE_USAGE.md)
- [Quickstart Guide](../QUICKSTART.md)
- [Security Guide](../security/README.md)

## 🤝 Вклад

При добавлении новых функций следуйте:

- Clean Architecture принципам
- Паттернам из существующих сервисов (catalog-server, bank-account-server)
- Документируйте изменения в этом README

## 📞 Поддержка

- **Email:** support@company.com
- **Issues:** [GitHub Issues](https://github.com/rusgainew/kkm-project-mks/issues)

---

**Версия:** 1.0.0  
**Статус:** ✅ Production Ready  
**Последнее обновление:** 8 января 2026
