# 🚀 Быстрый старт с новой архитектурой

## Как использовать готовые компоненты

### 1. Импорт и использование List компонента

```typescript
import { CompaniesList } from "@/features/companies/components";

export default function Page() {
  return (
    <div className="container">
      <CompaniesList />
    </div>
  );
}
```

### 2. Использование Form компонента

```typescript
import { CompanyForm } from "@/features/companies/components";
import { createCompany } from "@/lib/api/companies";

export default function CreatePage() {
  const handleSubmit = async (data) => {
    const company = await createCompany(data, token);
    router.push(`/companies/${company.id}`);
  };

  return <CompanyForm onSubmit={handleSubmit} />;
}
```

### 3. Использование API функций

```typescript
import {
  listCompanies,
  getCompanyById,
  createCompany,
  updateCompany,
  deleteCompany,
} from "@/lib/api/companies";

// Получить список
const companies = await listCompanies(token);

// Получить одну
const company = await getCompanyById("123", token);

// Создать
const newCompany = await createCompany(data, token);

// Обновить
const updated = await updateCompany("123", data, token);

// Удалить
await deleteCompany("123", token);
```

---

## Структура для нового feature

Если нужно добавить новый feature (например, "Expenses" - расходы):

### 1. Создать директорию

```bash
mkdir -p features/expenses/components
mkdir -p features/expenses/hooks
mkdir -p features/expenses/types
```

### 2. Создать компоненты (скопировать из Companies)

**features/expenses/components/ExpensesForm.tsx**

```typescript
"use client";

import React, { useState } from "react";
import { Loader2 } from "lucide-react";

interface ExpenseFormProps {
  initialData?: any;
  onSubmit?: (data: any) => Promise<void>;
}

export default function ExpenseForm({
  initialData,
  onSubmit,
}: ExpenseFormProps) {
  const [formData, setFormData] = useState({
    category: initialData?.category || "",
    amount: initialData?.amount || 0,
    date: initialData?.date || new Date().toISOString().split("T")[0],
    description: initialData?.description || "",
  });

  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors: Record<string, string> = {};
    if (!formData.category) newErrors.category = "Категория обязательна";
    if (formData.amount <= 0) newErrors.amount = "Сумма должна быть больше 0";
    return newErrors;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const newErrors = validateForm();

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setIsLoading(true);
    try {
      if (onSubmit) await onSubmit(formData);
    } catch (err) {
      setErrors({ submit: err instanceof Error ? err.message : "Ошибка" });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-gray-900 rounded-lg border border-gray-800 p-6"
    >
      {/* Форма здесь */}
    </form>
  );
}
```

### 3. Создать API клиент

**lib/api/expenses.ts**

```typescript
import axios from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";
const client = axios.create({ baseURL: API_BASE_URL, withCredentials: true });

export interface ExpenseData {
  category: string;
  amount: number;
  date: string;
  description: string;
}

export async function listExpenses(token: string) {
  return (
    await client.get("/api/expenses", {
      headers: { Authorization: `Bearer ${token}` },
    })
  ).data;
}

export async function createExpense(data: ExpenseData, token: string) {
  return (
    await client.post("/api/expenses", data, {
      headers: { Authorization: `Bearer ${token}` },
    })
  ).data;
}

// ... и т.д.
```

### 4. Создать pages

**app/expenses/page.tsx**

```typescript
"use client";
import ExpensesList from "@/features/expenses/components/ExpensesList";

export default function Page() {
  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4">
      <div className="max-w-7xl mx-auto">
        <ExpensesList />
      </div>
    </div>
  );
}
```

---

## Типичные пути в приложении

| Feature         | Путь                    | Компонент         |
| --------------- | ----------------------- | ----------------- |
| Список компаний | `/companies`            | CompaniesList     |
| Создать         | `/companies/create`     | CompanyForm       |
| Редактировать   | `/companies/123/edit`   | CompanyForm       |
| Удалить         | `/companies/123/delete` | DeleteCompanyPage |
| Счета           | `/invoices`             | InvoicesList      |
| Товары          | `/catalog`              | CatalogList       |
| Счета банка     | `/bank-accounts`        | BankAccountsList  |
| Документы       | `/documents/upload`     | DocumentsUpload   |

---

## Ошибки и их обработка

Все API функции используют единый формат ошибок:

```typescript
try {
  const data = await createCompany(formData, token);
} catch (error) {
  // error всегда будет Error объект
  console.log(error.message); // "Company with INN already exists"
}
```

**parseApiError()** автоматически преобразует:

- Ошибки Axios → Error с правильным message
- Response errors → Error message
- Network errors → "Unknown error"

---

## Валидация форм

Все формы используют локальную валидацию:

```typescript
const validateForm = () => {
  const newErrors: Record<string, string> = {};

  if (!formData.name) newErrors.name = "Требуется";
  if (formData.name.length < 2) newErrors.name = "Минимум 2 символа";
  if (!isValidEmail(formData.email)) newErrors.email = "Некорректный формат";

  return newErrors;
};

// Отображение ошибок
{
  errors.name && <p className="text-red-400 text-sm mt-1">{errors.name}</p>;
}
```

---

## Состояния загрузки и отключение

```typescript
const [isLoading, setIsLoading] = useState(false);

<button
  type="submit"
  disabled={isLoading}
  className="disabled:opacity-50 disabled:cursor-not-allowed"
>
  {isLoading ? (
    <>
      <Loader2 className="animate-spin" size={18} />
      Сохранение...
    </>
  ) : (
    "Сохранить"
  )}
</button>;
```

---

## Модальное подтверждение удаления

Все delete страницы используют паттерн:

```typescript
const [isConfirmed, setIsConfirmed] = useState(false);

<label className="flex items-center gap-3 p-4 bg-gray-800 rounded-lg">
  <input
    type="checkbox"
    checked={isConfirmed}
    onChange={(e) => setIsConfirmed(e.target.checked)}
  />
  <span>Я уверен, что хочу удалить</span>
</label>

<button
  disabled={!isConfirmed}
  onClick={handleDelete}
>
  Удалить
</button>
```

---

## Токен и аутентификация

Используется hook `useApiToken()`:

```typescript
import { useApiToken } from "@/lib/hooks/useApiToken";

const token = useApiToken();

if (!token) {
  setError("Требуется аутентификация");
  return;
}

const data = await listCompanies(token);
```

---

## Цвета по feature

Используйте эти цвета для согласованности:

```typescript
// Companies
bg-blue-600 hover:bg-blue-700 text-blue-400

// Invoices
bg-orange-600 hover:bg-orange-700 text-orange-400

// Catalog
bg-purple-600 hover:bg-purple-700 text-purple-400

// Bank Accounts
bg-green-600 hover:bg-green-700 text-green-400

// Documents
bg-indigo-600 hover:bg-indigo-700 text-indigo-400
```

---

## Полезные иконки (Lucide React)

```typescript
import {
  Plus, // Добавить
  Edit2, // Редактировать
  Trash2, // Удалить
  Search, // Поиск
  Loader2, // Загрузка (с animation)
  Building2, // Компании
  FileText, // Счета
  Package, // Товары
  CreditCard, // Банк
  FileUp, // Документы
  AlertTriangle, // Предупреждение
} from "lucide-react";
```

---

## Структура Form пропсов

Все Form компоненты принимают одинаковую сигнатуру:

```typescript
interface FormProps {
  initialData?: Entity;
  onSubmit?: (data: EntityData) => Promise<void>;
}
```

- `initialData` - для режима редактирования (null/undefined для создания)
- `onSubmit` - функция обработки отправки (optional)

---

## Примеры использования всех features

### Компании

```typescript
<CompaniesList />  // Список
<CompanyForm />    // Форма
```

### Счета

```typescript
<InvoicesList />   // Список
<InvoiceForm />    // Форма с товарами
```

### Товары

```typescript
<CatalogList />    // Список с фильтром
<CatalogForm />    // Форма с расчетом маржи
```

### Банк

```typescript
<BankAccountsList />  // Список счетов
<BankAccountForm />   // Форма с валидацией БИК
```

### Документы

```typescript
<DocumentsUpload /> // Загрузка и список
```

---

## Минимальный пример страницы

```typescript
"use client";

import React from "react";
import { CompanyForm } from "@/features/companies/components";
import { useRouter } from "next/navigation";
import { useApiToken } from "@/lib/hooks/useApiToken";
import { createCompany } from "@/lib/api/companies";

export default function CreatePage() {
  const router = useRouter();
  const token = useApiToken();

  const handleSubmit = async (data: any) => {
    const company = await createCompany(data, token!);
    router.push(`/companies/${company.id}`);
  };

  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4">
      <div className="max-w-2xl mx-auto">
        <CompanyForm onSubmit={handleSubmit} />
      </div>
    </div>
  );
}
```

---

✅ **Готово к разработке!**
