-- Migration: create_suppliers_table
-- Created at: Thu Dec 18 01:53:11 AM EET 2025

-- Add your UP migration SQL here
CREATE TABLE suppliers (
    id BIGSERIAL PRIMARY KEY,

    company_name VARCHAR(255) NOT NULL,
    contact_person_name VARCHAR(255) NOT NULL,
    contact_number VARCHAR(50) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_by BIGINT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_suppliers_created_by
        FOREIGN KEY (created_by)
        REFERENCES admins(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

CREATE INDEX idx_suppliers_deleted_at
    ON suppliers (deleted_at);