package requests

type AttributeRequest struct {
	NameAr          string                  `json:"name_ar" binding:"required"`
	NameEn          string                  `json:"name_en" binding:"required"`
	AttributeValues []AttributeValueRequest `binding:"required,dive" json:"attribute_values"`
}

type AttributeValueRequest struct {
	ValueAr     string `json:"value_ar" binding:"required"`
	ValueEn     string `json:"value_en" binding:"required"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active" binding:"required"`
}

type AttributeValuesResponse struct {
	ID          int64  `json:"id"`
	AttributeID int64  `json:"attribute_id"`
	ValueAr       string `json:"value_ar"`
	ValueEn       string `json:"value_en"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// AttributeDetail represents an attribute with its values.
type AttributeResponse struct {
	ID       int64                     `json:"id"`
	NameAr     string                    `json:"name_ar"`
	NameEn   string                    `json:"name_en"`
	Values   []AttributeValuesResponse `json:"attribute_values"`
}

type AttributeListItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}
