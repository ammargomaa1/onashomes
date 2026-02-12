package requests

type CreateSectionRequest struct {
	NameAr string `json:"name_ar" binding:"required"`
	NameEn string `json:"name_en" binding:"required"`
	IconID int64  `json:"icon_id" binding:"required"`
}
