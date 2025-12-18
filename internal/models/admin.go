package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        int64          `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	FirstName string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string         `gorm:"type:varchar(100)" json:"last_name"`
	RoleID    int64          `gorm:"type:bigint;not null" json:"role_id"`
	Role      Role           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

