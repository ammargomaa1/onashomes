package requests

type CreateBrandRequest struct {
	NameAr string `json:"name_ar" binding:"required"`
	NameEn string `json:"name_en" binding:"required"`
	LogoID int64  `json:"logo_id" binding:"required"`
	IconID int64  `json:"icon_id" binding:"required"`
}