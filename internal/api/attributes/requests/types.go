package requests

type AttributeRequest struct {
	Name            string                  `json:"name" binding:"required"`
	AttributeValues []AttributeValueRequest `binding:"required,dive" json:"attribute_values"`
}

type AttributeValueRequest struct {
	Value     string `json:"value" binding:"required"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active" binding:"required"`
}


type AttributeValuesResponse struct {
	ID          int64  `json:"id"`
	AttributeID int64  `json:"attribute_id"`
	Value       string `json:"value"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// AttributeDetail represents an attribute with its values.
type AttributeResponse struct {
	ID       int64            `json:"id"`
	Name     string           `json:"name"`
	IsActive bool             `json:"is_active"`
	Values   []AttributeValuesResponse `json:"values"`
}

type AttributeListItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}
