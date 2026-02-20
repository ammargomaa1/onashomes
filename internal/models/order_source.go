package models

import "time"

type OrderSource struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	NameEn    string    `gorm:"type:varchar(100);not null" json:"name_en"`
	NameAr    string    `gorm:"type:varchar(100);not null" json:"name_ar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
