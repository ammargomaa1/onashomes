# Phase 2 — GORM Models

All models follow existing conventions: `int64` IDs, `gorm` tags, `json` tags, soft delete via `gorm.DeletedAt`.

## StoreFront Model

**File:** `internal/models/store_front.go`

```go
package models

import "time"

type StoreFront struct {
    ID              int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
    Name            string    `gorm:"type:varchar(255);not null" json:"name"`
    Slug            string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
    Domain          string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"domain"`
    Currency        string    `gorm:"type:varchar(10);not null;default:'SAR'" json:"currency"`
    DefaultLanguage string    `gorm:"type:varchar(10);not null;default:'ar'" json:"default_language"`
    IsActive        bool      `gorm:"default:true" json:"is_active"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

## Updated Product Model

**File:** `internal/models/product.go` *(NEW — replaces raw SQL struct)*

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// ProductStatus enum
const (
    ProductStatusDraft    = "draft"
    ProductStatusActive   = "active"
    ProductStatusInactive = "inactive"
    ProductStatusArchived = "archived"
)

// AttributeType enum
const (
    AttributeTypeSize  = "size"
    AttributeTypeColor = "color"
)

type Product struct {
    ID                 int64          `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
    NameEn             string         `gorm:"type:varchar(255);not null" json:"name_en"`
    NameAr             string         `gorm:"type:varchar(255)" json:"name_ar"`
    Slug               string         `gorm:"type:varchar(255)" json:"slug"`
    DescriptionEn      string         `gorm:"type:text" json:"description_en"`
    DescriptionAr      string         `gorm:"type:text" json:"description_ar"`
    BrandID            *int64         `gorm:"type:bigint" json:"brand_id"`
    CategoryID         *int64         `gorm:"type:bigint" json:"category_id"`
    SupplierID         *int64         `gorm:"type:bigint" json:"supplier_id"`
    IsInternalSupplier bool           `gorm:"default:false" json:"is_internal_supplier"`
    Status             string         `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
    IsPublished        bool           `gorm:"default:false" json:"is_published"`
    IsFeatured         bool           `gorm:"default:false" json:"is_featured"`
    IsNew              bool           `gorm:"default:false" json:"is_new"`
    IsBestSeller       bool           `gorm:"default:false" json:"is_best_seller"`
    AttributeType      *string        `gorm:"type:varchar(20)" json:"attribute_type"`
    // Legacy field — keep for backward compat during migration
    Name               string         `gorm:"type:varchar(255)" json:"name"`
    Description        string         `gorm:"type:text" json:"description"`
    Price              float64        `gorm:"type:numeric(10,2)" json:"price"`
    IsActive           bool           `gorm:"default:true" json:"is_active"`
    CreatedAt          time.Time      `json:"created_at"`
    UpdatedAt          time.Time      `json:"updated_at"`
    DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

    // Relations (not loaded by default — use Preload)
    Brand       *Brand         `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
    Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    Supplier    *Supplier      `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
    Variants    []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
    StoreFronts []StoreFront   `gorm:"many2many:product_storefront;" json:"store_fronts,omitempty"`
    SEO         *ProductSEO    `gorm:"foreignKey:ProductID" json:"seo,omitempty"`
}
```

## Updated ProductVariant Model

**File:** `internal/models/product_variant.go` *(NEW)*

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type ProductVariant struct {
    ID              int64          `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
    ProductID       int64          `gorm:"type:bigint;not null" json:"product_id"`
    SKU             string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
    AttributeValue  string         `gorm:"type:varchar(255)" json:"attribute_value"` // "XL", "Red"
    Price           *float64       `gorm:"type:numeric(12,2)" json:"price"`
    CompareAtPrice  *float64       `gorm:"type:numeric(12,2)" json:"compare_at_price"`
    CostPrice       *float64       `gorm:"type:numeric(12,2)" json:"cost_price"`
    Barcode         *string        `gorm:"type:varchar(100)" json:"barcode"`
    Weight          *float64       `gorm:"type:numeric(10,3)" json:"weight"`
    IsActive        bool           `gorm:"default:true" json:"is_active"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

    // Relations
    Inventory []VariantInventory `gorm:"foreignKey:ProductVariantID" json:"inventory,omitempty"`
}
```

## VariantInventory Model

**File:** `internal/models/variant_inventory.go`

```go
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

// AvailableQuantity returns sellable quantity
func (v *VariantInventory) AvailableQuantity() int {
    return v.Quantity - v.ReservedQuantity
}

// IsLowStock returns true if available stock <= threshold
func (v *VariantInventory) IsLowStock() bool {
    return v.AvailableQuantity() <= v.LowStockThreshold
}
```

## InventoryAdjustment Model

**File:** `internal/models/inventory_adjustment.go`

```go
package models

import "time"

const (
    AdjustmentReasonAdjustment  = "adjustment"
    AdjustmentReasonCorrection  = "correction"
    AdjustmentReasonRestock     = "restock"
    AdjustmentReasonSale        = "sale"
    AdjustmentReasonReturn      = "return"
)

type InventoryAdjustment struct {
    ID                 int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
    VariantInventoryID int64     `gorm:"type:bigint;not null" json:"variant_inventory_id"`
    AdjustedBy         int64     `gorm:"type:bigint;not null" json:"adjusted_by"`
    PreviousQuantity   int       `gorm:"not null" json:"previous_quantity"`
    NewQuantity        int       `gorm:"not null" json:"new_quantity"`
    AdjustmentAmount   int       `gorm:"not null" json:"adjustment_amount"`
    Reason             string    `gorm:"type:varchar(50);not null" json:"reason"`
    Notes              string    `gorm:"type:text" json:"notes"`
    CreatedAt          time.Time `json:"created_at"`
}
```

## ProductSEO Model

**File:** `internal/models/product_seo.go`

```go
package models

import "time"

type ProductSEO struct {
    ID               int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
    ProductID        int64     `gorm:"type:bigint;uniqueIndex;not null" json:"product_id"`
    MetaTitleEn      string    `gorm:"type:varchar(70)" json:"meta_title_en"`
    MetaTitleAr      string    `gorm:"type:varchar(70)" json:"meta_title_ar"`
    MetaDescriptionEn string   `gorm:"type:varchar(170)" json:"meta_description_en"`
    MetaDescriptionAr string   `gorm:"type:varchar(170)" json:"meta_description_ar"`
    MetaKeywords     string    `gorm:"type:text" json:"meta_keywords"`
    CanonicalURL     string    `gorm:"type:varchar(500)" json:"canonical_url"`
    OgTitle          string    `gorm:"type:varchar(200)" json:"og_title"`
    OgDescription    string    `gorm:"type:varchar(300)" json:"og_description"`
    OgImage          string    `gorm:"type:varchar(500)" json:"og_image"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}
```

## Pivot Models (for explicit control)

**File:** `internal/models/pivots.go`

```go
package models

// ProductStorefront is the M2M pivot for products ↔ store_fronts
type ProductStorefront struct {
    ProductID    int64 `gorm:"primaryKey" json:"product_id"`
    StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (ProductStorefront) TableName() string { return "product_storefront" }

type CategoryStorefront struct {
    CategoryID   int64 `gorm:"primaryKey" json:"category_id"`
    StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (CategoryStorefront) TableName() string { return "category_storefront" }

type BrandStorefront struct {
    BrandID      int64 `gorm:"primaryKey" json:"brand_id"`
    StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (BrandStorefront) TableName() string { return "brand_storefront" }

type SupplierStorefront struct {
    SupplierID   int64 `gorm:"primaryKey" json:"supplier_id"`
    StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (SupplierStorefront) TableName() string { return "supplier_storefront" }
```
