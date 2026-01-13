package products

import (
	"context"
	"database/sql"
	"errors"

	"github.com/onas/ecommerce-api/internal/api/products/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type AdminProductListItem struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	IsActive bool    `json:"is_active"`
}

type AdminAttributeValue struct {
	ID            int64  `json:"id"`
	AttributeID   int64  `json:"attribute_id"`
	AttributeName string `json:"attribute_name"`
	Value         string `json:"value"`
	SortOrder     int    `json:"sort_order"`
}

type AdminProductAttribute struct {
	AttributeID   int64                 `json:"attribute_id"`
	Name          string                `json:"name"`
	SortOrder     int                   `json:"sort_order"`
	AllowedValues []AdminAttributeValue `json:"allowed_values"`
}

type AdminVariantAttributeValue struct {
	AttributeID   int64  `json:"attribute_id"`
	AttributeName string `json:"attribute_name"`
	ValueID       int64  `json:"value_id"`
	Value         string `json:"value"`
}

type AdminVariantImage struct {
	FileID   int64 `json:"file_id"`
	Position int   `json:"position"`
}

type AdminVariantAddOn struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
}

type AdminVariant struct {
	ID              int64                        `json:"id"`
	SKU             string                       `json:"sku"`
	IsActive        bool                         `json:"is_active"`
	AttributeValues []AdminVariantAttributeValue `json:"attribute_values"`
	Images          []AdminVariantImage          `json:"images"`
	AddOns          []AdminVariantAddOn          `json:"add_ons"`
}

type AdminProductDetail struct {
	ID          int64                   `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Price       float64                 `json:"price"`
	IsActive    bool                    `json:"is_active"`
	Attributes  []AdminProductAttribute `json:"attributes"`
	Variants    []AdminVariant          `json:"variants"`
}

type Service interface {
	CreateProduct(ctx context.Context, req requests.AdminProductRequest) utils.IResource
	UpdateProduct(ctx context.Context, id int64, req requests.AdminProductRequest) utils.IResource
	DeleteProduct(ctx context.Context, id int64) utils.IResource
	ListProducts(ctx context.Context, pagination *utils.Pagination) utils.IResource
	GetProductByID(ctx context.Context, id int64) utils.IResource
	CreateVariant(ctx context.Context, productID int64, req requests.AdminVariantRequest) utils.IResource
	UpdateVariant(ctx context.Context, productID, variantID int64, req requests.AdminVariantRequest) utils.IResource
	ListVariantAddOns(ctx context.Context, productID, variantID int64) utils.IResource
	CreateVariantAddOn(ctx context.Context, productID, variantID, addOnProductID int64) utils.IResource
	UpdateVariantAddOn(ctx context.Context, productID, variantID, addOnID, addOnProductID int64) utils.IResource
	DeleteVariantAddOn(ctx context.Context, productID, variantID, addOnID int64) utils.IResource

	// Public APIs
	PublicListProducts(ctx context.Context, pagination *utils.Pagination, q string) utils.IResource
	PublicGetProduct(ctx context.Context, id int64) utils.IResource
}

type service struct {
	db   *sql.DB
	repo Repository
}

func NewService(db *sql.DB, repo Repository) Service {
	return &service{
		db:   db,
		repo: repo,
	}
}

func (s *service) CreateProduct(ctx context.Context, req requests.AdminProductRequest) utils.IResource {
	if req.Name == "" {
		return utils.NewBadRequestResource("name is required", nil)
	}
	if req.Price < 0 {
		return utils.NewBadRequestResource("price must be greater than or equal to zero", nil)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	productID, err := s.repo.CreateProduct(ctx, tx, req)
	if err != nil {
		tx.Rollback()
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if len(req.Attributes) > 0 {
		if err := s.repo.ReplaceProductAttributes(ctx, tx, productID, req.Attributes); err != nil {
			tx.Rollback()
			return utils.NewBadRequestResource(err.Error(), nil)
		}
	}

	if err := tx.Commit(); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	detail, err := s.repo.GetAdminProductByID(ctx, productID)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewCreatedResource("Product created successfully", detail)
}

func (s *service) UpdateProduct(ctx context.Context, id int64, req requests.AdminProductRequest) utils.IResource {
	if req.Name == "" {
		return utils.NewBadRequestResource("name is required", nil)
	}
	if req.Price < 0 {
		return utils.NewBadRequestResource("price must be greater than or equal to zero", nil)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := s.repo.UpdateProduct(ctx, tx, id, req); err != nil {
		tx.Rollback()
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if req.Attributes != nil {
		if err := s.repo.ReplaceProductAttributes(ctx, tx, id, req.Attributes); err != nil {
			tx.Rollback()
			return utils.NewBadRequestResource(err.Error(), nil)
		}
	}

	if err := tx.Commit(); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	detail, err := s.repo.GetAdminProductByID(ctx, id)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Product updated successfully", detail)
}

func (s *service) DeleteProduct(ctx context.Context, id int64) utils.IResource {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := s.repo.SoftDeleteProduct(ctx, tx, id); err != nil {
		tx.Rollback()
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := tx.Commit(); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewNoContentResource()
}

func (s *service) ListProducts(ctx context.Context, pagination *utils.Pagination) utils.IResource {
	products, total, err := s.repo.ListAdminProducts(ctx, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve products", err.Error())
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Products retrieved successfully", products, pagination.GetMeta())
}

func (s *service) GetProductByID(ctx context.Context, id int64) utils.IResource {
	product, err := s.repo.GetAdminProductByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	return utils.NewOKResource("Product retrieved successfully", product)
}

// PublicListProducts lists only active, non-deleted products for public APIs.
func (s *service) PublicListProducts(ctx context.Context, pagination *utils.Pagination, q string) utils.IResource {
	// For now, reuse admin listing and filter via WHERE is_active = TRUE in the repository.
	products, total, err := s.repo.ListPublicProducts(ctx, pagination, q)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve products", err.Error())
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Products retrieved successfully", products, pagination.GetMeta())
}

// PublicGetProduct returns a product and its active variants for public APIs.
func (s *service) PublicGetProduct(ctx context.Context, id int64) utils.IResource {
	product, err := s.repo.GetAdminProductByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundResource("Product not found", nil)
	}
	if !product.IsActive {
		return utils.NewNotFoundResource("Product not found", nil)
	}

	// Filter out inactive variants for public response
	activeVariants := make([]AdminVariant, 0, len(product.Variants))
	for _, v := range product.Variants {
		if v.IsActive {
			activeVariants = append(activeVariants, v)
		}
	}
	product.Variants = activeVariants
	return utils.NewOKResource("Product retrieved successfully", product)
}

func (s *service) CreateVariant(ctx context.Context, productID int64, req requests.AdminVariantRequest) utils.IResource {
	if req.SKU == "" {
		return utils.NewBadRequestResource("sku is required", nil)
	}
	if len(req.AttributeValueIDs) == 0 {
		return utils.NewBadRequestResource("attribute_value_ids is required", nil)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	variantID, err := s.repo.CreateVariant(ctx, tx, productID, req)
	if err != nil {
		tx.Rollback()
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := tx.Commit(); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	variant, err := s.repo.GetAdminVariantByID(ctx, productID, variantID)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewCreatedResource("Variant created successfully", variant)
}

func (s *service) ListVariantAddOns(ctx context.Context, productID, variantID int64) utils.IResource {
	addOns, err := s.repo.ListVariantAddOns(ctx, productID, variantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewNotFoundResource("Product or variant not found", nil)
		}
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Variant add-ons retrieved successfully", addOns)
}

func (s *service) CreateVariantAddOn(ctx context.Context, productID, variantID, addOnProductID int64) utils.IResource {
	if addOnProductID <= 0 {
		return utils.NewBadRequestResource("add_on_product_id is required", nil)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	addOnID, err := s.repo.CreateVariantAddOn(ctx, tx, productID, variantID, addOnProductID)
	if err != nil {
		tx.Rollback()
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := tx.Commit(); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	addOns, err := s.repo.ListVariantAddOns(ctx, productID, variantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewNotFoundResource("Product, variant or add-on product not found", nil)
		}
		return utils.NewBadRequestResource(err.Error(), nil)
	}
	for _, ao := range addOns {
		if ao.ID == addOnID {
			return utils.NewCreatedResource("Variant add-on created successfully", &ao)
		}
	}
	return utils.NewNotFoundResource("Product, variant or add-on product not found", nil)
}

func (s *service) UpdateVariantAddOn(ctx context.Context, productID, variantID, addOnID, addOnProductID int64) utils.IResource {
	if addOnProductID <= 0 {
		return utils.NewBadRequestResource("add_on_product_id is required", nil)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := s.repo.UpdateVariantAddOn(ctx, tx, productID, variantID, addOnID, addOnProductID); err != nil {
		tx.Rollback()
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := tx.Commit(); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	addOns, err := s.repo.ListVariantAddOns(ctx, productID, variantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewNotFoundResource("Product, variant or add-on not found", nil)
		}
		return utils.NewBadRequestResource(err.Error(), nil)
	}
	for _, ao := range addOns {
		if ao.ID == addOnID {
			return utils.NewOKResource("Variant add-on updated successfully", &ao)
		}
	}
	return utils.NewNotFoundResource("Product, variant or add-on not found", nil)
}

func (s *service) DeleteVariantAddOn(ctx context.Context, productID, variantID, addOnID int64) utils.IResource {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := s.repo.DeleteVariantAddOn(ctx, tx, productID, variantID, addOnID); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewNotFoundResource("Product, variant or add-on not found", nil)
		}
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := tx.Commit(); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewNoContentResource()
}

func (s *service) UpdateVariant(ctx context.Context, productID, variantID int64, req requests.AdminVariantRequest) utils.IResource {
	if req.SKU == "" {
		return utils.NewBadRequestResource("sku is required", nil)
	}
	if len(req.AttributeValueIDs) == 0 {
		return utils.NewBadRequestResource("attribute_value_ids is required", nil)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := s.repo.UpdateVariant(ctx, tx, productID, variantID, req); err != nil {
		tx.Rollback()
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	if err := tx.Commit(); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	variant, err := s.repo.GetAdminVariantByID(ctx, productID, variantID)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Variant updated successfully", variant)
}
