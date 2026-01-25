-- Create organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    owner_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on owner_id for faster queries
CREATE INDEX idx_organizations_owner_id ON organizations(owner_id);

-- Create employees table
CREATE TABLE IF NOT EXISTS employees (
    id VARCHAR(36) PRIMARY KEY,
    organization_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    role VARCHAR(50) NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    UNIQUE(organization_id, user_id)
);

-- Create indexes for employees
CREATE INDEX idx_employees_organization_id ON employees(organization_id);
CREATE INDEX idx_employees_user_id ON employees(user_id);

-- Add comments
COMMENT ON TABLE organizations IS 'Organizations in the system';
COMMENT ON TABLE employees IS 'Organization members';
COMMENT ON COLUMN organizations.id IS 'Unique organization identifier';
COMMENT ON COLUMN organizations.name IS 'Organization name (unique)';
COMMENT ON COLUMN organizations.owner_id IS 'User ID of the organization owner';
COMMENT ON COLUMN employees.id IS 'Unique employee record identifier';
COMMENT ON COLUMN employees.role IS 'Employee role: owner, admin, member, viewer';
