-- Drop document versions table
DROP TABLE IF EXISTS document_versions;

-- Drop approval requests table
DROP TABLE IF EXISTS approval_requests;

-- Drop workflows table
DROP TABLE IF EXISTS document_workflows;

-- Drop document entries table
DROP TABLE IF EXISTS document_entries;

-- Drop indexes
DROP INDEX IF EXISTS idx_documents_organization_id;
DROP INDEX IF EXISTS idx_documents_status;
DROP INDEX IF EXISTS idx_documents_created_by;
DROP INDEX IF EXISTS idx_documents_assigned_to;

-- Drop documents table
DROP TABLE IF EXISTS documents;
