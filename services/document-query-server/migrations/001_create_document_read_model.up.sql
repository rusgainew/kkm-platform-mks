-- Document Query Server Read Model
-- Оптимизирована для быстрого чтения и фильтрации документов

CREATE TABLE IF NOT EXISTS document_read_model (
    id UUID PRIMARY KEY,
    document_number VARCHAR(100) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    document_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    company_id UUID NOT NULL,
    company_name VARCHAR(255),
    created_by_user_id UUID,
    created_by_user_name VARCHAR(255),
    approval_status VARCHAR(50),
    approved_by VARCHAR(255),
    approved_by_user_id UUID,
    approved_by_user_name VARCHAR(255),
    approved_at TIMESTAMP,
    sent_at TIMESTAMP,
    rejected_by VARCHAR(255),
    rejected_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    archived_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Индексы для оптимизации чтения
CREATE INDEX IF NOT EXISTS idx_document_number ON document_read_model(document_number);
CREATE INDEX IF NOT EXISTS idx_document_status ON document_read_model(status);
CREATE INDEX IF NOT EXISTS idx_document_type ON document_read_model(document_type);
CREATE INDEX IF NOT EXISTS idx_document_company_id ON document_read_model(company_id);
CREATE INDEX IF NOT EXISTS idx_document_created_at ON document_read_model(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_document_approval_status ON document_read_model(approval_status);
CREATE INDEX IF NOT EXISTS idx_document_archived_at ON document_read_model(archived_at);
CREATE INDEX IF NOT EXISTS idx_document_deleted_at ON document_read_model(deleted_at);

-- Композитный индекс для типичных фильтров
CREATE INDEX IF NOT EXISTS idx_document_status_created ON document_read_model(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_document_company_status ON document_read_model(company_id, status);

-- View для активных документов
CREATE OR REPLACE VIEW active_documents AS
SELECT * FROM document_read_model
WHERE deleted_at IS NULL AND archived_at IS NULL
ORDER BY created_at DESC;

-- View для требующих одобрения
CREATE OR REPLACE VIEW pending_approval_documents AS
SELECT * FROM document_read_model
WHERE approval_status = 'pending'
  AND deleted_at IS NULL
  AND archived_at IS NULL
ORDER BY created_at ASC;

-- Комментарии
COMMENT ON TABLE document_read_model IS 'Read-side model для document-query-server (CQRS)';
COMMENT ON VIEW active_documents IS 'Активные, неудаленные и неархивированные документы';
COMMENT ON VIEW pending_approval_documents IS 'Документы ожидающие одобрения, отсортированные по дате создания';
