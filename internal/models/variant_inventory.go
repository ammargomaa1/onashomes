package models

import "time"

type VariantInventory struct {
	ID                int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	ProductVariantID  int64     `gorm:"type:bigint;not null" json:"product_variant_id"`
	StoreFrontID      int64     `gorm:"type:bigint;not null" json:"store_front_id"`
	Quantity          int       `gorm:"not null;default:0" json:"quantity"`
	ReservedQuantity  int       `gorm:"not null;default:0" json:"reserved_quantity"`
	LowStockThreshold int       `gorm:"not null;default:5" json:"low_stock_threshold"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relations
	StoreFront     StoreFront     `gorm:"foreignKey:StoreFrontID" json:"store_front,omitempty"`
	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
}

func (VariantInventory) TableName() string { return "variant_inventory" }

// AvailableQuantity returns sellable quantity
func (v *VariantInventory) AvailableQuantity() int {
	return v.Quantity - v.ReservedQuantity
}

// IsLowStock returns true if available stock is at or below threshold
func (v *VariantInventory) IsLowStock() bool {
	return v.AvailableQuantity() <= v.LowStockThreshold
}
