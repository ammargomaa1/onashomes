package requests

type AttributeRequest struct {
	Name            string                  `json:"name" binding:"required"`
	AttributeValues []AttributeValueRequest `json:"attribute_values"`
}

type AttributeValueRequest struct {
	Value     string `json:"value" binding:"required"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}
