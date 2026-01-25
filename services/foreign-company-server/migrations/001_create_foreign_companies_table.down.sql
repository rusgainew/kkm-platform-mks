-- Удаление индексов
DROP INDEX IF EXISTS idx_foreign_companies_pin_unique;
DROP INDEX IF EXISTS idx_foreign_companies_created_at;
DROP INDEX IF EXISTS idx_foreign_companies_is_active;
DROP INDEX IF EXISTS idx_foreign_companies_country;
DROP INDEX IF EXISTS idx_foreign_companies_full_name;
DROP INDEX IF EXISTS idx_foreign_companies_pin;

-- Удаление таблицы
DROP TABLE IF EXISTS foreign_companies;
