-- Удаление индексов
DROP INDEX IF EXISTS idx_bank_accounts_org_default_unique;
DROP INDEX IF EXISTS idx_bank_accounts_org_number_unique;
DROP INDEX IF EXISTS idx_bank_accounts_created_at;
DROP INDEX IF EXISTS idx_bank_accounts_is_default;
DROP INDEX IF EXISTS idx_bank_accounts_is_active;
DROP INDEX IF EXISTS idx_bank_accounts_bank;
DROP INDEX IF EXISTS idx_bank_accounts_organization;

-- Удаление таблицы
DROP TABLE IF EXISTS bank_accounts;
