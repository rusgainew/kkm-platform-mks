# Company Server

Микросервис для управления организациями и их участниками в системе ККМ.

## Возможности

- ✅ Создание, обновление, удаление организаций
- ✅ Управление участниками организаций
- ✅ Список организаций с пагинацией и фильтрацией
- ✅ Event-driven архитектура (RabbitMQ)
- ✅ Полная observability (метрики, трассировка, логирование)
- ✅ Health checks
- ✅ PostgreSQL для хранения данных

## Архитектура

Проект следует принципам Clean Architecture:

```
company-server/
├── cmd/                    # Entry point
├── internal/
│   ├── domain/            # Бизнес-логика и entities
│   │   ├── organization.go
│   │   ├── employee.go
│   │   ├── errors.go
│   │   ├── events/
│   │   └── ports/
│   ├── application/       # Use cases
│   │   └── company/
│   ├── infrastructure/    # Внешние зависимости
│   │   ├── repository/
│   │   ├── messaging/
│   │   ├── config/
│   │   └── observability/
│   └── interfaces/        # gRPC handlers
│       └── grpc/
└── migrations/            # Database migrations
```

## Запуск

### С Docker Compose

```bash
docker-compose up -d company-server
```

### Локально

1. Запустите PostgreSQL и RabbitMQ
2. Скопируйте `.env.example` в `.env` и настройте переменные
3. Запустите миграции:
   ```bash
   psql -h localhost -U postgres -d company_db -f migrations/001_create_organizations_table.up.sql
   ```
4. Запустите сервис:
   ```bash
   go run cmd/main.go
   ```

## gRPC API

Сервис реализует `CompanyService` с методами:

- `GetOrganization` - получить организацию
- `CreateOrganization` - создать организацию
- `UpdateOrganization` - обновить организацию
- `DeleteOrganization` - удалить организацию
- `ListOrganizations` - список организаций
- `GetOrganizationMembers` - участники организации
- `AddMember` - добавить участника
- `RemoveMember` - удалить участника

## События

Сервис публикует события в RabbitMQ:

- `organization.created`
- `organization.updated`
- `organization.deleted`
- `employee.added`
- `employee.removed`

## Метрики

Prometheus метрики доступны на `:9092/metrics`

## Health Check

gRPC health check: `grpc.health.v1.Health/Check`
