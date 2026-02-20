package utils

// DropdownItem is a lightweight struct for dropdown/select lists.
type DropdownItem struct {
	ID     int64  `json:"id"`
	NameEn string `json:"name_en"`
	NameAr string `json:"name_ar"`
}
