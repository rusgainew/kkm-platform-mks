# Priority 9: Enhanced Health Checks - Implementation Summary

## ✅ Completed Tasks

Все 6 query серверов теперь имеют расширенные health checks с поддержкой Kubernetes probes.

### Обновлённые сервисы

1. ✅ **catalog-query-server**
2. ✅ **user-query-server**
3. ✅ **document-query-server**
4. ✅ **invoice-query-server**
5. ✅ **bank-account-query-server**
6. ✅ **foreign-company-query-server**

## Реализованные возможности

### 1. Health Check Manager

Создан универсальный пакет `internal/infrastructure/health/` для всех query серверов:

- **Менеджер проверок**: Централизованное управление health checks
- **Prometheus метрики**: `service_health_status` для каждой проверки
- **Гибкая архитектура**: Легко добавлять новые проверки

### 2. Endpoints

Каждый сервис теперь предоставляет 3 endpoint'а:

#### `/health` - Детальная проверка

- Выполняет все зарегистрированные проверки
- Возвращает детальную информацию о состоянии
- Коды ответа: 200 (healthy/degraded), 503 (unhealthy)

#### `/health/live` - Liveness Probe

- Проверяет, что процесс жив
- Используется Kubernetes для перезапуска при сбое
- Всегда возвращает 200 OK, если процесс запущен

#### `/health/ready` - Readiness Probe

- Проверяет готовность к обработке запросов
- Учитывает состояние кэша
- Используется для исключения из load balancer

### 3. Встроенные проверки

#### Cache Staleness Check

- **Цель**: Контроль свежести данных в кэше
- **Параметры**: Max age = 1 час
- **Статусы**:
  - `healthy`: Данные обновлялись < 1 часа назад
  - `degraded`: Данные устарели (> 1 часа) или не инициализированы
- **Метрики**: Возраст данных в секундах, время последнего обновления

#### Cache Size Check

- **Цель**: Контроль размера кэша
- **Параметры**: minSize = 10, maxSize = 100000
- **Статусы**:
  - `healthy`: 10 ≤ size ≤ 100000
  - `degraded`: size = 0, size < 10, или size > 100000
- **Метрики**: Текущий размер кэша

### 4. Repository Changes

Все InMemory репозитории обновлены:

```go
type InMemory{Service}Repository struct {
    // ... existing fields
    lastUpdate time.Time // NEW: Track last update
}

// NEW: Методы для health checks
func (r *InMemory{Service}Repository) GetLastUpdate() time.Time
func (r *InMemory{Service}Repository) GetCacheSize() int
```

Метод `UpsertXXX` обновляет `r.lastUpdate = time.Now()` при каждой вставке/обновлении.

## Изменённые файлы

### Новые файлы

```
services/catalog-query-server/internal/infrastructure/health/health.go
services/user-query-server/internal/infrastructure/health/health.go
services/document-query-server/internal/infrastructure/health/health.go
services/invoice-query-server/internal/infrastructure/health/health.go
services/bank-account-query-server/internal/infrastructure/health/health.go
services/foreign-company-query-server/internal/infrastructure/health/health.go
HEALTH_CHECKS.md (полная документация)
```

### Обновлённые файлы

#### cmd/main.go (все 6 серверов)

- Импорт health пакета с alias `apphealth`
- Инициализация `healthManager := apphealth.NewManager(logger)`
- Регистрация проверок staleness и size
- Обновление HTTP mux для новых endpoints
- Удаление старых `healthCheckHandler` функций

#### internal/infrastructure/repository/inmemory\_\*\_repository.go (все 6)

- Добавлено поле `lastUpdate time.Time`
- Методы `GetLastUpdate()` и `GetCacheSize()`
- Обновление `lastUpdate` в методах Upsert

## Kubernetes Integration

### Пример конфигурации

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: catalog-query-server
spec:
  template:
    spec:
      containers:
        - name: catalog-query-server
          ports:
            - name: grpc
              containerPort: 50051
            - name: health
              containerPort: 8081

          livenessProbe:
            httpGet:
              path: /health/live
              port: health
            initialDelaySeconds: 10
            periodSeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3

          readinessProbe:
            httpGet:
              path: /health/ready
              port: health
            initialDelaySeconds: 5
            periodSeconds: 10
            timeoutSeconds: 3
            failureThreshold: 3
```

## Prometheus Monitoring

### Метрика: service_health_status

**Тип**: Gauge  
**Значения**: 1.0 (healthy) / 0.0 (degraded/unhealthy)  
**Labels**: `check_type` (cache_staleness, cache_size)

### Пример запросов

```promql
# Все неуспешные проверки
service_health_status < 1

# Проверка staleness для всех серверов
service_health_status{check_type="cache_staleness"}

# Пустые кэши
service_health_status{check_type="cache_size"} == 0
```

### Рекомендуемые алерты

```yaml
- alert: CacheStale
  expr: service_health_status{check_type="cache_staleness"} < 1
  for: 5m
  severity: warning

- alert: CacheEmpty
  expr: service_health_status{check_type="cache_size"} < 1
  for: 2m
  severity: critical
```

## Тестирование

### Локальный запуск

```bash
# Запустить сервис
go run ./services/catalog-query-server/cmd/

# Проверить liveness
curl http://localhost:8081/health/live

# Проверить readiness
curl http://localhost:8081/health/ready

# Детальная проверка
curl http://localhost:8081/health | jq .

# Prometheus метрики
curl http://localhost:8081/metrics | grep service_health_status
```

### Ожидаемые ответы

#### Healthy (после загрузки данных)

```json
{
  "status": "healthy",
  "timestamp": "2026-01-27T08:30:00Z",
  "uptime": 3600000000000,
  "checks": {
    "cache_staleness": {
      "status": "healthy",
      "message": "Cache is up to date",
      "last_checked": "2026-01-27T08:30:00Z",
      "details": {
        "last_update": "2026-01-27T08:29:45Z",
        "age_seconds": 15
      }
    },
    "cache_size": {
      "status": "healthy",
      "message": "Cache size is within acceptable range",
      "last_checked": "2026-01-27T08:30:00Z",
      "details": {
        "size": 1234
      }
    }
  }
}
```

#### Degraded (устаревшие данные)

```json
{
  "status": "degraded",
  "checks": {
    "cache_staleness": {
      "status": "degraded",
      "message": "Cache data is stale",
      "details": {
        "age_seconds": 5400,
        "max_age_seconds": 3600
      }
    }
  }
}
```

## Статус компиляции

✅ Все 6 query серверов успешно скомпилированы:

```
✓ catalog-query-server
✓ user-query-server
✓ document-query-server
✓ invoice-query-server
✓ bank-account-query-server
✓ foreign-company-query-server
```

## Следующие шаги

Рекомендации для дальнейшего улучшения:

1. **Добавить проверки внешних зависимостей**
   - RabbitMQ connection check
   - Redis connection check
   - PostgreSQL connection check

2. **Unit тесты для health checks**
   - Тесты для CacheStalenessChecker
   - Тесты для CacheSizeChecker
   - Тесты для Manager.RunChecks

3. **Интеграционные тесты**
   - E2E тесты с Kubernetes probes
   - Тесты с симуляцией сбоев

4. **Дополнительные метрики**
   - Latency проверок
   - Количество failed checks за период
   - История состояний

## Документация

Полная документация доступна в:

- **[HEALTH_CHECKS.md](./HEALTH_CHECKS.md)** - Подробное руководство по health checks
- **[PROMETHEUS_METRICS.md](./PROMETHEUS_METRICS.md)** - Документация по метрикам

## Выводы

✅ **Priority 9 успешно реализован**

Все 6 query серверов теперь имеют:

- ✅ Liveness probes для Kubernetes
- ✅ Readiness probes для load balancing
- ✅ Детальные health checks с метриками
- ✅ Контроль staleness кэша
- ✅ Контроль размера кэша
- ✅ Prometheus интеграцию
- ✅ Production-ready конфигурацию

Сервисы готовы к развёртыванию в Kubernetes с автоматическим управлением жизненным циклом.
