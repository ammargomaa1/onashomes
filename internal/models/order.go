package models

import (
	"time"

	"gorm.io/gorm"
)

// Status Models
type OrderStatus struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	NameEn string `json:"name_en" gorm:"size:255;not null"`
	NameAr string `json:"name_ar" gorm:"size:255;not null"`
	Slug   string `json:"slug" gorm:"size:255;not null;uniqueIndex"`
}

type PaymentStatus struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	NameEn string `json:"name_en" gorm:"size:255;not null"`
	NameAr string `json:"name_ar" gorm:"size:255;not null"`
	Slug   string `json:"slug" gorm:"size:255;not null;uniqueIndex"`
}

type FulfillmentStatus struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	NameEn string `json:"name_en" gorm:"size:255;not null"`
	NameAr string `json:"name_ar" gorm:"size:255;not null"`
	Slug   string `json:"slug" gorm:"size:255;not null;uniqueIndex"`
}

type Currency struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	NameEn string `json:"name_en" gorm:"size:255;not null"`
	NameAr string `json:"name_ar" gorm:"size:255;not null"`
	Code   string `json:"code" gorm:"size:3;not null;uniqueIndex"`
	Symbol string `json:"symbol" gorm:"size:10"`
}

type Order struct {
	ID           int64  `json:"id" gorm:"primaryKey"`
	StoreFrontID int64  `json:"store_front_id" gorm:"index;not null"`
	OrderNumber  string `json:"order_number" gorm:"uniqueIndex;not null;size:64"`

	// Foreign Keys for Statuses
	OrderStatusID       int64  `json:"order_status_id" gorm:"index;not null;default:1"`
	PaymentStatusID     int64  `json:"payment_status_id" gorm:"index;not null;default:1"`
	PaymentMethodID     *int64 `json:"payment_method_id" gorm:"index"` // Nullable
	OrderSourceID       *int64 `json:"order_source_id" gorm:"index"`   // Nullable FK to OrderSources
	FulfillmentStatusID int64  `json:"fulfillment_status_id" gorm:"index;not null;default:1"`
	CurrencyID          int64  `json:"currency_id" gorm:"index;not null;default:1"`

	// Customer Info
	CustomerID    *int64 `json:"customer_id,omitempty" gorm:"index"` // Nullable FK to Customers
	CustomerName  string `json:"customer_name" gorm:"size:255"`
	CustomerEmail string `json:"customer_email" gorm:"size:255;index"`
	CustomerPhone string `json:"customer_phone" gorm:"size:50;index"`

	Subtotal       float64 `json:"subtotal" gorm:"not null;default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"not null;default:0"`
	TaxAmount      float64 `json:"tax_amount" gorm:"not null;default:0"`
	ShippingAmount float64 `json:"shipping_amount" gorm:"not null;default:0"`
	TotalAmount    float64 `json:"total_amount" gorm:"not null;default:0"`
	Notes          string  `json:"notes" gorm:"type:text"`

	CreatedByID int64          `json:"created_by_id" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Associations
	Items             []OrderItem        `json:"items" gorm:"foreignKey:OrderID"`
	StoreFront        *StoreFront        `json:"store_front,omitempty" gorm:"foreignKey:StoreFrontID"`
	OrderStatus       *OrderStatus       `json:"order_status,omitempty" gorm:"foreignKey:OrderStatusID"`
	PaymentStatus     *PaymentStatus     `json:"payment_status,omitempty" gorm:"foreignKey:PaymentStatusID"`
	PaymentMethod     *PaymentMethod     `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID"`
	OrderSource       *OrderSource       `json:"order_source,omitempty" gorm:"foreignKey:OrderSourceID"`
	FulfillmentStatus *FulfillmentStatus `json:"fulfillment_status,omitempty" gorm:"foreignKey:FulfillmentStatusID"`
	Currency          *Currency          `json:"currency,omitempty" gorm:"foreignKey:CurrencyID"`
	CreatedBy         *Admin             `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	Customer          *Customer          `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Address           *OrderAddress      `json:"address,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID                    int64   `json:"id" gorm:"primaryKey"`
	OrderID               int64   `json:"order_id" gorm:"index;not null"`
	ProductID             int64   `json:"product_id" gorm:"index;not null"`
	ProductVariantID      int64   `json:"product_variant_id" gorm:"index;not null"`
	SKU                   string  `json:"sku" gorm:"size:64;not null"`
	ProductNameSnapshotEn string  `json:"product_name_snapshot_en" gorm:"size:255"`
	ProductNameSnapshotAr string  `json:"product_name_snapshot_ar" gorm:"size:255"`
	UnitPrice             float64 `json:"unit_price" gorm:"not null"`
	CostPrice             float64 `json:"cost_price" gorm:"not null;default:0"`
	Quantity              int     `json:"quantity" gorm:"not null"`
	TotalPrice            float64 `json:"total_price" gorm:"not null"`

	// Associations
	Product        *Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ProductVariant *ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`
}
