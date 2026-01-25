# Структура проекта

## 📁 Организация файлов

```
kkm-platform/
├── app/                          # Next.js App Router
│   ├── layout.tsx               # Root layout с Inter шрифтом
│   ├── page.tsx                 # Главная страница POS
│   └── globals.css              # Глобальные стили
│
├── components/                   # Все React компоненты
│   ├── layout/                  # Компоненты макета
│   │   └── Sidebar.tsx          # Боковая навигационная панель
│   │
│   ├── features/                # Функциональные компоненты
│   │   ├── ProductGrid.tsx      # Сетка товаров с поиском
│   │   ├── Cart.tsx             # Корзина (текущий чек)
│   │   └── CheckoutPanel.tsx    # Панель оплаты
│   │
│   ├── ui/                      # UI компоненты
│   │   └── ProductCard.tsx      # Карточка товара
│   │
│   └── index.ts                 # Barrel экспорт компонентов
│
├── hooks/                        # Custom React hooks
│   ├── useCart.ts               # Хук для управления корзиной
│   ├── useOnlineStatus.ts       # Хук для отслеживания статуса online/offline
│   └── index.ts                 # Barrel экспорт hooks
│
├── types/                        # TypeScript типы и интерфейсы
│   ├── product.ts               # Product, ProductCategory
│   ├── cart.ts                  # CartItem, CartSummary, PaymentMethod
│   └── index.ts                 # Barrel экспорт типов
│
├── constants/                    # Константы приложения
│   ├── products.ts              # MOCK_PRODUCTS данные
│   ├── config.ts                # TAX_RATE, PAYMENT_METHODS
│   └── index.ts                 # Barrel экспорт констант
│
├── lib/                         # Утилиты и хелперы
│   ├── utils.ts                 # Вспомогательные функции
│   └── index.ts                 # Barrel экспорт утилит
│
└── public/                      # Статические файлы

```

## 🎯 Принципы организации

### 1. **Feature-based структура**

Компоненты разделены по назначению:

- `layout/` - компоненты макета и навигации
- `features/` - бизнес-логика и функциональность
- `ui/` - переиспользуемые UI компоненты

### 2. **Централизованные типы**

Все TypeScript типы в `types/`:

- Легко найти и переиспользовать
- Единый источник правды для типов
- Простое обслуживание

### 3. **Custom Hooks**

Бизнес-логика вынесена в hooks:

- `useCart` - управление корзиной
- `useOnlineStatus` - отслеживание статуса соединения
- Переиспользуемость и тестируемость

### 4. **Константы и конфигурация**

Все настраиваемые значения в `constants/`:

- Легко изменить налоговую ставку
- Централизованные mock данные
- Простая конфигурация

### 5. **Barrel Exports**

Каждая директория имеет `index.ts`:

```typescript
// Вместо:
import ProductCard from "@/components/ui/ProductCard";
import Cart from "@/components/features/Cart";

// Можно:
import { ProductCard, Cart } from "@/components";
```

### 6. **Абсолютные импорты**

Использование `@/` алиаса:

```typescript
import { Product } from "@/types";
import { MOCK_PRODUCTS } from "@/constants";
import { useCart } from "@/hooks";
```

## 📝 Соглашения по именованию

### Файлы

- **Компоненты**: PascalCase (`ProductCard.tsx`)
- **Hooks**: camelCase с префиксом use (`useCart.ts`)
- **Типы**: camelCase (`product.ts`)
- **Утилиты**: camelCase (`utils.ts`)

### Экспорты

- **Default export** для компонентов
- **Named exports** для hooks, типов, констант

## 🔄 Паттерны

### 1. Composition over Inheritance

Компоненты составляются из более мелких компонентов

### 2. Single Responsibility

Каждый компонент отвечает за одну задачу

### 3. Separation of Concerns

- Логика в hooks
- UI в компонентах
- Типы отдельно
- Данные в constants

### 4. DRY (Don't Repeat Yourself)

Переиспользуемые компоненты и утилиты

## 🚀 Расширение проекта

### Добавление нового компонента:

1. Определите категорию (`ui/`, `features/`, `layout/`)
2. Создайте файл компонента
3. Добавьте в `components/index.ts`

### Добавление нового типа:

1. Создайте файл в `types/`
2. Экспортируйте через `types/index.ts`

### Добавление нового hook:

1. Создайте файл в `hooks/`
2. Начните имя с `use`
3. Экспортируйте через `hooks/index.ts`

## 📦 Зависимости

- **Next.js 16** - React фреймворк
- **React 19** - UI библиотека
- **TypeScript** - Типизация
- **Tailwind CSS** - Стилизация
- **Lucide React** - Иконки

## 🎨 Стилизация

- Tailwind CSS для всех стилей
- Utility-first подход
- Responsive дизайн
- Темные/светлые темы через CSS переменные

## ✅ Лучшие практики

- ✅ TypeScript strict mode
- ✅ ESLint для контроля качества кода
- ✅ Абсолютные импорты
- ✅ Barrel exports
- ✅ Custom hooks для переиспользования логики
- ✅ Мемоизация в hooks (useMemo, useCallback)
- ✅ Типизация всех props и состояний
- ✅ Разделение concerns (UI, логика, данные)
