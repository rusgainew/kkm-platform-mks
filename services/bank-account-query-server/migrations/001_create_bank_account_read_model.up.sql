-- Bank Account Query Server Read Model
-- Simplified tables for fast queries

CREATE TABLE IF NOT EXISTS bank_accounts (
    id BIGSERIAL PRIMARY KEY,
    account_name VARCHAR(255) NOT NULL,
    bank_account VARCHAR(30) NOT NULL UNIQUE,
    contractor_tin VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for fast queries
CREATE INDEX IF NOT EXISTS idx_bank_accounts_account_name ON bank_accounts(account_name);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_bank_account ON bank_accounts(bank_account);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_active ON bank_accounts(is_active);
