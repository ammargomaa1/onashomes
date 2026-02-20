package models

import (
	"time"

	"gorm.io/gorm"
)

type Attribute struct {
	ID        int64          `gorm:"primaryKey"`
	NameAr    string         `gorm:"type:varchar(255);not null"`
	NameEn    string         `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
