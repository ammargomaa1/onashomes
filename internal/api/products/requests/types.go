package requests

type AdminProductAttributeRequest struct {
	AttributeID     int64   `json:"attribute_id" binding:"required"`
	SortOrder       int     `json:"sort_order"`
	AllowedValueIDs []int64 `json:"allowed_value_ids"`
}

type AdminProductRequest struct {
	Name        string                         `json:"name" binding:"required"`
	Description string                         `json:"description"`
	Price       float64                        `json:"price" binding:"required"`
	IsActive    bool                           `json:"is_active"`
	Attributes  []AdminProductAttributeRequest `json:"attributes"`
}

type AdminVariantRequest struct {
	SKU               string  `json:"sku" binding:"required"`
	IsActive          bool    `json:"is_active"`
	AttributeValueIDs []int64 `json:"attribute_value_ids" binding:"required"`
	ImageFileIDs      []int64 `json:"image_file_ids"`
	AddOnProductIDs   []int64 `json:"add_on_product_ids"`
}

type VariantAddOnRequest struct {
	AddOnProductID int64 `json:"add_on_product_id" binding:"required"`
}
