-- Создание таблицы для банковских счетов
CREATE SCHEMA IF NOT EXISTS bank_accounts;
SET search_path TO bank_accounts;

CREATE TABLE IF NOT EXISTS bank_accounts (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    bank_id UUID NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    iban VARCHAR(34),
    currency VARCHAR(3) NOT NULL DEFAULT 'KZT',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    bic VARCHAR(11),
    bank_name VARCHAR(300) NOT NULL,
    created_by UUID NOT NULL,
    updated_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_bank_accounts_organization ON bank_accounts(organization_id);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_bank ON bank_accounts(bank_id);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_active ON bank_accounts(is_active);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_default ON bank_accounts(organization_id, is_default);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_created_at ON bank_accounts(created_at DESC);

-- Уникальный индекс для номера счета внутри организации
CREATE UNIQUE INDEX IF NOT EXISTS idx_bank_accounts_org_number_unique ON bank_accounts(organization_id, account_number);

-- Уникальный индекс для счета по умолчанию (только один на организацию)
CREATE UNIQUE INDEX IF NOT EXISTS idx_bank_accounts_org_default_unique ON bank_accounts(organization_id) WHERE is_default = true;

-- Комментарии к таблице
COMMENT ON TABLE bank_accounts IS 'Банковские счета организаций';
COMMENT ON COLUMN bank_accounts.id IS 'Уникальный идентификатор счета';
COMMENT ON COLUMN bank_accounts.organization_id IS 'ID организации-владельца';
COMMENT ON COLUMN bank_accounts.bank_id IS 'ID банка из справочника';
COMMENT ON COLUMN bank_accounts.account_number IS 'Номер банковского счета';
COMMENT ON COLUMN bank_accounts.iban IS 'Международный номер счета (IBAN)';
COMMENT ON COLUMN bank_accounts.currency IS 'Код валюты (KZT, USD, EUR и т.д.)';
COMMENT ON COLUMN bank_accounts.is_active IS 'Активен ли счет';
COMMENT ON COLUMN bank_accounts.is_default IS 'Счет по умолчанию для организации';
COMMENT ON COLUMN bank_accounts.bic IS 'БИК банка';
COMMENT ON COLUMN bank_accounts.bank_name IS 'Название банка';
