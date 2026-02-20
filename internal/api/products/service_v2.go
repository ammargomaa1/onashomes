package products

import (
	"fmt"
	"strings"

	"github.com/onas/ecommerce-api/internal/api/inventory"
	"github.com/onas/ecommerce-api/internal/api/products/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

// ServiceV2 handles Phase 2 product operations
type ServiceV2 struct {
	db      *gorm.DB
	repo    *V2Repository
	invRepo *inventory.Repository
}

func NewServiceV2(db *gorm.DB, repo *V2Repository, invRepo *inventory.Repository) *ServiceV2 {
	return &ServiceV2{db: db, repo: repo, invRepo: invRepo}
}

// validateVariantAttributeRule checks that variant matches the product attribute type
func validateVariantAttributeRule(product *models.Product, attributeValue string) error {
	if product.AttributeType == nil {
		if attributeValue != "" {
			return fmt.Errorf("simple product cannot have attribute values")
		}
		return nil
	}
	if attributeValue == "" {
		return fmt.Errorf("variant must specify attribute_value for %s product", *product.AttributeType)
	}
	return nil
}

// resolveAttributeType fetches attribute name by ID
func (s *ServiceV2) resolveAttributeType(id *int64) (*string, error) {
	if id == nil {
		return nil, nil
	}
	var attr models.Attribute
	if err := s.db.First(&attr, *id).Error; err != nil {
		return nil, fmt.Errorf("invalid attribute_type id: %d", *id)
	}
	// Use NameEn as the key for now
	name := strings.ToLower(attr.NameEn)
	return &name, nil
}

func (s *ServiceV2) CreateProductV2(req requests.CreateProductV2Request) utils.IResource {
	// Resolve attribute type ID to string
	attrType, err := s.resolveAttributeType(req.AttributeType)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	// Generate slug
	slug := utils.Slugify(req.NameEn)
	if slug == "" {
		return utils.NewBadRequestResource("Could not generate slug from name_en", nil)
	}

	// Validate slug uniqueness per each store
	for _, sfID := range req.StoreFrontIDs {
		unique, err := s.repo.IsSlugUniqueForStore(slug, sfID, 0)
		if err != nil {
			return utils.NewInternalErrorResource("Failed to validate slug", err)
		}
		if !unique {
			return utils.NewBadRequestResource(fmt.Sprintf("Slug '%s' already exists in store %d", slug, sfID), nil)
		}
	}



	// Validate variants if any
	if len(req.Variants) > 0 {
		skuMap := make(map[string]bool)
		for _, vReq := range req.Variants {
			// 1. Check duplicates within the request
			if skuMap[vReq.SKU] {
				return utils.NewBadRequestResource(fmt.Sprintf("Duplicate SKU '%s' in request", vReq.SKU), nil)
			}
			skuMap[vReq.SKU] = true

			// 2. Check duplicates in DB
			unique, err := s.repo.IsSKUUnique(vReq.SKU, 0)
			if err != nil {
				return utils.NewInternalErrorResource("Failed to validate SKU uniqueness", err)
			}
			if !unique {
				return utils.NewBadRequestResource(fmt.Sprintf("SKU '%s' is already taken", vReq.SKU), nil)
			}
		}
	}

	// Create product in transaction
	var productID int64
	err = s.db.Transaction(func(tx *gorm.DB) error {
		product := &models.Product{
			NameEn:             req.NameEn,
			NameAr:             req.NameAr,
			Name:               req.NameEn, // backward compat
			Slug:               slug,
			DescriptionEn:      req.DescriptionEn,
			DescriptionAr:      req.DescriptionAr,
			Description:        req.DescriptionEn, // backward compat
			BrandID:            req.BrandID,
			CategoryID:         req.CategoryID,
			SupplierID:         req.SupplierID,
			IsInternalSupplier: req.IsInternalSupplier,
			AttributeType:      attrType,
			Status:             models.ProductStatusDraft,
			IsPublished:        false,
			IsFeatured:         req.IsFeatured,
			IsNew:              req.IsNew,
			IsActive:           true,
		}

		if err := tx.Create(product).Error; err != nil {
			return err
		}
		productID = product.ID

		// Assign to stores
		repoTx := &V2Repository{db: tx}
		if err := repoTx.AssignProductToStores(tx, productID, req.StoreFrontIDs); err != nil {
			return err
		}

		// Upsert SEO if provided
		if req.SEO != nil {
			if err := repoTx.UpsertProductSEO(tx, productID, *req.SEO); err != nil {
				return err
			}
		}

		// Create Variants
		if len(req.Variants) > 0 {
			for _, vReq := range req.Variants {
				if err := validateVariantAttributeRule(product, vReq.AttributeValue); err != nil {
					return fmt.Errorf("variant validation failed: %w", err)
				}
				variant := &models.ProductVariant{
					ProductID:      productID,
					SKU:            vReq.SKU,
					AttributeValue: vReq.AttributeValue,
					Price:          vReq.Price,
					CompareAtPrice: vReq.CompareAtPrice,
					CostPrice:      vReq.CostPrice,
					Weight:         vReq.Weight,
					Length:         vReq.Length,
					Width:          vReq.Width,
					Height:         vReq.Height,
					IsActive:       vReq.IsActive,
				}
				if vReq.Barcode != "" {
					variant.Barcode = &vReq.Barcode
				}

				if err := repoTx.CreateVariantV2(variant); err != nil {
					return fmt.Errorf("failed to create variant %s: %w", vReq.SKU, err)
				}

				// Create inventory for all storefronts
				initialStock := 0
				if vReq.Stock != nil {
					initialStock = *vReq.Stock
				}
				for _, sfID := range req.StoreFrontIDs {
					inv := models.VariantInventory{
						ProductVariantID:  variant.ID,
						StoreFrontID:      sfID,
						Quantity:          initialStock,
						LowStockThreshold: 5,
					}
					if err := repoTx.CreateVariantInventory(&inv); err != nil {
						return fmt.Errorf("failed to create inventory for variant %s: %w", vReq.SKU, err)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return utils.NewInternalErrorResource("Failed to create product", err)
	}

	detail, err := s.repo.GetAdminProductV2ByID(productID)
	if err != nil {
		return utils.NewInternalErrorResource("Product created but failed to retrieve", err)
	}

	return utils.NewCreatedResource("Product created successfully", detail)
}

func (s *ServiceV2) UpdateProductV2(id int64, req requests.UpdateProductV2Request) utils.IResource {
	product, err := s.repo.GetProductModelByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	attrType, err := s.resolveAttributeType(req.AttributeType)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	// Re-generate slug if name changed
	slug := utils.Slugify(req.NameEn)
	if slug != product.Slug {
		for _, sfID := range req.StoreFrontIDs {
			unique, err := s.repo.IsSlugUniqueForStore(slug, sfID, id)
			if err != nil {
				return utils.NewInternalErrorResource("Failed to validate slug", err)
			}
			if !unique {
				return utils.NewBadRequestResource(fmt.Sprintf("Slug '%s' already exists in store %d", slug, sfID), nil)
			}
		}
	}

	// Validate variants if any
	if len(req.Variants) > 0 {
		skuMap := make(map[string]bool)
		for _, vReq := range req.Variants {
			// 1. Check duplicates within the request
			if skuMap[vReq.SKU] {
				return utils.NewBadRequestResource(fmt.Sprintf("Duplicate SKU '%s' in request", vReq.SKU), nil)
			}
			skuMap[vReq.SKU] = true

			// 2. Check duplicates in DB (exclude current variant ID if updating)
			excludeID := int64(0)
			if vReq.ID != nil {
				excludeID = *vReq.ID
			}

			unique, err := s.repo.IsSKUUnique(vReq.SKU, excludeID)
			if err != nil {
				return utils.NewInternalErrorResource("Failed to validate SKU uniqueness", err)
			}
			if !unique {
				return utils.NewBadRequestResource(fmt.Sprintf("SKU '%s' is already taken", vReq.SKU), nil)
			}
		}
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		product.NameEn = req.NameEn
		product.NameAr = req.NameAr
		product.Name = req.NameEn
		product.Slug = slug
		product.DescriptionEn = req.DescriptionEn
		product.DescriptionAr = req.DescriptionAr
		product.Description = req.DescriptionEn
		product.BrandID = req.BrandID
		product.CategoryID = req.CategoryID
		product.SupplierID = req.SupplierID
		product.IsInternalSupplier = req.IsInternalSupplier
		product.AttributeType = attrType
		product.IsFeatured = req.IsFeatured
		product.IsNew = req.IsNew
		product.IsBestSeller = req.IsBestSeller

		if err := tx.Save(product).Error; err != nil {
			return err
		}

		repoTx := &V2Repository{db: tx}
		if err := repoTx.AssignProductToStores(tx, id, req.StoreFrontIDs); err != nil {
			return err
		}

		if req.SEO != nil {
			if err := repoTx.UpsertProductSEO(tx, id, *req.SEO); err != nil {
				return err
			}
		}

		// Handle Variants (Upsert)
			for _, vReq := range req.Variants {
				if vReq.ID != nil && *vReq.ID > 0 {
					// Update existing
					var existing models.ProductVariant
					if err := tx.Where("id = ? AND product_id = ?", *vReq.ID, id).First(&existing).Error; err != nil {
						return fmt.Errorf("variant %d not found for product", *vReq.ID)
					}

					if err := validateVariantAttributeRule(product, vReq.AttributeValue); err != nil {
						return fmt.Errorf("variant validation failed: %w", err)
					}

					existing.SKU = vReq.SKU
					existing.AttributeValue = vReq.AttributeValue
					existing.Price = vReq.Price
					existing.CompareAtPrice = vReq.CompareAtPrice
					existing.CostPrice = vReq.CostPrice
					existing.Weight = vReq.Weight
					existing.Length = vReq.Length
					existing.Width = vReq.Width
					existing.Height = vReq.Height
					existing.IsActive = vReq.IsActive
					if vReq.Barcode != "" {
						existing.Barcode = &vReq.Barcode
					}

					if err := repoTx.UpdateVariantV2(&existing); err != nil {
						return fmt.Errorf("failed to update variant %s: %w", vReq.SKU, err)
					}
				} else {
					// Create new
					if err := validateVariantAttributeRule(product, vReq.AttributeValue); err != nil {
						return fmt.Errorf("variant validation failed: %w", err)
					}

					variant := &models.ProductVariant{
						ProductID:      id,
						SKU:            vReq.SKU,
						AttributeValue: vReq.AttributeValue,
						Price:          vReq.Price,
						CompareAtPrice: vReq.CompareAtPrice,
						CostPrice:      vReq.CostPrice,
						Weight:         vReq.Weight,
						Length:         vReq.Length,
						Width:          vReq.Width,
						Height:         vReq.Height,
						IsActive:       vReq.IsActive,
					}
					if vReq.Barcode != "" {
						variant.Barcode = &vReq.Barcode
					}

					if err := repoTx.CreateVariantV2(variant); err != nil {
						return fmt.Errorf("failed to create new variant %s: %w", vReq.SKU, err)
					}

					// Create inventory for new variant
					initialStock := 0
					if vReq.Stock != nil {
						initialStock = *vReq.Stock
					}
					for _, sfID := range req.StoreFrontIDs {
						inv := models.VariantInventory{
							ProductVariantID:  variant.ID,
							StoreFrontID:      sfID,
							Quantity:          initialStock,
							LowStockThreshold: 5,
						}
						// Ignore error if exists (unlikely for new variant)
						_ = repoTx.CreateVariantInventory(&inv)
					}
				}
			}		

		return nil
	})

	if err != nil {
		return utils.NewInternalErrorResource("Failed to update product", err)
	}

	detail, err := s.repo.GetAdminProductV2ByID(id)
	if err != nil {
		return utils.NewInternalErrorResource("Product updated but failed to retrieve", err)
	}

	return utils.NewOKResource("Product updated successfully", detail)
}

func (s *ServiceV2) UpdateProductStatus(id int64, req requests.UpdateProductStatusRequest) utils.IResource {
	product, err := s.repo.GetProductModelByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	newStatus := req.Status

	// Validate transitions
	switch newStatus {
	case models.ProductStatusActive:
		// Must have at least 1 active variant
		count, err := s.repo.CountActiveVariants(id)
		if err != nil {
			return utils.NewInternalErrorResource("Failed to check variants", err)
		}
		if count == 0 {
			return utils.NewBadRequestResource("Cannot activate product: requires at least 1 active variant", nil)
		}

		// Must have inventory
		hasInv, err := s.invRepo.HasInventoryForProduct(id)
		if err != nil {
			return utils.NewInternalErrorResource("Failed to check inventory", err)
		}
		if !hasInv {
			return utils.NewBadRequestResource("Cannot activate product: requires inventory for at least 1 variant", nil)
		}

		// Activate and publish
		if err := s.repo.UpdateProductStatus(id, models.ProductStatusActive, true); err != nil {
			return utils.NewInternalErrorResource("Failed to update status", err)
		}

	case models.ProductStatusInactive:
		if product.Status != models.ProductStatusActive {
			return utils.NewBadRequestResource("Can only deactivate an active product", nil)
		}
		if err := s.repo.UpdateProductStatus(id, models.ProductStatusInactive, false); err != nil {
			return utils.NewInternalErrorResource("Failed to update status", err)
		}

	case models.ProductStatusArchived:
		if product.Status != models.ProductStatusActive && product.Status != models.ProductStatusInactive {
			return utils.NewBadRequestResource("Can only archive active or inactive products", nil)
		}
		if err := s.repo.UpdateProductStatus(id, models.ProductStatusArchived, false); err != nil {
			return utils.NewInternalErrorResource("Failed to update status", err)
		}

	case models.ProductStatusDraft:
		if product.Status != models.ProductStatusArchived {
			return utils.NewBadRequestResource("Can only re-open an archived product to draft", nil)
		}
		if err := s.repo.UpdateProductStatus(id, models.ProductStatusDraft, false); err != nil {
			return utils.NewInternalErrorResource("Failed to update status", err)
		}

	default:
		return utils.NewBadRequestResource("Invalid status", nil)
	}

	detail, err := s.repo.GetAdminProductV2ByID(id)
	if err != nil {
		return utils.NewInternalErrorResource("Status updated but failed to retrieve product", err)
	}

	return utils.NewOKResource("Product status updated successfully", detail)
}

func (s *ServiceV2) UpdateProductSEO(id int64, req requests.ProductSEORequest) utils.IResource {
	_, err := s.repo.GetProductModelByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	// Warn about SEO field lengths
	var warnings []string
	if len(req.MetaTitleEn) > 0 && (len(req.MetaTitleEn) < 50 || len(req.MetaTitleEn) > 60) {
		warnings = append(warnings, "meta_title_en recommended length: 50-60 chars")
	}
	if len(req.MetaDescriptionEn) > 0 && (len(req.MetaDescriptionEn) < 140 || len(req.MetaDescriptionEn) > 160) {
		warnings = append(warnings, "meta_description_en recommended length: 140-160 chars")
	}

	if err := s.repo.UpsertProductSEO(s.db, id, req); err != nil {
		return utils.NewInternalErrorResource("Failed to update SEO", err)
	}

	seo, _ := s.repo.GetProductSEO(id)

	result := map[string]interface{}{
		"seo": seo,
	}
	if len(warnings) > 0 {
		result["warnings"] = warnings
	}

	return utils.NewOKResource("SEO updated successfully", result)
}

func (s *ServiceV2) SuggestSEO(id int64) utils.IResource {
	detail, err := s.repo.GetAdminProductV2ByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	brandName := ""
	if detail.Brand != nil {
		brandName = detail.Brand.NameEn
	}
	categoryName := ""
	if detail.Category != nil {
		categoryName = detail.Category.NameEn
	}

	metaTitleEn := truncateStr(detail.NameEn+" | "+brandName, 60)
	metaDescEn := truncateStr("Shop "+detail.NameEn+" from "+brandName+" in "+categoryName+". "+detail.DescriptionEn, 160)

	keywords := []string{detail.NameEn}
	if brandName != "" {
		keywords = append(keywords, brandName)
	}
	if categoryName != "" {
		keywords = append(keywords, categoryName)
	}

	suggestion := requests.ProductSEORequest{
		MetaTitleEn:       metaTitleEn,
		MetaTitleAr:       truncateStr(detail.NameAr, 60),
		MetaDescriptionEn: metaDescEn,
		MetaDescriptionAr: truncateStr(detail.DescriptionAr, 160),
		MetaKeywords:      strings.Join(keywords, ","),
		CanonicalURL:      "/products/" + detail.Slug,
		OgTitle:           detail.NameEn,
		OgDescription:     truncateStr(detail.DescriptionEn, 200),
	}

	return utils.NewOKResource("SEO suggestions generated", suggestion)
}

func (s *ServiceV2) ListAdminProductsV2(filter requests.ProductFilterRequest, pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.ListAdminProductsV2(filter, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve products", err)
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Products retrieved successfully", items, pagination.GetMeta())
}

func (s *ServiceV2) GetProductV2ByID(id int64) utils.IResource {
	detail, err := s.repo.GetAdminProductV2ByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	return utils.NewOKResource("Product retrieved successfully", detail)
}

// CreateVariantV2 creates a new variant with pricing
func (s *ServiceV2) CreateVariantV2(productID int64, req requests.CreateVariantV2Request) utils.IResource {
	product, err := s.repo.GetProductModelByID(productID)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	// Validate attribute rule
	if err := validateVariantAttributeRule(product, req.AttributeValue); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	// Validate SKU uniqueness
	unique, err := s.repo.IsSKUUnique(req.SKU, 0)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate SKU", err)
	}
	if !unique {
		return utils.NewBadRequestResource("SKU already exists", nil)
	}

	variant := &models.ProductVariant{
		ProductID:      productID,
		SKU:            req.SKU,
		AttributeValue: req.AttributeValue,
		Price:          req.Price,
		CompareAtPrice: req.CompareAtPrice,
		CostPrice:      req.CostPrice,
		Weight:         req.Weight,
		IsActive:       req.IsActive,
	}
	if req.Barcode != "" {
		variant.Barcode = &req.Barcode
	}

	if err := s.repo.CreateVariantV2(variant); err != nil {
		return utils.NewInternalErrorResource("Failed to create variant", err)
	}

	// Create inventory for all storefronts
	sfIDs, err := s.repo.GetProductStoreFrontIDs(productID)
	if err == nil && len(sfIDs) > 0 {
		initialStock := 0
		if req.Stock != nil {
			initialStock = *req.Stock
		}
		for _, sfID := range sfIDs {
			inv := models.VariantInventory{
				ProductVariantID:  variant.ID,
				StoreFrontID:      sfID,
				Quantity:          initialStock,
				LowStockThreshold: 5, // Default
			}
			_ = s.repo.CreateVariantInventory(&inv)
		}
	}

	return utils.NewCreatedResource("Variant created successfully", AdminVariantV2{
		ID:             variant.ID,
		SKU:            variant.SKU,
		AttributeValue: variant.AttributeValue,
		Price:          variant.Price,
		CompareAtPrice: variant.CompareAtPrice,
		CostPrice:      variant.CostPrice,
		Barcode:        variant.Barcode,
		Weight:         variant.Weight,
		Length:         variant.Length,
		Width:          variant.Width,
		Height:         variant.Height,
		IsActive:       variant.IsActive,
	})
}

// UpdateVariantV2 updates a variant with pricing
func (s *ServiceV2) UpdateVariantV2(productID, variantID int64, req requests.CreateVariantV2Request) utils.IResource {
	product, err := s.repo.GetProductModelByID(productID)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	variant, err := s.repo.GetVariantByID(variantID)
	if err != nil {
		return utils.NewNotFoundResource("Variant not found", nil)
	}
	if variant.ProductID != productID {
		return utils.NewNotFoundResource("Variant does not belong to this product", nil)
	}

	if err := validateVariantAttributeRule(product, req.AttributeValue); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	// SKU uniqueness check (exclude self)
	unique, err := s.repo.IsSKUUnique(req.SKU, variantID)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate SKU", err)
	}
	if !unique {
		return utils.NewBadRequestResource("SKU already exists", nil)
	}

	variant.SKU = req.SKU
	variant.AttributeValue = req.AttributeValue
	variant.Price = req.Price
	variant.CompareAtPrice = req.CompareAtPrice
	variant.CostPrice = req.CostPrice
	variant.Weight = req.Weight
	variant.IsActive = req.IsActive
	if req.Barcode != "" {
		variant.Barcode = &req.Barcode
	}

	if err := s.repo.UpdateVariantV2(variant); err != nil {
		return utils.NewInternalErrorResource("Failed to update variant", err)
	}

	return utils.NewOKResource("Variant updated successfully", AdminVariantV2{
		ID:             variant.ID,
		SKU:            variant.SKU,
		AttributeValue: variant.AttributeValue,
		Price:          variant.Price,
		CompareAtPrice: variant.CompareAtPrice,
		CostPrice:      variant.CostPrice,
		Barcode:        variant.Barcode,
		Weight:         variant.Weight,
		IsActive:       variant.IsActive,
	})
}

// ============== Product Images ==============

func (s *ServiceV2) AddProductImages(productID int64, req requests.AddProductImagesRequest) utils.IResource {
	_, err := s.repo.GetProductModelByID(productID)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	// Validate file IDs exist
	var fileCount int64
	s.db.Model(&models.File{}).Where("id IN ? AND deleted_at IS NULL", req.FileIDs).Count(&fileCount)
	if fileCount != int64(len(req.FileIDs)) {
		return utils.NewBadRequestResource("One or more file IDs are invalid", nil)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		repoTx := &V2Repository{db: tx}
		return repoTx.AddProductImages(tx, productID, req.FileIDs)
	})
	if err != nil {
		return utils.NewInternalErrorResource("Failed to add images", err)
	}

	detail, _ := s.repo.GetAdminProductV2ByID(productID)
	return utils.NewOKResource("Images added successfully", detail)
}

func (s *ServiceV2) RemoveProductImage(productID, imageID int64) utils.IResource {
	_, err := s.repo.GetProductModelByID(productID)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		repoTx := &V2Repository{db: tx}
		return repoTx.RemoveProductImage(tx, productID, imageID)
	})
	if err != nil {
		return utils.NewInternalErrorResource("Failed to remove image", err)
	}

	detail, _ := s.repo.GetAdminProductV2ByID(productID)
	return utils.NewOKResource("Image removed successfully", detail)
}

func (s *ServiceV2) SetCoverImage(productID int64, req requests.SetCoverImageRequest) utils.IResource {
	_, err := s.repo.GetProductModelByID(productID)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		repoTx := &V2Repository{db: tx}
		return repoTx.SetCoverImage(tx, productID, req.FileID)
	})
	if err != nil {
		return utils.NewInternalErrorResource("Failed to set cover image", err)
	}

	detail, _ := s.repo.GetAdminProductV2ByID(productID)
	return utils.NewOKResource("Cover image set successfully", detail)
}

// Storefront public APIs
func (s *ServiceV2) StorefrontListProducts(storeFrontID int64, filter requests.ProductFilterRequest, pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.ListStorefrontProducts(storeFrontID, filter, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve products", err)
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Products retrieved successfully", items, pagination.GetMeta())
}

func (s *ServiceV2) StorefrontGetProduct(storeFrontID int64, slug string) utils.IResource {
	detail, err := s.repo.GetStorefrontProduct(storeFrontID, slug)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	return utils.NewOKResource("Product retrieved successfully", detail)
}

func (s *ServiceV2) GetProductStructuredData(storeFrontID int64, slug string) utils.IResource {
	jsonLD, err := s.repo.GetProductStructuredData(storeFrontID, slug)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	return utils.NewOKResource("Structured data generated", jsonLD)
}

// truncateStr truncates a string to maxLen
func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
