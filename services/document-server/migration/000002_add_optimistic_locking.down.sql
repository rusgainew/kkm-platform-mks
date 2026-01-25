-- Rollback optimistic locking migration
-- This migration reverts the changes made in 000002_add_optimistic_locking.up.sql

DROP INDEX IF EXISTS idx_document_version_history_changed_at;
DROP INDEX IF EXISTS idx_document_version_history_document_id;
DROP TABLE IF EXISTS document_version_history;

DROP INDEX IF EXISTS idx_documents_version;

ALTER TABLE IF EXISTS documents 
DROP CONSTRAINT IF EXISTS check_version_positive;
