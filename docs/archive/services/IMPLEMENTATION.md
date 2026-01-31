# Company Service Implementation

## 🎉 Реализован сервис CompanyService

### ✅ Что реализовано:

#### 1. **Domain Layer** (Бизнес-логика)

- ✅ Entity `Organization` - организации
- ✅ Entity `Employee` - участники организаций
- ✅ Domain events (OrganizationCreated, EmployeeAdded и др.)
- ✅ Domain errors с четкой семантикой
- ✅ Ports (интерфейсы для repository и event publisher)

#### 2. **Application Layer** (Use Cases)

- ✅ `CompanyService` со всеми методами:
  - `CreateOrganization` - создание организации
  - `GetOrganization` - получение организации
  - `UpdateOrganization` - обновление организации
  - `DeleteOrganization` - удаление организации
  - `ListOrganizations` - список с пагинацией
  - `AddMember` - добавление участника
  - `RemoveMember` - удаление участника
  - `GetOrganizationMembers` - список участников

#### 3. **Infrastructure Layer**

- ✅ PostgreSQL repositories (Organizations, Employees)
- ✅ RabbitMQ event publisher
- ✅ NoOp publisher для тестов
- ✅ Config management
- ✅ Prometheus metrics
- ✅ Structured logging (zap)

#### 4. **Interfaces Layer** (gRPC)

- ✅ `CompanyHandler` - полная реализация всех gRPC методов
- ✅ `HealthHandler` - health checks
- ✅ Маппинг domain ошибок в gRPC статусы
- ✅ Преобразование domain entities в protobuf

#### 5. **Database**

- ✅ Миграции для создания таблиц
- ✅ Индексы для оптимизации
- ✅ Foreign key constraints
- ✅ Unique constraints

#### 6. **DevOps**

- ✅ Dockerfile для сборки
- ✅ Docker Compose конфигурация
- ✅ Интеграция с общей инфраструктурой
- ✅ Prometheus мониторинг
- ✅ Jaeger трассировка
- ✅ Shared PostgreSQL с multiple databases

## 📁 Структура проекта

```
company-server/
├── cmd/
│   └── main.go                          # Entry point
├── internal/
│   ├── domain/                          # Бизнес-логика
│   │   ├── organization.go              # Entity Organization
│   │   ├── employee.go                  # Entity Employee
│   │   ├── errors.go                    # Domain errors
│   │   ├── events/
│   │   │   └── events.go                # Domain events
│   │   └── ports/
│   │       ├── repository.go            # Repository interfaces
│   │       └── event_publisher.go       # Event publisher interface
│   ├── application/
│   │   └── company/
│   │       └── service.go               # Business logic
│   ├── infrastructure/
│   │   ├── config/
│   │   │   └── config.go                # Configuration
│   │   ├── repository/
│   │   │   ├── postgres_organization_repo.go
│   │   │   └── postgres_employee_repo.go
│   │   ├── messaging/
│   │   │   ├── rabbitmq_publisher.go
│   │   │   └── noop_publisher.go
│   │   └── observability/
│   │       └── metrics.go               # Prometheus metrics
│   └── interfaces/
│       └── grpc/
│           ├── company_handler.go       # gRPC handler
│           └── health_handler.go        # Health check
├── migrations/
│   ├── 001_create_organizations_table.up.sql
│   └── 001_create_organizations_table.down.sql
├── Dockerfile
├── .env.example
├── go.mod
└── README.md
```

## 🚀 Запуск

### С Docker Compose (рекомендуется)

```bash
# Запуск всех сервисов
docker-compose up -d

# Проверка логов
docker-compose logs -f company-server

# Применение миграций
docker-compose exec postgres psql -U postgres -d company_db -f /migrations/001_create_organizations_table.up.sql
```

### Локально

```bash
cd company-server

# Установка зависимостей
go mod tidy

# Копирование конфига
cp .env.example .env

# Применение миграций
psql -h localhost -U postgres -d company_db -f migrations/001_create_organizations_table.up.sql

# Запуск
go run cmd/main.go
```

## 🔌 Endpoints

- **gRPC**: `:50052`
- **Metrics**: `:9092/metrics`
- **Health**: `:9092/health`

## 📊 Monitoring

- **Prometheus**: http://localhost:9091
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686

## 🔔 Events

Сервис публикует события в RabbitMQ:

- `organization.created` - организация создана
- `organization.updated` - организация обновлена
- `organization.deleted` - организация удалена
- `employee.added` - участник добавлен
- `employee.removed` - участник удален

## 🧪 Testing

```bash
# Unit тесты
go test ./internal/...

# Integration тесты
go test ./internal/... -tags=integration
```

## 📝 API Examples

### Создание организации

```bash
grpcurl -plaintext \
  -d '{"name": "My Company", "description": "Test org", "owner_id": "user-123"}' \
  localhost:50052 \
  api.company.CompanyService/CreateOrganization
```

### Получение организации

```bash
grpcurl -plaintext \
  -d '{"organization_id": "org-uuid"}' \
  localhost:50052 \
  api.company.CompanyService/GetOrganization
```

### Список организаций

```bash
grpcurl -plaintext \
  -d '{"page": 1, "per_page": 10}' \
  localhost:50052 \
  api.company.CompanyService/ListOrganizations
```

### Добавление участника

```bash
grpcurl -plaintext \
  -d '{"organization_id": "org-uuid", "user_id": "user-456", "role": "member"}' \
  localhost:50052 \
  api.company.CompanyService/AddMember
```

## 🏗️ Архитектура

Сервис следует принципам **Clean Architecture**:

1. **Domain** - независимая бизнес-логика
2. **Application** - use cases и orchestration
3. **Infrastructure** - внешние зависимости (БД, брокеры)
4. **Interfaces** - gRPC handlers

Применены паттерны:

- **Repository Pattern**
- **Dependency Injection**
- **Event-Driven Architecture**
- **CQRS-like** (команды и запросы разделены)

## 🔒 Security

- Валидация всех входных данных
- Проверка прав доступа на уровне сервиса
- Подготовленные SQL запросы (защита от SQL injection)
- Structured logging без sensitive данных

## 📈 Next Steps

1. Реализовать авторизацию (проверка прав через AuthService)
2. Добавить integration тесты
3. Реализовать rate limiting
4. Добавить cache слой (Redis)
5. Реализовать InvoiceService (следующий приоритет)
