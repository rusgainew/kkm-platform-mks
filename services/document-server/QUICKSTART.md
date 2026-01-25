# Document-Server Quickstart

## 🚀 Быстрый запуск

### 1. Через Docker Compose (рекомендуется)

```bash
# Из корня проекта
docker-compose up document-server

# Или с пересборкой контейнера
docker-compose up --build document-server

# Проверить логи
docker-compose logs -f document-server
```

### 2. Локальный запуск

```bash
# Убедитесь, что PostgreSQL и RabbitMQ запущены
docker-compose up -d postgres rabbitmq

# Перейти в папку сервиса
cd services/document-server

# Загрузить зависимости
go mod download

# Собрать
go build -o document-server ./cmd/main.go

# Запустить
./document-server
```

## 📋 Проверка работоспособности

### gRPC Health Check

```bash
grpcurl -plaintext localhost:50054 grpc.health.v1.Health/Check
```

Ожидаемый ответ:

```json
{
  "status": "SERVING"
}
```

### Prometheus Metrics

```bash
curl http://localhost:9094/metrics
```

### Создание документа (через grpcurl)

```bash
grpcurl -plaintext \
  -d '{
    "organization_id": "org-123",
    "title": "Sample Document",
    "content": "This is a sample document",
    "created_by": "user-456"
  }' \
  localhost:50054 api.document.DocumentService/CreateDocument
```

## 🔌 Портах

- **50054**: gRPC сервер
- **9094**: Prometheus metrics

## 📊 Статусы документов

- `draft` - черновик (по умолчанию)
- `sent` - отправлен на согласование
- `approved` - одобрен
- `rejected` - отклонен
- `archived` - архивирован

## 🗄️ База данных

```
Сервер: postgres:5432
База: document_db
Пользователь: postgres
Пароль: postgres
```

Миграции выполняются автоматически при старте сервиса.

## 🔐 Аутентификация

Все методы требуют JWT токен в заголовке:

```
Authorization: Bearer <JWT_TOKEN>
```

## 🐛 Отладка

### Увеличить уровень логирования

```bash
LOG_LEVEL=debug docker-compose up document-server
```

### Проверить подключение к БД

```bash
docker-compose exec postgres psql -U postgres -d document_db -c "SELECT * FROM documents LIMIT 1;"
```

### Проверить RabbitMQ

```
Открить http://localhost:15672/
Username: admin
Password: admin123
```

## 📝 Примеры использования

### Через API Gateway (когда интегрирован)

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "org-123",
    "title": "Invoice Document",
    "content": "Invoice content..."
  }'
```

## 🚨 Типичные ошибки

| Ошибка                          | Причина                 | Решение                                  |
| ------------------------------- | ----------------------- | ---------------------------------------- |
| `connection refused`            | Сервис не запущен       | `docker-compose up document-server`      |
| `failed to connect to database` | PostgreSQL не доступна  | `docker-compose up postgres`             |
| `invalid token`                 | Неправильный JWT_SECRET | Проверить `.env` и environment variables |
| `exchange not found`            | RabbitMQ не запущен     | `docker-compose up rabbitmq`             |

## 📚 Дополнительные ресурсы

- [IMPLEMENTATION_REPORT.md](./IMPLEMENTATION_REPORT.md) - Полная документация архитектуры
- [README.md](./README.md) - Описание сервиса
- [services/document-server](../document-server) - Исходный код
