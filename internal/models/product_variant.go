package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID             int64          `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	ProductID      int64          `gorm:"type:bigint;not null" json:"product_id"`
	SKU            string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	AttributeValue string         `gorm:"type:varchar(255)" json:"attribute_value"`
	Price          *float64       `gorm:"type:numeric(12,2)" json:"price"`
	CompareAtPrice *float64       `gorm:"type:numeric(12,2)" json:"compare_at_price"`
	CostPrice      *float64       `gorm:"type:numeric(12,2)" json:"cost_price"`
	Barcode        *string        `gorm:"type:varchar(100)" json:"barcode"`
	Weight         *float64       `gorm:"type:numeric(10,3)" json:"weight"`
	Length         *float64       `gorm:"type:numeric(10,3)" json:"length"`
	Width          *float64       `gorm:"type:numeric(10,3)" json:"width"`
	Height         *float64       `gorm:"type:numeric(10,3)" json:"height"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Product   *Product           `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Inventory []VariantInventory `gorm:"foreignKey:ProductVariantID" json:"inventory,omitempty"`
}

func (ProductVariant) TableName() string { return "product_variants" }
