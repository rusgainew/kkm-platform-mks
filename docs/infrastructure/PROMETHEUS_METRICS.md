# 📊 Prometheus Metrics - Query Services

**Дата:** 27 января 2026 г.  
**Статус:** ✅ Внедрено во всех 6 query серверах

---

## 🎯 Обзор

Все query-сервисы экспортируют метрики in-memory кэша на эндпоинте `/metrics` для мониторинга через Prometheus.

---

## 📈 Метрики по сервисам

### 1. **Catalog Query Server** (`:50061/metrics`)

**Cache Size:**

```promql
catalog_cache_size
# Текущее количество каталогов в памяти
```

**Operation Duration:**

```promql
catalog_operation_duration_seconds{operation="list"}
catalog_operation_duration_seconds{operation="get"}
catalog_operation_duration_seconds{operation="get_by_tnved"}
catalog_operation_duration_seconds{operation="upsert"}
catalog_operation_duration_seconds{operation="delete"}
# Histogram длительности операций (в секундах)
```

**Memory Evictions:**

```promql
catalog_memory_evictions_total
# Счетчик LRU evictions
```

**Index Size:**

```promql
catalog_index_size{index_type="tnved"}
catalog_index_size{index_type="gked"}
# Размер вторичных индексов
```

---

### 2. **User Query Server** (`:50051/metrics`)

**Cache Size:**

```promql
user_cache_size
# Текущее количество пользователей в памяти
```

**Operation Duration:**

```promql
user_operation_duration_seconds{operation="list"}
user_operation_duration_seconds{operation="get"}
user_operation_duration_seconds{operation="search"}
user_operation_duration_seconds{operation="upsert"}
user_operation_duration_seconds{operation="delete"}
```

**Memory Evictions:**

```promql
user_memory_evictions_total
```

---

### 3. **Document Query Server** (`:50062/metrics`)

**Cache Size:**

```promql
document_cache_size
# Текущее количество документов в памяти
```

**Operation Duration:**

```promql
document_operation_duration_seconds{operation="list"}
document_operation_duration_seconds{operation="get"}
document_operation_duration_seconds{operation="search"}
document_operation_duration_seconds{operation="get_pending"}
document_operation_duration_seconds{operation="upsert"}
document_operation_duration_seconds{operation="delete"}
```

**Memory Evictions:**

```promql
document_memory_evictions_total
```

---

### 4. **Invoice Query Server** (`:50063/metrics`)

**Cache Size:**

```promql
invoice_cache_size
```

**Operation Duration:**

```promql
invoice_operation_duration_seconds{operation="list"}
invoice_operation_duration_seconds{operation="get"}
invoice_operation_duration_seconds{operation="search"}
invoice_operation_duration_seconds{operation="upsert"}
invoice_operation_duration_seconds{operation="delete"}
```

**Memory Evictions:**

```promql
invoice_memory_evictions_total
```

---

### 5. **Bank Account Query Server** (`:50064/metrics`)

**Cache Size:**

```promql
bank_account_cache_size
```

**Operation Duration:**

```promql
bank_account_operation_duration_seconds{operation="list"}
bank_account_operation_duration_seconds{operation="get"}
bank_account_operation_duration_seconds{operation="search"}
bank_account_operation_duration_seconds{operation="upsert"}
bank_account_operation_duration_seconds{operation="delete"}
```

**Memory Evictions:**

```promql
bank_account_memory_evictions_total
```

---

### 6. **Foreign Company Query Server** (`:50065/metrics`)

**Cache Size:**

```promql
foreign_company_cache_size
```

**Operation Duration:**

```promql
foreign_company_operation_duration_seconds{operation="list"}
foreign_company_operation_duration_seconds{operation="get"}
foreign_company_operation_duration_seconds{operation="search"}
foreign_company_operation_duration_seconds{operation="upsert"}
foreign_company_operation_duration_seconds{operation="delete"}
```

**Memory Evictions:**

```promql
foreign_company_memory_evictions_total
```

---

## 🔥 Пример Grafana Queries

### Cache Hit Rate (если есть Redis fallback)

```promql
rate(catalog_operation_duration_seconds_count{operation="get"}[5m])
```

### Average Operation Latency

```promql
rate(catalog_operation_duration_seconds_sum[5m])
/
rate(catalog_operation_duration_seconds_count[5m])
```

### Memory Usage Trend

```promql
catalog_cache_size
```

### Eviction Rate (per minute)

```promql
rate(catalog_memory_evictions_total[1m]) * 60
```

### p95 Latency

```promql
histogram_quantile(0.95,
  rate(catalog_operation_duration_seconds_bucket[5m])
)
```

### Cache Size Percentage

```promql
(catalog_cache_size / 100000) * 100
# 100000 = defaultMaxEntries
```

---

## 🚨 Рекомендуемые алерты

### High Memory Usage

```yaml
- alert: CatalogCacheHighMemory
  expr: catalog_cache_size > 90000
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Catalog cache near memory limit"
    description: "Cache size is {{ $value }}, approaching max 100000"
```

### High Eviction Rate

```yaml
- alert: CatalogHighEvictionRate
  expr: rate(catalog_memory_evictions_total[5m]) > 10
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "High cache eviction rate"
    description: "Evicting {{ $value }} entries/sec"
```

### Slow Operations

```yaml
- alert: CatalogSlowOperations
  expr: histogram_quantile(0.95, rate(catalog_operation_duration_seconds_bucket[5m])) > 0.5
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Catalog operations are slow"
    description: "p95 latency is {{ $value }}s"
```

---

## 📊 Structured Logging

Все операции логируются с zap:

**Upsert:**

```go
r.logger.Debug("Catalog upserted",
    zap.String("action", "insert|update"),
    zap.String("number", catalog.Number),
    zap.String("tnved_code", catalog.TnvedCode),
    zap.Int("cache_size", len(r.catalogs)))
```

**Delete:**

```go
r.logger.Debug("Catalog deleted",
    zap.String("code", code),
    zap.Int("cache_size", len(r.catalogs)))
```

**Eviction:**

```go
r.logger.Info("Memory eviction performed",
    zap.Int("evicted_count", evicted),
    zap.Int("remaining_count", len(r.catalogs)),
    zap.Int("max_entries", r.maxEntries))
```

---

## 🎓 Best Practices

1. **Мониторинг cache_size** - следить за приближением к лимиту
2. **p95/p99 latency** - для SLA мониторинга
3. **Eviction rate** - индикатор нехватки памяти
4. **Index size** (только catalog) - для оценки overhead индексов
5. **Логи уровня Info** только для evictions - остальное Debug

---

## ✅ Результаты

- **6/6 query серверов** с метриками
- **5 операций** инструментированы (list, get, search, upsert, delete)
- **3-4 метрики** на сервис
- **Structured logging** с zap
- **Production-ready** для Grafana dashboards
