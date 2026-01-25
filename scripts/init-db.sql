-- Initialize KKM Database Schemas
-- This script creates the necessary schemas and extensions for the microservices

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create schemas for each microservice domain
CREATE SCHEMA IF NOT EXISTS users;
CREATE SCHEMA IF NOT EXISTS companies;
CREATE SCHEMA IF NOT EXISTS catalogs;
CREATE SCHEMA IF NOT EXISTS invoices;
CREATE SCHEMA IF NOT EXISTS bank_accounts;
CREATE SCHEMA IF NOT EXISTS foreign_companies;
CREATE SCHEMA IF NOT EXISTS documents;

-- Grant permissions to the application user
GRANT ALL PRIVILEGES ON SCHEMA users TO kkm_user;
GRANT ALL PRIVILEGES ON SCHEMA companies TO kkm_user;
GRANT ALL PRIVILEGES ON SCHEMA catalogs TO kkm_user;
GRANT ALL PRIVILEGES ON SCHEMA invoices TO kkm_user;
GRANT ALL PRIVILEGES ON SCHEMA bank_accounts TO kkm_user;
GRANT ALL PRIVILEGES ON SCHEMA foreign_companies TO kkm_user;
GRANT ALL PRIVILEGES ON SCHEMA documents TO kkm_user;

-- Grant usage on all existing tables (will be created by migrations)
ALTER DEFAULT PRIVILEGES IN SCHEMA users GRANT ALL ON TABLES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA companies GRANT ALL ON TABLES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA catalogs GRANT ALL ON TABLES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA invoices GRANT ALL ON TABLES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA bank_accounts GRANT ALL ON TABLES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA foreign_companies GRANT ALL ON TABLES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA documents GRANT ALL ON TABLES TO kkm_user;

-- Grant usage on sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA users GRANT ALL ON SEQUENCES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA companies GRANT ALL ON SEQUENCES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA catalogs GRANT ALL ON SEQUENCES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA invoices GRANT ALL ON SEQUENCES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA bank_accounts GRANT ALL ON SEQUENCES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA foreign_companies GRANT ALL ON SEQUENCES TO kkm_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA documents GRANT ALL ON SEQUENCES TO kkm_user;

-- Create a search path for easier access
ALTER DATABASE kkm_db SET search_path TO users, companies, catalogs, invoices, bank_accounts, foreign_companies, documents, public;

-- Log initialization
DO $$
BEGIN
    RAISE NOTICE 'KKM Database initialized successfully';
    RAISE NOTICE 'Schemas created: users, companies, catalogs, invoices, bank_accounts, foreign_companies, documents';
    RAISE NOTICE 'Extensions enabled: uuid-ossp, pg_trgm';
END $$;
