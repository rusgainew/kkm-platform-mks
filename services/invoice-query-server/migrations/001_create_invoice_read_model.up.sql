-- Invoice Query Server Read Model
-- Simplified tables for fast queries with pagination

CREATE TABLE IF NOT EXISTS invoices (
    document_uuid VARCHAR(36) PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    invoice_date VARCHAR(10),
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    note TEXT,
    created_date VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS invoice_details (
    id BIGSERIAL PRIMARY KEY,
    invoice_uuid VARCHAR(36) NOT NULL REFERENCES invoices(document_uuid),
    base_count DECIMAL(10, 3) NOT NULL DEFAULT 0,
    price DECIMAL(15, 2) NOT NULL DEFAULT 0,
    amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    amount_without_vat DECIMAL(15, 2) NOT NULL DEFAULT 0,
    amount_vat DECIMAL(15, 2) NOT NULL DEFAULT 0,
    amount_st DECIMAL(15, 2) NOT NULL DEFAULT 0,
    goods_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for fast queries
CREATE INDEX IF NOT EXISTS idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX IF NOT EXISTS idx_invoices_created_date ON invoices(created_date DESC);
CREATE INDEX IF NOT EXISTS idx_invoice_details_invoice_uuid ON invoice_details(invoice_uuid);
CREATE INDEX IF NOT EXISTS idx_invoice_details_id ON invoice_details(id DESC);
