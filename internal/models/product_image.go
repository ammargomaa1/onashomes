package models

import "time"

// ProductImage links product-level images to uploaded files
type ProductImage struct {
	ID        int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	ProductID int64     `gorm:"type:bigint;not null;index" json:"product_id"`
	FileID    int64     `gorm:"type:bigint;not null" json:"file_id"`
	Position  int       `gorm:"type:int;default:0" json:"position"`
	IsCover   bool      `gorm:"default:false" json:"is_cover"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	File *File `gorm:"foreignKey:FileID" json:"file,omitempty"`
}

func (ProductImage) TableName() string { return "product_images" }
