package models

import (
	"time"

	"gorm.io/gorm"
)

type Supplier struct {
	ID                int64          `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	CompanyName       string         `gorm:"type:varchar(255);not null" json:"company_name"`
	ContactPersonName string         `gorm:"type:varchar(255);not null" json:"contact_person_name"`
	ContactNumber     string         `gorm:"type:varchar(50);not null" json:"contact_number"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedBy         int64          `gorm:"type:bigint;not null" json:"created_by"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
