-- Migration: create_brands_table
-- Created at: Fri 16 Jan 2026 03:39:46 PM EET

-- Add your UP migration SQL here
CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    name_ar VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    icon_id INT,
    logo_id INT,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE brands
ADD CONSTRAINT fk_brand_logo_id FOREIGN KEY (logo_id) REFERENCES files(id) ON DELETE SET NULL;

ALTER TABLE brands
ADD CONSTRAINT fk_brand_icon_id FOREIGN KEY (icon_id) REFERENCES files(id) ON DELETE SET NULL;