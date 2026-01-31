# 🧹 План очистки и улучшения проекта

## 📊 Текущее состояние

### Статистика:
- **912 MD файлов** - избыточная документация
- **658 MB** - .next (build артефакты)
- **775 MB** - node_modules
- **263 тестовых файла**
- **5 backup файлов** (.backup, .bak, .old)
- **Дубликаты типов** (auth.ts, company.ts, enums.ts vs entities.ts)

---

## 🗑️ ЧТО УДАЛИТЬ

### 1. Backup и временные файлы (НЕМЕДЛЕННО)
```
./components/auth/ProtectedRoute.tsx.backup
./app/[locale]/auth/page.tsx.bak
./features/companies/components/CompanyForm.tsx.backup
./features/catalog/components/CatalogForm.tsx.backup
./lib/api/catalog.ts.backup
```

### 2. Устаревшие MD документы (30+ файлов)
**Оставить только:**
- README.md (главный)
- QUICK_START.md (быстрый старт)
- PROJECT_STRUCTURE.md (структура)

**Удалить:**
- PHASE_*_COMPLETE.md (история разработки)
- IMPLEMENTATION_REPORT_*.md (отчеты)
- *_QUICK_START.md (дубликаты)
- API_CLIENT_USAGE.md, INVOICE_API_EXAMPLES.md (устаревшие примеры)
- CATALOG_ERROR_DIAGNOSIS.md, DOCUMENTS_REFACTORING.md (временные)
- BACKEND_FIRST_PLAN.md, KKM_PLATFORM_ENHANCEMENT_PLAN.md, FORMS_IMPLEMENTATION_PLAN.md (планы)

### 3. Дублирующиеся типы
**Удалить:**
- types/auth.ts → мигрировано в entities.ts
- types/company.ts → мигрировано в entities.ts
- types/enums.ts → мигрировано в entities.ts

### 4. Build артефакты (добавить в .gitignore)
```
.next/
test-results/
playwright-report/
coverage/
```

### 5. Неиспользуемые компоненты
Нужна проверка:
- components/examples/ (демо компоненты)
- app/[locale]/audit-demo/ (демо страница)

---

## 🔧 ЧТО УЛУЧШИТЬ

### 1. Типы и импорты

**Проблема:** Импорты из устаревших типов
```typescript
// ПЛОХО
import { User } from '@/types/auth';
import { Company } from '@/types/company';

// ХОРОШО
import { User, Company } from '@/types/entities';
```

**Решение:**
- Обновить все импорты на @/types/entities
- Удалить старые файлы типов
- Оставить только специализированные (invoice.ts, party.ts, reference.ts)

### 2. API клиенты

**Проблема:** Дублирование логики в lib/api/*
- Каждый файл реализует свой apiRequest
- Нет единого error handling
- Дублирование Bearer auth логики

**Решение:**
```typescript
// lib/api/client.ts - единый клиент
export class ApiClient {
  private baseUrl: string;
  
  async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    // Единая логика: auth, retries, error handling
  }
}

// lib/api/resources/companies.ts
export const companiesAPI = {
  list: () => apiClient.request<ListResponse<Company>>('/companies'),
  create: (data) => apiClient.request<Company>('/companies', { method: 'POST', body: data })
}
```

### 3. Компоненты

**Проблема:** Дублирование стилей и логики
- Повторяющиеся className строки
- Копипаста модальных окон
- Дубликаты таблиц

**Решение:**
```typescript
// components/ui/Modal.tsx - базовый компонент
export function Modal({ title, children, onClose }) {
  // Единая реализация модального окна
}

// components/ui/DataTable.tsx - универсальная таблица
export function DataTable<T>({ columns, data, onRowClick }) {
  // Единая реализация таблицы с generic типами
}
```

### 4. Hooks

**Проблема:** Разрозненные хуки для API
- hooks/useCompanies.ts
- lib/hooks/useCompaniesApi.ts (дубликат!)
- Разная логика в разных местах

**Решение:**
```typescript
// hooks/api/useResource.ts - универсальный хук
export function useResource<T>(resource: string) {
  return {
    list: useQuery([resource, 'list'], () => api[resource].list()),
    create: useMutation((data) => api[resource].create(data)),
    // ...
  }
}

// Использование:
const companies = useResource<Company>('companies');
const invoices = useResource<Invoice>('invoices');
```

### 5. Конфигурация

**Проблема:** Дублирование getApiBase()
- В каждом API файле свой getApiBase
- Одинаковая логика в 10+ местах

**Решение:**
```typescript
// lib/config/api.ts
export const apiConfig = {
  baseUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost/api/v1',
  timeout: 10000,
  retries: 3
};
```

### 6. Тесты

**Проблема:** 2 ошибки в тестах
- e2e/bank-accounts.spec.ts:280 - неправильный синтаксис toHaveCount
- invoice-validation.test.ts:173 - отсутствует isPriceWithoutTaxes

**Решение:** Исправить тесты (низкий приоритет)

---

## 📁 Рекомендуемая структура

```
kkm-platform/
├── app/                    # Next.js app router
├── components/            
│   ├── ui/               # Базовые UI компоненты (Button, Modal, Table)
│   ├── features/         # Фича-специфичные компоненты
│   └── layout/           # Layout компоненты
├── lib/
│   ├── api/              # API клиенты (только resources)
│   ├── hooks/            # React hooks (api/, ui/)
│   ├── utils/            # Утилиты
│   └── config/           # Конфигурация
├── types/
│   ├── entities.ts       # Основные типы (User, Company, Invoice)
│   ├── invoice.ts        # ESF специфичные типы
│   ├── api-response.ts   # API response типы
│   └── index.ts          # Экспорты
├── store/                # Zustand stores
├── hooks/                # Общие hooks
├── features/             # Feature modules (companies/, invoices/)
└── tests/               # Тесты (unit/, e2e/)
```

---

## 🎯 План действий

### Фаза 1: Очистка (15 минут)
1. ✅ Удалить backup файлы (5 файлов)
2. ✅ Удалить 25+ устаревших MD файлов
3. ✅ Удалить старые типы (auth.ts, company.ts, enums.ts)
4. ✅ Обновить .gitignore

### Фаза 2: Рефакторинг типов (30 минут)
1. ✅ Обновить импорты (find & replace)
2. ✅ Удалить неиспользуемые типы
3. ✅ Очистить types/index.ts

### Фаза 3: Унификация API (1 час)
1. 🔄 Создать lib/api/client.ts
2. 🔄 Рефакторить API клиенты
3. 🔄 Убрать дублирование getApiBase()

### Фаза 4: Компоненты (1 час)
1. 🔄 Создать базовые UI компоненты
2. 🔄 Рефакторить модальные окна
3. 🔄 Унифицировать таблицы

### Фаза 5: Hooks (30 минут)
1. 🔄 Объединить дубликаты hooks
2. 🔄 Создать useResource generic hook
3. 🔄 Удалить старые hooks

---

## 💾 Ожидаемый результат

### До:
- 912 MD файлов
- Дубликаты типов в 3 файлах
- 10+ реализаций apiRequest
- Повторяющийся код в компонентах

### После:
- 3-5 MD файлов (документация)
- Единый источник типов (entities.ts)
- 1 API клиент с retry/error handling
- Переиспользуемые UI компоненты
- Сокращение кода на ~30%
- Улучшение производительности
- Легче поддерживать

---

## ⚠️ Риски

1. **Импорты** - после удаления файлов могут сломаться импорты
   - Решение: Использовать find & replace перед удалением
   
2. **Тесты** - могут сломаться после рефакторинга
   - Решение: Запускать тесты после каждого изменения
   
3. **Git история** - удаление файлов
   - Решение: Делать коммит перед каждой фазой

---

## 🚀 Начинаем?

Готов начать с Фазы 1 (Очистка) - это безопасно и даст быстрый результат.
