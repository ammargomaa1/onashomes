package products

import (
	"context"
	"database/sql"
	"fmt"

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
	CreateProduct(ctx context.Context, req requests.AdminProductRequest) (*AdminProductDetail, error)
	UpdateProduct(ctx context.Context, id int64, req requests.AdminProductRequest) (*AdminProductDetail, error)
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context, pagination *utils.Pagination) ([]AdminProductListItem, int64, error)
	GetProductByID(ctx context.Context, id int64) (*AdminProductDetail, error)
	CreateVariant(ctx context.Context, productID int64, req requests.AdminVariantRequest) (*AdminVariant, error)
	UpdateVariant(ctx context.Context, productID, variantID int64, req requests.AdminVariantRequest) (*AdminVariant, error)
	ListVariantAddOns(ctx context.Context, productID, variantID int64) ([]AdminVariantAddOn, error)
	CreateVariantAddOn(ctx context.Context, productID, variantID, addOnProductID int64) (*AdminVariantAddOn, error)
	UpdateVariantAddOn(ctx context.Context, productID, variantID, addOnID, addOnProductID int64) (*AdminVariantAddOn, error)
	DeleteVariantAddOn(ctx context.Context, productID, variantID, addOnID int64) error

	// Public APIs
	PublicListProducts(ctx context.Context, pagination *utils.Pagination, q string) ([]AdminProductListItem, int64, error)
	PublicGetProduct(ctx context.Context, id int64) (*AdminProductDetail, error)
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

func (s *service) CreateProduct(ctx context.Context, req requests.AdminProductRequest) (*AdminProductDetail, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("price must be greater than or equal to zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	productID, err := s.repo.CreateProduct(ctx, tx, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(req.Attributes) > 0 {
		if err := s.repo.ReplaceProductAttributes(ctx, tx, productID, req.Attributes); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.repo.GetAdminProductByID(ctx, productID)
}

func (s *service) UpdateProduct(ctx context.Context, id int64, req requests.AdminProductRequest) (*AdminProductDetail, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("price must be greater than or equal to zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateProduct(ctx, tx, id, req); err != nil {
		tx.Rollback()
		return nil, err
	}

	if req.Attributes != nil {
		if err := s.repo.ReplaceProductAttributes(ctx, tx, id, req.Attributes); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.repo.GetAdminProductByID(ctx, id)
}

func (s *service) DeleteProduct(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := s.repo.SoftDeleteProduct(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *service) ListProducts(ctx context.Context, pagination *utils.Pagination) ([]AdminProductListItem, int64, error) {
	return s.repo.ListAdminProducts(ctx, pagination)
}

func (s *service) GetProductByID(ctx context.Context, id int64) (*AdminProductDetail, error) {
	return s.repo.GetAdminProductByID(ctx, id)
}

// PublicListProducts lists only active, non-deleted products for public APIs.
func (s *service) PublicListProducts(ctx context.Context, pagination *utils.Pagination, q string) ([]AdminProductListItem, int64, error) {
	// For now, reuse admin listing and filter via WHERE is_active = TRUE in the repository.
	return s.repo.ListPublicProducts(ctx, pagination, q)
}

// PublicGetProduct returns a product and its active variants for public APIs.
func (s *service) PublicGetProduct(ctx context.Context, id int64) (*AdminProductDetail, error) {
	product, err := s.repo.GetAdminProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !product.IsActive {
		return nil, sql.ErrNoRows
	}

	// Filter out inactive variants for public response
	activeVariants := make([]AdminVariant, 0, len(product.Variants))
	for _, v := range product.Variants {
		if v.IsActive {
			activeVariants = append(activeVariants, v)
		}
	}
	product.Variants = activeVariants
	return product, nil
}

func (s *service) CreateVariant(ctx context.Context, productID int64, req requests.AdminVariantRequest) (*AdminVariant, error) {
	if req.SKU == "" {
		return nil, fmt.Errorf("sku is required")
	}
	if len(req.AttributeValueIDs) == 0 {
		return nil, fmt.Errorf("attribute_value_ids is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	variantID, err := s.repo.CreateVariant(ctx, tx, productID, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.repo.GetAdminVariantByID(ctx, productID, variantID)
}

func (s *service) ListVariantAddOns(ctx context.Context, productID, variantID int64) ([]AdminVariantAddOn, error) {
	return s.repo.ListVariantAddOns(ctx, productID, variantID)
}

func (s *service) CreateVariantAddOn(ctx context.Context, productID, variantID, addOnProductID int64) (*AdminVariantAddOn, error) {
	if addOnProductID <= 0 {
		return nil, fmt.Errorf("add_on_product_id is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	addOnID, err := s.repo.CreateVariantAddOn(ctx, tx, productID, variantID, addOnProductID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	addOns, err := s.repo.ListVariantAddOns(ctx, productID, variantID)
	if err != nil {
		return nil, err
	}
	for _, ao := range addOns {
		if ao.ID == addOnID {
			return &ao, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (s *service) UpdateVariantAddOn(ctx context.Context, productID, variantID, addOnID, addOnProductID int64) (*AdminVariantAddOn, error) {
	if addOnProductID <= 0 {
		return nil, fmt.Errorf("add_on_product_id is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateVariantAddOn(ctx, tx, productID, variantID, addOnID, addOnProductID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	addOns, err := s.repo.ListVariantAddOns(ctx, productID, variantID)
	if err != nil {
		return nil, err
	}
	for _, ao := range addOns {
		if ao.ID == addOnID {
			return &ao, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (s *service) DeleteVariantAddOn(ctx context.Context, productID, variantID, addOnID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteVariantAddOn(ctx, tx, productID, variantID, addOnID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateVariant(ctx context.Context, productID, variantID int64, req requests.AdminVariantRequest) (*AdminVariant, error) {
	if req.SKU == "" {
		return nil, fmt.Errorf("sku is required")
	}
	if len(req.AttributeValueIDs) == 0 {
		return nil, fmt.Errorf("attribute_value_ids is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateVariant(ctx, tx, productID, variantID, req); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.repo.GetAdminVariantByID(ctx, productID, variantID)
}
