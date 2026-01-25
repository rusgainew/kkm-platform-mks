# Document Query Server

CQRS Read-side service для документов.

## Архитектура

```
document-server (write-side)
    ↓ publishes document.* events
RabbitMQ (documents exchange)
    ↓ routes by routing key
document-query-server queue
    ↓ consumes events
document-query-server (read-side)
    ↓ updates read model
PostgreSQL (document_db - read optimized)
    ↓
API Gateway (GET /documents)
    ↓
Client
```

## Порты

- **gRPC**: 50062
- **Metrics (Prometheus)**: 9103

## Запуск

```bash
# Development
go run ./cmd/main.go

# Production
docker-compose -f docker-compose.prod.yml up document-query-server
```

## Миграции

```bash
# Создание таблиц read-model
psql -U postgres -d document_db -f migrations/001_create_document_read_model.up.sql
```

## API

### gRPC Methods

```protobuf
service DocumentQueryService {
  rpc GetDocument (GetDocumentRequest) returns (Document);
  rpc ListDocuments (ListDocumentsRequest) returns (ListDocumentsResponse);
  rpc SearchDocuments (SearchDocumentsRequest) returns (ListDocumentsResponse);
  rpc GetDocumentsByStatus (GetDocumentsByStatusRequest) returns (ListDocumentsResponse);
}
```

## Event Handling

Consumsит события из `documents` exchange:

- `document.created` - новый документ
- `document.updated` - обновление документа
- `document.sent` - отправка документа
- `document.approved` - утверждение
- `document.rejected` - отклонение
- `document.archived` - архивирование

## Синхронизация

Read-модель синхронизируется асинхронно из events.
Возможна временная рассинхронизация (eventual consistency).

## Мониторинг

- Metrics: http://localhost:9103/metrics
- Jaeger: http://localhost:16686
