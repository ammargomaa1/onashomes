package models

import (
	"time"

	"gorm.io/gorm"
)

// Section represents a product section in the system
type Section struct {
	ID        int64          `json:"id"`
	NameAr    string         `json:"name_ar"`
	NameEn    string         `json:"name_en"`
	IconID    int64          `json:"icon_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Section
func (Section) TableName() string {
	return "sections"
}
