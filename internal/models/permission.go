package models

import (
	"time"
)

type Permission struct {
	ID          int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

