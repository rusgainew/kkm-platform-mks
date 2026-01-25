# Catalog UI Implementation Guide

## ✅ Реализованные компоненты

### 1. **CatalogDashboard** (`features/catalog/components/CatalogDashboard.tsx`)

- ✅ Статистика товаров (всего, активные, стоимость, низкий запас)
- ✅ Алерты о низком запасе
- ✅ Список товаров с пагинацией

### 2. **CatalogList** (`features/catalog/components/CatalogList.tsx`)

- ✅ Таблица товаров с действиями
- ✅ Фильтрация и сортировка
- ✅ Диалоги Create/Edit/View
- ✅ Удаление товаров

### 3. **CatalogForm** (`features/catalog/components/CatalogFormNew.tsx`)

- ✅ Форма создания/редактирования товара
- ✅ Валидация полей
- ✅ Категории, цены, количество

### 4. **CatalogFilter** (`features/catalog/components/CatalogFilter.tsx`)

- ✅ Поиск по названию
- ✅ Фильтр по категории
- ✅ Фильтр по статусу (активные/неактивные)

### 5. **CatalogDashboardWidget** (`components/dashboard/CatalogDashboardWidget.tsx`)

- ✅ Мини-виджет в главном дашборде
- ✅ Статистика товаров
- ✅ Последние товары
- ✅ Товары с низким запасом

## 📍 Интеграция в Dashboard

Добавлено в `/components/dashboard/Dashboard.tsx`:

```tsx
// Catalog Section
<div className="rounded-lg border border-gray-800 bg-gray-900/50 backdrop-blur-sm p-6">
  <div className="flex items-center justify-between mb-6">
    <h2 className="text-xl font-bold text-white">Каталог товаров</h2>
    <a href="/catalog" className="text-sm text-blue-400 hover:text-blue-300">
      Перейти в каталог →
    </a>
  </div>
  <CatalogDashboardWidget />
</div>
```

## 🔌 API Integration

### Реализованные API хуки (`lib/api/useCatalogApi.ts`)

```typescript
// Query endpoints
useCatalogList(); // GET /catalogs-query
useCatalogFilter(filters); // GET /catalogs-query/filter

// Command endpoints
useCatalogById(id); // GET /catalog/{id}
useCreateCatalog(); // POST /catalog
useUpdateCatalog(); // PUT /catalog/{id}
useDeleteCatalog(); // DELETE /catalog/{id}
```

## 📄 Страницы

### /catalog - Главная страница каталога

```tsx
import { CatalogDashboard } from "@/features/catalog/components";

export default function CatalogPage() {
  return (
    <div className="container mx-auto py-8">
      <CatalogDashboard />
    </div>
  );
}
```

## 🎯 Основные операции

### 1. Создание товара

```tsx
const { createItem } = useCatalog();

await createItem({
  name: "Товар",
  sku: "SKU-001",
  category: "electronics",
  price: 100,
  costPrice: 50,
  stock: 50,
  active: true,
});
```

### 2. Получение списка

```tsx
const { fetchCatalog, items, filteredItems } = useCatalog();

await fetchCatalog({
  category: "electronics",
  search: "phone",
});
```

### 3. Фильтрация

```tsx
const { items } = useCatalogFilter({
  name: "electronics",
  active: true,
  sort_field: "price",
  sort_order: "ASC",
  page: 0,
  page_size: 20,
});
```

### 4. Редактирование товара

```tsx
const { updateItem } = useCatalog();

await updateItem("catalog-123", {
  name: "Updated Name",
  stock: 100,
  price: 120,
});
```

### 5. Удаление товара

```tsx
const { deleteItem } = useCatalog();

await deleteItem("catalog-123");
```

## 🔍 Zustand Store

Все операции управляются через `useCatalogStore`:

```typescript
const {
  items, // Все товары
  filteredItems, // Отфильтрованные товары
  loading, // Состояние загрузки
  error, // Ошибки
  pagination, // Пагинация

  // CRUD actions
  fetchCatalog, // Получить товары
  createItem, // Создать товар
  updateItem, // Обновить товар
  deleteItem, // Удалить товар

  // Utilities
  lowStockItems, // Товары с низким запасом
  totalInventoryValue, // Общая стоимость склада
} = useCatalogStore();
```

## 📊 Доступные компоненты в Dashboard

1. **Stats Cards** - Статистика (всего, активные, стоимость, низкий запас)
2. **Low Stock Alert** - Алерт о товарах с низким запасом
3. **Catalog List** - Таблица всех товаров
4. **Filter Panel** - Панель фильтров
5. **Dialogs** - Модальные окна для Create/Edit/View

## 🎨 UI Features

- ✅ Темная тема (matches dashboard styling)
- ✅ Responsive grid layouts
- ✅ Модальные диалоги
- ✅ Иконки Lucide React
- ✅ Toast уведомления об успехе/ошибке
- ✅ Пагинация
- ✅ Фильтры и сортировка

## 📱 Endpoints реализованные в UI

| Метод  | Endpoint                 | Компонент                     |
| ------ | ------------------------ | ----------------------------- |
| POST   | `/catalog`               | CatalogForm (Create)          |
| GET    | `/catalog/{id}`          | CatalogDetail (View)          |
| PUT    | `/catalog/{id}`          | CatalogForm (Edit)            |
| DELETE | `/catalog/{id}`          | CatalogList (Delete action)   |
| GET    | `/catalogs-query`        | CatalogList, Dashboard Widget |
| GET    | `/catalogs-query/filter` | CatalogFilter                 |

## 🔄 CQRS Pattern в Frontend

- **Write operations** → useCatalog hooks → POST/PUT/DELETE
- **Read operations** → useCatalogFilter/List → GET /catalogs-query
- **Caching** → React Query handles automatic cache invalidation

## ✨ Готово к использованию!

Все компоненты синхронизированы с API Gateway endpoints:

- `/catalog` - Command operations
- `/catalogs-query` - Query operations with Redis caching
- Полная CQRS реализация на фронтенде
