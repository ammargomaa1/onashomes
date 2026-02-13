package models

import (
	"time"

	"gorm.io/gorm"
)

// Product status constants
const (
	ProductStatusDraft    = "draft"
	ProductStatusActive   = "active"
	ProductStatusInactive = "inactive"
	ProductStatusArchived = "archived"
)

// Attribute type constants
const (
	AttributeTypeSize  = "size"
	AttributeTypeColor = "color"
)

type Product struct {
	ID                 int64          `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	Name               string         `gorm:"type:varchar(255)" json:"name"`
	NameEn             string         `gorm:"type:varchar(255)" json:"name_en"`
	NameAr             string         `gorm:"type:varchar(255)" json:"name_ar"`
	Slug               string         `gorm:"type:varchar(255)" json:"slug"`
	Description        string         `gorm:"type:text" json:"description"`
	DescriptionEn      string         `gorm:"type:text" json:"description_en"`
	DescriptionAr      string         `gorm:"type:text" json:"description_ar"`
	Price              float64        `gorm:"type:numeric(10,2)" json:"price"`
	BrandID            *int64         `gorm:"type:bigint" json:"brand_id"`
	CategoryID         *int64         `gorm:"type:bigint" json:"category_id"`
	SupplierID         *int64         `gorm:"type:bigint" json:"supplier_id"`
	IsInternalSupplier bool           `gorm:"default:false" json:"is_internal_supplier"`
	Status             string         `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	IsPublished        bool           `gorm:"default:false" json:"is_published"`
	IsFeatured         bool           `gorm:"default:false" json:"is_featured"`
	IsNew              bool           `gorm:"default:false" json:"is_new"`
	IsBestSeller       bool           `gorm:"default:false" json:"is_best_seller"`
	AttributeType      *string        `gorm:"type:varchar(20)" json:"attribute_type"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Brand       *Brand           `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Category    *Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Supplier    *Supplier        `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Variants    []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	StoreFronts []StoreFront     `gorm:"many2many:product_storefront;" json:"store_fronts,omitempty"`
	Images      []ProductImage   `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	SEO         *ProductSEO      `gorm:"foreignKey:ProductID" json:"seo,omitempty"`
}

func (Product) TableName() string { return "products" }
