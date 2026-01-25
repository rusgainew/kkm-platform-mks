-- Create invoices table
CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(36) PRIMARY KEY,
    document_uuid VARCHAR(36) NOT NULL UNIQUE,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    number VARCHAR(50) NOT NULL,
    corrected_receipt_uuid VARCHAR(36),
    invoice_date TIMESTAMP NOT NULL,
    created_date TIMESTAMP NOT NULL,
    delivery_date TIMESTAMP,
    corrected_receipt_creation_date TIMESTAMP,
    total_amount NUMERIC(15, 2) NOT NULL CHECK (total_amount >= 0),
    is_resident BOOLEAN NOT NULL DEFAULT true,
    note TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'sent', 'signed', 'accepted', 'rejected', 'revoked')),
    legal_person_id VARCHAR(36) NOT NULL,
    contractor_id VARCHAR(36) NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    signed_at TIMESTAMP,
    signed_by VARCHAR(36),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create invoice_details table
CREATE TABLE IF NOT EXISTS invoice_details (
    id SERIAL PRIMARY KEY,
    invoice_uuid VARCHAR(36) NOT NULL REFERENCES invoices(document_uuid) ON DELETE CASCADE,
    base_count NUMERIC(15, 3) NOT NULL CHECK (base_count > 0),
    price NUMERIC(15, 2) NOT NULL CHECK (price >= 0),
    amount NUMERIC(15, 2) NOT NULL CHECK (amount >= 0),
    amount_without_vat NUMERIC(15, 2) NOT NULL CHECK (amount_without_vat >= 0),
    amount_vat NUMERIC(15, 2) NOT NULL CHECK (amount_vat >= 0),
    amount_st NUMERIC(15, 2) NOT NULL CHECK (amount_st >= 0),
    goods_name VARCHAR(500) NOT NULL,
    tnved_code VARCHAR(50),
    gked_code VARCHAR(50),
    fcd_number VARCHAR(50),
    catalog_id VARCHAR(36) NOT NULL,
    unit_type_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create invoice_financials table
CREATE TABLE IF NOT EXISTS invoice_financials (
    invoice_uuid VARCHAR(36) PRIMARY KEY REFERENCES invoices(document_uuid) ON DELETE CASCADE,
    total_amount NUMERIC(15, 2) NOT NULL,
    opening_balances NUMERIC(15, 2) NOT NULL DEFAULT 0,
    assessed_contributions_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    paid_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    penalties_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    fines_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    closing_balances NUMERIC(15, 2) NOT NULL DEFAULT 0,
    amount_to_be_paid NUMERIC(15, 2) NOT NULL DEFAULT 0,
    personal_account_number VARCHAR(50),
    legal_person_bank_account VARCHAR(100),
    contractor_bank_account VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for invoices
CREATE INDEX idx_invoices_document_uuid ON invoices(document_uuid);
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_legal_person_id ON invoices(legal_person_id);
CREATE INDEX idx_invoices_contractor_id ON invoices(contractor_id);
CREATE INDEX idx_invoices_created_date ON invoices(created_date DESC);
CREATE INDEX idx_invoices_invoice_date ON invoices(invoice_date DESC);

-- Create indexes for invoice_details
CREATE INDEX idx_invoice_details_invoice_uuid ON invoice_details(invoice_uuid);
CREATE INDEX idx_invoice_details_catalog_id ON invoice_details(catalog_id);
CREATE INDEX idx_invoice_details_goods_name ON invoice_details(goods_name);

-- Create indexes for invoice_financials
CREATE INDEX idx_invoice_financials_invoice_uuid ON invoice_financials(invoice_uuid);
