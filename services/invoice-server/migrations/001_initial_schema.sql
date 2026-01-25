-- Invoice tables schema

-- Main invoice table
CREATE TABLE IF NOT EXISTS invoices (
    id SERIAL PRIMARY KEY,
    document_uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) NOT NULL,
    invoice_date TIMESTAMP,
    delivery_date TIMESTAMP,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    is_resident BOOLEAN NOT NULL DEFAULT true,
    note TEXT,
    organization_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL
);

-- Invoice details table
CREATE TABLE IF NOT EXISTS invoice_details (
    id SERIAL PRIMARY KEY,
    invoice_uuid UUID NOT NULL REFERENCES invoices(document_uuid) ON DELETE CASCADE,
    catalog_id VARCHAR(255),
    goods_name VARCHAR(500) NOT NULL,
    base_count DECIMAL(15,3) NOT NULL,
    price DECIMAL(15,2) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    amount_without_vat DECIMAL(15,2),
    amount_vat DECIMAL(15,2),
    amount_st DECIMAL(15,2),
    unit_type_id VARCHAR(50),
    tnved_code VARCHAR(50),
    gked_code VARCHAR(50),
    fcd_number VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Financial data table
CREATE TABLE IF NOT EXISTS invoice_financial_data (
    id SERIAL PRIMARY KEY,
    invoice_uuid UUID UNIQUE NOT NULL REFERENCES invoices(document_uuid) ON DELETE CASCADE,
    total_amount DECIMAL(15,2) NOT NULL,
    opening_balances DECIMAL(15,2),
    assessed_contributions_amount DECIMAL(15,2),
    paid_amount DECIMAL(15,2),
    penalties_amount DECIMAL(15,2),
    fines_amount DECIMAL(15,2),
    closing_balances DECIMAL(15,2),
    amount_to_be_paid DECIMAL(15,2),
    personal_account_number VARCHAR(100),
    legal_person_bank_account VARCHAR(100),
    contractor_bank_account VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_invoices_organization ON invoices(organization_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_date ON invoices(invoice_date);
CREATE INDEX IF NOT EXISTS idx_invoices_created_by ON invoices(created_by);
CREATE INDEX IF NOT EXISTS idx_invoice_details_invoice ON invoice_details(invoice_uuid);
