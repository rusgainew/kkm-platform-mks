# 🔄 Миграция User-Query и Document-Query на Redis + In-Memory

**Дата:** 27 января 2026 г.  
**Цель:** Удалить PostgreSQL из user-query-server и document-query-server, использовать только Redis + RabbitMQ + In-memory

---

## 📊 Текущее состояние

### ✅ Уже мигрированы (используют Redis + RabbitMQ + In-memory):

- catalog-query-server
- invoice-query-server
- bank-account-query-server
- foreign-company-query-server

### ⚠️ Требуют миграции (используют PostgreSQL):

- **user-query-server** - читает из `USER_QUERY_DB_URL`
- **document-query-server** - читает из `DOCUMENT_QUERY_DB_URL`

---

## 🎯 План миграции

### user-query-server

#### 1. Обновить config.go

**Файл:** `services/user-query-server/internal/infrastructure/config/config.go`

**Удалить:**

```go
DatabaseURL      string
DatabaseDriver   string
```

**Обновить Validate():**

```go
func (c *Config) Validate() error {
    if c.GRPCPort == 0 {
        return errors.New("GRPC_PORT is required")
    }
    if c.RedisURL == "" {
        return errors.New("USER_QUERY_REDIS_URL is required")
    }
    if c.RabbitMQURL == "" {
        return errors.New("USER_QUERY_RABBITMQ_URL is required")
    }
    return nil
}
```

#### 2. Обновить cmd/main.go

**Файл:** `services/user-query-server/cmd/main.go`

**Удалить:**

```go
// Connect to database
db, err := sqlx.Connect(cfg.DatabaseDriver, cfg.DatabaseURL)
if err != nil {
    logger.Fatal("Failed to connect to database", zap.Error(err))
}
defer db.Close()

// Verify connection
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
if err = db.PingContext(ctx); err != nil {
    cancel()
    logger.Fatal("Failed to ping database", zap.Error(err))
}
cancel()
```

**Заменить на:**

```go
// Create in-memory repository (data will come from RabbitMQ events)
repo := repository.NewInMemoryUserRepository()
```

#### 3. Создать InMemoryUserRepository

**Файл:** `services/user-query-server/internal/infrastructure/repository/inmemory_user_repository.go`

```go
package repository

import (
    "context"
    "strings"
    "sync"

    "github.com/rusgainew/kkm-project-mks/user-query-server/internal/domain/ports"
    pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
)

type InMemoryUserRepository struct {
    mu    sync.RWMutex
    users map[string]*pb.User // key: user_id
}

func NewInMemoryUserRepository() ports.UserQueryRepository {
    return &InMemoryUserRepository{
        users: make(map[string]*pb.User),
    }
}

func (r *InMemoryUserRepository) GetUserByID(ctx context.Context, userID string) (*pb.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    user, ok := r.users[userID]
    if !ok {
        return nil, nil
    }
    return user, nil
}

func (r *InMemoryUserRepository) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for _, user := range r.users {
        if user.Email == email {
            return user, nil
        }
    }
    return nil, nil
}

func (r *InMemoryUserRepository) ListUsers(ctx context.Context, page, size int32) ([]*pb.User, int32, error) {
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 10
    }

    r.mu.RLock()
    defer r.mu.RUnlock()

    allUsers := make([]*pb.User, 0, len(r.users))
    for _, user := range r.users {
        allUsers = append(allUsers, user)
    }

    totalCount := int32(len(allUsers))
    offset := (page - 1) * size

    if offset >= totalCount {
        return []*pb.User{}, totalCount, nil
    }

    end := offset + size
    if end > totalCount {
        end = totalCount
    }

    return allUsers[offset:end], totalCount, nil
}

func (r *InMemoryUserRepository) SearchUsers(ctx context.Context, query string, page, size int32) ([]*pb.User, int32, error) {
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 10
    }

    r.mu.RLock()
    defer r.mu.RUnlock()

    query = strings.ToLower(query)
    var filtered []*pb.User

    for _, user := range r.users {
        if strings.Contains(strings.ToLower(user.Email), query) ||
           strings.Contains(strings.ToLower(user.FirstName), query) ||
           strings.Contains(strings.ToLower(user.LastName), query) {
            filtered = append(filtered, user)
        }
    }

    totalCount := int32(len(filtered))
    offset := (page - 1) * size

    if offset >= totalCount {
        return []*pb.User{}, totalCount, nil
    }

    end := offset + size
    if end > totalCount {
        end = totalCount
    }

    return filtered[offset:end], totalCount, nil
}

// UpsertUser добавляет или обновляет пользователя (вызывается из RabbitMQ consumer)
func (r *InMemoryUserRepository) UpsertUser(ctx context.Context, user *pb.User) error {
    if user == nil || user.UserId == "" {
        return nil
    }

    r.mu.Lock()
    defer r.mu.Unlock()
    r.users[user.UserId] = user
    return nil
}

// DeleteUser удаляет пользователя
func (r *InMemoryUserRepository) DeleteUser(ctx context.Context, userID string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.users, userID)
    return nil
}
```

#### 4. Обновить RabbitMQ Consumer

**Файл:** `services/user-query-server/internal/infrastructure/messaging/rabbitmq_consumer.go`

Заменить все вызовы `db.ExecContext()` и `db.QueryContext()` на вызовы `repo.UpsertUser()`.

#### 5. Обновить docker-compose.yml

**Удалить (если есть):**

```yaml
USER_QUERY_DB_URL: postgres://...
```

**Оставить только:**

```yaml
environment:
  USER_QUERY_REDIS_URL: redis://redis:6379/0
  USER_QUERY_RABBITMQ_URL: amqp://kkm_user:kkm_password@rabbitmq:5672/
  GRPC_PORT: 50061
  METRICS_PORT: 9102
```

---

### document-query-server

#### 1. Обновить config.go

**Файл:** `services/document-query-server/internal/infrastructure/config/config.go`

**Удалить:**

```go
DatabaseURL      string
DatabaseDriver   string
```

**Обновить Validate():**

```go
func (c *Config) Validate() error {
    if c.GRPCPort == 0 {
        return errors.New("GRPC_PORT is required")
    }
    if c.RedisURL == "" {
        return errors.New("DOCUMENT_QUERY_REDIS_URL is required")
    }
    if c.RabbitMQURL == "" {
        return errors.New("DOCUMENT_QUERY_RABBITMQ_URL is required")
    }
    return nil
}
```

#### 2. Обновить cmd/main.go

**Удалить подключение к PostgreSQL, заменить на:**

```go
// Create in-memory repository
repo := repository.NewInMemoryDocumentRepository()
```

#### 3. Создать InMemoryDocumentRepository

**Файл:** `services/document-query-server/internal/infrastructure/repository/inmemory_document_repository.go`

```go
package repository

import (
    "context"
    "strings"
    "sync"

    "github.com/rusgainew/kkm-project-mks/document-query-server/internal/domain/ports"
    pb "github.com/rusgainew/kkm-project-mks/proto-lib/document"
)

type InMemoryDocumentRepository struct {
    mu        sync.RWMutex
    documents map[string]*pb.Document // key: document_id
}

func NewInMemoryDocumentRepository() ports.DocumentQueryRepository {
    return &InMemoryDocumentRepository{
        documents: make(map[string]*pb.Document),
    }
}

func (r *InMemoryDocumentRepository) GetDocumentByID(ctx context.Context, documentID string) (*pb.Document, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    doc, ok := r.documents[documentID]
    if !ok {
        return nil, nil
    }
    return doc, nil
}

func (r *InMemoryDocumentRepository) ListDocuments(ctx context.Context, page, size int32) ([]*pb.Document, int32, error) {
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 10
    }

    r.mu.RLock()
    defer r.mu.RUnlock()

    allDocs := make([]*pb.Document, 0, len(r.documents))
    for _, doc := range r.documents {
        allDocs = append(allDocs, doc)
    }

    totalCount := int32(len(allDocs))
    offset := (page - 1) * size

    if offset >= totalCount {
        return []*pb.Document{}, totalCount, nil
    }

    end := offset + size
    if end > totalCount {
        end = totalCount
    }

    return allDocs[offset:end], totalCount, nil
}

func (r *InMemoryDocumentRepository) SearchDocuments(ctx context.Context, query string, page, size int32) ([]*pb.Document, int32, error) {
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 10
    }

    r.mu.RLock()
    defer r.mu.RUnlock()

    query = strings.ToLower(query)
    var filtered []*pb.Document

    for _, doc := range r.documents {
        if strings.Contains(strings.ToLower(doc.Title), query) ||
           strings.Contains(strings.ToLower(doc.Content), query) {
            filtered = append(filtered, doc)
        }
    }

    totalCount := int32(len(filtered))
    offset := (page - 1) * size

    if offset >= totalCount {
        return []*pb.Document{}, totalCount, nil
    }

    end := offset + size
    if end > totalCount {
        end = totalCount
    }

    return filtered[offset:end], totalCount, nil
}

// UpsertDocument добавляет или обновляет документ
func (r *InMemoryDocumentRepository) UpsertDocument(ctx context.Context, doc *pb.Document) error {
    if doc == nil || doc.DocumentId == "" {
        return nil
    }

    r.mu.Lock()
    defer r.mu.Unlock()
    r.documents[doc.DocumentId] = doc
    return nil
}

// DeleteDocument удаляет документ
func (r *InMemoryDocumentRepository) DeleteDocument(ctx context.Context, documentID string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.documents, documentID)
    return nil
}
```

#### 4. Обновить RabbitMQ Consumer

Заменить все SQL запросы на вызовы `repo.UpsertDocument()`.

#### 5. Обновить docker-compose.yml

**Удалить:**

```yaml
DOCUMENT_QUERY_DB_URL: postgres://...
```

**Оставить:**

```yaml
environment:
  REDIS_URL: redis://redis:6379/5
  RABBITMQ_URL: amqp://kkm_user:kkm_password@rabbitmq:5672/
  GRPC_PORT: 50067
  METRICS_PORT: 9114
```

---

## ✅ Checklist миграции

### user-query-server

- [ ] Удалить `DatabaseURL`, `DatabaseDriver` из config.go
- [ ] Обновить Validate() в config.go
- [ ] Удалить `sqlx.Connect()` из main.go
- [ ] Создать `inmemory_user_repository.go`
- [ ] Обновить RabbitMQ consumer для работы с in-memory
- [ ] Удалить `USER_QUERY_DB_URL` из docker-compose.yml
- [ ] Удалить зависимость от postgres в docker-compose.yml
- [ ] Пересобрать: `docker compose build user-query-server`
- [ ] Перезапустить: `docker compose up -d user-query-server`

### document-query-server

- [ ] Удалить `DatabaseURL`, `DatabaseDriver` из config.go
- [ ] Обновить Validate() в config.go
- [ ] Удалить `sqlx.Connect()` из main.go
- [ ] Создать `inmemory_document_repository.go`
- [ ] Обновить RabbitMQ consumer для работы с in-memory
- [ ] Удалить `DOCUMENT_QUERY_DB_URL` из docker-compose.yml (если есть)
- [ ] Удалить зависимость от postgres в docker-compose.yml
- [ ] Пересобрать: `docker compose build document-query-server`
- [ ] Перезапустить: `docker compose up -d document-query-server`

---

## 🎯 Ожидаемый результат

После миграции **ВСЕ 6 query-серверов** будут использовать:

- ✅ **Redis** для кэширования
- ✅ **RabbitMQ** для получения событий
- ✅ **In-memory** для хранения read-моделей
- ❌ **PostgreSQL НЕ ИСПОЛЬЗУЕТСЯ**

### Архитектура query-серверов:

```
                    ┌──────────────┐
                    │  RabbitMQ    │
                    │  (события)   │
                    └──────┬───────┘
                           │
                           ▼
                  ┌─────────────────┐
                  │  Query Server   │
                  │                 │
                  │  ┌───────────┐  │
                  │  │In-Memory  │  │
                  │  │Repository │  │
                  │  └───────────┘  │
                  │        │        │
                  │        ▼        │
                  │  ┌───────────┐  │
                  │  │Redis Cache│  │
                  │  └───────────┘  │
                  └─────────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   gRPC API   │
                    │   Clients    │
                    └──────────────┘
```

---

## 📊 Сравнение до/после

| Сервис                       | До миграции   | После миграции                    |
| ---------------------------- | ------------- | --------------------------------- |
| catalog-query-server         | PostgreSQL    | ✅ Redis + In-memory              |
| invoice-query-server         | PostgreSQL    | ✅ Redis + In-memory              |
| bank-account-query-server    | PostgreSQL    | ✅ Redis + In-memory              |
| foreign-company-query-server | PostgreSQL    | ✅ Redis + In-memory              |
| **user-query-server**        | ⚠️ PostgreSQL | 🔄 Redis + In-memory (в процессе) |
| **document-query-server**    | ⚠️ PostgreSQL | 🔄 Redis + In-memory (в процессе) |

---

## ⚡ Преимущества миграции

1. **Производительность:** In-memory в 100x быстрее PostgreSQL
2. **Масштабируемость:** Легко горизонтальное масштабирование
3. **Простота:** Нет миграций баз данных для query-side
4. **CQRS правильный:** Полное разделение write/read моделей
5. **Отказоустойчивость:** При падении query-сервера данные восстанавливаются из событий

---

## 🚀 Команды для миграции

```bash
# 1. Остановить текущие query-серверы
docker compose stop user-query-server document-query-server

# 2. Применить изменения в коде (см. выше)

# 3. Пересобрать образы
docker compose build user-query-server document-query-server

# 4. Запустить
docker compose up -d user-query-server document-query-server

# 5. Проверить логи
docker compose logs -f user-query-server
docker compose logs -f document-query-server

# 6. Проверить статус
docker compose ps | grep query
```

---

## 📝 Примечания

- Данные из PostgreSQL НЕ мигрируются автоматически
- In-memory кэш заполняется из RabbitMQ событий
- При первом запуске кэш может быть пустым - это нормально
- События постепенно заполнят read-модель

**Если нужна начальная загрузка данных:**

1. Trigger re-publish событий из command-серверов
2. Или использовать event replay механизм
