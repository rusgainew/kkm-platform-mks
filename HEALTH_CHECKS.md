# Enhanced Health Checks

## Обзор

Все query-серверы (catalog, user, document, invoice, bank-account, foreign-company) теперь имеют расширенные health checks с поддержкой Kubernetes readiness и liveness проб, а также мониторинга staleness кэша.

## Endpoints

### `/health` - Детальная проверка здоровья

Выполняет все зарегистрированные проверки и возвращает детальную информацию.

**Ответ при успехе (200 OK):**

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

**Ответ при деградации (200 OK):**

```json
{
  "status": "degraded",
  "timestamp": "2026-01-27T08:30:00Z",
  "uptime": 3600000000000,
  "checks": {
    "cache_staleness": {
      "status": "degraded",
      "message": "Cache data is stale",
      "last_checked": "2026-01-27T08:30:00Z",
      "details": {
        "last_update": "2026-01-27T07:00:00Z",
        "age_seconds": 5400,
        "max_age_seconds": 3600
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

**Ответ при проблемах (503 Service Unavailable):**

```json
{
  "status": "unhealthy",
  "timestamp": "2026-01-27T08:30:00Z",
  "uptime": 3600000000000,
  "checks": {
    "cache_size": {
      "status": "unhealthy",
      "message": "Cache is empty",
      "last_checked": "2026-01-27T08:30:00Z",
      "details": {
        "size": 0
      }
    }
  }
}
```

### `/health/live` - Liveness проба

Проверяет, что сервис запущен и работает. Используется Kubernetes для перезапуска pod при сбое.

**Ответ (200 OK):**

```json
{
  "status": "alive",
  "uptime": 3600.5
}
```

**Настройка в Kubernetes:**

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8081
  initialDelaySeconds: 10
  periodSeconds: 30
  timeoutSeconds: 5
  failureThreshold: 3
```

### `/health/ready` - Readiness проба

Проверяет готовность сервиса принимать трафик. Учитывает состояние кэша.

**Ответ при готовности (200 OK):**

```json
{
  "status": "healthy",
  "timestamp": "2026-01-27T08:30:00Z",
  "uptime": 3600000000000,
  "checks": { ... }
}
```

**Ответ при деградации (200 OK):**
Сервис продолжает обслуживать запросы, но с предупреждением о деградации.

**Ответ при неготовности (503 Service Unavailable):**

```json
{
  "status": "unhealthy",
  "timestamp": "2026-01-27T08:30:00Z",
  "uptime": 3600000000000,
  "checks": { ... }
}
```

**Настройка в Kubernetes:**

```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 3
  successThreshold: 1
  failureThreshold: 3
```

## Health Checks

### Cache Staleness Check

Проверяет, когда последний раз обновлялся кэш. Считает данные устаревшими, если последнее обновление было более 1 часа назад.

**Статусы:**

- `healthy`: Данные обновлялись в течение последнего часа
- `degraded`: Данные не обновлялись более часа
- `degraded`: Кэш не инициализирован (lastUpdate = zero time)

**Параметры:**

- `maxAge`: 1 час (по умолчанию)

### Cache Size Check

Проверяет размер кэша в памяти.

**Статусы:**

- `healthy`: Размер кэша в допустимых пределах (10 ≤ size ≤ 100000)
- `degraded`: Кэш пустой (size = 0)
- `degraded`: Размер кэша меньше минимума (size < 10)
- `degraded`: Размер кэша превышает максимум (size > 100000)

**Параметры:**

- `minSize`: 10 записей
- `maxSize`: 100000 записей

## Prometheus Метрики

### service_health_status

Gauge метрика со статусом каждой health проверки.

**Labels:**

- `check_type`: Тип проверки (cache_staleness, cache_size)

**Значения:**

- `1.0`: Проверка прошла успешно (healthy)
- `0.0`: Проверка не прошла (degraded или unhealthy)

**Пример запроса:**

```promql
# Все неуспешные проверки
service_health_status{check_type=~"cache_staleness|cache_size"} < 1

# Проверка staleness для catalog-query-server
service_health_status{check_type="cache_staleness", job="catalog-query-server"}
```

## Интеграция с Docker Compose

Обновите `docker-compose.yml` для использования новых health checks:

```yaml
services:
  catalog-query-server:
    image: catalog-query-server:latest
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--quiet",
          "--tries=1",
          "--spider",
          "http://localhost:8081/health/ready",
        ]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    depends_on:
      rabbitmq:
        condition: service_healthy
```

## Интеграция с Kubernetes

### Полный пример Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: catalog-query-server
  namespace: kkm-platform
spec:
  replicas: 3
  selector:
    matchLabels:
      app: catalog-query-server
  template:
    metadata:
      labels:
        app: catalog-query-server
    spec:
      containers:
        - name: catalog-query-server
          image: catalog-query-server:latest
          ports:
            - name: grpc
              containerPort: 50051
              protocol: TCP
            - name: health
              containerPort: 8081
              protocol: TCP

          # Liveness probe - перезапуск при сбое
          livenessProbe:
            httpGet:
              path: /health/live
              port: health
            initialDelaySeconds: 10
            periodSeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3

          # Readiness probe - исключение из балансировки при неготовности
          readinessProbe:
            httpGet:
              path: /health/ready
              port: health
            initialDelaySeconds: 5
            periodSeconds: 10
            timeoutSeconds: 3
            successThreshold: 1
            failureThreshold: 3

          # Startup probe - дать время на инициализацию
          startupProbe:
            httpGet:
              path: /health/live
              port: health
            initialDelaySeconds: 0
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 30 # 30 * 5s = 150s на старт

          env:
            - name: HEALTH_PORT
              value: "8081"
            - name: GRPC_PORT
              value: "50051"

          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
```

### Service для Health Checks

```yaml
apiVersion: v1
kind: Service
metadata:
  name: catalog-query-server-health
  namespace: kkm-platform
spec:
  selector:
    app: catalog-query-server
  ports:
    - name: health
      port: 8081
      targetPort: health
  type: ClusterIP
```

## Мониторинг и Алерты

### Grafana Dashboard

Добавьте панели для мониторинга health checks:

```json
{
  "title": "Health Check Status",
  "targets": [
    {
      "expr": "service_health_status",
      "legendFormat": "{{job}} - {{check_type}}"
    }
  ],
  "yAxes": [
    {
      "min": 0,
      "max": 1
    }
  ]
}
```

### Prometheus Alerts

Создайте алерты для критичных проблем:

```yaml
groups:
  - name: health_checks
    interval: 30s
    rules:
      - alert: CacheStale
        expr: service_health_status{check_type="cache_staleness"} < 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Cache is stale for {{ $labels.job }}"
          description: "Cache has not been updated for more than 1 hour"

      - alert: CacheEmpty
        expr: service_health_status{check_type="cache_size"} < 1
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Cache is empty for {{ $labels.job }}"
          description: "Query server cache is empty, service cannot handle queries"

      - alert: ServiceUnhealthy
        expr: up{job=~".*-query-server"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Query server {{ $labels.job }} is down"
          description: "Service is not responding to health checks"
```

## Тестирование

### Локальное тестирование

```bash
# Проверка liveness
curl http://localhost:8081/health/live

# Проверка readiness
curl http://localhost:8081/health/ready

# Детальная проверка
curl http://localhost:8081/health | jq .

# Проверка метрик
curl http://localhost:8081/metrics | grep service_health_status
```

### Тестирование в Kubernetes

```bash
# Проверка health через pod
kubectl exec -it catalog-query-server-xyz -- wget -qO- http://localhost:8081/health/ready

# Проверка статуса probes
kubectl describe pod catalog-query-server-xyz | grep -A 10 "Liveness\|Readiness"

# Логи health checks
kubectl logs -f catalog-query-server-xyz | grep "health"
```

## Расширение Health Checks

### Добавление Custom Checker

```go
// В cmd/main.go
healthManager.RegisterCheck("custom_check",
    apphealth.CheckFunc(func(ctx context.Context) apphealth.CheckResult {
        // Ваша логика проверки
        return apphealth.CheckResult{
            Status:      apphealth.StatusHealthy,
            Message:     "Custom check passed",
            LastChecked: time.Now(),
            Details: map[string]interface{}{
                "custom_metric": 42,
            },
        }
    }))
```

### Проверка внешних зависимостей

```go
// RabbitMQ connection check
healthManager.RegisterCheck("rabbitmq",
    apphealth.CheckFunc(func(ctx context.Context) apphealth.CheckResult {
        if rabbitConn == nil || rabbitConn.IsClosed() {
            return apphealth.CheckResult{
                Status:      apphealth.StatusUnhealthy,
                Message:     "RabbitMQ connection is closed",
                LastChecked: time.Now(),
            }
        }
        return apphealth.CheckResult{
            Status:      apphealth.StatusHealthy,
            Message:     "RabbitMQ connection is active",
            LastChecked: time.Now(),
        }
    }))
```

## Best Practices

1. **Liveness vs Readiness**
   - Liveness: Проверяет, что процесс жив (не deadlock, не hang)
   - Readiness: Проверяет, что сервис готов обрабатывать запросы

2. **Timeout значения**
   - Liveness: Более длительные таймауты (30s+), меньше рисков ложного срабатывания
   - Readiness: Короткие таймауты (5-10s), быстрее исключить из балансировки

3. **Статусы**
   - `healthy`: Всё в порядке
   - `degraded`: Проблемы, но сервис работает (используется для мониторинга, не влияет на readiness)
   - `unhealthy`: Критичные проблемы (сервис не готов к работе)

4. **Мониторинг**
   - Следите за метриками `service_health_status`
   - Настройте алерты на длительные degraded состояния
   - Анализируйте Details в логах для диагностики

## Troubleshooting

### Проблема: Cache Staleness всегда degraded

**Причина:** RabbitMQ не отправляет события, или consumer не работает.

**Решение:**

```bash
# Проверьте подключение к RabbitMQ
kubectl logs catalog-query-server-xyz | grep "rabbitmq\|consumer"

# Проверьте очереди в RabbitMQ
kubectl exec -it rabbitmq-0 -- rabbitmqctl list_queues
```

### Проблема: Cache Size = 0

**Причина:** Данные не загрузились после старта, или была ошибка при обработке событий.

**Решение:**

```bash
# Проверьте логи на ошибки
kubectl logs catalog-query-server-xyz | grep "error\|ERROR"

# Проверьте метрики операций
curl http://localhost:8081/metrics | grep catalog_operation
```

### Проблема: Pod постоянно перезапускается

**Причина:** Liveness probe не проходит из-за слишком строгих параметров.

**Решение:**

```yaml
# Увеличьте таймауты и failureThreshold
livenessProbe:
  httpGet:
    path: /health/live
    port: health
  initialDelaySeconds: 30 # Было: 10
  periodSeconds: 60 # Было: 30
  timeoutSeconds: 10 # Было: 5
  failureThreshold: 5 # Было: 3
```

## См. также

- [PROMETHEUS_METRICS.md](./PROMETHEUS_METRICS.md) - Документация по метрикам
- [DOCKER_SETUP_COMPLETE.md](./DOCKER_SETUP_COMPLETE.md) - Docker конфигурация
- [monitoring/prometheus-alerts.yml](./monitoring/prometheus-alerts.yml) - Примеры алертов
