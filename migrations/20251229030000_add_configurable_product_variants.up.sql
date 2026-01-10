-- Migration: add_configurable_product_variants
-- Created at: Mon Dec 29 02:57:00 AM EET 2025

-- Core products table
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_products_name ON products (name);
CREATE INDEX idx_products_is_active ON products (is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_deleted_at ON products (deleted_at);

-- Global attributes
CREATE TABLE attributes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ux_attributes_name ON attributes (name);
CREATE INDEX idx_attributes_is_active ON attributes (is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_attributes_deleted_at ON attributes (deleted_at);

-- Global attribute values
CREATE TABLE attribute_values (
    id BIGSERIAL PRIMARY KEY,
    attribute_id BIGINT NOT NULL,
    value VARCHAR(100) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_attribute_values_attribute_id
        FOREIGN KEY (attribute_id)
        REFERENCES attributes(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX ux_attribute_values_attribute_id_value
    ON attribute_values(attribute_id, value);
CREATE INDEX idx_attribute_values_attribute_id
    ON attribute_values(attribute_id);
CREATE INDEX idx_attribute_values_is_active
    ON attribute_values(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_attribute_values_deleted_at
    ON attribute_values(deleted_at);

-- Attributes assigned to a product
CREATE TABLE product_attributes (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    attribute_id BIGINT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_attributes_product_id
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_product_attributes_attribute_id
        FOREIGN KEY (attribute_id)
        REFERENCES attributes(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX ux_product_attributes_product_id_attribute_id
    ON product_attributes(product_id, attribute_id);
CREATE INDEX idx_product_attributes_product_id
    ON product_attributes(product_id);
CREATE INDEX idx_product_attributes_attribute_id
    ON product_attributes(attribute_id);

-- Allowed attribute values per product
CREATE TABLE product_attribute_values (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    attribute_value_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_attribute_values_product_id
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_product_attribute_values_attribute_value_id
        FOREIGN KEY (attribute_value_id)
        REFERENCES attribute_values(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX ux_product_attribute_values_product_id_attribute_value_id
    ON product_attribute_values(product_id, attribute_value_id);
CREATE INDEX idx_product_attribute_values_attribute_value_id
    ON product_attribute_values(attribute_value_id);

-- Product variants (SKU-level)
CREATE TABLE product_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    sku VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_product_variants_product_id
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX ux_product_variants_sku
    ON product_variants(sku);
CREATE INDEX idx_product_variants_product_id
    ON product_variants(product_id);
CREATE INDEX idx_product_variants_is_active
    ON product_variants(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_variants_deleted_at
    ON product_variants(deleted_at);

-- Attribute values that define each variant's combination
CREATE TABLE product_variant_attribute_values (
    id BIGSERIAL PRIMARY KEY,
    product_variant_id BIGINT NOT NULL,
    attribute_value_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_variant_attribute_values_variant_id
        FOREIGN KEY (product_variant_id)
        REFERENCES product_variants(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_product_variant_attribute_values_attribute_value_id
        FOREIGN KEY (attribute_value_id)
        REFERENCES attribute_values(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX ux_product_variant_attribute_values_variant_id_attribute_value_id
    ON product_variant_attribute_values(product_variant_id, attribute_value_id);
CREATE INDEX idx_product_variant_attribute_values_attribute_value_id
    ON product_variant_attribute_values(attribute_value_id);

-- Variant images (ordered list of file_ids per variant)
CREATE TABLE product_variant_images (
    id BIGSERIAL PRIMARY KEY,
    product_variant_id BIGINT NOT NULL,
    file_id BIGINT NOT NULL,
    position INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_variant_images_variant_id
        FOREIGN KEY (product_variant_id)
        REFERENCES product_variants(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_product_variant_images_file_id
        FOREIGN KEY (file_id)
        REFERENCES files(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX ux_product_variant_images_variant_id_position
    ON product_variant_images(product_variant_id, position);
CREATE UNIQUE INDEX ux_product_variant_images_variant_id_file_id
    ON product_variant_images(product_variant_id, file_id);
CREATE INDEX idx_product_variant_images_file_id
    ON product_variant_images(file_id);

-- Variant-based add-on compatibility (add-ons are products)
CREATE TABLE product_variant_add_ons (
    id BIGSERIAL PRIMARY KEY,
    product_variant_id BIGINT NOT NULL,
    add_on_product_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_variant_add_ons_variant_id
        FOREIGN KEY (product_variant_id)
        REFERENCES product_variants(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_product_variant_add_ons_add_on_product_id
        FOREIGN KEY (add_on_product_id)
        REFERENCES products(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX ux_product_variant_add_ons_variant_id_add_on_product_id
    ON product_variant_add_ons(product_variant_id, add_on_product_id);
CREATE INDEX idx_product_variant_add_ons_add_on_product_id
    ON product_variant_add_ons(add_on_product_id);
