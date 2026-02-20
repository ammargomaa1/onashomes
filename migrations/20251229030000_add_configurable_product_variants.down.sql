-- Migration: add_configurable_product_variants (rollback)
-- Created at: Mon Dec 29 02:57:00 AM EET 2025

-- Drop indexes and tables in reverse dependency order

-- Variant-based add-on compatibility
DROP INDEX IF EXISTS ux_product_variant_add_ons_variant_id_add_on_product_id;
DROP INDEX IF EXISTS idx_product_variant_add_ons_add_on_product_id;
DROP TABLE IF EXISTS product_variant_add_ons;

-- Variant images
DROP INDEX IF EXISTS ux_product_variant_images_variant_id_position;
DROP INDEX IF EXISTS ux_product_variant_images_variant_id_file_id;
DROP INDEX IF EXISTS idx_product_variant_images_file_id;
DROP TABLE IF EXISTS product_variant_images;

-- Variant attribute values
DROP INDEX IF EXISTS ux_product_variant_attribute_values_variant_id_attribute_value_id;
DROP INDEX IF EXISTS idx_product_variant_attribute_values_attribute_value_id;
DROP TABLE IF EXISTS product_variant_attribute_values;

-- Product variants
DROP INDEX IF EXISTS ux_product_variants_sku;
DROP INDEX IF EXISTS idx_product_variants_product_id;
DROP INDEX IF EXISTS idx_product_variants_is_active;
DROP INDEX IF EXISTS idx_product_variants_deleted_at;
DROP TABLE IF EXISTS product_variants;

-- Product attribute values
DROP INDEX IF EXISTS ux_product_attribute_values_product_id_attribute_value_id;
DROP INDEX IF EXISTS idx_product_attribute_values_attribute_value_id;
DROP TABLE IF EXISTS product_attribute_values;

-- Product attributes
DROP INDEX IF EXISTS ux_product_attributes_product_id_attribute_id;
DROP INDEX IF EXISTS idx_product_attributes_product_id;
DROP INDEX IF EXISTS idx_product_attributes_attribute_id;
DROP TABLE IF EXISTS product_attributes;

-- Global attribute values
DROP INDEX IF EXISTS ux_attribute_values_attribute_id_value;
DROP INDEX IF EXISTS idx_attribute_values_attribute_id;
DROP INDEX IF EXISTS idx_attribute_values_is_active;
DROP INDEX IF EXISTS idx_attribute_values_deleted_at;
DROP TABLE IF EXISTS attribute_values;

-- Global attributes
DROP INDEX IF EXISTS ux_attributes_name;
DROP INDEX IF EXISTS idx_attributes_is_active;
DROP INDEX IF EXISTS idx_attributes_deleted_at;
DROP TABLE IF EXISTS attributes;

-- Products
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_deleted_at;
DROP TABLE IF EXISTS products;
