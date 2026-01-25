# Catalog Server

Микросервис для управления каталогом товаров и услуг в системе ККМ.

## 🎯 Возможности

- ✅ Создание, обновление, удаление элементов каталога
- ✅ Список элементов с пагинацией
- ✅ Поиск по имени, коду и описанию
- ✅ Поддержка ТНВЭД/ГКЭД классификации
- ✅ Управление ценами и ставками НДС
- ✅ Event-driven архитектура (RabbitMQ)
- ✅ Полная observability (метрики, трассировка, логирование)
- ✅ Health checks
- ✅ Clean Architecture (domain/application/infrastructure/interfaces)

## 📊 Entities

### CatalogItem (Элемент каталога)

Товар или услуга в каталоге организации:

- **ID** - уникальный идентификатор
- **OrganizationID** - владелец элемента
- **Name** - наименование товара/услуги
- **Code** - артикул/код (уникален в рамках организации)
- **Description** - подробное описание
- **UnitType** - единица измерения (шт, кг, л, м и т.д.)
- **Price** - цена за единицу
- **VATRate** - ставка НДС (0%, 12%, 20%)
- **Category** - категория товара/услуги
- **IsActive** - активен ли элемент
- **TNVED** - код ТНВЭД (для товаров)
- **GKED** - код ГКЭД (для услуг)
- **Barcode** - штрих-код

## 🚀 Запуск

### С Docker Compose

```bash
# Запуск всего стека
docker-compose up -d catalog-server

# Просмотр логов
docker-compose logs -f catalog-server
```

### Локальная разработка

```bash
# Установка зависимостей
go mod download

# Запуск PostgreSQL и RabbitMQ
docker-compose up -d postgres rabbitmq

# Создание базы данных
psql -U postgres -h localhost -c "CREATE DATABASE catalogdb;"

# Применение миграций
psql -U postgres -h localhost -d catalogdb -f migrations/001_create_catalog_items_table.up.sql

# Запуск сервера
go run cmd/main.go
```

### Настройка переменных окружения

Скопируйте `.env.example` в `.env` и настройте:

```bash
cp .env.example .env
```

Основные переменные:

- `CATALOG_SERVER_PORT=50053` - gRPC порт
- `CATALOG_SERVER_METRICS_PORT=9093` - Prometheus метрики
- `CATALOG_SERVER_DB_URL` - строка подключения к PostgreSQL
- `CATALOG_SERVER_RABBITMQ_URL` - строка подключения к RabbitMQ

## 🧪 Тестирование

### С grpcurl

```bash
# Проверка health check
grpcurl -plaintext localhost:50053 grpc.health.v1.Health/Check

# Создание элемента каталога
grpcurl -plaintext -d '{
  "organization_id": "123e4567-e89b-12d3-a456-426614174000",
  "created_by": "123e4567-e89b-12d3-a456-426614174001",
  "name": "Товар #1",
  "code": "TOV001",
  "description": "Описание товара",
  "unit_type": "шт",
  "price": 1000.0,
  "vat_rate": 20.0
}' localhost:50053 api.CatalogCommandService.CreateCatalogItem

# Получение списка элементов
grpcurl -plaintext -d '{
  "organization_id": "123e4567-e89b-12d3-a456-426614174000",
  "page": 1,
  "per_page": 10
}' localhost:50053 api.CatalogQueryService.GetCatalogPage

# Поиск элементов
grpcurl -plaintext -d '{
  "organization_id": "123e4567-e89b-12d3-a456-426614174000",
  "search_query": "товар",
  "page": 1,
  "per_page": 10
}' localhost:50053 api.CatalogQueryService.GetCatalogPage
```

## 📊 Observability

### Метрики Prometheus

Доступны на `http://localhost:9093/metrics`:

- `catalog_items_total` - всего создано элементов
- `catalog_items_active` - количество активных элементов
- `catalog_operation_duration_seconds` - время выполнения операций
- `catalog_operation_errors_total` - количество ошибок
- `catalog_search_duration_seconds` - время поиска

### Логирование

Структурированные JSON логи (production) или цветные консольные логи (development).

Уровни: debug, info, warn, error, fatal

## 🏗️ Архитектура

```
catalog-server/
├── cmd/
│   └── main.go                          # Точка входа
├── internal/
│   ├── domain/                          # Бизнес-логика
│   │   ├── entities.go                  # CatalogItem entity
│   │   ├── errors.go                    # Доменные ошибки
│   │   ├── events/events.go             # Доменные события
│   │   └── ports/                       # Интерфейсы
│   │       ├── repository.go            # Репозиторий
│   │       └── publisher.go             # Event publisher
│   ├── application/
│   │   └── catalog/service.go           # Бизнес-логика сервиса
│   ├── infrastructure/
│   │   ├── config/config.go             # Конфигурация
│   │   ├── repository/                  # PostgreSQL реализация
│   │   ├── messaging/                   # RabbitMQ publisher
│   │   └── observability/               # Метрики, логи
│   └── interfaces/
│       └── grpc/                        # gRPC обработчики
│           ├── catalog_handler.go       # Command/Query handlers
│           └── health_handler.go        # Health check
├── migrations/                          # SQL миграции
├── .env.example                         # Пример конфигурации
├── Dockerfile                           # Docker образ
├── go.mod                               # Go зависимости
└── README.md                            # Эта документация
```

## 🔄 События RabbitMQ

Сервис публикует события в exchange `catalog.events`:

### catalog.item.created

```json
{
  "event_id": "uuid",
  "event_type": "catalog.item.created",
  "timestamp": "2026-01-08T12:00:00Z",
  "item_id": "uuid",
  "organization_id": "uuid",
  "name": "Товар #1",
  "code": "TOV001",
  "price": 1000.0,
  "created_by": "uuid"
}
```

### catalog.item.updated

```json
{
  "event_id": "uuid",
  "event_type": "catalog.item.updated",
  "timestamp": "2026-01-08T12:00:00Z",
  "item_id": "uuid",
  "organization_id": "uuid",
  "name": "Товар #1 (обновлен)",
  "price": 1200.0,
  "updated_by": "uuid"
}
```

### catalog.item.deleted

```json
{
  "event_id": "uuid",
  "event_type": "catalog.item.deleted",
  "timestamp": "2026-01-08T12:00:00Z",
  "item_id": "uuid",
  "organization_id": "uuid",
  "deleted_by": "uuid"
}
```

## 🔐 Безопасность

- Валидация всех входных данных
- Проверка прав доступа (organization_id)
- SQL-инъекции защищены (prepared statements)
- Структурированные ошибки (без раскрытия внутренней информации)

## 🤝 Интеграция с другими сервисами

### С invoice-server

Invoice-server использует catalog-server для получения информации о товарах/услугах при создании счетов-фактур.

### С company-server

Проверка принадлежности элементов каталога к организациям.

## 📝 TODO

- [ ] Добавить интеграционные тесты
- [ ] Добавить кэширование (Redis)
- [ ] Добавить импорт/экспорт каталога (Excel/CSV)
- [ ] Добавить историю изменений цен
- [ ] Добавить групповые операции (массовое обновление)

## 📚 Документация

- [Proto файлы](../proto/api/)
- [OWASP Security Checklist](../security/owasp-top10-checklist.md)
- [Copilot Instructions](../.github/copilot-instructions.md)

---

**Версия:** 1.0.0  
**Статус:** ✅ Готов к разработке
**Контакт:** support@company.com
