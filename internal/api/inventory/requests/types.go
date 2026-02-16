package requests

type AdjustInventoryRequest struct {
	ProductVariantID int64  `json:"product_variant_id" binding:"required"`
	StoreFrontID     int64  `json:"store_front_id" binding:"required"`
	Adjustment       int    `json:"adjustment" binding:"required"` // Can be negative
	Reason           string `json:"reason" binding:"required,oneof=restock correction sale return transfer"`
	Notes            string `json:"notes"`
}

type BulkInventoryUpdateItem struct {
	VariantInventoryID int64  `json:"variant_inventory_id" binding:"required"`
	NewQuantity        int    `json:"new_quantity" binding:"required,min=0"`
	Reason             string `json:"reason" binding:"required,oneof=restock correction sale return transfer"`
	Notes              string `json:"notes"`
}

type BulkInventoryUpdateRequest struct {
	Items []BulkInventoryUpdateItem `json:"items" binding:"required,dive"`
}

type InventoryFilterRequest struct {
	StoreFrontID int64   `form:"store_front_id"`
	BrandID      *int64  `form:"brand_id"`
	CategoryID   *int64  `form:"category_id"`
	ProductID    *int64  `form:"product_id"`
	SupplierID   *int64  `form:"supplier_id"`
	LowStockOnly *bool   `form:"low_stock_only"`
	Search       *string `form:"search"`
}

type ReserveStockRequest struct {
	VariantID int64 `json:"variant_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}
