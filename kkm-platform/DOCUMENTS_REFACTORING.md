# Documents Module Refactoring

## Overview

Проведен полный рефакторинг модуля документов для улучшения масштабируемости, переиспользования кода и типизации.

## Changes Made

### 1. **Константы и утилиты** (`lib/documents/`)

#### `constants.ts`

- Извлечены все константы в отдельный файл
- Добавлены типизированные статусы документов
- STATUS_LABELS и STATUS_COLORS для единой обработки отображения

#### `formatting.ts`

- Утилиты форматирования (дата, ID документа, текст)
- Переиспользуемые функции для форматирования

### 2. **Компоненты документов** (`components/documents/`)

#### `CreateDocumentForm.tsx`

- ✅ Удалены неиспользуемые константы (DOCUMENT_TYPES)
- ✅ Используются импортированные ORGANIZATIONS
- Упрощение логики формы

#### `DocumentsList.tsx`

- ✅ Миграция на хук `useDocuments`
- ✅ Удалена дублированная логика загрузки
- ✅ Упрощена логика рендера (используется DocumentRow)
- ✅ Улучшена читаемость кода

#### `DocumentRow.tsx` (новый)

- Новый компонент для строки таблицы
- Изолирует логику рендера одной строки
- Упрощает повторное использование в других местах

### 3. **Хуки** (`hooks/`)

#### `useDocuments.ts` (новый)

- Централизованная логика загрузки документов
- Обработка ошибок и состояния загрузки
- Переиспользуемый хук для компонентов

### 4. **Типы** (`types/`)

#### `documents.ts` (новый)

- Централизованные типы документов
- DocumentStatus, DocumentFilters и т.д.

## File Structure

```
lib/
  documents/
    constants.ts        # Константы и перечисления
    formatting.ts       # Функции форматирования

hooks/
  useDocuments.ts      # Хук загрузки документов

components/documents/
  DocumentsList.tsx    # Список документов (рефакторен)
  DocumentRow.tsx      # Строка таблицы (новый)
  CreateDocumentForm.tsx # Форма создания (упрощен)

types/
  documents.ts         # TypeScript типы
```

## Benefits

✅ **Уменьшение дублирования кода** - логика загрузки в одном месте  
✅ **Лучшая типизация** - централизованные типы  
✅ **Переиспользуемость** - компоненты и хуки могут использоваться в других местах  
✅ **Упрощение поддержки** - четкая разделение ответственности  
✅ **Лучшая производительность** - мемоизация в DocumentRow  
✅ **Более чистый код** - DocumentsList стал проще и понятнее

## Usage Examples

### Использование хука useDocuments

```typescript
const { documents, isLoading, error } = useDocuments({
  page: 1,
  pageSize: 10,
  status: "draft",
});
```

### Использование форматирования

```typescript
import { formatDate, formatDocumentId } from "@/lib/documents/formatting";

formatDate(doc.created_at); // "16.01.2026 12:30"
formatDocumentId(doc.id); // "ABC12345"
```

### Использование констант

```typescript
import { STATUS_COLORS, STATUS_LABELS } from "@/lib/documents/constants";

const color = STATUS_COLORS["draft"]; // "bg-gray-500/20 text-gray-400"
const label = STATUS_LABELS["draft"]; // "Черновик"
```

## Next Steps

1. [ ] Добавить обработчики действий в DocumentRow (onView, onDownload)
2. [ ] Создать компонент DocumentDetail для просмотра
3. [ ] Добавить пагинацию в DocumentsList
4. [ ] Добавить фильтры по организации
5. [ ] Создать тесты для компонентов и хуков
