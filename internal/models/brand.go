package models

import (
	"time"

	"gorm.io/gorm"
)

// Brand represents a product brand in the system
type Brand struct {
	ID        int64          `json:"id"`
	NameAr    string         `json:"name_ar"`
	NameEn    string         `json:"name_en"`
	LogoID    int64          `json:"logo_id"`
	IconID    int64          `json:"icon_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}
