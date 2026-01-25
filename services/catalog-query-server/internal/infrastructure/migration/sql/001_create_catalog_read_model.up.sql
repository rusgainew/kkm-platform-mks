-- Catalog Query Server Read Model
-- Simplified tables for fast queries

CREATE TABLE IF NOT EXISTS catalogs (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    number VARCHAR(50),
    tnved_code VARCHAR(10),
    gked_code VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for fast queries
CREATE INDEX IF NOT EXISTS idx_catalogs_name ON catalogs(name);
CREATE INDEX IF NOT EXISTS idx_catalogs_number ON catalogs(number);
