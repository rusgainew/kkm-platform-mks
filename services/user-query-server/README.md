# User Query Server

CQRS Read-side service для пользователей.

## Архитектура

```
user-server (write-side)
    ↓ publishes user.* events
RabbitMQ (user-events exchange)
    ↓ routes by routing key
user-query-server queue
    ↓ consumes events
user-query-server (read-side)
    ↓ updates read model
PostgreSQL (userdb - read optimized)
    ↓
API Gateway (GET /users)
    ↓
Client
```

## Порты

- **gRPC**: 50061
- **Metrics (Prometheus)**: 9102

## Запуск

```bash
# Development
go run ./cmd/main.go

# Production
docker-compose -f docker-compose.prod.yml up user-query-server
```

## Миграции

```bash
# Создание таблиц read-model
psql -U postgres -d userdb -f migrations/001_create_user_read_model.up.sql
```

## API

### gRPC Methods

```protobuf
service UserQueryService {
  rpc GetUser (GetUserRequest) returns (User);
  rpc ListUsers (ListUsersRequest) returns (ListUsersResponse);
  rpc SearchUsers (SearchUsersRequest) returns (ListUsersResponse);
}
```

## Event Handling

Consumsит события из `user-events` exchange:

- `user.created` - новый пользователь
- `user.updated` - обновление профиля
- `user.deleted` - удаление пользователя
- `user.role_changed` - изменение роли

## Синхронизация

Read-модель синхронизируется асинхронно из events.
Возможна временная рассинхронизация (eventual consistency).

## Мониторинг

- Metrics: http://localhost:9102/metrics
- Jaeger: http://localhost:16686
