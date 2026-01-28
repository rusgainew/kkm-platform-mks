# Backend Analytics для Dashboard

## 📊 Обзор

Backend API для Dashboard аналитики, предоставляющий агрегированные данные о счетах-фактурах, продажах и контрагентах.

## 🔗 API Эндпоинты

Все эндпоинты требуют JWT аутентификации и доступны по префиксу `/api/v1/analytics`

### 1. Статистика Dashboard

```http
GET /api/v1/analytics/stats
```

**Query Parameters:**

- `period` (optional): `today`, `week`, `month`, `quarter`, `year`, `custom` (default: `month`)
- `startDate` (optional): Начальная дата для custom период (YYYY-MM-DD)
- `endDate` (optional): Конечная дата для custom период (YYYY-MM-DD)

**Response:**

```json
{
  "success": true,
  "data": {
    "totalRevenue": 1500000.5,
    "totalInvoices": 245,
    "averageInvoiceAmount": 6122.45,
    "activeContractors": 38,
    "pendingInvoices": 12,
    "approvedInvoices": 220,
    "rejectedInvoices": 13,
    "revenueChange": 15.7,
    "invoiceCountChange": 8.3,
    "averageAmountChange": 3.2,
    "contractorsChange": 5.5,
    "period": {
      "startDate": "2024-01-01",
      "endDate": "2024-01-31"
    }
  }
}
```

### 2. График продаж

```http
GET /api/v1/analytics/sales-chart
```

**Query Parameters:**

- `period`, `startDate`, `endDate` (как выше)

**Response:**

```json
{
  "success": true,
  "data": {
    "data": [
      {
        "date": "2024-01-01",
        "value": 45000
      },
      {
        "date": "2024-01-02",
        "value": 46500
      }
    ],
    "period": {
      "startDate": "2024-01-01",
      "endDate": "2024-01-31"
    }
  }
}
```

### 3. Статистика по статусам

```http
GET /api/v1/analytics/status-stats
```

**Response:**

```json
{
  "success": true,
  "data": {
    "data": [
      {
        "name": "Утверждено",
        "value": 220,
        "percentage": 89.8,
        "color": "#10b981"
      },
      {
        "name": "На рассмотрении",
        "value": 12,
        "percentage": 4.9,
        "color": "#f59e0b"
      }
    ]
  }
}
```

### 4. Статистика по типам операций

```http
GET /api/v1/analytics/operation-stats
```

**Response:**

```json
{
  "success": true,
  "data": {
    "data": [
      {
        "name": "Продажа",
        "value": 180,
        "percentage": 73.5,
        "color": "#3b82f6"
      }
    ]
  }
}
```

### 5. Топ контрагенты

```http
GET /api/v1/analytics/top-contractors
```

**Query Parameters:**

- `period`, `startDate`, `endDate` (как выше)
- `limit` (optional): Количество контрагентов (default: 10)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "contractorId": "1",
      "contractorName": "ТОО \"Рога и копыта\"",
      "totalAmount": 450000,
      "invoiceCount": 45
    }
  ]
}
```

### 6. Выручка по месяцам

```http
GET /api/v1/analytics/revenue-by-month
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "date": "Январь",
      "value": 120000
    }
  ]
}
```

## 🏗️ Архитектура

### Текущая реализация (Mock Data)

```
┌─────────────┐
│   Frontend  │
│  (Next.js)  │
└──────┬──────┘
       │ HTTP/JSON
       ▼
┌─────────────────────┐
│   API Gateway       │
│ Analytics Handler   │
│   (Mock Data)       │
└─────────────────────┘
```

### Будущая реализация (с БД)

```
┌─────────────┐
│   Frontend  │
└──────┬──────┘
       │
       ▼
┌─────────────────────┐
│   API Gateway       │
│ Analytics Handler   │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│ Analytics Service   │
│   (Repository)      │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│    PostgreSQL       │
│  (Aggregations)     │
└──────┬──────────────┘
       │ Cache
       ▼
┌─────────────────────┐
│      Redis          │
│  (5 min TTL)        │
└─────────────────────┘
```

## 📝 Файлы

### API Gateway

- `services/api-gateway/internal/interfaces/http/analytics_handler.go` (470 строк)
  - Все 6 эндпоинтов
  - Вычисление диапазонов дат
  - Mock данные для тестирования

- `services/api-gateway/internal/interfaces/http/analytics_handler_test.go` (243 строки)
  - 8 unit тестов
  - Покрытие всех эндпоинтов
  - Тесты calculateDateRange

- `services/api-gateway/internal/interfaces/http/router_configurator.go`
  - Роуты `/api/v1/analytics/*`
  - JWT middleware
  - Rate limiting

## 🔧 Установка и запуск

### Локальный запуск

```bash
# Перейти в API Gateway
cd services/api-gateway

# Установить зависимости
go mod download

# Запустить сервер
go run main.go
```

### Docker

```bash
# Собрать образ
docker-compose build api-gateway

# Запустить
docker-compose up -d api-gateway
```

## 🧪 Тестирование

### Unit тесты

```bash
cd services/api-gateway
go test -v ./internal/interfaces/http -run TestAnalytics
```

### Manual тестирование с curl

```bash
# Получить JWT токен
TOKEN=$(curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' \
  | jq -r '.data.token')

# Получить статистику
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/stats?period=month"

# График продаж
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/sales-chart?period=week"

# Топ контрагенты
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/analytics/top-contractors?limit=5"
```

## 📊 Следующие шаги

### 1. Database Integration (Priority: High)

Создать SQL queries для агрегации данных:

```sql
-- Total revenue and invoice count
SELECT
  SUM(total_amount) as total_revenue,
  COUNT(*) as total_invoices,
  AVG(total_amount) as average_amount
FROM invoices
WHERE created_date BETWEEN $1 AND $2
  AND status = 'APPROVED';

-- Revenue by date (for charts)
SELECT
  DATE(created_date) as date,
  SUM(total_amount) as value
FROM invoices
WHERE created_date BETWEEN $1 AND $2
GROUP BY DATE(created_date)
ORDER BY date;

-- Top contractors
SELECT
  contractor_tin as contractor_id,
  contractor_name,
  SUM(total_amount) as total_amount,
  COUNT(*) as invoice_count
FROM invoices
WHERE created_date BETWEEN $1 AND $2
GROUP BY contractor_tin, contractor_name
ORDER BY total_amount DESC
LIMIT $3;

-- Status distribution
SELECT
  status as name,
  COUNT(*) as value,
  ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage
FROM invoices
WHERE created_date BETWEEN $1 AND $2
GROUP BY status;
```

### 2. Redis Caching (Priority: High)

```go
// Пример кеширования
func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
    period := c.DefaultQuery("period", "month")
    cacheKey := fmt.Sprintf("analytics:stats:%s", period)

    // Try cache first
    var stats DashboardStats
    if err := h.cache.Get(ctx, cacheKey, &stats); err == nil {
        c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
        return
    }

    // Query database
    stats = h.queryDashboardStats(period)

    // Cache for 5 minutes
    h.cache.Set(ctx, cacheKey, stats, 5*time.Minute)

    c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}
```

### 3. Performance Optimization

- **Materialized Views** для предрассчитанных агрегаций
- **Batch queries** для уменьшения количества запросов
- **Connection pooling** для PostgreSQL
- **Query timeout** (5 секунд максимум)

### 4. Мониторинг

- Prometheus метрики для каждого эндпоинта
- Отслеживание времени ответа
- Alert на медленные запросы (> 3 сек)

## 🔐 Безопасность

- ✅ JWT аутентификация обязательна
- ✅ Rate limiting через middleware
- ✅ SQL injection protection (подготовленные запросы)
- ✅ Input validation для дат

## 📈 Performance Requirements

- Время ответа: < 500ms (95 percentile)
- Cache hit rate: > 80%
- Database queries: < 3 per request
- Concurrent users: 100+

## 🎯 Текущий статус

- ✅ API эндпоинты созданы (6 штук)
- ✅ Mock данные для тестирования
- ✅ Unit тесты (8 тестов)
- ✅ JWT аутентификация
- ✅ Swagger документация
- ⏳ Database integration (TODO)
- ⏳ Redis caching (TODO)
- ⏳ Real aggregations (TODO)

## 🤝 Интеграция с Frontend

Frontend использует TypeScript API client:

```typescript
import { analyticsAPI } from "@/lib/api/analytics";

// Получить статистику
const stats = await analyticsAPI.getStats({ period: "month" });

// График продаж
const chart = await analyticsAPI.getSalesChart({
  period: "custom",
  startDate: "2024-01-01",
  endDate: "2024-01-31",
});
```

Полная документация frontend: `/kkm-platform/PHASE_8_DASHBOARD_COMPLETE.md`

## 📞 Контакты

Для вопросов и предложений создавайте issue в репозитории.
