-- Add optimistic locking support
-- This migration ensures proper version handling for concurrent updates

-- Add NOT NULL constraint to version if not already present
ALTER TABLE IF EXISTS documents 
ADD CONSTRAINT check_version_positive CHECK (version > 0);

-- Create index on version for faster lookups
CREATE INDEX IF NOT EXISTS idx_documents_version ON documents(version);

-- Add comment explaining optimistic locking strategy
COMMENT ON COLUMN documents.version IS 'Version counter for optimistic locking. Incremented on each update. Used to detect concurrent modifications.';

-- Create audit table to track version changes
CREATE TABLE IF NOT EXISTS document_version_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    old_version INT NOT NULL,
    new_version INT NOT NULL,
    changed_by UUID,
    changed_at BIGINT NOT NULL,
    operation VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_document_version_history_document_id 
ON document_version_history(document_id);

CREATE INDEX IF NOT EXISTS idx_document_version_history_changed_at 
ON document_version_history(changed_at);
