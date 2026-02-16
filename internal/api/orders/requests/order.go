package requests

type CreateOrderItemRequest struct {
	ProductVariantID int64 `json:"product_variant_id" binding:"required"`
	Quantity         int   `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	StoreFrontID int64                    `json:"store_front_id" binding:"required"`
	Items        []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`

	// New Customer Fields
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`

	ShippingAmount float64 `json:"shipping_amount" binding:"min=0"`
	TaxAmount      float64 `json:"tax_amount" binding:"min=0"`
	DiscountAmount float64 `json:"discount_amount" binding:"min=0"`
	Notes          string  `json:"notes"`
}

type OrderFilterRequest struct {
	StoreFrontID  int64  `form:"store_front_id"`
	Status        string `form:"status"`         // Order Status Slug
	PaymentStatus string `form:"payment_status"` // Payment Status Slug
	ProductIDs    string `form:"product_ids"`    // Comma separated product IDs
	DateFrom      string `form:"date_from"`      // YYYY-MM-DD
	DateTo        string `form:"date_to"`        // YYYY-MM-DD
	Search        string `form:"search"`
}
