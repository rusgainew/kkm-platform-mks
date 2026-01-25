# Document Service

Микросервис для управления документами в системе.

## Описание

`document-server` реализует следующий функционал:

- Создание, чтение, обновление и удаление документов
- Управление статусами документов (draft, sent, approved, rejected, archived)
- Workflow согласования документов
- История версий документов
- Пулинг событий через RabbitMQ
- JWT аутентификация
- Prometheus metrics и Jaeger трассировка

## Стек технологий

- Go 1.24
- PostgreSQL
- gRPC
- Protocol Buffers
- RabbitMQ
- Jaeger для трассировки
- Prometheus для метрик

## Структура проекта

```
document-server/
├── cmd/main.go                          # Точка входа
├── internal/
│   ├── domain/                          # Business logic
│   │   ├── errors.go                    # Domain ошибки
│   │   ├── events/                      # События
│   │   └── ports/                       # Интерфейсы (interfaces)
│   ├── application/                     # Use cases / Application services
│   │   └── document/
│   │       └── service.go               # Business логика документов
│   ├── infrastructure/                  # Детали реализации
│   │   ├── config/                      # Конфигурация
│   │   ├── repository/                  # Data access
│   │   ├── messaging/                   # Event publishing
│   │   ├── middleware/                  # Middleware (auth, etc)
│   │   ├── migration/                   # Database migrations
│   │   └── observability/               # Трассировка и логирование
│   └── interfaces/                      # Handlers (gRPC, HTTP, etc)
│       └── grpc/                        # gRPC handlers
├── migrations/                          # SQL миграции БД
├── go.mod                               # Go зависимости
└── Dockerfile                           # Контейнеризация
```

## API

### gRPC Методы DocumentService

- `GetDocument(GetDocumentRequest)` - Получить документ
- `CreateDocument(CreateDocumentRequest)` - Создать документ
- `UpdateDocument(UpdateDocumentRequest)` - Обновить документ
- `SendDocument(SendDocumentRequest)` - Отправить на согласование
- `ApproveDocument(ApproveDocumentRequest)` - Одобрить
- `RejectDocument(RejectDocumentRequest)` - Отклонить
- `ArchiveDocument(ArchiveDocumentRequest)` - Архивировать
- `ListDocuments(ListDocumentsRequest)` - Список документов

## Запуск

### Локально

```bash
# Установка зависимостей
go mod download

# Запуск миграций
go run cmd/main.go

# Или через docker-compose (из корня проекта)
docker-compose -f docker-compose.dev.yml up document-server
```

### Docker

```bash
docker build -t document-server:latest .
docker run -e DB_HOST=postgres document-server:latest
```

## Переменные окружения

См. `.env.example`

## Авторы

KKM Project
