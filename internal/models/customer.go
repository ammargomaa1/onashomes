package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID           int64          `json:"id" gorm:"primaryKey"`
	UserID       *int64         `json:"user_id,omitempty" gorm:"index"` // Nullable, links to auth user if they have an account
	StoreFrontID int64          `json:"store_front_id" gorm:"index;not null"`
	StoreFront   StoreFront     `json:"store_front,omitempty" gorm:"foreignKey:StoreFrontID"`
	FirstName    string         `json:"first_name" gorm:"size:255;not null"`
	LastName     string         `json:"last_name" gorm:"size:255;not null"`
	Email        string         `json:"email" gorm:"size:255;index"`
	Phone        string         `json:"phone" gorm:"size:50;index;not null"` // Phone is the primary identifier for guest/admin-created customers
	CreatedAt    time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relations
	Orders []Order `json:"orders,omitempty" gorm:"foreignKey:CustomerID"`
}
