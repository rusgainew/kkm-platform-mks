# Быстрый гайд: Использование ESF API типов

**Файл:** `/kkm-platform/QUICK_START_ESF_TYPES.md`  
**Версия:** 1.0  
**Дата:** 23 января 2026 г.

---

## 📦 Импорт типов

```typescript
// Все типы доступны из одного места
import type {
  // Справочники
  ReferenceItem,
  VatTaxType,
  TaxInfo,
  CurrencyReference,

  // Стороны сделки
  Party,
  LegalPerson,
  ForeignParty,
  ContractParty,

  // Счета-фактуры
  Invoice,
  InvoiceDetail,
  CreateInvoiceRequest,
  UpdateInvoiceRequest,
  AcceptOrRejectInvoiceRequest,
  SignInvoiceRequest,
  RevokeInvoiceRequest,
  InvoiceFilters,
  InvoiceOperationResponse,

  // ESF API ответы
  ESFAPIResponse,
  ESFCreateInvoiceResponse,
  ESFInvoiceResponse,
  ESFInvoiceListResponse,

  // Перечисления
  CurrencyCode,
  ESFDocumentStatus,
  ESFDocumentType,
  ESFOperationType,
  ESFDeliveryType,
  ESFPaymentType,
  ESFTaxType,
  ESFInvoiceAction,
  ESFUnitClassification,
} from "@/types";

// Вспомогательные функции
import {
  createParty,
  createForeignParty,
  isPartyResident,
  isForeignParty,
  getPartyId,
  createReferenceItem,
  createTaxInfo,
  createVatTaxType,
  createInvoiceDetail,
  createCatalogEntry,
  adaptESFResponse,
} from "@/types";

// Маппинги и константы
import {
  ESF_STATUS_LABELS,
  OPERATION_TYPE_LABELS,
  ESF_DOCUMENT_TYPE_LABELS,
  DELIVERY_TYPE_LABELS,
  PAYMENT_TYPE_LABELS,
  UNIT_CLASSIFICATION_LABELS,
  ESF_CURRENCY_CODES,
  ESF_COUNTRY_CODES,
  CURRENCY_CODES,
} from "@/types";
```

---

## 🎯 Примеры использования

### 1. Создание сторон сделки

```typescript
// Резидент - местная компания
const seller: Party = createParty(
  "01206200110146", // ИНН
  "ООО Экспортер", // Наименование
  "г.Бишкек, ул.Чуй 123", // Адрес
  "01206200110100", // ИНН главной организации (если филиал)
);

// Покупатель (резидент)
const buyer: Party = createParty("01206200110100", "АО Импортер");

// Иностранная компания
const foreignBuyer: ForeignParty = createForeignParty(
  "LLC International Trade",
  "US", // Код страны
  "123 Main Street, NY",
);

// Определить тип стороны
if (isPartyResident(seller)) {
  console.log("Это резидент КР");
}

if (isForeignParty(foreignBuyer)) {
  console.log("Это иностранная компания");
}

// Получить идентификатор
const id = getPartyId(seller); // '01206200110146'
```

---

### 2. Создание деталей счета-фактуры

```typescript
// Способ 1: Используя вспомогательную функцию (с автоматическим расчетом НДС)
const detail1: InvoiceDetail = createInvoiceDetail(
  "invoice-uuid-xxx", // UUID счета
  "Подшипник 180200 FAG", // Название
  100, // Количество
  250, // Цена за единицу
  "796", // Единица (штука)
  "12", // НДС % (опционально)
);
// Результат: НДС автоматически рассчитано

// Способ 2: Полное определение
const detail2: InvoiceDetail = {
  invoiceUuid: "invoice-uuid-xxx",
  baseCount: 50,
  price: 1000,
  amount: 50000,
  amountWithoutVAT: 44642.86,
  amountVAT: 5357.14,
  goodsName: "Услуга консультирования",
  tnvedCode: "8482109008",
  gked: "7210Z",
  unitClassification: {
    code: "255", // Час
    name: "Час",
  },
  goodsType: {
    code: "2",
    name: "Услуга",
  },
};
```

---

### 3. Создание счета-фактуры

```typescript
// Подготовка товаров из каталога
const catalogEntries = [
  createCatalogEntry(
    2047, // ID из каталога
    100, // Количество
    250, // Цена
    "796", // Единица (штука)
    "10", // Код налога
    12, // НДС %
  ),
  createCatalogEntry(
    2048,
    50,
    1000,
    "255", // Час
    "10",
    12,
  ),
];

// Создание счета-фактуры
const createRequest: CreateInvoiceRequest = {
  // Основные данные
  operationTypeCode: ESFOperationType.SALES, // Реализация
  invoiceNumber: "СФ-001-2026",
  invoiceDate: "2026-01-23",
  deliveryDate: "2026-01-25",

  // Стороны сделки
  contractorTin: "01206200110100", // ИНН покупателя
  supplierBankAccount: "KG1234567890",
  contractorBankAccount: "KG0987654321",

  // Справочники
  deliveryTypeCode: ESFDeliveryType.DIRECT, // Прямая доставка
  paymentCode: ESFPaymentType.BANK_TRANSFER, // Банковский перевод
  currencyCode: CurrencyCode.KGS, // Киргизский сом
  countryCode: "KG",
  taxRateVATCode: "12",

  // Флаги
  isResident: true,
  isPriceWithoutTaxes: false,
  isIndustry: false,

  // Товары
  catalogEntries,
};

// Отправка на создание
try {
  // const response = await invoiceAPI.createInvoice(createRequest);
  console.log("Счет создан:", createRequest);
} catch (error) {
  console.error("Ошибка создания:", error);
}
```

---

### 4. Работа со статусами

```typescript
// Использование enum
const status: ESFDocumentStatus = ESFDocumentStatus.SENT;

// Получение названия для UI
const label = ESF_STATUS_LABELS[status]; // "Отправлен"

// Маппинг всех статусов
const allStatuses = [
  {
    code: ESFDocumentStatus.DRAFT,
    label: ESF_STATUS_LABELS[ESFDocumentStatus.DRAFT],
  },
  {
    code: ESFDocumentStatus.SENT,
    label: ESF_STATUS_LABELS[ESFDocumentStatus.SENT],
  },
  {
    code: ESFDocumentStatus.ACCEPTED,
    label: ESF_STATUS_LABELS[ESFDocumentStatus.ACCEPTED],
  },
  {
    code: ESFDocumentStatus.REJECTED,
    label: ESF_STATUS_LABELS[ESFDocumentStatus.REJECTED],
  },
];

// Использование в фильтрах
const filters: InvoiceFilters = {
  status: ESFDocumentStatus.SENT,
  dateFrom: "2026-01-01",
  dateTo: "2026-01-31",
  currency: CurrencyCode.KGS,
  limit: 50,
};
```

---

### 5. Операции со счетами

```typescript
// Принять счет
const acceptRequest: AcceptOrRejectInvoiceRequest = {
  invoiceUuid: "invoice-uuid-xxx",
  action: "accept",
  comment: "Принят без замечаний",
};

// Отклонить счет
const rejectRequest: AcceptOrRejectInvoiceRequest = {
  invoiceUuid: "invoice-uuid-xxx",
  action: "reject",
  reason: "Ошибка в сумме",
  comment: "Требуется коррекция",
};

// Подписать счет электронной подписью
const signRequest: SignInvoiceRequest = {
  invoiceUuid: "invoice-uuid-xxx",
  signature: "base64-encoded-signature",
  certificateData: "cert-data-base64",
  timestamp: Date.now(),
};

// Отозвать счет
const revokeRequest: RevokeInvoiceRequest = {
  invoiceUuid: "invoice-uuid-xxx",
  reason: "Отправлен по ошибке",
  comment: "Требуется переоформление",
};
```

---

### 6. Обработка ESF API ответов

```typescript
import { adaptESFResponse } from "@/types";

// ESF API возвращает ответ в своём формате
const esfResponse: ESFAPIResponse<Invoice> = {
  responseId: "resp-123",
  requestId: "req-456",
  data: {
    documentUuid: "doc-uuid",
    invoiceNumber: "СФ-001",
    // ... остальные поля
  },
};

// Преобразовать в наш стандартный формат
const standardResponse = adaptESFResponse(esfResponse);

// Теперь работать с обычными функциями
if (isSuccessResponse(standardResponse)) {
  console.log("Успех:", standardResponse.data);
} else {
  console.error("Ошибка:", standardResponse.error);
}
```

---

### 7. React Hook пример

```typescript
import { useState } from 'react';
import type { Invoice, InvoiceFilters } from '@/types';
import { ESFDocumentStatus, ESF_STATUS_LABELS } from '@/types';

export function InvoiceListComponent() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [filters, setFilters] = useState<InvoiceFilters>({
    status: ESFDocumentStatus.SENT,
    limit: 20
  });

  const handleStatusChange = (newStatus: ESFDocumentStatus) => {
    setFilters(prev => ({
      ...prev,
      status: newStatus
    }));
  };

  return (
    <div>
      <select onChange={(e) => handleStatusChange(e.target.value as ESFDocumentStatus)}>
        <option value={ESFDocumentStatus.DRAFT}>
          {ESF_STATUS_LABELS[ESFDocumentStatus.DRAFT]}
        </option>
        <option value={ESFDocumentStatus.SENT}>
          {ESF_STATUS_LABELS[ESFDocumentStatus.SENT]}
        </option>
        <option value={ESFDocumentStatus.ACCEPTED}>
          {ESF_STATUS_LABELS[ESFDocumentStatus.ACCEPTED]}
        </option>
      </select>

      {invoices.map(invoice => (
        <div key={invoice.documentUuid}>
          <h3>{invoice.invoiceNumber}</h3>
          <p>
            Статус: {ESF_STATUS_LABELS[invoice.status?.code as ESFDocumentStatus] || 'Неизвестно'}
          </p>
          <p>Сумма: {invoice.totalAmount} {invoice.currency?.code}</p>
        </div>
      ))}
    </div>
  );
}
```

---

### 8. Форма создания счета

```typescript
import { useForm } from 'react-hook-form';
import type { CreateInvoiceRequest } from '@/types';
import { CurrencyCode, ESFDeliveryType, ESFPaymentType, ESFOperationType } from '@/types';

export function CreateInvoiceForm() {
  const { register, handleSubmit, formState: { errors } } = useForm<CreateInvoiceRequest>({
    defaultValues: {
      operationTypeCode: ESFOperationType.SALES,
      currencyCode: CurrencyCode.KGS,
      deliveryTypeCode: ESFDeliveryType.DIRECT,
      paymentCode: ESFPaymentType.BANK_TRANSFER,
      isResident: true,
      isPriceWithoutTaxes: false,
      taxRateVATCode: '12',
      catalogEntries: []
    }
  });

  const onSubmit = (data: CreateInvoiceRequest) => {
    // Отправить на сервер
    console.log('Отправка:', data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input
        {...register('invoiceNumber', { required: true })}
        placeholder="Номер счета"
      />

      <input
        {...register('contractorTin', { required: true })}
        placeholder="ИНН покупателя"
      />

      <input
        type="date"
        {...register('deliveryDate', { required: true })}
      />

      <select {...register('currencyCode')}>
        <option value={CurrencyCode.KGS}>KGS - Киргизский сом</option>
        <option value={CurrencyCode.RUB}>RUB - Российский рубль</option>
        <option value={CurrencyCode.USD}>USD - Американский доллар</option>
      </select>

      <button type="submit">Создать счет</button>
    </form>
  );
}
```

---

## 🔗 Связанные файлы

- **Документация:** `/doc-bec-api/TYPESCRIPT_TYPES_AUDIT.md`
- **API инструкции:** `/doc-bec-api/ESF-API-INSTRUCTION.md`
- **Маппинг типов:** `/doc-bec-api/TYPE_MAPPING.md`
- **Отчет реализации:** `/kkm-platform/IMPLEMENTATION_REPORT_PRIORITY_1.md`

---

## 📝 Заметки

- Валюта `KGS` (Киргизский сом) - **основная валюта для всех операций в КР**
- Все коды статусов используют строки, но значения - числовые коды ESF
- `Party` для местных компаний (резидентов), `ForeignParty` для иностранцев
- Все суммы в денежных единицах (сомах или другой валюте)
- НДС обычно 12% для стандартной ставки (0% для некоторых категорий)

---

**Версия:** 1.0  
**Последнее обновление:** 23 января 2026 г.
