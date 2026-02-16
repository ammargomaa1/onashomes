ALTER TABLE product_variants
    DROP COLUMN IF EXISTS length,
    DROP COLUMN IF EXISTS width,
    DROP COLUMN IF EXISTS height;
