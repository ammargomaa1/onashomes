package models

import "time"

type ProductSEO struct {
	ID                int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	ProductID         int64     `gorm:"type:bigint;uniqueIndex;not null" json:"product_id"`
	MetaTitleEn       string    `gorm:"type:varchar(70)" json:"meta_title_en"`
	MetaTitleAr       string    `gorm:"type:varchar(70)" json:"meta_title_ar"`
	MetaDescriptionEn string    `gorm:"type:varchar(170)" json:"meta_description_en"`
	MetaDescriptionAr string    `gorm:"type:varchar(170)" json:"meta_description_ar"`
	MetaKeywords      string    `gorm:"type:text" json:"meta_keywords"`
	CanonicalURL      string    `gorm:"type:varchar(500)" json:"canonical_url"`
	OgTitle           string    `gorm:"type:varchar(200)" json:"og_title"`
	OgDescription     string    `gorm:"type:varchar(300)" json:"og_description"`
	OgImage           string    `gorm:"type:varchar(500)" json:"og_image"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (ProductSEO) TableName() string { return "product_seo" }
