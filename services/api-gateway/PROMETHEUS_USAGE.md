# Prometheus Использование в API Gateway

## Описание

API Gateway имеет встроенную поддержку Prometheus для мониторинга и сбора метрик. Отдельный HTTP сервер на порту 9097 экспортирует метрики в формате Prometheus.

## Архитектура метрик

### Структура

```
api-gateway (порт 8080)
├── REST endpoints (/api/v1/...)
└── Health checks (/health, /ready)

metrics-server (порт 9097)
└── /metrics (Prometheus формат)
```

### Компоненты

1. **Metrics Registry** - центральное хранилище метрик
2. **HTTP Middleware** - сбор метрик для каждого запроса
3. **Prometheus Handler** - экспорт метрик
4. **Prometheus** - скрейпинг и хранение данных

## Запуск

### 1. Через Docker Compose

```bash
# Запустить весь стек
docker-compose up -d

# Проверить что все работает
curl http://localhost:9097/metrics | head -50
```

### 2. Локально

```bash
# Убедиться что .env правильно настроен
cat .env | grep METRICS
# METRICS_PORT=9097
# ENABLE_METRICS=true

# Запустить API Gateway
cd services/api-gateway
go run cmd/api/main.go

# В другом терминале
curl http://localhost:9097/metrics
```

## Доступ к метрикам

### Напрямую через HTTP

```bash
# Получить все метрики
curl http://localhost:9097/metrics

# Получить конкретную метрику
curl http://localhost:9097/metrics | grep api_gateway_http_requests_total
```

### Через Prometheus Dashboard

```
http://localhost:9090
```

**Примеры запросов:**

```promql
# 1. Количество запросов в секунду
rate(api_gateway_http_requests_total[1m])

# 2. HTTP запросы по статусу
sum(rate(api_gateway_http_requests_total[5m])) by (status)

# 3. Время ответа (95-й перцентиль)
histogram_quantile(0.95, api_gateway_http_request_duration_seconds)

# 4. Активные соединения
api_gateway_active_connections

# 5. Ошибки gRPC
rate(api_gateway_grpc_calls_total{status="error"}[1m])

# 6. Cache hit rate
rate(api_gateway_cache_hits_total[5m]) /
(rate(api_gateway_cache_hits_total[5m]) + rate(api_gateway_cache_misses_total[5m]))

# 7. Размер кэша в памяти
api_gateway_cache_memory_bytes_used

# 8. Активные ключи JWT
api_gateway_active_keys_count

# 9. Среднее время валидации JWT
rate(api_gateway_key_validation_latency_sum[5m]) /
rate(api_gateway_key_validation_latency_count[5m])

# 10. Отозванные токены
api_gateway_revoked_tokens_total
```

### Через Grafana

```
http://localhost:3000
```

**Настройка Data Source:**

1. Нажать `Configuration` → `Data Sources`
2. Нажать `Add data source`
3. Выбрать `Prometheus`
4. URL: `http://prometheus:9090`
5. Нажать `Save & test`

**Создание Dashboard:**

1. Нажать `+` → `Dashboard`
2. Нажать `Add panel`
3. Выбрать метрику из списка
4. Настроить график

## Доступные метрики

### HTTP Метрики

| Метрика                                     | Тип       | Описание                | Labels               |
| ------------------------------------------- | --------- | ----------------------- | -------------------- |
| `api_gateway_http_requests_total`           | Counter   | Всего HTTP запросов     | method, path, status |
| `api_gateway_http_request_duration_seconds` | Histogram | Время ответа в секундах | method, path         |
| `api_gateway_active_connections`            | Gauge     | Активные соединения     | -                    |

### gRPC Метрики

| Метрика                                  | Тип       | Описание           | Labels                  |
| ---------------------------------------- | --------- | ------------------ | ----------------------- |
| `api_gateway_grpc_calls_total`           | Counter   | Всего gRPC вызовов | service, method, status |
| `api_gateway_grpc_call_duration_seconds` | Histogram | Время gRPC вызова  | service, method         |

### Cache Метрики

| Метрика                               | Тип       | Описание                    | Labels                |
| ------------------------------------- | --------- | --------------------------- | --------------------- |
| `api_gateway_cache_hits_total`        | Counter   | Попадания в кэш             | cache_name            |
| `api_gateway_cache_misses_total`      | Counter   | Промахи кэша                | cache_name            |
| `api_gateway_cache_latency_seconds`   | Histogram | Задержка кэша               | cache_name, operation |
| `api_gateway_cache_evictions_total`   | Counter   | Вытеснения из кэша          | cache_name            |
| `api_gateway_cache_entries_count`     | Gauge     | Количество элементов в кэше | cache_name            |
| `api_gateway_cache_memory_bytes_used` | Gauge     | Память используемая кэшем   | -                     |

### JWT Метрики

| Метрика                                      | Тип       | Описание                 | Labels |
| -------------------------------------------- | --------- | ------------------------ | ------ |
| `api_gateway_key_rotations_total`            | Counter   | Всего ротаций ключей     | key_id |
| `api_gateway_active_keys_count`              | Gauge     | Активные ключи JWT       | -      |
| `api_gateway_key_validation_latency_seconds` | Histogram | Задержка валидации ключа | -      |

### Token Blacklist Метрики

| Метрика                                       | Тип       | Описание                    | Labels |
| --------------------------------------------- | --------- | --------------------------- | ------ |
| `api_gateway_revoked_tokens_total`            | Counter   | Отозванные токены           | reason |
| `api_gateway_blacklist_size_gauge`            | Gauge     | Размер чёрного списка       | -      |
| `api_gateway_blacklist_check_latency_seconds` | Histogram | Задержка проверки blacklist | -      |

### Error Метрики

| Метрика                    | Тип     | Описание     | Labels              |
| -------------------------- | ------- | ------------ | ------------------- |
| `api_gateway_errors_total` | Counter | Всего ошибок | error_type, service |

## Примеры использования

### Мониторинг производительности

```promql
# Средняя задержка API
avg(rate(api_gateway_http_request_duration_seconds_sum[5m]) /
    rate(api_gateway_http_request_duration_seconds_count[5m]))

# P95 задержка
histogram_quantile(0.95,
  rate(api_gateway_http_request_duration_seconds_bucket[5m]))

# P99 задержка
histogram_quantile(0.99,
  rate(api_gateway_http_request_duration_seconds_bucket[5m]))
```

### Анализ ошибок

```promql
# Процент ошибок
sum(rate(api_gateway_http_requests_total{status=~"5.."}[5m])) /
sum(rate(api_gateway_http_requests_total[5m])) * 100

# Ошибки по статусу
sum(rate(api_gateway_http_requests_total[1m])) by (status)
```

### Cache анализ

```promql
# Hit rate в процентах
sum(rate(api_gateway_cache_hits_total[5m])) /
(sum(rate(api_gateway_cache_hits_total[5m])) +
 sum(rate(api_gateway_cache_misses_total[5m]))) * 100

# Операции в секунду
rate(api_gateway_cache_hits_total[1m]) +
rate(api_gateway_cache_misses_total[1m])

# Средняя задержка кэша
avg(rate(api_gateway_cache_latency_seconds_sum[5m]) /
    rate(api_gateway_cache_latency_seconds_count[5m]))
```

## Конфигурация

### .env переменные

```bash
# Включить/отключить метрики
ENABLE_METRICS=true

# Порт для метрик
METRICS_PORT=9097
```

### docker-compose.yml

```yaml
# API Gateway экспортирует метрики на порту 9097
ports:
  - "9097:9097" # Metrics

# Prometheus скрейпит метрики
prometheus:
  image: prom/prometheus:latest
  volumes:
    - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
```

### monitoring/prometheus.yml

```yaml
scrape_configs:
  - job_name: "api-gateway"
    static_configs:
      - targets: ["api-gateway:9097"]
    scrape_interval: 10s
    scrape_timeout: 5s
```

## Портет

| Сервис      | Порт | Endpoint   |
| ----------- | ---- | ---------- |
| API Gateway | 8080 | /api/v1/\* |
| Metrics     | 9097 | /metrics   |
| Prometheus  | 9090 | /          |
| Grafana     | 3000 | /          |

## Troubleshooting

### Метрики не экспортируются

```bash
# Проверить что метрики включены
curl -s http://localhost:9097/metrics | head -5

# Проверить что сервис запущен
curl -s http://localhost:8080/health
```

### Prometheus не видит метрики

```bash
# Проверить конфиг prometheus.yml
docker exec prometheus cat /etc/prometheus/prometheus.yml | grep api-gateway

# Проверить targets в Prometheus UI
http://localhost:9090/targets
```

### Высокое использование памяти

```promql
# Проверить размер кэша
api_gateway_cache_memory_bytes_used

# Проверить количество метрик
count(count by(__name__) ({job="api-gateway"}))
```

## Best Practices

1. **Scrape Interval** - установить 10-30 секунд в зависимости от нужд
2. **Retention** - 15 дней по умолчанию (достаточно для анализа)
3. **Alerting** - создавать алерты на аномалии (см. prometheus-alerts.yml)
4. **Dashboard** - создавать персонализированные графики в Grafana
5. **Capacity Planning** - мониторить тренды для планирования ёмкости

## Дополнительно

- [Prometheus Документация](https://prometheus.io/docs/)
- [PromQL Guide](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
