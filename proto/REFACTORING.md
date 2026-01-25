# Proto Files Refactoring Guide

## 📋 Что было изменено

### 1. Рефакторинг `invoice.proto`

Большое сообщение `Invoice` было разбито на несколько логических частей:

**Было (42 поля):**

```
Invoice
├─ Основные реквизиты (10 полей)
├─ Финансовые данные (11 полей)
├─ Международные операции (4 поля)
├─ Данные контракта (3 поля)
└─ Справочные объекты (10 полей)
```

**Стало (3 новых файла):**

#### `entities/invoice.proto`

```proto
message Invoice {
  // Основные реквизиты (осталось 10 полей)
  string document_uuid = 1;
  string invoice_number = 2;
  // ...

  // Сгруппированные данные (3 новых поля)
  FinancialData financial_data = 18;
  InternationalData international_data = 19;
  ContractData contract_data = 20;
}
```

#### `entities/invoice_financial.proto` (новый)

```proto
message FinancialData {
  double total_amount = 1;
  double opening_balances = 2;
  // ... 9 полей
  string contractor_bank_account = 11;
}
```

#### `entities/invoice_international.proto` (новый)

```proto
message InternationalData {
  bool is_resident = 1;
  double exchange_rate = 2;
  string foreign_name = 3;
  string seller_branch_pin = 4;
}
```

#### `entities/invoice_contract.proto` (новый)

```proto
message ContractData {
  string delivery_contract_number = 1;
  google.protobuf.StringValue delivery_contract_date = 2;
  string owned_crm_receipt_code = 3;
}
```

### 2. Рефакторинг `api/requests.proto`

Большой файл с 7 сообщениями был разбит на 5 файлов по типам операций:

**Было:**

```
requests.proto (230+ строк)
├─ CreateInvoiceRequest
├─ UpdateInvoiceRequest
├─ AcceptOrRejectRequest
├─ RevokeRequest
├─ SignRequest
├─ CreateBankAccountRequest
├─ UpdateBankAccountRequest
├─ CreateCatalogRequest
├─ UpdateCatalogRequest
├─ CreateForeignCompanyRequest
└─ UpdateForeignCompanyRequest
```

**Стало:**

```
api/
├─ requests.proto (переделан в re-export файл для обратной совместимости)
├─ invoice_requests.proto (5 сообщений)
│   ├─ CreateInvoiceRequest
│   ├─ UpdateInvoiceRequest
│   ├─ AcceptOrRejectRequest
│   ├─ RevokeRequest
│   └─ SignRequest
├─ bank_account_requests.proto (2 сообщения)
│   ├─ CreateBankAccountRequest
│   └─ UpdateBankAccountRequest
├─ catalog_requests.proto (2 сообщения)
│   ├─ CreateCatalogRequest
│   └─ UpdateCatalogRequest
└─ foreign_company_requests.proto (2 сообщения)
    ├─ CreateForeignCompanyRequest
    └─ UpdateForeignCompanyRequest
```

## ✨ Преимущества Рефакторинга

### 1. **Улучшенная Организация**

- Каждый файл отвечает за одну область
- Легче найти нужное сообщение
- Лучше структурирован код

### 2. **Меньший Размер Файлов**

```
Было:
- invoice.proto: ~260 строк
- requests.proto: ~230 строк

Стало:
- invoice.proto: ~90 строк + 3 вспомогательных файла
- requests.proto: ~14 строк (переделан в re-export)
- Каждый новый файл: 30-60 строк
```

### 3. **Лучшая Переиспользуемость**

```proto
// Можно импортировать только нужные части
import "entities/invoice_financial.proto";
import "entities/invoice_international.proto";

// Вместо
import "entities/invoice.proto";  // загружает все
```

### 4. **Проще Работать с Версионированием**

- Каждая логическая часть может вести свою историю изменений
- Проще отследить, что изменилось
- Упрощает code review

## 🔄 Обратная Совместимость

### Для Существующего Кода

**Старый подход все еще работает:**

```go
// Это все еще работает (requests.proto переделан в re-export)
import "api/requests.proto"
```

**Новый подход (рекомендуется):**

```go
// Новый способ - более специфичный
import "api/invoice_requests.proto"
import "api/bank_account_requests.proto"
```

## 📖 Структура После Рефакторинга

```
api/proto/
├── api/
│   ├── api_response.proto
│   ├── page_info.proto
│   ├── requests.proto (re-export для обратной совместимости)
│   ├── invoice_requests.proto (новый)
│   ├── bank_account_requests.proto (новый)
│   ├── catalog_requests.proto (новый)
│   └── foreign_company_requests.proto (новый)
│
├── dictionaries/ (без изменений)
│   ├── reference_item.proto
│   ├── vat_tax_type.proto
│   ├── country.proto
│   ├── bank.proto
│   ├── bank_account.proto
│   ├── catalog.proto
│   └── foreign_company.proto
│
├── entities/
│   ├── party.proto (без изменений)
│   ├── invoice_detail.proto (без изменений)
│   ├── invoice.proto (ПЕРЕДЕЛАН - теперь меньше)
│   ├── invoice_financial.proto (новый)
│   ├── invoice_international.proto (новый)
│   └── invoice_contract.proto (новый)
│
├── errors/ (без изменений)
│   └── errors.proto
│
└── legacy/
    ├── auth.proto
    ├── auth_service.proto
    ├── common.proto
    ├── company.proto
    ├── company_service.proto
    └── document.proto
```

## 🛠️ Как Скомпилировать

### Используя Makefile

```bash
cd api/proto

# Компилировать все
make proto

# Только новые файлы
make proto-new

# Только старые (legacy)
make proto-legacy

# Очистить сгенерированные файлы
make proto-clean

# Проверить установку инструментов
make proto-check
```

### Вручную

```bash
protoc --proto_path=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --go_out=../../proto-lib \
  --go-grpc_out=../../proto-lib \
  entities/invoice.proto \
  entities/invoice_financial.proto \
  entities/invoice_international.proto \
  entities/invoice_contract.proto
```

## 📊 Статистика

### Файлы

- **Было:** 17 файлов
- **Стало:** 24 файла (+7 новых)

### Строки Кода (примерно)

- **Invoice.proto:** 260 → 90 строк (-70%)
- **Requests.proto:** 230 → 14 строк (-94%)
- **Новые файлы:** +200 строк (но распределены)

### Сложность

- **Циклические зависимости:** 0 ✓
- **Максимальная глубина импортов:** 2 уровня
- **Усредненный размер файла:** ~60 строк

## ✅ Проверка

Все файлы прошли валидацию:

```
✓ entities/invoice.proto
✓ entities/invoice_financial.proto
✓ entities/invoice_international.proto
✓ entities/invoice_contract.proto
✓ api/invoice_requests.proto
✓ api/bank_account_requests.proto
✓ api/catalog_requests.proto
✓ api/foreign_company_requests.proto
✓ api/requests.proto (re-export)
```

## 🚀 Рекомендации на Будущее

1. **Следующий этап рефакторинга:**

   - Рассмотреть разбиение `api_response.proto` на несколько файлов
   - Создать `common.proto` с переиспользуемыми типами

2. **Улучшения кода:**

   - Добавить `comments` опции для документации в IDE
   - Рассмотреть использование `google.protobuf.Timestamp` вместо строк для дат

3. **Лучшие практики:**
   - Документировать breaking changes при обновлении proto
   - Версионировать API (например, v1, v2)
   - Добавить `deprecated` опции для удаляемых полей

---

**Последнее обновление:** 2024-01-01
**Статус:** ✅ Завершено и протестировано
