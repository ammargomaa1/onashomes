package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category in the system
type Category struct {
	ID        int64          `json:"id"`
	SectionID int64          `json:"section_id"`
	NameAr    string         `json:"name_ar"`
	NameEn    string         `json:"name_en"`
	IconID    int64          `json:"icon_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}