-- Migration: phase2_storefront_products
-- Created at: 2026-02-12

-- ============================================================
-- 1. STORE FRONTS
-- ============================================================
CREATE TABLE store_fronts (
    id               BIGSERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    slug             VARCHAR(255) NOT NULL UNIQUE,
    domain           VARCHAR(255) NOT NULL UNIQUE,
    currency         VARCHAR(10)  NOT NULL DEFAULT 'SAR',
    default_language VARCHAR(10)  NOT NULL DEFAULT 'ar',
    is_active        BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_store_fronts_domain ON store_fronts (domain);
CREATE INDEX idx_store_fronts_slug   ON store_fronts (slug);
CREATE INDEX idx_store_fronts_active ON store_fronts (is_active);

-- ============================================================
-- 2. PIVOT TABLES (Store <-> Catalog)
-- ============================================================
CREATE TABLE product_storefront (
    product_id     BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    store_front_id BIGINT NOT NULL REFERENCES store_fronts(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, store_front_id)
);
CREATE INDEX idx_product_storefront_store ON product_storefront (store_front_id);

CREATE TABLE category_storefront (
    category_id    BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    store_front_id BIGINT NOT NULL REFERENCES store_fronts(id) ON DELETE CASCADE,
    PRIMARY KEY (category_id, store_front_id)
);
CREATE INDEX idx_category_storefront_store ON category_storefront (store_front_id);

CREATE TABLE brand_storefront (
    brand_id       BIGINT NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    store_front_id BIGINT NOT NULL REFERENCES store_fronts(id) ON DELETE CASCADE,
    PRIMARY KEY (brand_id, store_front_id)
);
CREATE INDEX idx_brand_storefront_store ON brand_storefront (store_front_id);

CREATE TABLE supplier_storefront (
    supplier_id    BIGINT NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    store_front_id BIGINT NOT NULL REFERENCES store_fronts(id) ON DELETE CASCADE,
    PRIMARY KEY (supplier_id, store_front_id)
);
CREATE INDEX idx_supplier_storefront_store ON supplier_storefront (store_front_id);

CREATE TABLE section_storefront (
    section_id     BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    store_front_id BIGINT NOT NULL REFERENCES store_fronts(id) ON DELETE CASCADE,
    PRIMARY KEY (section_id, store_front_id)
);
CREATE INDEX idx_section_storefront_store ON section_storefront (store_front_id);

-- ============================================================
-- 3. ALTER PRODUCTS TABLE
-- ============================================================
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS name_en          VARCHAR(255),
    ADD COLUMN IF NOT EXISTS name_ar          VARCHAR(255),
    ADD COLUMN IF NOT EXISTS description_en   TEXT,
    ADD COLUMN IF NOT EXISTS description_ar   TEXT,
    ADD COLUMN IF NOT EXISTS slug             VARCHAR(255),
    ADD COLUMN IF NOT EXISTS brand_id         BIGINT REFERENCES brands(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS category_id      BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS supplier_id      BIGINT REFERENCES suppliers(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS is_internal_supplier BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_published     BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS status           VARCHAR(20) NOT NULL DEFAULT 'draft',
    ADD COLUMN IF NOT EXISTS is_featured      BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_new           BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_best_seller   BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS attribute_type   VARCHAR(20);

-- Backfill name_en from existing name column
UPDATE products SET name_en = name WHERE name_en IS NULL AND name IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_products_slug        ON products (slug);
CREATE INDEX IF NOT EXISTS idx_products_brand       ON products (brand_id);
CREATE INDEX IF NOT EXISTS idx_products_category    ON products (category_id);
CREATE INDEX IF NOT EXISTS idx_products_supplier    ON products (supplier_id);
CREATE INDEX IF NOT EXISTS idx_products_status      ON products (status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_published   ON products (is_published) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_featured    ON products (is_featured) WHERE deleted_at IS NULL AND is_published = TRUE;
CREATE INDEX IF NOT EXISTS idx_products_attr_type   ON products (attribute_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_admin_filter ON products (status, brand_id, category_id, supplier_id) WHERE deleted_at IS NULL;

-- ============================================================
-- 4. ALTER PRODUCT VARIANTS
-- ============================================================
ALTER TABLE product_variants
    ADD COLUMN IF NOT EXISTS price              NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS compare_at_price   NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS cost_price         NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS barcode            VARCHAR(100),
    ADD COLUMN IF NOT EXISTS weight             NUMERIC(10,3),
    ADD COLUMN IF NOT EXISTS attribute_value    VARCHAR(255);

CREATE UNIQUE INDEX IF NOT EXISTS ux_product_variants_barcode
    ON product_variants (barcode) WHERE barcode IS NOT NULL AND deleted_at IS NULL;

-- ============================================================
-- 5. VARIANT INVENTORY (Per Store)
-- ============================================================
CREATE TABLE variant_inventory (
    id                  BIGSERIAL PRIMARY KEY,
    product_variant_id  BIGINT NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    store_front_id      BIGINT NOT NULL REFERENCES store_fronts(id) ON DELETE CASCADE,
    quantity            INT    NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved_quantity   INT    NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    low_stock_threshold INT    NOT NULL DEFAULT 5,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ux_variant_inventory UNIQUE (product_variant_id, store_front_id)
);

CREATE INDEX idx_variant_inventory_store   ON variant_inventory (store_front_id);
CREATE INDEX idx_variant_inventory_variant ON variant_inventory (product_variant_id);

-- ============================================================
-- 6. INVENTORY ADJUSTMENTS
-- ============================================================
CREATE TABLE inventory_adjustments (
    id                   BIGSERIAL PRIMARY KEY,
    variant_inventory_id BIGINT      NOT NULL REFERENCES variant_inventory(id) ON DELETE CASCADE,
    adjusted_by          BIGINT      NOT NULL,
    previous_quantity    INT         NOT NULL,
    new_quantity         INT         NOT NULL,
    adjustment_amount    INT         NOT NULL,
    reason               VARCHAR(50) NOT NULL,
    notes                TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inventory_adj_inventory ON inventory_adjustments (variant_inventory_id);
CREATE INDEX idx_inventory_adj_date      ON inventory_adjustments (created_at);

-- ============================================================
-- 7. PRODUCT SEO
-- ============================================================
CREATE TABLE product_seo (
    id                  BIGSERIAL PRIMARY KEY,
    product_id          BIGINT NOT NULL UNIQUE REFERENCES products(id) ON DELETE CASCADE,
    meta_title_en       VARCHAR(70),
    meta_title_ar       VARCHAR(70),
    meta_description_en VARCHAR(170),
    meta_description_ar VARCHAR(170),
    meta_keywords       TEXT,
    canonical_url       VARCHAR(500),
    og_title            VARCHAR(200),
    og_description      VARCHAR(300),
    og_image            VARCHAR(500),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_seo_product ON product_seo (product_id);
