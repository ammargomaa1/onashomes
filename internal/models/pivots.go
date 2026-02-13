package models

// ProductStorefront is the M2M pivot for products <-> store_fronts
type ProductStorefront struct {
	ProductID    int64 `gorm:"primaryKey" json:"product_id"`
	StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (ProductStorefront) TableName() string { return "product_storefront" }

// CategoryStorefront is the M2M pivot for categories <-> store_fronts
type CategoryStorefront struct {
	CategoryID   int64 `gorm:"primaryKey" json:"category_id"`
	StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (CategoryStorefront) TableName() string { return "category_storefront" }

// BrandStorefront is the M2M pivot for brands <-> store_fronts
type BrandStorefront struct {
	BrandID      int64 `gorm:"primaryKey" json:"brand_id"`
	StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (BrandStorefront) TableName() string { return "brand_storefront" }

// SupplierStorefront is the M2M pivot for suppliers <-> store_fronts
type SupplierStorefront struct {
	SupplierID   int64 `gorm:"primaryKey" json:"supplier_id"`
	StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (SupplierStorefront) TableName() string { return "supplier_storefront" }

// SectionStorefront is the M2M pivot for sections <-> store_fronts
type SectionStorefront struct {
	SectionID    int64 `gorm:"primaryKey" json:"section_id"`
	StoreFrontID int64 `gorm:"primaryKey" json:"store_front_id"`
}

func (SectionStorefront) TableName() string { return "section_storefront" }
