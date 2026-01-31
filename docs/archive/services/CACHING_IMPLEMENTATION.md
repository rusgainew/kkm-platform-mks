# Реализация кэширования в API Gateway

## 📋 Обзор

Интегрировано **Redis кэширование** для оптимизации производительности API Gateway. Кэш автоматически использует часто запрашиваемые данные (компании, товары, счета), снижая нагрузку на backend сервисы.

## 🏗️ Архитектура

### Компоненты

```
API Gateway
├── HTTP Request
│   └── Handler
│       └── CompanyService (+ кэш)
│           ├── Check Cache → Cache HIT ✅ (быстрый ответ)
│           └── Cache MISS → gRPC Backend → Store in Cache → Response
```

### Кэширование работает на уровне Application Services

- ✅ **CompanyService** — кэширует GetCompany() запросы
- 🔄 **InvoiceService** — может расширить в будущем
- 🔄 **CatalogService** — может расширить в будущем
- 🔄 **BankAccountService** — может расширить в будущем

## 🔧 Конфигурация

### 1. Переменные окружения (`.env`)

```bash
# Redis подключение
REDIS_ADDR=redis:6379           # Адрес Redis сервера
REDIS_PASSWORD=                 # Пароль (если требуется)
REDIS_TTL=300                   # TTL в секундах (по умолчанию 5 минут)
```

### 2. Docker Compose интеграция

```yaml
# docker-compose.yml
services:
  api-gateway:
    environment:
      REDIS_ADDR: redis:6379
      REDIS_TTL: 300
    depends_on:
      - redis

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
```

## 📝 Реализация

### CompanyService с кэшированием

```go
// internal/application/services/company_service.go

type CompanyService struct {
    connManager *client.ConnectionManager
    serviceAddr string
    metrics     *observability.Metrics
    tracer      *observability.Tracer
    logger      *zap.Logger
    cache       *cache.RedisCache          // ← Добавлен кэш
}

// GetCompany с кэшированием
func (s *CompanyService) GetCompany(ctx context.Context, id string) (*models.Company, error) {
    // 1️⃣ Проверяем кэш
    if s.cache != nil && s.cache.IsAvailable(ctx) {
        var cachedCompany models.Company
        cacheKey := fmt.Sprintf("company:%s", id)
        if err := s.cache.Get(ctx, cacheKey, &cachedCompany); err == nil {
            return &cachedCompany  // ← CACHE HIT ✅
        }
    }

    // 2️⃣ Вызываем backend через gRPC
    resp, err := client.GetOrganization(ctx, ...)
    if err != nil {
        return nil, err
    }

    company := &models.Company{...}

    // 3️⃣ Сохраняем в кэш
    if s.cache != nil && s.cache.IsAvailable(ctx) {
        cacheKey := fmt.Sprintf("company:%s", id)
        s.cache.Set(ctx, cacheKey, company)
    }

    return company, nil
}
```

### Инвалидация кэша при обновлении

```go
// UpdateCompany инвалидирует кэш
func (s *CompanyService) UpdateCompany(ctx context.Context, company *models.Company) (*models.Company, error) {
    // ... вызов backend ...

    // ❌ Инвалидируем кэш после обновления
    if s.cache != nil {
        cacheKey := fmt.Sprintf("company:%s", company.ID)
        s.cache.Delete(ctx, cacheKey)
    }

    return updatedCompany, nil
}
```

## 🔑 Стратегия ключей кэша

```
company:<id>                    # Отдельная компания (TTL: 5 мин)
catalog_item:<id>              # Отдельный товар (TTL: 10 мин)
invoice:<id>                    # Отдельный счет (TTL: 5 мин)
bank_account:<id>              # Отдельный счет (TTL: 5 мин)
```

## 📊 Производительность

### Ожидаемые улучшения

```
Без кэша:
- GetCompany: ~50-100ms (gRPC + БД)
- Нагрузка на БД: 100%

С кэшем (hit):
- GetCompany: ~1-5ms (Redis)
- Нагрузка на БД: 20-30%

Ожидаемый hit rate: 70-80%
```

## ⚠️ Graceful Degradation

Если Redis недоступен:

- ❌ Кэш полностью отключается
- ✅ Приложение продолжает работать
- ✅ Все запросы идут прямо в backend

```go
if s.cache != nil && s.cache.IsAvailable(ctx) {
    // Используем кэш
} else {
    // Fail-open: идем в backend
}
```

## 🔗 Связанные файлы

- [company_service.go](./internal/application/services/company_service.go) — Сервис с кэшем
- [redis_cache.go](./internal/infrastructure/cache/redis_cache.go) — Клиент Redis
- [initializer.go](./pkg/di_container/initializer.go) — Инициализация кэша
- [factories.go](./pkg/di_container/factories.go) — Фабрики сервисов
