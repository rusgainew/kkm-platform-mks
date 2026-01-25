# ESF Invoice API Client Usage Guide

**Файл:** `/kkm-platform/lib/api/invoice.ts`  
**Hooks:** `/kkm-platform/hooks/useInvoice.ts`  
**Версия:** 1.0  
**Дата:** 23 января 2026 г.

---

## 📚 Обзор

Полнофункциональный Axios клиент для работы с ESF API счетами-фактурами через API Gateway.

**Особенности:**

- ✅ Type-safe типизация (TypeScript)
- ✅ Автоматическое управление токеном
- ✅ Error handling с типизацией
- ✅ React Query интеграция (через hooks)
- ✅ Кэширование результатов
- ✅ Retry логика

---

## 🚀 Быстрый старт

### 1. Использование клиента напрямую

```typescript
import { invoiceAPI } from '@/lib/api/invoice';
import { useAuthStore } from '@/store/authStore';

// В компоненте или сервисе
const { token } = useAuthStore();

// Установить токен (один раз при инициализации)
invoiceAPI.setToken(token);

// Создать счет
const newInvoice = await invoiceAPI.createInvoice({
  operationTypeCode: ESFOperationType.SALES,
  invoiceNumber: 'СФ-001-2026',
  deliveryDate: '2026-01-25',
  contractorTin: '01206200110100',
  deliveryTypeCode: ESFDeliveryType.DIRECT,
  paymentCode: ESFPaymentType.BANK_TRANSFER,
  currencyCode: CurrencyCode.KGS,
  isResident: true,
  isPriceWithoutTaxes: false,
  taxRateVATCode: '12',
  catalogEntries: [...]
});
```

### 2. Использование React Hooks (рекомендуется)

```typescript
import {
  useInvoiceList,
  useInvoiceDetail,
  useCreateInvoice,
  useAcceptInvoice
} from '@/hooks/useInvoice';

function MyComponent() {
  // Получить список счетов
  const { data, isLoading, error } = useInvoiceList({
    status: ESFDocumentStatus.SENT,
    limit: 20
  });

  // Создать счет
  const { mutate: createInvoice, isPending } = useCreateInvoice();

  const handleCreate = async (request: CreateInvoiceRequest) => {
    createInvoice(request, {
      onSuccess: (newInvoice) => {
        console.log('Created:', newInvoice.invoiceNumber);
      },
      onError: (error) => {
        console.error('Failed:', error.message);
      }
    });
  };

  return (
    // ... JSX
  );
}
```

---

## 📖 Полная документация API

### Основной клиент (`InvoiceAPIClient`)

#### Конструктор

```typescript
const client = new InvoiceAPIClient(baseURL?: string);

// baseURL - URL API Gateway (по умолчанию http://localhost:8080/api)
```

#### Методы

##### 1. `createInvoice(request: CreateInvoiceRequest): Promise<Invoice>`

Создать новый счет-фактуру.

```typescript
const invoice = await invoiceAPI.createInvoice({
  operationTypeCode: ESFOperationType.SALES, // Обязательно
  invoiceNumber: "СФ-001", // Опционально
  deliveryDate: "2026-01-25", // Обязательно (YYYY-MM-DD)
  contractorTin: "01206200110100", // Обязательно (ИНН покупателя)
  deliveryTypeCode: ESFDeliveryType.DIRECT, // Обязательно
  paymentCode: ESFPaymentType.BANK_TRANSFER, // Обязательно
  currencyCode: CurrencyCode.KGS, // Обязательно
  isResident: true, // Обязательно
  isPriceWithoutTaxes: false, // Обязательно
  taxRateVATCode: "12", // Обязательно
  catalogEntries: [
    // Обязательно (минимум 1)
    {
      id: 2047,
      quantity: 100,
      price: 250,
      unitClassificationCode: "796",
      salesTaxCode: "10",
      vatAmount: 5357.14,
      amountWithoutTaxes: 25000,
      totalAmount: 30357.14,
    },
  ],
});
```

**Ошибки:**

- `400` - Bad Request (валидация)
- `401` - Unauthorized (токен истек)
- `409` - Conflict (дублирование номера)

---

##### 2. `getInvoice(invoiceUuid: string): Promise<Invoice>`

Получить счет по UUID.

```typescript
const invoice = await invoiceAPI.getInvoice("uuid-xxx-xxx");

console.log(invoice.invoiceNumber); // СФ-001
console.log(invoice.totalAmount); // 30357.14
console.log(invoice.status?.code); // 10 (новый)
```

**Ошибки:**

- `404` - Invoice not found

---

##### 3. `listInvoices(filters?: InvoiceFilters): Promise<ListResponse<Invoice>>`

Получить список счетов с пагинацией.

```typescript
const response = await invoiceAPI.listInvoices({
  status: ESFDocumentStatus.SENT,
  dateFrom: "2026-01-01",
  dateTo: "2026-01-31",
  currency: CurrencyCode.KGS,
  page: 1,
  limit: 20, // Max 100
});

console.log(response.data.length); // Массив счетов
console.log(response.meta.total_count); // Всего найдено
console.log(response.meta.total_pages); // Всего страниц
```

**Фильтры:**

```typescript
interface InvoiceFilters {
  status?: string; // ESFDocumentStatus код
  invoiceNumber?: string; // Номер счета
  contractorTin?: string; // ИНН контрагента
  dateFrom?: string; // С даты (YYYY-MM-DD)
  dateTo?: string; // По дату (YYYY-MM-DD)
  createdBy?: string; // ID создателя
  isResident?: boolean; // Только резиденты
  currency?: string; // КОД валюты
  page?: number; // Страница (по умолчанию 1)
  limit?: number; // На странице (по умолчанию 20, max 100)
}
```

---

##### 4. `updateInvoice(request: UpdateInvoiceRequest): Promise<Invoice>`

Обновить существующий счет.

```typescript
const updated = await invoiceAPI.updateInvoice({
  documentUuid: "uuid-xxx", // Обязательно
  deliveryDate: "2026-01-26", // Опционально
  comment: "Обновлено", // Опционально
  // ... другие поля
});
```

**Ограничения:**

- Можно обновлять только черновики (статус 10)
- После отправки - обновление невозможно

---

##### 5. `acceptInvoice(invoiceUuid: string, comment?: string): Promise<InvoiceOperationResponse>`

Принять счет-фактуру.

```typescript
const result = await invoiceAPI.acceptInvoice(
  "uuid-xxx",
  "Принято без замечаний", // Опционально
);

console.log(result.success); // true
console.log(result.status); // '40' (принят)
```

---

##### 6. `rejectInvoice(invoiceUuid: string, reason: string, comment?: string): Promise<InvoiceOperationResponse>`

Отклонить счет-фактуру.

```typescript
const result = await invoiceAPI.rejectInvoice(
  "uuid-xxx",
  "Ошибка в сумме", // Обязательно
  "Требуется коррекция", // Опционально
);

console.log(result.status); // '50' (отклонен)
```

**Ошибки:**

- `400` - Reason is required

---

##### 7. `signInvoice(invoiceUuid: string, signature: string, certificateData?: string): Promise<InvoiceOperationResponse>`

Подписать счет электронной подписью.

```typescript
const result = await invoiceAPI.signInvoice(
  "uuid-xxx",
  "base64-encoded-signature",
  "certificate-data-base64", // Опционально
);
```

---

##### 8. `revokeInvoice(invoiceUuid: string, reason: string, comment?: string): Promise<InvoiceOperationResponse>`

Отозвать (отменить) счет.

```typescript
const result = await invoiceAPI.revokeInvoice(
  "uuid-xxx",
  "Отправлен по ошибке",
  "Требуется переоформление",
);

console.log(result.status); // '20' (отозван)
```

---

##### 9. `getInvoiceDetails(invoiceUuid: string): Promise<InvoiceDetail[]>`

Получить детали (строки) счета.

```typescript
const details = await invoiceAPI.getInvoiceDetails("uuid-xxx");

details.forEach((detail) => {
  console.log(detail.goodsName); // Название товара
  console.log(detail.baseCount); // Количество
  console.log(detail.price); // Цена
  console.log(detail.amount); // Сумма
});
```

---

##### 10. `downloadInvoicePDF(invoiceUuid: string): Promise<Blob>`

Скачать счет в PDF.

```typescript
const blob = await invoiceAPI.downloadInvoicePDF("uuid-xxx");

// Автоматически скачается (в компоненте используйте hook)
```

---

##### 11. `downloadInvoiceXML(invoiceUuid: string): Promise<Blob>`

Скачать счет в XML.

```typescript
const blob = await invoiceAPI.downloadInvoiceXML("uuid-xxx");
```

---

##### 12. `getInvoiceStatusHistory(invoiceUuid: string): Promise<Array>`

Получить историю статусов.

```typescript
const history = await invoiceAPI.getInvoiceStatusHistory("uuid-xxx");

history.forEach((entry) => {
  console.log(entry.status); // '10', '30', '40'
  console.log(entry.changedAt); // '2026-01-23T10:30:00Z'
  console.log(entry.changedBy); // 'user-id'
});
```

---

### Управление токеном

```typescript
// Установить токен
invoiceAPI.setToken("jwt-token-xxx");

// Очистить токен
invoiceAPI.clearToken();
```

---

## 🪝 React Hooks API

### `useInvoiceList(filters?, enabled?): UseQueryResult`

Получить список счетов с реактивностью.

```typescript
function InvoiceListPage() {
  const { data, isLoading, error, refetch } = useInvoiceList(
    { status: ESFDocumentStatus.SENT, limit: 20 },
    true  // enabled
  );

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      {data?.data.map(invoice => (
        <div key={invoice.documentUuid}>
          {invoice.invoiceNumber}
        </div>
      ))}
      <button onClick={() => refetch()}>Refresh</button>
    </div>
  );
}
```

---

### `useInvoiceDetail(invoiceUuid, enabled?): UseQueryResult`

Получить один счет.

```typescript
function InvoiceDetailPage({ invoiceId }: { invoiceId: string }) {
  const { data: invoice, isLoading } = useInvoiceDetail(invoiceId);

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>{invoice?.invoiceNumber}</h1>
      <p>Amount: {invoice?.totalAmount}</p>
    </div>
  );
}
```

---

### `useCreateInvoice(): UseMutationResult`

Создать счет с оптимистичными обновлениями.

```typescript
function CreateInvoiceForm() {
  const { mutate: createInvoice, isPending } = useCreateInvoice();

  const onSubmit = (formData: CreateInvoiceRequest) => {
    createInvoice(formData, {
      onSuccess: (invoice) => {
        toast.success(`Created: ${invoice.invoiceNumber}`);
        navigate(`/invoices/${invoice.documentUuid}`);
      },
      onError: (error) => {
        toast.error(error.message);
      }
    });
  };

  return (
    <form onSubmit={(e) => {
      e.preventDefault();
      onSubmit(getFormData());
    }}>
      {/* Form fields */}
      <button type="submit" disabled={isPending}>
        {isPending ? 'Creating...' : 'Create'}
      </button>
    </form>
  );
}
```

---

### `useAcceptInvoice(): UseMutationResult`

Принять счет.

```typescript
const { mutate: accept, isPending } = useAcceptInvoice();

await accept({
  invoiceUuid: "uuid-xxx",
  comment: "Accepted",
});
```

---

### `useRejectInvoice(): UseMutationResult`

Отклонить счет.

```typescript
const { mutate: reject, isPending } = useRejectInvoice();

await reject({
  invoiceUuid: "uuid-xxx",
  reason: "Wrong amount",
  comment: "Requires correction",
});
```

---

### `useSignInvoice(): UseMutationResult`

Подписать счет.

```typescript
const { mutate: sign } = useSignInvoice();

await sign({
  invoiceUuid: "uuid-xxx",
  signature: signatureBase64,
  certificateData: certBase64,
});
```

---

### `useRevokeInvoice(): UseMutationResult`

Отозвать счет.

```typescript
const { mutate: revoke } = useRevokeInvoice();

await revoke({
  invoiceUuid: "uuid-xxx",
  reason: "Sent by mistake",
});
```

---

### `useDownloadInvoicePDF(): UseMutationResult`

Скачать PDF (автоматически).

```typescript
const { mutate: downloadPDF, isPending } = useDownloadInvoicePDF();

const handleDownload = () => {
  downloadPDF("uuid-xxx");
  // Автоматически скачается файл invoice-uuid-xxx.pdf
};
```

---

### `useDownloadInvoiceXML(): UseMutationResult`

Скачать XML (автоматически).

```typescript
const { mutate: downloadXML } = useDownloadInvoiceXML();

downloadXML("uuid-xxx");
// Скачается файл invoice-uuid-xxx.xml
```

---

## ❌ Обработка ошибок

```typescript
import { InvoiceAPIError } from "@/lib/api/invoice";

try {
  const invoice = await invoiceAPI.getInvoice("uuid-xxx");
} catch (error) {
  if (error instanceof InvoiceAPIError) {
    console.log(error.statusCode); // 404
    console.log(error.message); // Invoice not found
    console.log(error.details); // Полные детали ошибки
  }
}
```

**Типичные коды:**

- `400` - Bad Request (валидация данных)
- `401` - Unauthorized (токен не установлен)
- `404` - Not Found (счет не найден)
- `409` - Conflict (уже существует)
- `500` - Server Error

---

## 🔐 Аутентификация

Токен устанавливается через `useAuthStore` автоматически в hooks:

```typescript
// hooks/useInvoice.ts
const { token } = useAuthStore();

invoiceAPI.setToken(token); // Автоматически в каждом hook
```

Или установить вручную:

```typescript
import { invoiceAPI } from "@/lib/api/invoice";

const token = localStorage.getItem("token");
invoiceAPI.setToken(token);
```

---

## 💾 Кэширование (React Query)

Автоматическое кэширование в hooks:

```typescript
// Первый запрос - с сервера
const { data: invoices1 } = useInvoiceList();

// Второй запрос - из кэша (5 минут)
const { data: invoices2 } = useInvoiceList();

// Одинаковые данные - без нового запроса!
```

Инвалидировать кэш:

```typescript
const queryClient = useQueryClient();

queryClient.invalidateQueries({
  queryKey: invoiceQueryKeys.lists(), // Очистить список
});
```

---

## 📝 Примеры

### Пример 1: Список счетов со статусом ОТПРАВЛЕНО

```typescript
function SentInvoicesList() {
  const { data: response } = useInvoiceList({
    status: ESFDocumentStatus.SENT,
    limit: 50
  });

  return (
    <div>
      <h2>Sent Invoices ({response?.meta.total_count})</h2>
      <table>
        <thead>
          <tr>
            <th>Number</th>
            <th>Contractor</th>
            <th>Amount</th>
            <th>Date</th>
          </tr>
        </thead>
        <tbody>
          {response?.data.map(invoice => (
            <tr key={invoice.documentUuid}>
              <td>{invoice.invoiceNumber}</td>
              <td>{invoice.contractor?.fullName}</td>
              <td>{invoice.totalAmount} {invoice.currency?.code}</td>
              <td>{invoice.invoiceDate}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

### Пример 2: Форма создания счета

```typescript
function CreateInvoiceForm() {
  const { mutate: create, isPending, error } = useCreateInvoice();
  const [formData, setFormData] = useState<CreateInvoiceRequest>({
    operationTypeCode: ESFOperationType.SALES,
    deliveryDate: new Date().toISOString().split('T')[0],
    contractorTin: '',
    deliveryTypeCode: ESFDeliveryType.DIRECT,
    paymentCode: ESFPaymentType.BANK_TRANSFER,
    currencyCode: CurrencyCode.KGS,
    isResident: true,
    isPriceWithoutTaxes: false,
    taxRateVATCode: '12',
    catalogEntries: [],
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    create(formData);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Contractor TIN"
        onChange={(e) =>
          setFormData({ ...formData, contractorTin: e.target.value })
        }
      />
      {error && <div className="error">{error.message}</div>}
      <button disabled={isPending}>{isPending ? 'Creating...' : 'Create'}</button>
    </form>
  );
}
```

### Пример 3: Действия над счетом

```typescript
function InvoiceActions({ invoiceUuid }: { invoiceUuid: string }) {
  const { mutate: accept } = useAcceptInvoice();
  const { mutate: reject } = useRejectInvoice();
  const { mutate: sign } = useSignInvoice();
  const { mutate: revoke } = useRevokeInvoice();

  return (
    <div className="actions">
      <button onClick={() => accept({ invoiceUuid })}>
        ✓ Accept
      </button>
      <button onClick={() => reject({
        invoiceUuid,
        reason: 'Incorrect data'
      })}>
        ✗ Reject
      </button>
      <button onClick={() => sign({
        invoiceUuid,
        signature: 'base64-sig'
      })}>
        🔐 Sign
      </button>
      <button onClick={() => revoke({
        invoiceUuid,
        reason: 'Sent by mistake'
      })}>
        🔄 Revoke
      </button>
    </div>
  );
}
```

---

## 🧪 Тестирование

```typescript
import { vi } from "vitest";
import { invoiceAPI } from "@/lib/api/invoice";

describe("Invoice API", () => {
  it("should create invoice", async () => {
    invoiceAPI.setToken("test-token");

    const invoice = await invoiceAPI.createInvoice({
      operationTypeCode: ESFOperationType.SALES,
      // ... другие поля
    });

    expect(invoice.documentUuid).toBeDefined();
  });

  it("should handle errors", async () => {
    invoiceAPI.setToken("invalid-token");

    expect(async () => {
      await invoiceAPI.getInvoice("invalid-uuid");
    }).rejects.toThrow("404");
  });
});
```

---

## 📋 Чеклист интеграции

- [ ] Установить токен при инициализации
- [ ] Обработать ошибки в UI
- [ ] Добавить loading состояния
- [ ] Настроить кэширование (React Query)
- [ ] Добавить toast уведомления
- [ ] Добавить валидацию форм
- [ ] Тестировать все операции
- [ ] Обновить типы при необходимости

---

**Версия:** 1.0  
**Последнее обновление:** 23 января 2026 г.
