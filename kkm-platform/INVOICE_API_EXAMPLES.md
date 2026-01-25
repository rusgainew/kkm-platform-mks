## Примеры Интеграции Invoice API

### ✅ Готовые файлы:

1. **lib/api/invoice.ts** — Invoice API Client (533 строк)
2. **hooks/useInvoice.ts** — React Query Hooks (383 строк)
3. **lib/api/index.ts** — Exports (обновлено)
4. **types/invoice.ts** — TypeScript типы (410 строк)
5. **types/api-response.ts** — ESF Response типы (515 строк)

---

## Базовый Пример: Список Счетов

```typescript
// features/invoices/pages/InvoiceListPage.tsx
'use client';

import { useInvoiceList } from '@/hooks/useInvoice';
import { InvoiceFilters } from '@/types';
import { useState } from 'react';

export default function InvoiceListPage() {
  const [filters, setFilters] = useState<InvoiceFilters>({
    status: 'DRAFT',
    limit: 20,
    page: 1,
  });

  const { data, isLoading, error, refetch } = useInvoiceList(filters);

  if (isLoading) return <div className="p-4">Загрузка...</div>;
  if (error) return <div className="text-red-500">Ошибка: {error.message}</div>;

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Счета-фактуры</h1>

      {/* Фильтры */}
      <div className="mb-4 space-x-2">
        <button
          onClick={() => setFilters({...filters, status: 'DRAFT'})}
          className="px-4 py-2 bg-blue-500 text-white rounded"
        >
          Черновики ({data?.data.filter(i => i.status.code === 'DRAFT').length})
        </button>
        <button
          onClick={() => setFilters({...filters, status: 'ACCEPTED'})}
          className="px-4 py-2 bg-green-500 text-white rounded"
        >
          Принятые ({data?.data.filter(i => i.status.code === 'ACCEPTED').length})
        </button>
      </div>

      {/* Таблица */}
      <table className="w-full border">
        <thead className="bg-gray-100">
          <tr>
            <th className="border p-2">Номер</th>
            <th className="border p-2">Статус</th>
            <th className="border p-2">Сумма</th>
            <th className="border p-2">Действия</th>
          </tr>
        </thead>
        <tbody>
          {data?.data.map(invoice => (
            <tr key={invoice.documentUuid} className="hover:bg-gray-50">
              <td className="border p-2">{invoice.invoiceNumber}</td>
              <td className="border p-2">
                <span className={`px-2 py-1 rounded text-sm ${
                  invoice.status.code === 'DRAFT' ? 'bg-yellow-100' : 'bg-green-100'
                }`}>
                  {invoice.status.name}
                </span>
              </td>
              <td className="border p-2">{invoice.totalAmount} сом</td>
              <td className="border p-2">
                <button className="text-blue-500">Просмотр</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Пагинация */}
      <div className="mt-4 flex gap-2">
        <button
          onClick={() => setFilters({...filters, page: Math.max(1, (filters.page || 1) - 1)})}
          className="px-4 py-2 bg-gray-300 rounded"
        >
          ← Предыдущая
        </button>
        <span className="px-4 py-2">
          Страница {filters.page || 1}
        </span>
        <button
          onClick={() => setFilters({...filters, page: (filters.page || 1) + 1})}
          className="px-4 py-2 bg-gray-300 rounded"
        >
          Следующая →
        </button>
      </div>
    </div>
  );
}
```

---

## Пример 2: Создание Счета

```typescript
// features/invoices/components/CreateInvoiceForm.tsx
'use client';

import { useCreateInvoice } from '@/hooks/useInvoice';
import { useInvoiceList } from '@/hooks/useInvoice';
import { CreateInvoiceRequest } from '@/types';
import { useState } from 'react';

export function CreateInvoiceForm() {
  const [formData, setFormData] = useState<CreateInvoiceRequest>({
    invoiceNumber: '',
    contractorTin: '',
    totalAmount: 0,
    description: '',
    isResident: true,
  });

  const { mutate, isPending, error } = useCreateInvoice();
  const { refetch } = useInvoiceList(); // Обновить список после создания

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutate(formData, {
      onSuccess: (invoice) => {
        alert(`Счет создан: ${invoice.documentUuid}`);
        setFormData({ invoiceNumber: '', contractorTin: '', totalAmount: 0, description: '', isResident: true });
        refetch(); // Обновить список
      },
      onError: (error) => {
        alert(`Ошибка: ${error.message}`);
      },
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-w-md">
      <h2 className="text-xl font-bold">Создать счет-фактуру</h2>

      <div>
        <label className="block font-semibold mb-1">Номер счета</label>
        <input
          type="text"
          value={formData.invoiceNumber}
          onChange={(e) => setFormData({...formData, invoiceNumber: e.target.value})}
          className="w-full border rounded px-3 py-2"
          required
        />
      </div>

      <div>
        <label className="block font-semibold mb-1">ИНН контрагента</label>
        <input
          type="text"
          value={formData.contractorTin}
          onChange={(e) => setFormData({...formData, contractorTin: e.target.value})}
          className="w-full border rounded px-3 py-2"
          required
        />
      </div>

      <div>
        <label className="block font-semibold mb-1">Сумма (сом)</label>
        <input
          type="number"
          value={formData.totalAmount}
          onChange={(e) => setFormData({...formData, totalAmount: Number(e.target.value)})}
          className="w-full border rounded px-3 py-2"
          required
        />
      </div>

      <div>
        <label className="block font-semibold mb-1">Описание</label>
        <textarea
          value={formData.description}
          onChange={(e) => setFormData({...formData, description: e.target.value})}
          className="w-full border rounded px-3 py-2"
          rows={3}
        />
      </div>

      <div className="flex items-center gap-2">
        <input
          type="checkbox"
          checked={formData.isResident}
          onChange={(e) => setFormData({...formData, isResident: e.target.checked})}
          id="resident"
        />
        <label htmlFor="resident">Резидент</label>
      </div>

      {error && <div className="text-red-500">{error.message}</div>}

      <button
        type="submit"
        disabled={isPending}
        className="w-full bg-blue-500 text-white font-semibold py-2 rounded disabled:bg-gray-400"
      >
        {isPending ? 'Создание...' : 'Создать счет'}
      </button>
    </form>
  );
}
```

---

## Пример 3: Детали Счета с Действиями

```typescript
// features/invoices/components/InvoiceDetail.tsx
'use client';

import {
  useInvoiceDetail,
  useInvoiceStatusHistory,
  useAcceptInvoice,
  useRejectInvoice,
  useSignInvoice,
} from '@/hooks/useInvoice';

export function InvoiceDetail({ uuid }: { uuid: string }) {
  const { data: invoice, isLoading } = useInvoiceDetail(uuid);
  const { data: history } = useInvoiceStatusHistory(uuid);

  const { mutate: accept, isPending: acceptPending } = useAcceptInvoice();
  const { mutate: reject, isPending: rejectPending } = useRejectInvoice();
  const { mutate: sign, isPending: signPending } = useSignInvoice();

  if (isLoading) return <div>Загрузка...</div>;
  if (!invoice) return <div>Счет не найден</div>;

  return (
    <div className="p-6 max-w-4xl">
      <h1 className="text-2xl font-bold mb-4">{invoice.invoiceNumber}</h1>

      {/* Основная информация */}
      <div className="grid grid-cols-2 gap-4 mb-6 border p-4 rounded">
        <div>
          <strong>UUID:</strong> {invoice.documentUuid}
        </div>
        <div>
          <strong>Статус:</strong>
          <span className="ml-2 px-2 py-1 bg-blue-100 rounded">
            {invoice.status.name}
          </span>
        </div>
        <div>
          <strong>Сумма:</strong> {invoice.totalAmount} сом
        </div>
        <div>
          <strong>Резидент:</strong> {invoice.isResident ? 'Да' : 'Нет'}
        </div>
      </div>

      {/* Действия */}
      <div className="flex gap-2 mb-6">
        {invoice.status.code === 'DRAFT' && (
          <>
            <button
              onClick={() => accept({ invoiceUuid: uuid })}
              disabled={acceptPending}
              className="px-4 py-2 bg-green-500 text-white rounded disabled:bg-gray-400"
            >
              {acceptPending ? 'Принимаю...' : 'Принять'}
            </button>
            <button
              onClick={() => reject({ invoiceUuid: uuid, reason: 'Отклонено' })}
              disabled={rejectPending}
              className="px-4 py-2 bg-red-500 text-white rounded disabled:bg-gray-400"
            >
              {rejectPending ? 'Отклоняю...' : 'Отклонить'}
            </button>
          </>
        )}
        {invoice.status.code === 'ACCEPTED' && (
          <button
            onClick={() => sign({ invoiceUuid: uuid, signature: '' })}
            disabled={signPending}
            className="px-4 py-2 bg-purple-500 text-white rounded disabled:bg-gray-400"
          >
            {signPending ? 'Подписываю...' : 'Подписать'}
          </button>
        )}
      </div>

      {/* История статусов */}
      <div className="border-t pt-4">
        <h3 className="font-bold mb-2">История статусов:</h3>
        <div className="space-y-2">
          {history?.data.map((item: any, idx: number) => (
            <div key={idx} className="text-sm text-gray-600">
              <strong>{item.status}</strong> - {item.timestamp}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
```

---

## Пример 4: Error Handling

```typescript
// Proper error handling in components
import { useInvoiceList } from '@/hooks/useInvoice';
import { InvoiceAPIError } from '@/lib/api/invoice';

export function SafeInvoiceList() {
  const { data, error, isError } = useInvoiceList();

  if (isError && error instanceof InvoiceAPIError) {
    return (
      <div className="border border-red-300 bg-red-50 p-4 rounded">
        <strong>API Ошибка {error.statusCode}:</strong> {error.message}
        {error.details && (
          <pre className="mt-2 text-sm overflow-auto">
            {JSON.stringify(error.details, null, 2)}
          </pre>
        )}
      </div>
    );
  }

  if (isError) {
    return <div className="text-red-500">Неизвестная ошибка</div>;
  }

  return <div>{data?.data.length} счетов найдено</div>;
}
```

---

## Пример 5: Custom Hook для Бизнес-Логики

```typescript
// hooks/useInvoiceWorkflow.ts - Custom hook для workflow

import { useInvoiceDetail } from "./useInvoice";
import { useAcceptInvoice, useSignInvoice } from "./useInvoice";
import { useCallback } from "react";

export function useInvoiceWorkflow(invoiceUuid: string) {
  const { data: invoice } = useInvoiceDetail(invoiceUuid);
  const { mutate: accept } = useAcceptInvoice();
  const { mutate: sign } = useSignInvoice();

  const canAccept = invoice?.status.code === "DRAFT";
  const canSign = invoice?.status.code === "ACCEPTED";

  const acceptAndThen = useCallback(async () => {
    return new Promise((resolve, reject) => {
      accept(
        { invoiceUuid },
        {
          onSuccess: () => {
            // После принятия можно сделать что-то еще
            console.log("Счет принят, готов к подписанию");
            resolve(true);
          },
          onError: reject,
        },
      );
    });
  }, [accept, invoiceUuid]);

  return {
    invoice,
    canAccept,
    canSign,
    acceptAndThen,
  };
}
```

---

## Технические Детали

### Query Key Structure

```typescript
// invoiceQueryKeys.all = ['invoices']
// invoiceQueryKeys.list() = ['invoices', 'list']
// invoiceQueryKeys.list(filters) = ['invoices', 'list', filters]
// invoiceQueryKeys.detail(id) = ['invoices', 'detail', id]
```

### Caching Strategy

- List queries: 5 минут stale time (часто меняются)
- Detail queries: 2 минуты stale time
- Status history: 1 минута stale time
- Все с 2x retry logic

### Token Management

Автоматический refresh через authStore:

```typescript
const { tokens } = useAuthStore();
const token = tokens?.accessToken;

// Если токена нет - query не запустится
// enabled: enabled && !!token
```

---

## Что можно делать дальше

1. **Unit Tests** для hooks и API client
2. **E2E Tests** для полного workflow'а
3. **React Components** для UI
4. **Zustand Store** для локального состояния
5. **Error Boundaries** для graceful error handling
6. **Loading Skeletons** для лучшего UX
7. **Pagination Controls** компоненты
8. **Search/Filter UI** компоненты

---

**Все файлы готовы к использованию! 🚀**
