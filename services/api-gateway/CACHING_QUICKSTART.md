# Redis Кэширование: Быстрый старт

## 🚀 Включить кэширование в 3 шага

### 1. Добавить Redis в docker-compose.yml

```yaml
services:
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

  api-gateway:
    environment:
      REDIS_ADDR: redis:6379
      REDIS_TTL: 300
    depends_on:
      - redis

volumes:
  redis_data:
```

### 2. Обновить .env (опционально)

```bash
# Если Redis на другом хосте
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=optional_password
REDIS_TTL=300  # 5 минут
```

### 3. Запустить

```bash
docker-compose -f docker-compose.dev.yml up
```

## ✅ Проверить что работает

```bash
# Проверить Redis доступен
redis-cli ping
# Output: PONG

# Смотреть ключи кэша
redis-cli KEYS "company:*"

# Проверить TTL
redis-cli TTL company:uuid-here
```

## 📊 Мониторинг

```bash
# Размер кэша
redis-cli INFO memory

# Количество ключей
redis-cli DBSIZE

# Статистика
redis-cli INFO stats
```

## 🔧 Расширить кэширование на другие сервисы

### InvoiceService (аналогично CompanyService)

```go
// internal/application/services/invoice_service.go

type InvoiceService struct {
    // ... существующие поля ...
    cache *cache.RedisCache  // ← Добавить
}

func (s *InvoiceService) GetInvoice(ctx context.Context, id string) (*models.Invoice, error) {
    // Проверить кэш
    if s.cache != nil && s.cache.IsAvailable(ctx) {
        var cached models.Invoice
        if err := s.cache.Get(ctx, fmt.Sprintf("invoice:%s", id), &cached); err == nil {
            return &cached
        }
    }

    // ... вызвать backend ...

    // Кэшировать результат
    if s.cache != nil {
        s.cache.Set(ctx, fmt.Sprintf("invoice:%s", id), invoice)
    }

    return invoice, nil
}
```

### В factories.go

```go
func createInvoiceService(..., redisCache *cache.RedisCache) *services.InvoiceService {
    return services.NewInvoiceService(
        // ... остальные параметры ...
        redisCache,  // ← Добавить
    )
}
```

### В initializer.go

```go
invoiceService := createInvoiceService(
    cfg, connManager, metrics, tracer, logger, redisCache,  // ← Добавить redisCache
)
```

## 📈 Ожидаемые результаты

| Метрика             | До        | После      |
| ------------------- | --------- | ---------- |
| Latency (cache hit) | 50-100ms  | 1-5ms      |
| DB Load             | 100%      | 20-30%     |
| Throughput          | 100 req/s | 300+ req/s |
| Hit Rate            | N/A       | 70-80%     |

## 🐛 Отладка

### Кэш не работает?

```bash
# 1. Проверить Redis доступен
docker logs redis

# 2. Проверить переменные окружения
echo $REDIS_ADDR
echo $REDIS_TTL

# 3. Проверить логи API Gateway
docker logs api-gateway | grep -i cache

# 4. Проверить ключи в Redis
redis-cli KEYS "*"
```

### Ключи не появляются?

- Убедитесь что это GET запрос (не POST/PUT/DELETE)
- Проверьте TTL истёк ли: `redis-cli TTL key-name`
- Посмотрите логи: `docker logs api-gateway | grep "cache"`

## 📚 Документация

- [CACHING_IMPLEMENTATION.md](./CACHING_IMPLEMENTATION.md) — Полное описание
- [internal/infrastructure/cache/redis_cache.go](./internal/infrastructure/cache/redis_cache.go) — Реализация кэша
- [internal/application/services/company_service.go](./internal/application/services/company_service.go) — Пример использования

## 🎯 Next Steps

1. ✅ Кэширование CompanyService (DONE)
2. ⏳ Кэширование InvoiceService
3. ⏳ Кэширование CatalogService
4. ⏳ Cache warming для популярных данных
5. ⏳ Adaptive TTL на основе access patterns
