package models

type PaymentMethod struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	NameEn   string `json:"name_en" gorm:"size:255;not null"`
	NameAr   string `json:"name_ar" gorm:"size:255;not null"`
	Slug     string `json:"slug" gorm:"size:255;not null;uniqueIndex"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}
