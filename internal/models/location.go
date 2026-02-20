package models

import "time"

// Location Models

type Country struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	NameEn    string    `json:"name_en" gorm:"size:255;not null"`
	NameAr    string    `json:"name_ar" gorm:"size:255;not null"`
	Code      string    `json:"code" gorm:"size:3;not null;uniqueIndex"` // ISO Code
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Governorate struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	CountryID int64     `json:"country_id" gorm:"index;not null"`
	NameEn    string    `json:"name_en" gorm:"size:255;not null"`
	NameAr    string    `json:"name_ar" gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	Country *Country `json:"country,omitempty" gorm:"foreignKey:CountryID"`
}

type City struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	GovernorateID int64     `json:"governorate_id" gorm:"index;not null"`
	NameEn        string    `json:"name_en" gorm:"size:255;not null"`
	NameAr        string    `json:"name_ar" gorm:"size:255;not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Associations
	Governorate *Governorate `json:"governorate,omitempty" gorm:"foreignKey:GovernorateID"`
}
