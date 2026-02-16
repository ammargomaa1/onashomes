package requests

type UpdateOrderRequest struct {
	StoreFrontID   int64             `json:"store_front_id"`
	CustomerName   string            `json:"customer_name"`
	CustomerEmail  string            `json:"customer_email"`
	CustomerPhone  string            `json:"customer_phone"`
	Notes          string            `json:"notes"`
	Items          []OrderItemUpdate `json:"items"`
	ShippingAmount float64           `json:"shipping_amount"`
	TaxAmount      float64           `json:"tax_amount"`
	DiscountAmount float64           `json:"discount_amount"`
}

type OrderItemUpdate struct {
	ID               int64 `json:"id"` // 0 for new items
	ProductVariantID int64 `json:"product_variant_id"`
	Quantity         int   `json:"quantity"`
	IsRemoved        bool  `json:"is_removed"` // Flag to mark for deletion
}
