# Database Migrations

Эта директория содержит миграции базы данных для company-server.

## Структура миграций

Миграции именуются в формате: `{version}_{description}.{up|down}.sql`

- `up.sql` - применение миграции
- `down.sql` - откат миграции

## Применение миграций

### Вручную через psql

```bash
# Применить миграцию
psql -h localhost -U postgres -d company_db -f migrations/001_create_organizations_table.up.sql

# Откатить миграцию
psql -h localhost -U postgres -d company_db -f migrations/001_create_organizations_table.down.sql
```

### Создание базы данных

```bash
# Создать базу данных
createdb -h localhost -U postgres company_db

# Или через psql
psql -h localhost -U postgres -c "CREATE DATABASE company_db;"
```

## Список миграций

1. **001_create_organizations_table** - Создание таблиц organizations и employees
   - Таблица `organizations` для хранения организаций
   - Таблица `employees` для хранения участников организаций
   - Индексы для оптимизации запросов
   - Foreign key constraints

## Схема базы данных

### organizations

- `id` (VARCHAR(36)) - UUID организации
- `name` (VARCHAR(255)) - название организации (уникальное)
- `description` (TEXT) - описание
- `owner_id` (VARCHAR(36)) - ID владельца
- `created_at` (TIMESTAMP) - дата создания
- `updated_at` (TIMESTAMP) - дата обновления

### employees

- `id` (VARCHAR(36)) - UUID записи
- `organization_id` (VARCHAR(36)) - ID организации (FK)
- `user_id` (VARCHAR(36)) - ID пользователя
- `role` (VARCHAR(50)) - роль (owner, admin, member, viewer)
- `joined_at` (TIMESTAMP) - дата добавления
- Уникальное ограничение: (organization_id, user_id)
