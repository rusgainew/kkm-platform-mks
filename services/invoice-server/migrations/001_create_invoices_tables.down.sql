-- Drop indexes for invoice_financials
DROP INDEX IF EXISTS idx_invoice_financials_invoice_uuid;

-- Drop indexes for invoice_details
DROP INDEX IF EXISTS idx_invoice_details_goods_name;
DROP INDEX IF EXISTS idx_invoice_details_catalog_id;
DROP INDEX IF EXISTS idx_invoice_details_invoice_uuid;

-- Drop indexes for invoices
DROP INDEX IF EXISTS idx_invoices_invoice_date;
DROP INDEX IF EXISTS idx_invoices_created_date;
DROP INDEX IF EXISTS idx_invoices_contractor_id;
DROP INDEX IF EXISTS idx_invoices_legal_person_id;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_invoice_number;
DROP INDEX IF EXISTS idx_invoices_document_uuid;

-- Drop tables in reverse order
DROP TABLE IF EXISTS invoice_financials;
DROP TABLE IF EXISTS invoice_details;
DROP TABLE IF EXISTS invoices;
