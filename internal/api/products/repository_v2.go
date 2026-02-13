package products

import (
	"fmt"
	"strings"
	"time"

	"github.com/onas/ecommerce-api/internal/api/products/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

// V2Repository handles V2 product operations using GORM
type V2Repository struct {
	db *gorm.DB
}

func NewV2Repository(db *gorm.DB) *V2Repository {
	return &V2Repository{db: db}
}

// ============== DTOs ==============

type AdminProductV2ListItem struct {
	ID            int64     `json:"id"`
	NameEn        string    `json:"name_en"`
	NameAr        string    `json:"name_ar"`
	Slug          string    `json:"slug"`
	Status        string    `json:"status"`
	IsPublished   bool      `json:"is_published"`
	IsFeatured    bool      `json:"is_featured"`
	AttributeType *string   `json:"attribute_type"`
	BrandName     *string   `json:"brand_name"`
	CategoryName  *string   `json:"category_name"`
	VariantCount  int64     `json:"variant_count"`
	MinPrice      *float64  `json:"min_price"`
	MaxPrice      *float64  `json:"max_price"`
	CreatedAt     time.Time `json:"created_at"`
}

type AdminProductV2Detail struct {
	ID                 int64              `json:"id"`
	NameEn             string             `json:"name_en"`
	NameAr             string             `json:"name_ar"`
	Slug               string             `json:"slug"`
	DescriptionEn      string             `json:"description_en"`
	DescriptionAr      string             `json:"description_ar"`
	Status             string             `json:"status"`
	IsPublished        bool               `json:"is_published"`
	IsFeatured         bool               `json:"is_featured"`
	IsNew              bool               `json:"is_new"`
	IsBestSeller       bool               `json:"is_best_seller"`
	IsInternalSupplier bool               `json:"is_internal_supplier"`
	AttributeType      *string            `json:"attribute_type"`
	Brand              *BrandInfo         `json:"brand"`
	Category           *CategoryInfo      `json:"category"`
	Supplier           *SupplierInfo      `json:"supplier"`
	StoreFronts        []StoreFrontInfo   `json:"store_fronts"`
	Variants           []AdminVariantV2   `json:"variants"`
	Images             []ProductImageInfo `json:"images"`
	SEO                *models.ProductSEO `json:"seo"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

type ProductImageInfo struct {
	ID       int64  `json:"id"`
	FileID   int64  `json:"file_id"`
	URL      string `json:"url"`
	Position int    `json:"position"`
	IsCover  bool   `json:"is_cover"`
}

type BrandInfo struct {
	ID     int64  `json:"id"`
	NameEn string `json:"name_en"`
}

type CategoryInfo struct {
	ID     int64  `json:"id"`
	NameEn string `json:"name_en"`
}

type SupplierInfo struct {
	ID          int64  `json:"id"`
	CompanyName string `json:"company_name"`
}

type StoreFrontInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type AdminVariantV2 struct {
	ID             int64    `json:"id"`
	SKU            string   `json:"sku"`
	AttributeValue string   `json:"attribute_value"`
	Price          *float64 `json:"price"`
	CompareAtPrice *float64 `json:"compare_at_price"`
	CostPrice      *float64 `json:"cost_price"`
	Barcode        *string  `json:"barcode"`
	Weight         *float64 `json:"weight"`
	IsActive       bool     `json:"is_active"`
}

type StorefrontProductItem struct {
	ID             int64    `json:"id"`
	NameEn         string   `json:"name_en"`
	NameAr         string   `json:"name_ar"`
	Slug           string   `json:"slug"`
	BrandName      *string  `json:"brand_name"`
	CategoryName   *string  `json:"category_name"`
	IsFeatured     bool     `json:"is_featured"`
	IsNew          bool     `json:"is_new"`
	IsBestSeller   bool     `json:"is_best_seller"`
	MinPrice       *float64 `json:"min_price"`
	MaxPrice       *float64 `json:"max_price"`
	CompareAtPrice *float64 `json:"compare_at_price"`
	InStock        bool     `json:"in_stock"`
	VariantCount   int64    `json:"variant_count"`
}

type StorefrontProductDetail struct {
	ID            int64               `json:"id"`
	NameEn        string              `json:"name_en"`
	NameAr        string              `json:"name_ar"`
	Slug          string              `json:"slug"`
	DescriptionEn string              `json:"description_en"`
	DescriptionAr string              `json:"description_ar"`
	AttributeType *string             `json:"attribute_type"`
	IsFeatured    bool                `json:"is_featured"`
	IsNew         bool                `json:"is_new"`
	Brand         *BrandInfo          `json:"brand"`
	Category      *CategoryInfo       `json:"category"`
	SEO           *models.ProductSEO  `json:"seo"`
	Variants      []StorefrontVariant `json:"variants"`
}

type StorefrontVariant struct {
	ID             int64    `json:"id"`
	SKU            string   `json:"sku"`
	AttributeValue string   `json:"attribute_value"`
	Price          *float64 `json:"price"`
	CompareAtPrice *float64 `json:"compare_at_price"`
	InStock        bool     `json:"in_stock"`
}

// ============== Repository Methods ==============

func (r *V2Repository) CreateProductV2(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *V2Repository) UpdateProductV2(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *V2Repository) GetProductModelByID(id int64) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *V2Repository) IsSlugUniqueForStore(slug string, storeFrontID int64, excludeProductID int64) (bool, error) {
	var count int64
	query := r.db.Table("product_storefront ps").
		Joins("JOIN products p ON p.id = ps.product_id").
		Where("p.slug = ? AND ps.store_front_id = ? AND p.deleted_at IS NULL", slug, storeFrontID)
	if excludeProductID > 0 {
		query = query.Where("p.id <> ?", excludeProductID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *V2Repository) AssignProductToStores(tx *gorm.DB, productID int64, storeFrontIDs []int64) error {
	// Remove existing assignments
	if err := tx.Where("product_id = ?", productID).Delete(&models.ProductStorefront{}).Error; err != nil {
		return err
	}
	// Insert new assignments
	for _, sfID := range storeFrontIDs {
		pivot := models.ProductStorefront{
			ProductID:    productID,
			StoreFrontID: sfID,
		}
		if err := tx.Create(&pivot).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *V2Repository) UpsertProductSEO(tx *gorm.DB, productID int64, seo requests.ProductSEORequest) error {
	var existing models.ProductSEO
	err := tx.Where("product_id = ?", productID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		// Create
		newSEO := models.ProductSEO{
			ProductID:         productID,
			MetaTitleEn:       seo.MetaTitleEn,
			MetaTitleAr:       seo.MetaTitleAr,
			MetaDescriptionEn: seo.MetaDescriptionEn,
			MetaDescriptionAr: seo.MetaDescriptionAr,
			MetaKeywords:      seo.MetaKeywords,
			CanonicalURL:      seo.CanonicalURL,
			OgTitle:           seo.OgTitle,
			OgDescription:     seo.OgDescription,
			OgImage:           seo.OgImage,
		}
		return tx.Create(&newSEO).Error
	}
	if err != nil {
		return err
	}

	// Update
	return tx.Model(&existing).Updates(map[string]interface{}{
		"meta_title_en":       seo.MetaTitleEn,
		"meta_title_ar":       seo.MetaTitleAr,
		"meta_description_en": seo.MetaDescriptionEn,
		"meta_description_ar": seo.MetaDescriptionAr,
		"meta_keywords":       seo.MetaKeywords,
		"canonical_url":       seo.CanonicalURL,
		"og_title":            seo.OgTitle,
		"og_description":      seo.OgDescription,
		"og_image":            seo.OgImage,
		"updated_at":          time.Now(),
	}).Error
}

func (r *V2Repository) GetProductSEO(productID int64) (*models.ProductSEO, error) {
	var seo models.ProductSEO
	err := r.db.Where("product_id = ?", productID).First(&seo).Error
	if err != nil {
		return nil, err
	}
	return &seo, nil
}

func (r *V2Repository) GetAdminProductV2ByID(productID int64) (*AdminProductV2Detail, error) {
	var product models.Product
	err := r.db.Where("id = ? AND deleted_at IS NULL", productID).First(&product).Error
	if err != nil {
		return nil, err
	}

	detail := &AdminProductV2Detail{
		ID:                 product.ID,
		NameEn:             product.NameEn,
		NameAr:             product.NameAr,
		Slug:               product.Slug,
		DescriptionEn:      product.DescriptionEn,
		DescriptionAr:      product.DescriptionAr,
		Status:             product.Status,
		IsPublished:        product.IsPublished,
		IsFeatured:         product.IsFeatured,
		IsNew:              product.IsNew,
		IsBestSeller:       product.IsBestSeller,
		IsInternalSupplier: product.IsInternalSupplier,
		AttributeType:      product.AttributeType,
		CreatedAt:          product.CreatedAt,
		UpdatedAt:          product.UpdatedAt,
	}

	// Load brand
	if product.BrandID != nil {
		var brand models.Brand
		if err := r.db.Select("id, name_en").First(&brand, *product.BrandID).Error; err == nil {
			detail.Brand = &BrandInfo{ID: brand.ID, NameEn: brand.NameEn}
		}
	}

	// Load category
	if product.CategoryID != nil {
		var cat models.Category
		if err := r.db.Select("id, name_en").First(&cat, *product.CategoryID).Error; err == nil {
			detail.Category = &CategoryInfo{ID: cat.ID, NameEn: cat.NameEn}
		}
	}

	// Load supplier
	if product.SupplierID != nil {
		var sup models.Supplier
		if err := r.db.Select("id, company_name").First(&sup, *product.SupplierID).Error; err == nil {
			detail.Supplier = &SupplierInfo{ID: sup.ID, CompanyName: sup.CompanyName}
		}
	}

	// Load store fronts
	var sfPivots []models.ProductStorefront
	r.db.Where("product_id = ?", productID).Find(&sfPivots)
	if len(sfPivots) > 0 {
		sfIDs := make([]int64, len(sfPivots))
		for i, p := range sfPivots {
			sfIDs[i] = p.StoreFrontID
		}
		var sfs []models.StoreFront
		r.db.Where("id IN ?", sfIDs).Find(&sfs)
		for _, sf := range sfs {
			detail.StoreFronts = append(detail.StoreFronts, StoreFrontInfo{
				ID: sf.ID, Name: sf.Name, Domain: sf.Domain,
			})
		}
	}
	if detail.StoreFronts == nil {
		detail.StoreFronts = []StoreFrontInfo{}
	}

	// Load variants
	var variants []models.ProductVariant
	r.db.Where("product_id = ? AND deleted_at IS NULL", productID).Order("id ASC").Find(&variants)
	for _, v := range variants {
		detail.Variants = append(detail.Variants, AdminVariantV2{
			ID:             v.ID,
			SKU:            v.SKU,
			AttributeValue: v.AttributeValue,
			Price:          v.Price,
			CompareAtPrice: v.CompareAtPrice,
			CostPrice:      v.CostPrice,
			Barcode:        v.Barcode,
			Weight:         v.Weight,
			IsActive:       v.IsActive,
		})
	}
	if detail.Variants == nil {
		detail.Variants = []AdminVariantV2{}
	}

	// Load images
	images, _ := r.GetProductImages(productID)
	detail.Images = images
	if detail.Images == nil {
		detail.Images = []ProductImageInfo{}
	}

	// Load SEO
	detail.SEO, _ = r.GetProductSEO(productID)

	return detail, nil
}

func (r *V2Repository) ListAdminProductsV2(filter requests.ProductFilterRequest, pagination *utils.Pagination) ([]AdminProductV2ListItem, int64, error) {
	query := r.db.Table("products p").
		Joins("LEFT JOIN brands b ON b.id = p.brand_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id").
		Where("p.deleted_at IS NULL")

	// Apply filters
	if filter.StoreFrontID != nil {
		query = query.Joins("JOIN product_storefront ps ON ps.product_id = p.id").
			Where("ps.store_front_id = ?", *filter.StoreFrontID)
	}
	if filter.CategoryID != nil {
		query = query.Where("p.category_id = ?", *filter.CategoryID)
	}
	if filter.BrandID != nil {
		query = query.Where("p.brand_id = ?", *filter.BrandID)
	}
	if filter.SupplierID != nil {
		query = query.Where("p.supplier_id = ?", *filter.SupplierID)
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("p.status = ?", *filter.Status)
	}
	if filter.AttributeType != nil {
		if *filter.AttributeType == "simple" {
			query = query.Where("p.attribute_type IS NULL")
		} else {
			query = query.Where("p.attribute_type = ?", *filter.AttributeType)
		}
	}
	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("(LOWER(p.name_en) LIKE ? OR LOWER(p.name_ar) LIKE ?)", search, search)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Select with variant aggregation
	var items []AdminProductV2ListItem
	offset := (pagination.Page - 1) * pagination.Limit
	err := query.
		Select(`
			p.id, p.name_en, p.name_ar, p.slug, p.status, p.is_published, p.is_featured,
			p.attribute_type, b.name_en as brand_name, c.name_en as category_name,
			p.created_at,
			(SELECT COUNT(*) FROM product_variants pv WHERE pv.product_id = p.id AND pv.deleted_at IS NULL) as variant_count,
			(SELECT MIN(pv.price) FROM product_variants pv WHERE pv.product_id = p.id AND pv.deleted_at IS NULL) as min_price,
			(SELECT MAX(pv.price) FROM product_variants pv WHERE pv.product_id = p.id AND pv.deleted_at IS NULL) as max_price
		`).
		Order("p.created_at DESC").
		Offset(offset).Limit(pagination.Limit).
		Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ============== Variant V2 ==============

func (r *V2Repository) CreateVariantV2(variant *models.ProductVariant) error {
	return r.db.Create(variant).Error
}

func (r *V2Repository) UpdateVariantV2(variant *models.ProductVariant) error {
	return r.db.Save(variant).Error
}

func (r *V2Repository) GetVariantByID(variantID int64) (*models.ProductVariant, error) {
	var v models.ProductVariant
	err := r.db.Where("id = ? AND deleted_at IS NULL", variantID).First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *V2Repository) IsSKUUnique(sku string, excludeVariantID int64) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProductVariant{}).Where("sku = ? AND deleted_at IS NULL", sku)
	if excludeVariantID > 0 {
		query = query.Where("id <> ?", excludeVariantID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *V2Repository) CountActiveVariants(productID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.ProductVariant{}).
		Where("product_id = ? AND is_active = true AND deleted_at IS NULL", productID).
		Count(&count).Error
	return count, err
}

func (r *V2Repository) UpdateProductStatus(productID int64, status string, isPublished bool) error {
	updates := map[string]interface{}{
		"status":       status,
		"is_published": isPublished,
		"updated_at":   time.Now(),
	}
	return r.db.Model(&models.Product{}).Where("id = ?", productID).Updates(updates).Error
}

// ============== Product Images ==============

func (r *V2Repository) GetProductImages(productID int64) ([]ProductImageInfo, error) {
	var images []models.ProductImage
	err := r.db.Where("product_id = ?", productID).
		Preload("File").
		Order("position ASC, id ASC").
		Find(&images).Error
	if err != nil {
		return nil, err
	}

	result := make([]ProductImageInfo, 0, len(images))
	for _, img := range images {
		info := ProductImageInfo{
			ID:       img.ID,
			FileID:   img.FileID,
			Position: img.Position,
			IsCover:  img.IsCover,
		}
		if img.File != nil {
			info.URL = img.File.FilePath
		}
		result = append(result, info)
	}
	return result, nil
}

func (r *V2Repository) AddProductImages(tx *gorm.DB, productID int64, fileIDs []int64) error {
	// Get current max position
	var maxPos int
	tx.Model(&models.ProductImage{}).Where("product_id = ?", productID).
		Select("COALESCE(MAX(position), -1)").Scan(&maxPos)

	// Check if product has any images (first image becomes cover)
	var imageCount int64
	tx.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&imageCount)

	for i, fileID := range fileIDs {
		img := models.ProductImage{
			ProductID: productID,
			FileID:    fileID,
			Position:  maxPos + 1 + i,
			IsCover:   imageCount == 0 && i == 0, // first image of new product = cover
		}
		if err := tx.Create(&img).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *V2Repository) RemoveProductImage(tx *gorm.DB, productID, imageID int64) error {
	var img models.ProductImage
	err := tx.Where("id = ? AND product_id = ?", imageID, productID).First(&img).Error
	if err != nil {
		return err
	}

	wasCover := img.IsCover
	if err := tx.Delete(&img).Error; err != nil {
		return err
	}

	// If we removed the cover, promote the first remaining image
	if wasCover {
		var firstImg models.ProductImage
		err := tx.Where("product_id = ?", productID).Order("position ASC, id ASC").First(&firstImg).Error
		if err == nil {
			tx.Model(&firstImg).Update("is_cover", true)
		}
	}
	return nil
}

func (r *V2Repository) SetCoverImage(tx *gorm.DB, productID, fileID int64) error {
	// Unset all covers for this product
	tx.Model(&models.ProductImage{}).Where("product_id = ?", productID).Update("is_cover", false)

	// Set the one with matching file_id
	result := tx.Model(&models.ProductImage{}).
		Where("product_id = ? AND file_id = ?", productID, fileID).
		Update("is_cover", true)
	if result.RowsAffected == 0 {
		return fmt.Errorf("image with file_id %d not found for product %d", fileID, productID)
	}
	return nil
}

// ============== Storefront Queries ==============

func (r *V2Repository) ListStorefrontProducts(storeFrontID int64, filter requests.ProductFilterRequest, pagination *utils.Pagination) ([]StorefrontProductItem, int64, error) {
	query := r.db.Table("products p").
		Joins("JOIN product_storefront ps ON ps.product_id = p.id").
		Joins("LEFT JOIN brands b ON b.id = p.brand_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id").
		Where("ps.store_front_id = ? AND p.status = 'active' AND p.is_published = true AND p.deleted_at IS NULL", storeFrontID)

	if filter.CategoryID != nil {
		query = query.Where("p.category_id = ?", *filter.CategoryID)
	}
	if filter.BrandID != nil {
		query = query.Where("p.brand_id = ?", *filter.BrandID)
	}
	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("(LOWER(p.name_en) LIKE ? OR LOWER(p.name_ar) LIKE ?)", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []StorefrontProductItem
	offset := (pagination.Page - 1) * pagination.Limit
	err := query.
		Select(`
			p.id, p.name_en, p.name_ar, p.slug, p.is_featured, p.is_new, p.is_best_seller,
			b.name_en as brand_name, c.name_en as category_name,
			(SELECT MIN(pv.price) FROM product_variants pv WHERE pv.product_id = p.id AND pv.is_active = true AND pv.deleted_at IS NULL) as min_price,
			(SELECT MAX(pv.price) FROM product_variants pv WHERE pv.product_id = p.id AND pv.is_active = true AND pv.deleted_at IS NULL) as max_price,
			(SELECT MAX(pv.compare_at_price) FROM product_variants pv WHERE pv.product_id = p.id AND pv.is_active = true AND pv.deleted_at IS NULL) as compare_at_price,
			(SELECT COUNT(*) FROM product_variants pv WHERE pv.product_id = p.id AND pv.is_active = true AND pv.deleted_at IS NULL) as variant_count,
			EXISTS(SELECT 1 FROM variant_inventory vi JOIN product_variants pv ON pv.id = vi.product_variant_id WHERE pv.product_id = p.id AND vi.store_front_id = ps.store_front_id AND vi.quantity > vi.reserved_quantity) as in_stock
		`).
		Order("p.created_at DESC").
		Offset(offset).Limit(pagination.Limit).
		Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *V2Repository) GetStorefrontProduct(storeFrontID int64, slug string) (*StorefrontProductDetail, error) {
	var product models.Product
	err := r.db.Table("products p").
		Joins("JOIN product_storefront ps ON ps.product_id = p.id").
		Where("ps.store_front_id = ? AND p.slug = ? AND p.status = 'active' AND p.is_published = true AND p.deleted_at IS NULL", storeFrontID, slug).
		Select("p.*").
		First(&product).Error
	if err != nil {
		return nil, err
	}

	detail := &StorefrontProductDetail{
		ID:            product.ID,
		NameEn:        product.NameEn,
		NameAr:        product.NameAr,
		Slug:          product.Slug,
		DescriptionEn: product.DescriptionEn,
		DescriptionAr: product.DescriptionAr,
		AttributeType: product.AttributeType,
		IsFeatured:    product.IsFeatured,
		IsNew:         product.IsNew,
	}

	// Brand
	if product.BrandID != nil {
		var brand models.Brand
		if err := r.db.Select("id, name_en").First(&brand, *product.BrandID).Error; err == nil {
			detail.Brand = &BrandInfo{ID: brand.ID, NameEn: brand.NameEn}
		}
	}

	// Category
	if product.CategoryID != nil {
		var cat models.Category
		if err := r.db.Select("id, name_en").First(&cat, *product.CategoryID).Error; err == nil {
			detail.Category = &CategoryInfo{ID: cat.ID, NameEn: cat.NameEn}
		}
	}

	// SEO
	detail.SEO, _ = r.GetProductSEO(product.ID)

	// Active variants with stock info
	var variants []models.ProductVariant
	r.db.Where("product_id = ? AND is_active = true AND deleted_at IS NULL", product.ID).Order("id ASC").Find(&variants)
	for _, v := range variants {
		// Check stock for this store
		var invCount int64
		r.db.Table("variant_inventory").
			Where("product_variant_id = ? AND store_front_id = ? AND quantity > reserved_quantity", v.ID, storeFrontID).
			Count(&invCount)

		detail.Variants = append(detail.Variants, StorefrontVariant{
			ID:             v.ID,
			SKU:            v.SKU,
			AttributeValue: v.AttributeValue,
			Price:          v.Price,
			CompareAtPrice: v.CompareAtPrice,
			InStock:        invCount > 0,
		})
	}
	if detail.Variants == nil {
		detail.Variants = []StorefrontVariant{}
	}

	return detail, nil
}

// GetProductStructuredData generates JSON-LD structured data
func (r *V2Repository) GetProductStructuredData(storeFrontID int64, slug string) (map[string]interface{}, error) {
	detail, err := r.GetStorefrontProduct(storeFrontID, slug)
	if err != nil {
		return nil, err
	}

	// Get store domain for URLs
	var storeFront models.StoreFront
	r.db.First(&storeFront, storeFrontID)

	// Calculate price range
	var lowPrice, highPrice float64
	for _, v := range detail.Variants {
		if v.Price != nil {
			if lowPrice == 0 || *v.Price < lowPrice {
				lowPrice = *v.Price
			}
			if *v.Price > highPrice {
				highPrice = *v.Price
			}
		}
	}

	// Determine availability
	availability := "https://schema.org/OutOfStock"
	for _, v := range detail.Variants {
		if v.InStock {
			availability = "https://schema.org/InStock"
			break
		}
	}

	jsonLD := map[string]interface{}{
		"@context":    "https://schema.org",
		"@type":       "Product",
		"name":        detail.NameEn,
		"description": detail.DescriptionEn,
		"offers": map[string]interface{}{
			"@type":         "AggregateOffer",
			"lowPrice":      lowPrice,
			"highPrice":     highPrice,
			"priceCurrency": storeFront.Currency,
			"availability":  availability,
			"offerCount":    len(detail.Variants),
		},
	}

	if detail.Brand != nil {
		jsonLD["brand"] = map[string]interface{}{
			"@type": "Brand",
			"name":  detail.Brand.NameEn,
		}
	}

	if len(detail.Variants) > 0 {
		jsonLD["sku"] = detail.Variants[0].SKU
	}

	// Breadcrumb
	breadcrumbItems := []map[string]interface{}{
		{
			"@type":    "ListItem",
			"position": 1,
			"name":     "Home",
			"item":     fmt.Sprintf("https://%s", storeFront.Domain),
		},
	}
	if detail.Category != nil {
		breadcrumbItems = append(breadcrumbItems, map[string]interface{}{
			"@type":    "ListItem",
			"position": 2,
			"name":     detail.Category.NameEn,
			"item":     fmt.Sprintf("https://%s/categories/%s", storeFront.Domain, utils.Slugify(detail.Category.NameEn)),
		})
	}
	breadcrumbItems = append(breadcrumbItems, map[string]interface{}{
		"@type":    "ListItem",
		"position": len(breadcrumbItems) + 1,
		"name":     detail.NameEn,
	})

	jsonLD["breadcrumb"] = map[string]interface{}{
		"@type":           "BreadcrumbList",
		"itemListElement": breadcrumbItems,
	}

	if detail.SEO != nil {
		if detail.SEO.CanonicalURL != "" {
			jsonLD["url"] = detail.SEO.CanonicalURL
		}
		if detail.SEO.OgImage != "" {
			jsonLD["image"] = detail.SEO.OgImage
		}
	}

	return jsonLD, nil
}
