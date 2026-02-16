package inventory

import (
	"time"

	"github.com/onas/ecommerce-api/internal/api/inventory/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// VariantInventoryItem is the DTO for listing inventory
type VariantInventoryItem struct {
	ID                int64    `json:"id"`
	ProductVariantID  int64    `json:"product_variant_id"`
	StoreFrontID      int64    `json:"store_front_id"`
	Quantity          int      `json:"quantity"`
	ReservedQuantity  int      `json:"reserved_quantity"`
	AvailableQuantity int      `json:"available_quantity"`
	LowStockThreshold int      `json:"low_stock_threshold"`
	IsLowStock        bool     `json:"is_low_stock"`
	SKU               string   `json:"sku"`
	ProductName       string   `json:"product_name"`
	AttributeValue    string   `json:"attribute_value"`
	Price             *float64 `json:"price"`
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetVariantInventory(variantID, storeFrontID int64) (*models.VariantInventory, error) {
	var inv models.VariantInventory
	err := r.db.Where("product_variant_id = ? AND store_front_id = ?", variantID, storeFrontID).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *Repository) EnsureInventoryRecord(tx *gorm.DB, variantID, storeFrontID int64) (*models.VariantInventory, error) {
	var inv models.VariantInventory
	err := tx.Where("product_variant_id = ? AND store_front_id = ?", variantID, storeFrontID).First(&inv).Error
	if err == nil {
		return &inv, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new inventory record
	inv = models.VariantInventory{
		ProductVariantID:  variantID,
		StoreFrontID:      storeFrontID,
		Quantity:          0,
		ReservedQuantity:  0,
		LowStockThreshold: 5,
	}
	if err := tx.Create(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *Repository) AdjustInventory(tx *gorm.DB, inventoryID int64, newQuantity int) error {
	return tx.Model(&models.VariantInventory{}).
		Where("id = ?", inventoryID).
		Updates(map[string]interface{}{
			"quantity":   newQuantity,
			"updated_at": time.Now(),
		}).Error
}

func (r *Repository) UpdateStock(tx *gorm.DB, inventoryID int64, quantity, reservedQuantity int) error {
	return tx.Model(&models.VariantInventory{}).
		Where("id = ?", inventoryID).
		Updates(map[string]interface{}{
			"quantity":          quantity,
			"reserved_quantity": reservedQuantity,
			"updated_at":        time.Now(),
		}).Error
}

func (r *Repository) LockInventory(tx *gorm.DB, inventoryID int64) (*models.VariantInventory, error) {
	var inv models.VariantInventory
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", inventoryID).
		First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *Repository) CreateAdjustment(tx *gorm.DB, adj *models.InventoryAdjustment) error {
	return tx.Create(adj).Error
}

func (r *Repository) ListInventory(filter requests.InventoryFilterRequest, pagination *utils.Pagination) ([]VariantInventoryItem, int64, error) {
	var total int64

	// Base query
	query := r.db.Table("variant_inventory vi").
		Joins("JOIN product_variants pv ON pv.id = vi.product_variant_id").
		Joins("JOIN products p ON p.id = pv.product_id")

	if filter.StoreFrontID > 0 {
		query = query.Where("vi.store_front_id = ?", filter.StoreFrontID)
	}

	query = query.Where("pv.deleted_at IS NULL").
		Where("p.deleted_at IS NULL")

	// Apply filters
	if filter.BrandID != nil {
		query = query.Where("p.brand_id = ?", *filter.BrandID)
	}
	if filter.CategoryID != nil {
		query = query.Where("p.category_id = ?", *filter.CategoryID)
	}
	if filter.SupplierID != nil {
		query = query.Where("p.supplier_id = ?", *filter.SupplierID)
	}
	if filter.ProductID != nil {
		query = query.Where("pv.product_id = ?", *filter.ProductID)
	}
	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query = query.Where("(vi.quantity - vi.reserved_quantity) <= vi.low_stock_threshold")
	}
	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		query = query.Where("(p.name_en ILIKE ? OR p.name_ar ILIKE ? OR pv.sku ILIKE ?)", search, search, search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Select and Paging
	var items []VariantInventoryItem
	offset := (pagination.Page - 1) * pagination.Limit

	err := query.Select(`
			vi.id, vi.product_variant_id, vi.store_front_id,
			vi.quantity, vi.reserved_quantity, vi.low_stock_threshold,
			(vi.quantity - vi.reserved_quantity) as available_quantity,
			CASE WHEN (vi.quantity - vi.reserved_quantity) <= vi.low_stock_threshold THEN true ELSE false END as is_low_stock,
			pv.sku, p.name_en as product_name, pv.attribute_value, pv.price
		`).
		Order("vi.quantity ASC").
		Offset(offset).Limit(pagination.Limit).
		Scan(&items).Error

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *Repository) ListAdjustmentHistory(inventoryID int64, pagination *utils.Pagination) ([]models.InventoryAdjustment, int64, error) {
	var total int64
	if err := r.db.Model(&models.InventoryAdjustment{}).Where("variant_inventory_id = ?", inventoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []models.InventoryAdjustment
	offset := (pagination.Page - 1) * pagination.Limit
	err := r.db.Where("variant_inventory_id = ?", inventoryID).
		Order("created_at DESC").
		Offset(offset).Limit(pagination.Limit).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// HasInventoryForProduct checks if any variant of a product has inventory in any store
func (r *Repository) HasInventoryForProduct(productID int64) (bool, error) {
	var count int64
	err := r.db.Table("variant_inventory vi").
		Joins("JOIN product_variants pv ON pv.id = vi.product_variant_id").
		Where("pv.product_id = ? AND vi.quantity > 0 AND pv.deleted_at IS NULL", productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
