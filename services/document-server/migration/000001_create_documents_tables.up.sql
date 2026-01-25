-- Create documents table
CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL,
    assigned_to UUID,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    status_changed_at BIGINT NOT NULL,
    version INT DEFAULT 1,
    created_at_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_documents_organization_id ON documents(organization_id);
CREATE INDEX IF NOT EXISTS idx_documents_status ON documents(status);
CREATE INDEX IF NOT EXISTS idx_documents_created_by ON documents(created_by);
CREATE INDEX IF NOT EXISTS idx_documents_assigned_to ON documents(assigned_to);

-- Create document_entries table for key-value entries
CREATE TABLE IF NOT EXISTS document_entries (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    created_at_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_document_entries_document_id ON document_entries(document_id);

-- Create document_workflows table
CREATE TABLE IF NOT EXISTS document_workflows (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL UNIQUE REFERENCES documents(id) ON DELETE CASCADE,
    current_step VARCHAR(100),
    rejection_reason TEXT,
    rejected_at BIGINT,
    approved_at BIGINT
);

-- Create approval_requests table
CREATE TABLE IF NOT EXISTS approval_requests (
    id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL REFERENCES document_workflows(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    comments TEXT,
    created_at BIGINT NOT NULL,
    responded_at BIGINT
);

CREATE INDEX IF NOT EXISTS idx_approval_requests_workflow_id ON approval_requests(workflow_id);
CREATE INDEX IF NOT EXISTS idx_approval_requests_approver_id ON approval_requests(approver_id);

-- Create document_versions table for history
CREATE TABLE IF NOT EXISTS document_versions (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version INT NOT NULL,
    title VARCHAR(255),
    content TEXT,
    status VARCHAR(50),
    created_at BIGINT NOT NULL,
    created_by UUID
);

CREATE INDEX IF NOT EXISTS idx_document_versions_document_id ON document_versions(document_id);
CREATE INDEX IF NOT EXISTS idx_document_versions_version ON document_versions(version);
