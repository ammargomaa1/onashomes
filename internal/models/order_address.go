package models

import "time"

type OrderAddress struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrderID        int64     `json:"order_id" gorm:"index;not null"`
	CountryID      int64     `json:"country_id" gorm:"index;not null"`
	GovernorateID  int64     `json:"governorate_id" gorm:"index;not null"`
	CityID         int64     `json:"city_id" gorm:"index;not null"`
	Street         string    `json:"street" gorm:"size:255;not null"`
	BuildingNumber string    `json:"building_number" gorm:"size:50"`
	Floor          string    `json:"floor" gorm:"size:50"`
	Apartment      string    `json:"apartment" gorm:"size:50"`
	SpecialMark    string    `json:"special_mark" gorm:"size:255"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Associations
	Country     *Country     `json:"country,omitempty" gorm:"foreignKey:CountryID"`
	Governorate *Governorate `json:"governorate,omitempty" gorm:"foreignKey:GovernorateID"`
	City        *City        `json:"city,omitempty" gorm:"foreignKey:CityID"`
}
