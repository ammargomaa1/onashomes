package requests

type AdjustInventoryRequest struct {
	ProductVariantID int64  `json:"product_variant_id" binding:"required"`
	StoreFrontID     int64  `json:"store_front_id" binding:"required"`
	Adjustment       int    `json:"adjustment" binding:"required"`
	Reason           string `json:"reason" binding:"required,oneof=adjustment correction restock"`
	Notes            string `json:"notes"`
}
