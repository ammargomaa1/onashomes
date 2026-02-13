package requests

// CreateProductV2Request is the V2 product creation request with bilingual fields and store assignment
type CreateProductV2Request struct {
	NameEn             string             `json:"name_en" binding:"required"`
	NameAr             string             `json:"name_ar" binding:"required"`
	DescriptionEn      string             `json:"description_en"`
	DescriptionAr      string             `json:"description_ar"`
	BrandID            *int64             `json:"brand_id"`
	CategoryID         *int64             `json:"category_id"`
	SupplierID         *int64             `json:"supplier_id"`
	IsInternalSupplier bool               `json:"is_internal_supplier"`
	AttributeType      *string            `json:"attribute_type"`
	StoreFrontIDs      []int64            `json:"store_front_ids" binding:"required,min=1"`
	IsFeatured         bool               `json:"is_featured"`
	IsNew              bool               `json:"is_new"`
	SEO                *ProductSEORequest `json:"seo"`
}

// UpdateProductV2Request is the V2 product update request
type UpdateProductV2Request struct {
	NameEn             string             `json:"name_en" binding:"required"`
	NameAr             string             `json:"name_ar" binding:"required"`
	DescriptionEn      string             `json:"description_en"`
	DescriptionAr      string             `json:"description_ar"`
	BrandID            *int64             `json:"brand_id"`
	CategoryID         *int64             `json:"category_id"`
	SupplierID         *int64             `json:"supplier_id"`
	IsInternalSupplier bool               `json:"is_internal_supplier"`
	AttributeType      *string            `json:"attribute_type"`
	StoreFrontIDs      []int64            `json:"store_front_ids" binding:"required,min=1"`
	IsFeatured         bool               `json:"is_featured"`
	IsNew              bool               `json:"is_new"`
	IsBestSeller       bool               `json:"is_best_seller"`
	SEO                *ProductSEORequest `json:"seo"`
}

// ProductSEORequest is the SEO metadata request
type ProductSEORequest struct {
	MetaTitleEn       string `json:"meta_title_en" binding:"max=60"`
	MetaTitleAr       string `json:"meta_title_ar" binding:"max=60"`
	MetaDescriptionEn string `json:"meta_description_en" binding:"max=160"`
	MetaDescriptionAr string `json:"meta_description_ar" binding:"max=160"`
	MetaKeywords      string `json:"meta_keywords"`
	CanonicalURL      string `json:"canonical_url"`
	OgTitle           string `json:"og_title"`
	OgDescription     string `json:"og_description"`
	OgImage           string `json:"og_image"`
}

// CreateVariantV2Request is the V2 variant creation request with pricing
type CreateVariantV2Request struct {
	SKU            string   `json:"sku" binding:"required"`
	AttributeValue string   `json:"attribute_value"`
	Price          *float64 `json:"price" binding:"required"`
	CompareAtPrice *float64 `json:"compare_at_price"`
	CostPrice      *float64 `json:"cost_price"`
	Barcode        string   `json:"barcode"`
	Weight         *float64 `json:"weight"`
	IsActive       bool     `json:"is_active"`
	ImageFileIDs   []int64  `json:"image_file_ids"`
}

// UpdateProductStatusRequest is the status change request
type UpdateProductStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft active inactive archived"`
}

// ProductFilterRequest holds admin filtering params
type ProductFilterRequest struct {
	StoreFrontID  *int64  `form:"store_front_id"`
	CategoryID    *int64  `form:"category_id"`
	BrandID       *int64  `form:"brand_id"`
	SupplierID    *int64  `form:"supplier_id"`
	Status        *string `form:"status"`
	AttributeType *string `form:"attribute_type"`
	StockStatus   *string `form:"stock_status"`
	DateFrom      *string `form:"date_from"`
	DateTo        *string `form:"date_to"`
	Search        string  `form:"search"`
}

// AddProductImagesRequest is the request to add images to a product
type AddProductImagesRequest struct {
	FileIDs []int64 `json:"file_ids" binding:"required,min=1"`
}

// SetCoverImageRequest sets the cover image for a product
type SetCoverImageRequest struct {
	FileID int64 `json:"file_id" binding:"required"`
}
