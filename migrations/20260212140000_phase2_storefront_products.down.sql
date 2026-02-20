-- Down migration: phase2_storefront_products

DROP TABLE IF EXISTS inventory_adjustments;
DROP TABLE IF EXISTS variant_inventory;
DROP TABLE IF EXISTS product_seo;
DROP TABLE IF EXISTS product_storefront;
DROP TABLE IF EXISTS category_storefront;
DROP TABLE IF EXISTS brand_storefront;
DROP TABLE IF EXISTS supplier_storefront;
DROP TABLE IF EXISTS section_storefront;
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
