package inventory

import (
	"fmt"

	"github.com/onas/ecommerce-api/internal/api/inventory/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *Repository
}

func NewService(db *gorm.DB, repo *Repository) *Service {
	return &Service{db: db, repo: repo}
}

func (s *Service) AdjustInventory(req requests.AdjustInventoryRequest, adminID int64) utils.IResource {
	var result *models.VariantInventory
	var adjustment *models.InventoryAdjustment

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Ensure inventory record exists
		invRepo := &Repository{db: tx}
		inv, err := invRepo.EnsureInventoryRecord(tx, req.ProductVariantID, req.StoreFrontID)
		if err != nil {
			return fmt.Errorf("failed to ensure inventory record: %w", err)
		}

		// Lock the row for update
		locked, err := invRepo.LockInventory(tx, inv.ID)
		if err != nil {
			return fmt.Errorf("failed to lock inventory: %w", err)
		}

		// Validate no negative stock
		newQty := locked.Quantity + req.Adjustment
		if newQty < 0 {
			return fmt.Errorf("adjustment would result in negative stock (current: %d, adjustment: %d)", locked.Quantity, req.Adjustment)
		}

		// Update quantity
		if err := invRepo.AdjustInventory(tx, locked.ID, newQty); err != nil {
			return fmt.Errorf("failed to update inventory: %w", err)
		}

		// Create audit record
		adj := &models.InventoryAdjustment{
			VariantInventoryID: locked.ID,
			AdjustedBy:         adminID,
			PreviousQuantity:   locked.Quantity,
			NewQuantity:        newQty,
			AdjustmentAmount:   req.Adjustment,
			Reason:             req.Reason,
			Notes:              req.Notes,
		}
		if err := invRepo.CreateAdjustment(tx, adj); err != nil {
			return fmt.Errorf("failed to create adjustment record: %w", err)
		}

		adjustment = adj
		locked.Quantity = newQty
		result = locked

		return nil
	})

	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Inventory adjusted successfully", map[string]interface{}{
		"id":                  result.ID,
		"product_variant_id":  result.ProductVariantID,
		"store_front_id":      result.StoreFrontID,
		"quantity":            result.Quantity,
		"reserved_quantity":   result.ReservedQuantity,
		"available_quantity":  result.AvailableQuantity(),
		"low_stock_threshold": result.LowStockThreshold,
		"is_low_stock":        result.IsLowStock(),
		"adjustment": map[string]interface{}{
			"previous_quantity": adjustment.PreviousQuantity,
			"new_quantity":      adjustment.NewQuantity,
			"adjustment_amount": adjustment.AdjustmentAmount,
			"reason":            adjustment.Reason,
		},
	})
}

func (s *Service) GetVariantInventory(variantID, storeFrontID int64) utils.IResource {
	inv, err := s.repo.GetVariantInventory(variantID, storeFrontID)
	if err != nil {
		return utils.NewNotFoundResource("Inventory not found", nil)
	}

	return utils.NewOKResource("Inventory retrieved successfully", map[string]interface{}{
		"id":                  inv.ID,
		"product_variant_id":  inv.ProductVariantID,
		"store_front_id":      inv.StoreFrontID,
		"quantity":            inv.Quantity,
		"reserved_quantity":   inv.ReservedQuantity,
		"available_quantity":  inv.AvailableQuantity(),
		"low_stock_threshold": inv.LowStockThreshold,
		"is_low_stock":        inv.IsLowStock(),
	})
}

func (s *Service) ListInventoryByStore(storeFrontID int64, pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.ListInventoryByStore(storeFrontID, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve inventory", err)
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Inventory retrieved successfully", items, pagination.GetMeta())
}

func (s *Service) GetLowStockAlerts(storeFrontID int64, pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.GetLowStockItems(storeFrontID, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve low stock items", err)
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Low stock items retrieved successfully", items, pagination.GetMeta())
}

func (s *Service) GetAdjustmentHistory(inventoryID int64, pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.ListAdjustmentHistory(inventoryID, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve adjustment history", err)
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Adjustment history retrieved successfully", items, pagination.GetMeta())
}
