# Bank Account Server

Микросервис для управления банковскими счетами организаций в системе ККМ.

## 🎯 Возможности

- ✅ Создание, обновление, удаление банковских счетов
- ✅ Список счетов организации с пагинацией
- ✅ Установка счета по умолчанию
- ✅ Валидация номеров счетов, IBAN, BIC
- ✅ Поддержка мультивалютных счетов (KZT, USD, EUR и др.)
- ✅ Event-driven архитектура (RabbitMQ)
- ✅ Полная observability (метрики, трассировка, логирование)
- ✅ Health checks
- ✅ Clean Architecture (domain/application/infrastructure/interfaces)

## 📊 Entities

### BankAccount (Банковский счет)

Банковский счет организации:

- **ID** - уникальный идентификатор
- **OrganizationID** - владелец счета
- **BankID** - ID банка из справочника
- **AccountNumber** - номер банковского счета (20 символов)
- **IBAN** - международный номер счета
- **Currency** - валюта счета (KZT, USD, EUR)
- **IsActive** - активен ли счет
- **IsDefault** - счет по умолчанию (только один на организацию)
- **BIC** - БИК банка
- **BankName** - название банка (денормализация)

## 🚀 Запуск

### С Docker Compose

```bash
# Запуск всего стека
docker-compose up -d bank-account-server

# Просмотр логов
docker-compose logs -f bank-account-server
```

### Локальная разработка

```bash
# Установка зависимостей
go mod download

# Запуск PostgreSQL и RabbitMQ
docker-compose up -d postgres rabbitmq

# Создание базы данных
psql -U postgres -h localhost -c "CREATE DATABASE bankaccountdb;"

# Применение миграций
psql -U postgres -h localhost -d bankaccountdb -f migrations/001_create_bank_accounts_table.up.sql

# Запуск сервера
go run cmd/main.go
```

### Настройка переменных окружения

Скопируйте `.env.example` в `.env` и настройте:

```bash
cp .env.example .env
```

Основные переменные:

- `BANK_ACCOUNT_SERVER_PORT=50054` - gRPC порт
- `BANK_ACCOUNT_SERVER_METRICS_PORT=9094` - Prometheus метрики
- `BANK_ACCOUNT_SERVER_DB_URL` - строка подключения к PostgreSQL
- `BANK_ACCOUNT_SERVER_RABBITMQ_URL` - строка подключения к RabbitMQ

## 🧪 Тестирование

### С grpcurl

```bash
# Проверка health check
grpcurl -plaintext localhost:50054 grpc.health.v1.Health/Check

# Создание банковского счета
grpcurl -plaintext -d '{
  "organization_id": "123e4567-e89b-12d3-a456-426614174000",
  "bank_id": "123e4567-e89b-12d3-a456-426614174002",
  "created_by": "123e4567-e89b-12d3-a456-426614174001",
  "account_number": "KZ86125KZT1001100100",
  "iban": "KZ86125KZT1001100100",
  "currency": "KZT",
  "bic": "HSBKKZKX",
  "bank_name": "Halyk Bank"
}' localhost:50054 api.BankAccountCommandService.CreateBankAccount

# Получение списка счетов
grpcurl -plaintext -d '{
  "organization_id": "123e4567-e89b-12d3-a456-426614174000",
  "page": 1,
  "per_page": 10
}' localhost:50054 api.BankAccountQueryService.GetBankAccountPage
```

## 📊 Observability

### Метрики Prometheus

Доступны на `http://localhost:9094/metrics`:

- `bank_accounts_total` - всего создано счетов
- `bank_accounts_active` - количество активных счетов
- `bank_account_operation_duration_seconds` - время выполнения операций
- `bank_account_operation_errors_total` - количество ошибок
- `bank_account_default_changes_total` - изменения счета по умолчанию

### Логирование

Структурированные JSON логи (production) или цветные консольные логи (development).

Уровни: debug, info, warn, error, fatal

## 🏗️ Архитектура

```
bank-account-server/
├── cmd/
│   └── main.go                          # Точка входа
├── internal/
│   ├── domain/                          # Бизнес-логика
│   │   ├── entities.go                  # BankAccount entity
│   │   ├── errors.go                    # Доменные ошибки
│   │   ├── events/events.go             # Доменные события
│   │   └── ports/                       # Интерфейсы
│   │       ├── repository.go            # Репозиторий
│   │       └── publisher.go             # Event publisher
│   ├── application/
│   │   └── bankaccount/service.go       # Бизнес-логика сервиса
│   ├── infrastructure/
│   │   ├── config/config.go             # Конфигурация
│   │   ├── repository/                  # PostgreSQL реализация
│   │   ├── messaging/                   # RabbitMQ publisher
│   │   └── observability/               # Метрики, логи
│   └── interfaces/
│       └── grpc/                        # gRPC обработчики
│           ├── bankaccount_handler.go   # Command/Query handlers
│           └── health_handler.go        # Health check
├── migrations/                          # SQL миграции
├── .env.example                         # Пример конфигурации
├── Dockerfile                           # Docker образ
├── go.mod                               # Go зависимости
└── README.md                            # Эта документация
```

## 🔄 События RabbitMQ

Сервис публикует события в exchange `bank_account.events`:

### bank_account.created

```json
{
  "event_id": "uuid",
  "event_type": "bank_account.created",
  "timestamp": "2026-01-08T12:00:00Z",
  "account_id": "uuid",
  "organization_id": "uuid",
  "account_number": "KZ86125KZT1001100100",
  "currency": "KZT",
  "bank_name": "Halyk Bank",
  "created_by": "uuid"
}
```

### bank_account.updated

```json
{
  "event_id": "uuid",
  "event_type": "bank_account.updated",
  "timestamp": "2026-01-08T12:00:00Z",
  "account_id": "uuid",
  "organization_id": "uuid",
  "account_number": "KZ86125KZT1001100100",
  "currency": "KZT",
  "updated_by": "uuid"
}
```

### bank_account.deleted

```json
{
  "event_id": "uuid",
  "event_type": "bank_account.deleted",
  "timestamp": "2026-01-08T12:00:00Z",
  "account_id": "uuid",
  "organization_id": "uuid",
  "deleted_by": "uuid"
}
```

### bank_account.set_default

```json
{
  "event_id": "uuid",
  "event_type": "bank_account.set_default",
  "timestamp": "2026-01-08T12:00:00Z",
  "account_id": "uuid",
  "organization_id": "uuid",
  "updated_by": "uuid"
}
```

## 🔐 Безопасность

- Валидация форматов номеров счетов (Казахстанский стандарт)
- Валидация IBAN (международный стандарт)
- Валидация BIC кодов
- Проверка прав доступа (organization_id)
- SQL-инъекции защищены (prepared statements)
- Структурированные ошибки (без раскрытия внутренней информации)

## ✅ Бизнес-правила

1. **Номер счета должен быть уникальным** в рамках организации
2. **Только один счет может быть по умолчанию** для организации
3. **Нельзя удалить счет по умолчанию**, если есть другие активные счета
4. **При установке нового счета по умолчанию** флаг снимается со старого

## 🤝 Интеграция с другими сервисами

### С invoice-server

Invoice-server использует bank-account-server для получения банковских реквизитов при создании счетов-фактур.

### С company-server

Проверка принадлежности банковских счетов к организациям.

## 📝 TODO

- [ ] Добавить интеграционные тесты
- [ ] Добавить кэширование (Redis)
- [ ] Добавить историю изменений счетов
- [ ] Добавить валидацию через внешний банковский API
- [ ] Добавить проверку остатков на счетах

## 📚 Документация

- [Proto файлы](../proto/api/)
- [OWASP Security Checklist](../security/owasp-top10-checklist.md)
- [Copilot Instructions](../.github/copilot-instructions.md)

---

**Версия:** 1.0.0  
**Статус:** ✅ Готов к разработке  
**Контакт:** support@company.com
