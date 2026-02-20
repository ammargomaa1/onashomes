# Phase 2 — Database Schema Design

## ER Diagram (Textual)

```
store_fronts (1) ─────────< product_storefront >─────────── (N) products
store_fronts (1) ─────────< category_storefront >────────── (N) categories
store_fronts (1) ─────────< brand_storefront >───────────── (N) brands
store_fronts (1) ─────────< supplier_storefront >────────── (N) suppliers

products (1) ─────────< product_variants >───────────────── (N) variants
products (1) ─────────< product_seo >────────────────────── (1) seo_metadata

product_variants (1) ──< variant_inventory >─────────────── (N) per-store inventory
product_variants (1) ──< product_variant_attribute_values > (N) attribute combos

suppliers (1) ─────────< products (supplier_id) ──────────  (N) products

brands (1) ────────────< products (brand_id) ─────────────  (N) products
categories (1) ────────< products (category_id) ──────────  (N) products

inventory_adjustments (N) ──> variant_inventory (1)
```

## Full SQL Migration

```sql
-- Migration: 20260212140000_phase2_storefront_products.up.sql

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

CREATE INDEX idx_store_fronts_domain   ON store_fronts (domain);
CREATE INDEX idx_store_fronts_slug     ON store_fronts (slug);
CREATE INDEX idx_store_fronts_active   ON store_fronts (is_active);

-- ============================================================
-- 2. PIVOT TABLES (Store ↔ Catalog Entities)
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

-- ============================================================
-- 3. ALTER PRODUCTS TABLE (New Fields)
-- ============================================================
ALTER TABLE products
    ADD COLUMN name_en          VARCHAR(255),
    ADD COLUMN name_ar          VARCHAR(255),
    ADD COLUMN description_en   TEXT,
    ADD COLUMN description_ar   TEXT,
    ADD COLUMN slug             VARCHAR(255),
    ADD COLUMN brand_id         BIGINT REFERENCES brands(id) ON DELETE SET NULL,
    ADD COLUMN category_id      BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    ADD COLUMN supplier_id      BIGINT REFERENCES suppliers(id) ON DELETE SET NULL,
    ADD COLUMN is_internal_supplier BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN is_published     BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN status           VARCHAR(20) NOT NULL DEFAULT 'draft',
    ADD COLUMN is_featured      BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN is_new           BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN is_best_seller   BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN attribute_type   VARCHAR(20); -- 'size', 'color', or NULL (simple)

-- Backfill name_en from existing name column
UPDATE products SET name_en = name WHERE name_en IS NULL;
ALTER TABLE products ALTER COLUMN name_en SET NOT NULL;

-- Slug uniqueness is per-store, enforced via unique index on pivot
-- Global slug index for fast lookups
CREATE INDEX idx_products_slug        ON products (slug);
CREATE INDEX idx_products_brand       ON products (brand_id);
CREATE INDEX idx_products_category    ON products (category_id);
CREATE INDEX idx_products_supplier    ON products (supplier_id);
CREATE INDEX idx_products_status      ON products (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_published   ON products (is_published) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_featured    ON products (is_featured) WHERE deleted_at IS NULL AND is_published = TRUE;
CREATE INDEX idx_products_attr_type   ON products (attribute_type) WHERE deleted_at IS NULL;

-- Composite index for admin filtering
CREATE INDEX idx_products_admin_filter ON products (status, brand_id, category_id, supplier_id)
    WHERE deleted_at IS NULL;

-- Per-store slug uniqueness
CREATE UNIQUE INDEX ux_product_storefront_slug
    ON product_storefront (store_front_id, (
        SELECT slug FROM products WHERE products.id = product_storefront.product_id
    ));
-- NOTE: Slug uniqueness per store is enforced at APPLICATION LEVEL
-- because PostgreSQL doesn't support cross-table unique indexes elegantly.
-- We use a service-layer check before insert/update.

-- ============================================================
-- 4. ALTER PRODUCT VARIANTS (New Pricing/Inventory Fields)
-- ============================================================
ALTER TABLE product_variants
    ADD COLUMN price              NUMERIC(12,2),
    ADD COLUMN compare_at_price   NUMERIC(12,2),
    ADD COLUMN cost_price         NUMERIC(12,2),
    ADD COLUMN barcode            VARCHAR(100),
    ADD COLUMN weight             NUMERIC(10,3),
    ADD COLUMN attribute_value    VARCHAR(255); -- e.g. "XL", "Red"

CREATE UNIQUE INDEX ux_product_variants_barcode
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
CREATE INDEX idx_variant_inventory_low     ON variant_inventory (quantity, low_stock_threshold)
    WHERE quantity <= low_stock_threshold;

-- ============================================================
-- 6. INVENTORY ADJUSTMENTS (Audit Trail)
-- ============================================================
CREATE TABLE inventory_adjustments (
    id                   BIGSERIAL PRIMARY KEY,
    variant_inventory_id BIGINT      NOT NULL REFERENCES variant_inventory(id) ON DELETE CASCADE,
    adjusted_by          BIGINT      NOT NULL, -- admin_id
    previous_quantity    INT         NOT NULL,
    new_quantity         INT         NOT NULL,
    adjustment_amount    INT         NOT NULL, -- can be negative
    reason               VARCHAR(50) NOT NULL, -- 'adjustment', 'correction', 'restock', 'sale', 'return'
    notes                TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inventory_adj_inventory ON inventory_adjustments (variant_inventory_id);
CREATE INDEX idx_inventory_adj_date      ON inventory_adjustments (created_at);

-- ============================================================
-- 7. PRODUCT SEO (Separate table, 1:1 with product)
-- ============================================================
CREATE TABLE product_seo (
    id                  BIGSERIAL PRIMARY KEY,
    product_id          BIGINT NOT NULL UNIQUE REFERENCES products(id) ON DELETE CASCADE,
    meta_title_en       VARCHAR(70),
    meta_title_ar       VARCHAR(70),
    meta_description_en VARCHAR(170),
    meta_description_ar VARCHAR(170),
    meta_keywords       TEXT,       -- comma-separated
    canonical_url       VARCHAR(500),
    og_title            VARCHAR(200),
    og_description      VARCHAR(300),
    og_image            VARCHAR(500),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_seo_product ON product_seo (product_id);
```

## Why Separate Table for SEO vs JSONB

| Approach | Pros | Cons |
|----------|------|------|
| **Separate `product_seo` table** ✅ | Indexable, type-safe, clear schema, easy to validate at DB level | Extra JOIN |
| JSONB column on products | Flexible, no JOIN | Not indexable for individual fields, no type safety, harder validation |

**Decision: Separate table.** SEO fields are known, stable, and need validation. The 1:1 relationship with product makes JOINs trivial.

## Indexing Strategy Summary

| Index | Purpose |
|-------|---------|
| `idx_products_admin_filter` | Composite for admin list filtering (status + brand + category + supplier) |
| `idx_products_published` | Storefront queries—only published products |
| `idx_products_featured` | Homepage featured product queries |
| `idx_variant_inventory_low` | Low stock alerts dashboard |
| `idx_store_fronts_domain` | Domain-based store resolution (SSR) |
| `ux_product_variants_sku` | Global SKU uniqueness (already exists) |
| `ux_product_variants_barcode` | Barcode uniqueness |

## Down Migration

```sql
-- 20260212140000_phase2_storefront_products.down.sql

DROP TABLE IF EXISTS inventory_adjustments;
DROP TABLE IF EXISTS variant_inventory;
DROP TABLE IF EXISTS product_seo;
DROP TABLE IF EXISTS product_storefront;
DROP TABLE IF EXISTS category_storefront;
DROP TABLE IF EXISTS brand_storefront;
DROP TABLE IF EXISTS supplier_storefront;
DROP TABLE IF EXISTS store_fronts;

ALTER TABLE products
    DROP COLUMN IF EXISTS name_en,
    DROP COLUMN IF EXISTS name_ar,
    DROP COLUMN IF EXISTS description_en,
    DROP COLUMN IF EXISTS description_ar,
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS brand_id,
    DROP COLUMN IF EXISTS category_id,
    DROP COLUMN IF EXISTS supplier_id,
    DROP COLUMN IF EXISTS is_internal_supplier,
    DROP COLUMN IF EXISTS is_published,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS is_featured,
    DROP COLUMN IF EXISTS is_new,
    DROP COLUMN IF EXISTS is_best_seller,
    DROP COLUMN IF EXISTS attribute_type;

ALTER TABLE product_variants
    DROP COLUMN IF EXISTS price,
    DROP COLUMN IF EXISTS compare_at_price,
    DROP COLUMN IF EXISTS cost_price,
    DROP COLUMN IF EXISTS barcode,
    DROP COLUMN IF EXISTS weight,
    DROP COLUMN IF EXISTS attribute_value;
```
