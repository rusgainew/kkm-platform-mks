# 📋 План доработки: Backend API → Frontend Integration

**Дата:** 27 января 2026 г.  
**Версия:** 2.0  
**Подход:** Backend-First Development

---

## 🎯 Цель

Реализовать бизнес-логику в api-gateway (Go бэкенд), затем интегрировать её в kkm-platform (Next.js фронтенд). Каждая фича разрабатывается полностью на бэкенде, тестируется, и только потом внедряется во фронтенд.

## 📐 Архитектурный подход

```
┌─────────────────────────────────────────────────────┐
│  ЭТАП 1: Backend Implementation (Go)                │
│  ├─ Business Logic (Services Layer)                 │
│  ├─ API Handlers (HTTP Layer)                       │
│  ├─ Infrastructure (Cache, Metrics, etc)            │
│  └─ Testing (Unit + Integration)                    │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  ЭТАП 2: API Contract (OpenAPI/Swagger)             │
│  ├─ Swagger Documentation                           │
│  ├─ Request/Response Types                          │
│  └─ API Testing (Postman/curl)                      │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  ЭТАП 3: Frontend Integration (TypeScript/React)    │
│  ├─ API Client Generation                           │
│  ├─ React Hooks                                     │
│  ├─ UI Components                                   │
│  └─ E2E Testing                                     │
└─────────────────────────────────────────────────────┘
```

---

## 📊 Текущее состояние

### ✅ Что уже есть в kkm-platform

- [x] JWT авторизация (токены + refresh)
- [x] Базовое кэширование (memory, localStorage, IndexedDB)
- [x] Performance оптимизации (code splitting, lazy loading)
- [x] Audit логирование (clientside)
- [x] API клиенты (invoice, catalog, companies)
- [x] Интернационализация (i18n middleware)
- [x] Zustand state management
- [x] TypeScript типизация

### ❌ Чего не хватает (есть в api-gateway)

- [ ] **JWT Key Rotation** - ротация ключей для токенов
- [ ] **Rate Limiting** - контроль частоты запросов
- [ ] **Distributed Tracing** - OpenTelemetry/Jaeger интеграция
- [ ] **Prometheus Metrics** - экспорт метрик производительности
- [ ] **Request/Response Logging** - детальное логирование HTTP
- [ ] **Circuit Breaker** - паттерн для отказоустойчивости
- [ ] **Health Checks Dashboard** - мониторинг статуса backend
- [ ] **Graceful Shutdown** - корректное завершение при деплое
- [ ] **Redis Cache Integration** - централизованный кэш
- [ ] **API Versioning** - версионирование API endpoints
- [ ] **Swagger/OpenAPI Client** - авто-генерация типов из OpenAPI

---

## 🚀 Приоритетный план доработки

### **Фаза 1: Безопасность и надёжность (Priority: HIGH)** 🔥

#### 1.1. Rate Limiting для API запросов

**Проблема:**

- Нет защиты от DoS атак
- Можно отправить неограниченное кол-во запросов

**Решение:**

```typescript
// lib/middleware/rateLimiter.ts
import { RateLimiter } from "@/lib/rate-limiter";

const limiter = new RateLimiter({
  maxRequests: 100,
  windowMs: 60 * 1000, // 100 запросов в минуту
  strategy: "sliding-window",
});

// Интеграция в API клиент
export async function apiCallWithRateLimit<T>(
  fn: () => Promise<T>,
): Promise<T> {
  await limiter.checkLimit();
  return fn();
}
```

**Файлы:**

- `lib/middleware/rateLimiter.ts` _(новый)_
- `lib/api/client.ts` _(обновить)_

**Время:** 4-6 часов

---

#### 1.2. JWT Token Validation Enhancement

**Проблема:**

- Нет проверки на blacklist (отозванные токены)
- Нет валидации key rotation

**Решение:**

```typescript
// lib/auth/tokenValidator.ts
interface TokenValidationConfig {
  checkBlacklist: boolean;
  supportKeyRotation: boolean;
  maxTokenAge: number;
}

export class TokenValidator {
  async validateToken(token: string): Promise<ValidationResult> {
    // 1. Базовая JWT валидация
    const decoded = jwt.verify(token, publicKey);

    // 2. Проверка blacklist
    if (await this.isBlacklisted(token)) {
      throw new Error("Token revoked");
    }

    // 3. Проверка ротации ключей
    if (this.isOldKey(decoded.kid)) {
      // Token с старым ключом - требуется refresh
      return { valid: true, requiresRefresh: true };
    }

    return { valid: true };
  }
}
```

**Файлы:**

- `lib/auth/tokenValidator.ts` _(новый)_
- `lib/auth/tokenBlacklist.ts` _(новый)_
- `store/authStore.ts` _(обновить)_

**Время:** 6-8 часов

---

#### 1.3. Request/Response Logging Infrastructure

**Проблема:**

- Нет централизованного логирования API запросов
- Сложно дебажить проблемы в production

**Решение:**

```typescript
// lib/logging/apiLogger.ts
interface APILog {
  requestId: string;
  method: string;
  url: string;
  status: number;
  duration: number;
  userId?: string;
  error?: string;
  timestamp: number;
}

export class APILogger {
  async logRequest(config: RequestConfig): Promise<void> {
    const log: APILog = {
      requestId: generateId(),
      method: config.method,
      url: config.url,
      timestamp: Date.now(),
      userId: getCurrentUserId(),
    };

    // Отправка в backend через API Gateway
    await fetch("/api/v1/logs", {
      method: "POST",
      body: JSON.stringify(log),
    });
  }
}

// Интеграция с Axios
axios.interceptors.request.use(async (config) => {
  config.metadata = { startTime: Date.now() };
  await apiLogger.logRequest(config);
  return config;
});

axios.interceptors.response.use(
  async (response) => {
    const duration = Date.now() - response.config.metadata.startTime;
    await apiLogger.logResponse(response, duration);
    return response;
  },
  async (error) => {
    await apiLogger.logError(error);
    throw error;
  },
);
```

**Файлы:**

- `lib/logging/apiLogger.ts` _(новый)_
- `lib/api/client.ts` _(обновить)_
- `components/admin/APILogsViewer.tsx` _(новый)_

**Время:** 8-10 часов

---

### **Фаза 2: Мониторинг и наблюдаемость (Priority: HIGH)** 📊

#### 2.1. Health Checks Dashboard

**Проблема:**

- Нет визуализации статуса backend сервисов
- Пользователь не знает когда сервис недоступен

**Решение:**

```typescript
// lib/health/healthChecker.ts
interface ServiceHealth {
  name: string;
  status: 'healthy' | 'degraded' | 'down';
  responseTime: number;
  lastCheck: number;
  error?: string;
}

export class HealthChecker {
  private services = [
    'api-gateway',
    'user-service',
    'invoice-service',
    'catalog-service'
  ];

  async checkAll(): Promise<ServiceHealth[]> {
    return Promise.all(
      this.services.map(service => this.checkService(service))
    );
  }

  private async checkService(name: string): Promise<ServiceHealth> {
    const startTime = Date.now();
    try {
      const response = await fetch(`/api/v1/health/${name}`);
      return {
        name,
        status: response.ok ? 'healthy' : 'degraded',
        responseTime: Date.now() - startTime,
        lastCheck: Date.now()
      };
    } catch (error) {
      return {
        name,
        status: 'down',
        responseTime: Date.now() - startTime,
        lastCheck: Date.now(),
        error: error.message
      };
    }
  }
}

// React компонент
export function HealthDashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ['health'],
    queryFn: () => healthChecker.checkAll(),
    refetchInterval: 30000 // Проверка каждые 30 сек
  });

  return (
    <div className="grid grid-cols-2 gap-4">
      {data?.map(service => (
        <ServiceHealthCard key={service.name} service={service} />
      ))}
    </div>
  );
}
```

**Файлы:**

- `lib/health/healthChecker.ts` _(новый)_
- `components/admin/HealthDashboard.tsx` _(новый)_
- `app/admin/health/page.tsx` _(новый)_

**Время:** 6-8 часов

---

#### 2.2. Client-side Metrics Collection (Prometheus-compatible)

**Проблема:**

- Нет сбора метрик производительности frontend
- Невозможно отследить проблемы пользователей

**Решение:**

```typescript
// lib/metrics/metricsCollector.ts
interface Metric {
  name: string;
  value: number;
  labels: Record<string, string>;
  timestamp: number;
}

export class MetricsCollector {
  private metrics: Metric[] = [];

  // HTTP метрики
  recordHTTPRequest(
    method: string,
    url: string,
    status: number,
    duration: number,
  ) {
    this.metrics.push({
      name: "http_requests_total",
      value: 1,
      labels: { method, url, status: String(status) },
      timestamp: Date.now(),
    });

    this.metrics.push({
      name: "http_request_duration_seconds",
      value: duration / 1000,
      labels: { method, url },
      timestamp: Date.now(),
    });
  }

  // Метрики UI
  recordPageLoad(page: string, duration: number) {
    this.metrics.push({
      name: "page_load_duration_seconds",
      value: duration / 1000,
      labels: { page },
      timestamp: Date.now(),
    });
  }

  // Периодическая отправка в API Gateway
  async flush() {
    if (this.metrics.length === 0) return;

    await fetch("/api/v1/metrics/client", {
      method: "POST",
      body: JSON.stringify(this.metrics),
    });

    this.metrics = [];
  }
}

// Web Vitals интеграция
import { onCLS, onFID, onLCP } from "web-vitals";

onCLS((metric) => metricsCollector.recordWebVital("CLS", metric.value));
onFID((metric) => metricsCollector.recordWebVital("FID", metric.value));
onLCP((metric) => metricsCollector.recordWebVital("LCP", metric.value));
```

**Файлы:**

- `lib/metrics/metricsCollector.ts` _(новый)_
- `lib/metrics/webVitals.ts` _(новый)_
- `app/layout.tsx` _(обновить - добавить MetricsProvider)_

**Время:** 8-10 часов

---

#### 2.3. Distributed Tracing (OpenTelemetry)

**Проблема:**

- Невозможно отследить цепочку вызовов frontend → api-gateway → backend
- Сложно найти узкие места в распределённой системе

**Решение:**

```typescript
// lib/tracing/tracer.ts
import { WebTracerProvider } from "@opentelemetry/sdk-trace-web";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";

export class FrontendTracer {
  private provider: WebTracerProvider;

  init() {
    this.provider = new WebTracerProvider({
      resource: new Resource({
        "service.name": "kkm-platform-frontend",
        "service.version": "1.0.0",
      }),
    });

    const exporter = new OTLPTraceExporter({
      url: "/api/v1/traces", // Прокси через API Gateway
    });

    this.provider.addSpanProcessor(new BatchSpanProcessor(exporter));

    this.provider.register();
  }

  startSpan(name: string, attributes?: Record<string, string>) {
    const tracer = this.provider.getTracer("kkm-platform");
    return tracer.startSpan(name, {
      attributes: {
        "user.id": getCurrentUserId(),
        ...attributes,
      },
    });
  }
}

// Использование
const span = tracer.startSpan("api.createInvoice", {
  "invoice.type": "ESF",
  "invoice.amount": "10000",
});

try {
  const result = await invoiceAPI.createInvoice(data);
  span.setStatus({ code: SpanStatusCode.OK });
  return result;
} catch (error) {
  span.setStatus({
    code: SpanStatusCode.ERROR,
    message: error.message,
  });
  throw error;
} finally {
  span.end();
}
```

**Файлы:**

- `lib/tracing/tracer.ts` _(новый)_
- `lib/tracing/spanProcessor.ts` _(новый)_
- `lib/api/client.ts` _(обновить - добавить трейсинг)_

**Время:** 10-12 часов

---

### **Фаза 3: Производительность и масштабирование (Priority: MEDIUM)** ⚡

#### 3.1. Redis Cache Integration

**Проблема:**

- Кэш не синхронизирован между пользователями
- Нет инвалидации кэша при изменениях

**Решение:**

```typescript
// lib/cache/redisCache.ts
interface CacheConfig {
  enabled: boolean;
  ttl: number;
  prefix: string;
}

export class RedisCacheClient {
  constructor(private apiUrl: string) {}

  async get<T>(key: string): Promise<T | null> {
    try {
      const response = await fetch(`/api/v1/cache/${this.buildKey(key)}`);

      if (!response.ok) return null;

      const data = await response.json();
      return data.value as T;
    } catch (error) {
      console.error("Cache get error:", error);
      return null;
    }
  }

  async set<T>(key: string, value: T, ttl?: number): Promise<void> {
    await fetch(`/api/v1/cache/${this.buildKey(key)}`, {
      method: "PUT",
      body: JSON.stringify({ value, ttl }),
    });
  }

  async delete(key: string): Promise<void> {
    await fetch(`/api/v1/cache/${this.buildKey(key)}`, {
      method: "DELETE",
    });
  }

  private buildKey(key: string): string {
    return `kkm-platform:${key}`;
  }
}

// React Query интеграция
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 минут
      cacheTime: 10 * 60 * 1000, // 10 минут

      // Кастомный кэш адаптер
      queryCache: new QueryCache({
        onSuccess: async (data, query) => {
          // Сохраняем в Redis через API Gateway
          await redisCache.set(query.queryKey, data);
        },
      }),
    },
  },
});
```

**Файлы:**

- `lib/cache/redisCache.ts` _(новый)_
- `lib/providers/QueryProvider.tsx` _(обновить)_
- `lib/api/client.ts` _(обновить)_

**Время:** 8-10 часов

---

#### 3.2. API Response Caching Strategy

**Проблема:**

- Нет умного кэширования ответов API
- Повторные запросы за одними и теми же данными

**Решение:**

```typescript
// lib/cache/apiCacheStrategy.ts
interface CacheStrategy {
  key: string;
  ttl: number;
  invalidateOn: string[]; // События для инвалидации
  staleWhileRevalidate?: boolean;
}

const strategies: Record<string, CacheStrategy> = {
  "invoices.list": {
    key: "invoices:list:{userId}:{filters}",
    ttl: 5 * 60, // 5 минут
    invalidateOn: ["invoice.created", "invoice.updated"],
    staleWhileRevalidate: true,
  },
  "companies.detail": {
    key: "company:{id}",
    ttl: 10 * 60, // 10 минут
    invalidateOn: ["company.updated"],
    staleWhileRevalidate: true,
  },
  "catalog.items": {
    key: "catalog:items:{categoryId}",
    ttl: 15 * 60, // 15 минут
    invalidateOn: ["catalog.updated"],
    staleWhileRevalidate: false,
  },
};

export class APICacheManager {
  async getWithStrategy<T>(
    strategyKey: string,
    fetcher: () => Promise<T>,
    params?: Record<string, string>,
  ): Promise<T> {
    const strategy = strategies[strategyKey];
    const cacheKey = this.buildKey(strategy.key, params);

    // Проверяем кэш
    const cached = await redisCache.get<T>(cacheKey);

    if (cached && !strategy.staleWhileRevalidate) {
      return cached;
    }

    // Stale-while-revalidate
    if (cached && strategy.staleWhileRevalidate) {
      // Возвращаем старые данные, но обновляем в фоне
      this.revalidateInBackground(cacheKey, fetcher, strategy.ttl);
      return cached;
    }

    // Нет в кэше - запрашиваем
    const fresh = await fetcher();
    await redisCache.set(cacheKey, fresh, strategy.ttl);
    return fresh;
  }

  // Инвалидация при событиях
  async invalidateOnEvent(event: string, params?: Record<string, string>) {
    for (const [key, strategy] of Object.entries(strategies)) {
      if (strategy.invalidateOn.includes(event)) {
        const cacheKey = this.buildKey(strategy.key, params);
        await redisCache.delete(cacheKey);
      }
    }
  }
}
```

**Файлы:**

- `lib/cache/apiCacheStrategy.ts` _(новый)_
- `lib/cache/cacheInvalidation.ts` _(новый)_
- `hooks/useInvoice.ts` _(обновить)_
- `hooks/useCompany.ts` _(обновить)_

**Время:** 10-12 часов

---

#### 3.3. Circuit Breaker Pattern

**Проблема:**

- При падении backend продолжаем слать запросы
- Пользователь видит бесконечный спиннер

**Решение:**

```typescript
// lib/patterns/circuitBreaker.ts
enum CircuitState {
  CLOSED = "closed", // Всё работает
  OPEN = "open", // Сервис недоступен
  HALF_OPEN = "half_open", // Проверяем восстановление
}

export class CircuitBreaker {
  private state: CircuitState = CircuitState.CLOSED;
  private failures = 0;
  private successCount = 0;
  private nextAttempt = 0;

  constructor(
    private threshold: number = 5, // Ошибок до OPEN
    private timeout: number = 60000, // Время в OPEN (1 минута)
    private successThreshold: number = 2, // Успехов для закрытия
  ) {}

  async execute<T>(fn: () => Promise<T>): Promise<T> {
    // Проверяем состояние
    if (this.state === CircuitState.OPEN) {
      if (Date.now() < this.nextAttempt) {
        throw new Error("Circuit breaker is OPEN");
      }
      // Переходим в HALF_OPEN
      this.state = CircuitState.HALF_OPEN;
    }

    try {
      const result = await fn();
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure();
      throw error;
    }
  }

  private onSuccess() {
    this.failures = 0;

    if (this.state === CircuitState.HALF_OPEN) {
      this.successCount++;
      if (this.successCount >= this.successThreshold) {
        this.state = CircuitState.CLOSED;
        this.successCount = 0;
      }
    }
  }

  private onFailure() {
    this.failures++;
    this.successCount = 0;

    if (this.failures >= this.threshold) {
      this.state = CircuitState.OPEN;
      this.nextAttempt = Date.now() + this.timeout;
    }
  }

  getState() {
    return this.state;
  }
}

// Интеграция с API клиентом
const breakers = {
  invoiceService: new CircuitBreaker(5, 60000, 2),
  catalogService: new CircuitBreaker(5, 60000, 2),
  companyService: new CircuitBreaker(5, 60000, 2),
};

export async function apiCallWithCircuitBreaker<T>(
  service: keyof typeof breakers,
  fn: () => Promise<T>,
): Promise<T> {
  return breakers[service].execute(fn);
}
```

**Файлы:**

- `lib/patterns/circuitBreaker.ts` _(новый)_
- `lib/api/client.ts` _(обновить)_
- `components/fallbacks/ServiceUnavailable.tsx` _(новый)_

**Время:** 6-8 часов

---

### **Фаза 4: Developer Experience (Priority: MEDIUM)** 🛠️

#### 4.1. OpenAPI/Swagger Client Generator

**Проблема:**

- Типы API вручную синхронизируются с backend
- Риск расхождения между frontend и backend контрактами

**Решение:**

```bash
# package.json scripts
"generate:api": "openapi-generator-cli generate -i http://localhost:8080/docs/swagger.json -g typescript-axios -o lib/generated/api"
"generate:types": "openapi-typescript http://localhost:8080/docs/swagger.json --output types/api.d.ts"
```

```typescript
// scripts/generateAPIClient.ts
import SwaggerParser from "@apidevtools/swagger-parser";
import { generateApi } from "swagger-typescript-api";

async function generate() {
  const spec = await SwaggerParser.validate(
    "http://localhost:8080/docs/swagger.json",
  );

  await generateApi({
    name: "APIClient.ts",
    output: path.resolve(process.cwd(), "lib/generated"),
    spec,
    httpClientType: "axios",
    generateClient: true,
    generateRouteTypes: true,
    extractRequestParams: true,
    extractRequestBody: true,
    extractResponseBody: true,
  });

  console.log("✅ API client generated successfully");
}

generate();
```

**Файлы:**

- `scripts/generateAPIClient.ts` _(новый)_
- `lib/generated/APIClient.ts` _(генерируется)_
- `types/api.d.ts` _(генерируется)_
- `package.json` _(обновить)_

**Время:** 4-6 часов

---

#### 4.2. API Versioning Support

**Проблема:**

- Нет поддержки нескольких версий API
- Сложно мигрировать между версиями

**Решение:**

```typescript
// lib/api/versionedClient.ts
enum APIVersion {
  V1 = "v1",
  V2 = "v2",
}

export class VersionedAPIClient {
  constructor(
    private baseURL: string,
    private defaultVersion: APIVersion = APIVersion.V1,
  ) {}

  getClient(version?: APIVersion) {
    const ver = version || this.defaultVersion;
    return axios.create({
      baseURL: `${this.baseURL}/api/${ver}`,
      headers: {
        "X-API-Version": ver,
      },
    });
  }

  // Миграция между версиями
  async migrateRequest<V1, V2>(
    v1Request: V1,
    transformer: (v1: V1) => V2,
  ): Promise<V2> {
    return transformer(v1Request);
  }
}

// Использование
const client = new VersionedAPIClient("http://localhost:8080");

// V1 запрос
const v1Data = await client.getClient(APIVersion.V1).get("/invoices");

// V2 запрос (новые поля)
const v2Data = await client.getClient(APIVersion.V2).get("/invoices");
```

**Файлы:**

- `lib/api/versionedClient.ts` _(новый)_
- `lib/api/migrations/` _(новая папка для миграций)_

**Время:** 4-6 часов

---

#### 4.3. Development Mock Server

**Проблема:**

- Зависимость от запущенного backend при разработке
- Сложно тестировать edge cases

**Решение:**

```typescript
// lib/mocks/mockServer.ts (MSW - Mock Service Worker)
import { setupWorker, rest } from "msw";

const handlers = [
  // Моки endpoints из swagger
  rest.get("/api/v1/invoices", (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        items: mockInvoices,
        total: 100,
        page: 1,
      }),
    );
  }),

  // Симуляция ошибок
  rest.get("/api/v1/companies/:id", (req, res, ctx) => {
    if (req.params.id === "error") {
      return res(ctx.status(500), ctx.json({ error: "Internal Server Error" }));
    }
    return res(ctx.json(mockCompany));
  }),

  // Симуляция задержек
  rest.get("/api/v1/catalog", async (req, res, ctx) => {
    await delay(2000); // 2 секунды задержки
    return res(ctx.json(mockCatalog));
  }),
];

export const worker = setupWorker(...handlers);

// В development режиме
if (
  process.env.NODE_ENV === "development" &&
  process.env.NEXT_PUBLIC_ENABLE_MOCKS === "true"
) {
  worker.start();
}
```

**Файлы:**

- `lib/mocks/mockServer.ts` _(новый)_
- `lib/mocks/handlers/` _(новая папка)_
- `lib/mocks/data/` _(моковые данные)_

**Время:** 6-8 часов

---

### **Фаза 5: DevOps и мониторинг (Priority: LOW)** 🔧

#### 5.1. Frontend Error Boundary с отчётностью

**Проблема:**

- Ошибки React не отслеживаются
- Нет автоматических отчётов об ошибках

**Решение:**

```typescript
// components/errors/ErrorBoundary.tsx
import { Component, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: React.ErrorInfo) => void;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // Логируем в backend
    this.reportError(error, errorInfo);

    // Кастомный обработчик
    this.props.onError?.(error, errorInfo);
  }

  async reportError(error: Error, errorInfo: React.ErrorInfo) {
    await fetch('/api/v1/errors/frontend', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: error.message,
        stack: error.stack,
        componentStack: errorInfo.componentStack,
        userId: getCurrentUserId(),
        url: window.location.href,
        timestamp: Date.now(),
        userAgent: navigator.userAgent,
      })
    });
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <ErrorFallback error={this.state.error!} />
      );
    }

    return this.props.children;
  }
}
```

**Файлы:**

- `components/errors/ErrorBoundary.tsx` _(новый)_
- `components/errors/ErrorFallback.tsx` _(новый)_
- `app/layout.tsx` _(обновить)_

**Время:** 4-6 часов

---

#### 5.2. Build-time Health Checks

**Проблема:**

- Деплой может сломать приложение
- Нет проверок перед production

**Решение:**

```typescript
// scripts/preDeployChecks.ts
import { exec } from "child_process";
import { promisify } from "util";

const execAsync = promisify(exec);

interface CheckResult {
  name: string;
  passed: boolean;
  message: string;
}

async function runPreDeployChecks(): Promise<CheckResult[]> {
  const checks: CheckResult[] = [];

  // 1. TypeScript проверка
  try {
    await execAsync("tsc --noEmit");
    checks.push({
      name: "TypeScript",
      passed: true,
      message: "No type errors",
    });
  } catch (error) {
    checks.push({
      name: "TypeScript",
      passed: false,
      message: error.message,
    });
  }

  // 2. Lint
  try {
    await execAsync("eslint . --max-warnings 0");
    checks.push({
      name: "ESLint",
      passed: true,
      message: "No linting errors",
    });
  } catch (error) {
    checks.push({
      name: "ESLint",
      passed: false,
      message: error.message,
    });
  }

  // 3. Unit тесты
  try {
    await execAsync("vitest run");
    checks.push({
      name: "Unit Tests",
      passed: true,
      message: "All tests passed",
    });
  } catch (error) {
    checks.push({
      name: "Unit Tests",
      passed: false,
      message: error.message,
    });
  }

  // 4. Build проверка
  try {
    await execAsync("next build");
    checks.push({
      name: "Build",
      passed: true,
      message: "Build successful",
    });
  } catch (error) {
    checks.push({
      name: "Build",
      passed: false,
      message: error.message,
    });
  }

  // 5. API Gateway доступность
  try {
    const response = await fetch("http://localhost:8080/health");
    checks.push({
      name: "API Gateway",
      passed: response.ok,
      message: response.ok ? "Available" : "Unavailable",
    });
  } catch (error) {
    checks.push({
      name: "API Gateway",
      passed: false,
      message: "Not accessible",
    });
  }

  return checks;
}

async function main() {
  console.log("🔍 Running pre-deploy checks...\n");

  const results = await runPreDeployChecks();

  results.forEach((result) => {
    const icon = result.passed ? "✅" : "❌";
    console.log(`${icon} ${result.name}: ${result.message}`);
  });

  const failed = results.filter((r) => !r.passed);

  if (failed.length > 0) {
    console.error(`\n❌ ${failed.length} checks failed. Deploy blocked.`);
    process.exit(1);
  }

  console.log("\n✅ All checks passed. Ready to deploy!");
}

main();
```

**Файлы:**

- `scripts/preDeployChecks.ts` _(новый)_
- `package.json` _(обновить скрипт predeploy)_

**Время:** 4-6 часов

---

## 📅 Временная оценка

### Итого по фазам:

| Фаза                             | Приоритет | Время            | Сложность |
| -------------------------------- | --------- | ---------------- | --------- |
| **Фаза 1: Безопасность**         | HIGH 🔥   | 18-24 часа       | Средняя   |
| **Фаза 2: Мониторинг**           | HIGH 📊   | 24-30 часов      | Высокая   |
| **Фаза 3: Производительность**   | MEDIUM ⚡ | 24-30 часов      | Средняя   |
| **Фаза 4: Developer Experience** | MEDIUM 🛠️ | 14-20 часов      | Низкая    |
| **Фаза 5: DevOps**               | LOW 🔧    | 8-12 часов       | Низкая    |
| **ИТОГО**                        | -         | **88-116 часов** | -         |

---

## 🎯 Рекомендуемая последовательность

### Неделя 1 (Фаза 1): Безопасность

1. Rate Limiting (день 1-2)
2. JWT Token Validation (день 3-4)
3. Request/Response Logging (день 5)

### Неделя 2 (Фаза 2, часть 1): Мониторинг

1. Health Checks Dashboard (день 1-2)
2. Client-side Metrics (день 3-4)

### Неделя 3 (Фаза 2, часть 2 + Фаза 3, часть 1): Трейсинг + Кэш

1. Distributed Tracing (день 1-3)
2. Redis Cache Integration (день 4-5)

### Неделя 4 (Фаза 3, часть 2): Производительность

1. API Cache Strategy (день 1-3)
2. Circuit Breaker (день 4-5)

### Неделя 5 (Фаза 4): DevEx

1. OpenAPI Generator (день 1-2)
2. API Versioning (день 3)
3. Mock Server (день 4-5)

### Неделя 6 (Фаза 5): DevOps

1. Error Boundary (день 1-2)
2. Pre-deploy Checks (день 3)
3. Тестирование и документация (день 4-5)

---

## 📦 Необходимые зависимости

```json
{
  "dependencies": {
    "@opentelemetry/api": "^1.7.0",
    "@opentelemetry/sdk-trace-web": "^1.19.0",
    "@opentelemetry/exporter-trace-otlp-http": "^0.47.0",
    "@opentelemetry/instrumentation": "^0.47.0",
    "web-vitals": "^3.5.0"
  },
  "devDependencies": {
    "@apidevtools/swagger-parser": "^10.1.0",
    "swagger-typescript-api": "^13.0.3",
    "openapi-typescript": "^6.7.3",
    "msw": "^2.0.0",
    "@faker-js/faker": "^8.3.1"
  }
}
```

---

## ✅ Критерии успеха

### Безопасность

- [ ] Rate limiting работает для всех API endpoints
- [ ] JWT токены валидируются с blacklist
- [ ] Все API запросы логируются

### Мониторинг

- [ ] Dashboard показывает статус всех сервисов
- [ ] Метрики отправляются в API Gateway
- [ ] Трейсы видны в Jaeger UI

### Производительность

- [ ] Redis кэш работает с TTL и инвалидацией
- [ ] Circuit breaker защищает от падений backend
- [ ] API cache hit rate > 70%

### DevEx

- [ ] Типы генерируются из Swagger спецификации
- [ ] Mock server работает в development
- [ ] API versioning поддерживает v1 и v2

### DevOps

- [ ] Error boundary ловит все React ошибки
- [ ] Pre-deploy checks блокируют плохой код
- [ ] Build проходит без warnings

---

## 🔄 Continuous Improvement

После завершения всех фаз рекомендуется:

1. **Настроить мониторинг метрик**
   - Prometheus + Grafana dashboard
   - Алерты на критические метрики

2. **Провести нагрузочное тестирование**
   - k6 скрипты для API endpoints
   - Проверка Circuit Breaker в action

3. **Документация**
   - Обновить README с новыми возможностями
   - Создать runbooks для операционных задач

4. **Обучение команды**
   - Сессии по новым инструментам
   - Code review guidelines

---

## 📚 Дополнительные ресурсы

- [OpenTelemetry Docs](https://opentelemetry.io/docs/)
- [MSW Documentation](https://mswjs.io/)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
- [Web Vitals](https://web.dev/vitals/)
- [API Gateway Swagger](http://localhost:8080/api/v1/docs/)

---

**Конец плана доработки**
