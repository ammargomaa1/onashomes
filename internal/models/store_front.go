package models

import "time"

type StoreFront struct {
	ID              int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	Name            string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug            string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Domain          string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"domain"`
	Currency        string    `gorm:"type:varchar(10);not null;default:'SAR'" json:"currency"`
	DefaultLanguage string    `gorm:"type:varchar(10);not null;default:'ar'" json:"default_language"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (StoreFront) TableName() string { return "store_fronts" }
