package requests

type UpdateOrderRequest struct {
	StoreFrontID    int64             `json:"store_front_id"`
	CustomerName    string            `json:"customer_name"`
	CustomerEmail   string            `json:"customer_email"`
	CustomerPhone   string            `json:"customer_phone"`
	Notes           string            `json:"notes"`
	Items           []OrderItemUpdate `json:"items"`
	ShippingAmount  float64           `json:"shipping_amount" binding:"min=0"`
	TaxAmount       float64           `json:"tax_amount" binding:"min=0"`
	DiscountAmount  float64           `json:"discount_amount" binding:"min=0"`
	OrderSourceID   *int64            `json:"order_source_id"`
	PaymentMethodID *int64            `json:"payment_method_id"`

	// Address
	CountryID      *int64 `json:"country_id"`
	GovernorateID  *int64 `json:"governorate_id"`
	CityID         *int64 `json:"city_id"`
	Street         string `json:"street"`
	BuildingNumber string `json:"building_number"`
	Floor          string `json:"floor"`
	Apartment      string `json:"apartment"`
	SpecialMark    string `json:"special_mark"`
}

type OrderItemUpdate struct {
	ID               int64 `json:"id"` // 0 for new items
	ProductVariantID int64 `json:"product_variant_id"`
	Quantity         int   `json:"quantity"`
	IsRemoved        bool  `json:"is_removed"` // Flag to mark for deletion
}
