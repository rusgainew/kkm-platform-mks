-- Удаление индексов
DROP INDEX IF EXISTS idx_catalog_items_org_code_unique;
DROP INDEX IF EXISTS idx_catalog_items_created_at;
DROP INDEX IF EXISTS idx_catalog_items_is_active;
DROP INDEX IF EXISTS idx_catalog_items_category;
DROP INDEX IF EXISTS idx_catalog_items_name;
DROP INDEX IF EXISTS idx_catalog_items_code;
DROP INDEX IF EXISTS idx_catalog_items_organization;

-- Удаление таблицы
DROP TABLE IF EXISTS catalog_items;
