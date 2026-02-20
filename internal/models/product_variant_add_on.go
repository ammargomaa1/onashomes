package models

import "time"

// ProductVariantAddOn represents a mapping between a product variant and an add-on product.
type ProductVariantAddOn struct {
	ID               int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	ProductVariantID int64     `gorm:"type:bigint;not null" json:"product_variant_id"`
	AddOnProductID   int64     `gorm:"type:bigint;not null" json:"add_on_product_id"`
	CreatedAt        time.Time `json:"created_at"`
}
