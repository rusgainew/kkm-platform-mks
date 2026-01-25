# Performance Optimization Guide

## 📊 Создано для Session 5 - Performance Phase

### Tier 3 Optimization Features Implemented

#### 1. **Code Splitting & Dynamic Imports**

```typescript
// lib/dynamic.ts - Динамическая загрузка компонентов
import { DynamicCharts, DynamicAdvancedFilter } from "@/lib/dynamic";

// Компоненты загружаются при необходимости
<DynamicCharts />; // Загрузится только когда видим на экране
```

#### 2. **Performance Utilities**

```typescript
// lib/performance.ts - 8 optimization techniques

// Debounce для поиска (300ms задержка)
const debouncedSearch = useDebounce(searchTerm, 300);

// Throttle для скролла (100ms интервал)
const throttledScroll = useThrottle(scrollPos, 100);

// Lazy load при скролле
const isVisible = useIntersectionObserver(ref);

// Виртуализация больших таблиц (10000+ строк)
const { visibleItems, offsetY } = useVirtualization(data, 40, 600);
```

#### 3. **Lazy Image Loading**

```typescript
// components/optimizations/LazyImage.tsx
import {
  LazyImage,
  LazyImageGallery,
} from "@/components/optimizations/LazyImage";

// Изображение загружается при скролле
<LazyImage
  src="/image.jpg"
  alt="..."
  width={400}
  height={300}
  quality={75} // Оптимизированное качество
/>;
```

#### 4. **Optimized Dashboard**

```typescript
// components/dashboard/OptimizedDashboard.tsx
// Динамическая загрузка с Suspense
<DynamicKPIMetrics /> // Загружаются первыми
<DynamicRevenueChart /> // Загружаются отложенно
<DynamicInvoiceStatusChart /> // С Skeleton loaders
```

#### 5. **Virtual Table for Large Data**

```typescript
// components/tables/VirtualTable.tsx
// Эффективно отрендеривает 10000+ строк
<VirtualTable
  data={largeDataset}
  itemHeight={40}
  containerHeight={600}
  columns={[...]}
/>
```

#### 6. **Caching Strategies**

```typescript
// lib/cache.ts - 3 уровня кеширования

// Memory Cache (5 минут по умолчанию)
apiCache.set("users", data, 5 * 60 * 1000);
const users = apiCache.get("users");

// LocalStorage Cache (7 дней)
setLocalStorageCache("preferences", data);
const prefs = getLocalStorageCache("preferences");

// IndexedDB Cache (для больших данных)
await indexedDBCache.set("catalog", largeData, 24 * 60 * 60 * 1000);
const catalog = await indexedDBCache.get("catalog");
```

#### 7. **Performance Monitoring**

```typescript
// lib/monitoring/performance.ts
import {
  usePerformanceMonitoring,
  useRenderTime,
} from "@/lib/monitoring/performance";

// Автоматический мониторинг Web Vitals
usePerformanceMonitoring();

// Отслеживание времени рендера компонента
useRenderTime("MyComponent");

// Профилирование функций
profile("complexCalculation", () => {
  // ... код
});
```

#### 8. **Font Optimization**

```typescript
// lib/fonts.ts - Оптимизированная загрузка шрифтов
import { inter, jetBrainsMono, getFontClassNames } from '@/lib/fonts';

// В layout.tsx
<html className={getFontClassNames()}>
```

#### 9. **Next.js Configuration**

```typescript
// next.config.optimization.ts
// - Image optimization (avif, webp)
// - Font loading strategy (swap)
// - Bundle splitting (react, ui, data chunks)
// - Webpack code splitting
// - Production optimizations
```

---

## 🎯 Performance Best Practices

### 1. **Code Splitting Strategy**

```typescript
// ❌ Плохо - всё загружается сразу
import Charts from "@/components/dashboard/Charts";
import AdvancedFilter from "@/components/filters/AdvancedFilter";

// ✅ Хорошо - загружается при необходимости
const DynamicCharts = dynamic(() => import("@/components/dashboard/Charts"));
const DynamicFilter = dynamic(
  () => import("@/components/filters/AdvancedFilter")
);
```

### 2. **Optimization Hooks Pattern**

```typescript
// ✅ Оптимизировать поиск
const [search, setSearch] = useState("");
const debouncedSearch = useDebounce(search, 300);

useEffect(() => {
  // API запрос только когда пользователь завершит набор
  fetchUsers(debouncedSearch);
}, [debouncedSearch]);

// ✅ Оптимизировать скролл
const throttledScroll = useThrottle(scrollPos, 100);

// ✅ Мемоизировать дорогие операции
const memoizedList = useMemo(() => {
  return items.filter((x) => x.price > threshold);
}, [items, threshold]);
```

### 3. **Lazy Loading Pattern**

```typescript
// ✅ Ленивая загрузка изображений
<LazyImage src={url} alt="..." width={400} height={300} />;

// ✅ Ленивая загрузка компонентов при скролле
const ref = useRef(null);
const isVisible = useIntersectionObserver(ref, { rootMargin: "100px" });
if (isVisible) {
  return <HeavyComponent />;
}
```

### 4. **Virtualization for Large Lists**

```typescript
// ❌ Плохо - рендерит все 10000 строк
return (
  <div>
    {data.map(item => (
      <div key={item.id}>{item.name}</div>
    ))}
  </div>
);

// ✅ Хорошо - рендерит только видимые строки
<VirtualTable
  data={data}
  itemHeight={40}
  containerHeight={600}
  columns={[...]}
/>
```

### 5. **Caching Strategy**

```typescript
// Сочетание разных уровней кеша:

// 1. React Query автоматический кеш
useQuery(["users"], fetchUsers);

// 2. Memory кеш для быстрого доступа
apiCache.get("users");

// 3. LocalStorage для сохранения между сессиями
getLocalStorageCache("userPreferences");

// 4. IndexedDB для больших данных
indexedDBCache.get("largeDataset");
```

---

## 📈 Performance Metrics to Monitor

### Core Web Vitals

- **LCP** (Largest Contentful Paint) < 2.5s
- **FID** (First Input Delay) < 100ms
- **CLS** (Cumulative Layout Shift) < 0.1
- **INP** (Interaction to Next Paint) < 200ms

### Custom Metrics

- API call duration
- Component render time
- Cache hit rate
- Bundle size

### Usage

```typescript
import { performanceMonitor } from "@/lib/monitoring/performance";

// Запустить мониторинг
performanceMonitor.startWebVitalsMonitoring();

// Получить метрики
const metrics = performanceMonitor.getMetrics();
```

---

## 🚀 Deployment Optimizations

### Build Time

```bash
npm run build
# Использует Turbopack для быстрой сборки
```

### Bundle Analysis

```bash
npm install --save-dev @next/bundle-analyzer
# См. next.config.optimization.ts для конфигурации
```

### Compression

- Gzip enabled
- Brotli compression
- Image optimization (avif, webp)

### Caching Headers

```typescript
// next.config.optimization.ts
// API: s-maxage=10, stale-while-revalidate=59
// Static: max-age=31536000, immutable
```

---

## 📚 Files Created

| File                                          | Purpose                     | Size       |
| --------------------------------------------- | --------------------------- | ---------- |
| `lib/dynamic.ts`                              | Code splitting utilities    | 100+ lines |
| `lib/performance.ts`                          | Performance hooks           | 180+ lines |
| `lib/fonts.ts`                                | Font optimization           | 40 lines   |
| `lib/cache.ts`                                | Multi-level caching         | 200+ lines |
| `lib/monitoring/performance.ts`               | Performance monitoring      | 200+ lines |
| `components/optimizations/LazyImage.tsx`      | Lazy loading images         | 50 lines   |
| `components/dashboard/OptimizedDashboard.tsx` | Optimized dashboard         | 70 lines   |
| `components/tables/VirtualTable.tsx`          | Virtual scrolling           | 150 lines  |
| `next.config.optimization.ts`                 | Next.js optimization config | 120 lines  |

---

## ✅ Session 5 Performance Phase - COMPLETE

**Tier 3 Optimization Fully Implemented:**

- ✅ Code splitting with dynamic imports
- ✅ Lazy loading (images, components, data)
- ✅ Virtual scrolling for large lists
- ✅ Performance monitoring & metrics
- ✅ Multi-level caching strategy
- ✅ Font optimization
- ✅ Next.js build optimization
- ✅ Performance utilities (debounce, throttle, intersection observer)

**Total Session 5 Progress: ~85% of Architecture**

**Ready for:**

- Production deployment
- Real-time WebSocket integration (Tier 3 - next)
- Audit logging
- Image optimization at scale
