package main

import (
	"log"

	"github.com/onas/ecommerce-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedLocations(db *gorm.DB) {
	// 1. Seed Country (Egypt)
	egypt := models.Country{
		Code:   "EG",
		NameEn: "Egypt",
		NameAr: "مصر",
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name_en", "name_ar"}),
	}).FirstOrCreate(&egypt).Error; err != nil {
		log.Printf("Failed to seed country Egypt: %v", err)
		return
	}

	// 2. Seed Governorate (Damietta)
	damiettaGov := models.Governorate{
		CountryID: egypt.ID,
		NameEn:    "Damietta",
		NameAr:    "دمياط",
	}

	// Use NameEn + CountryID as unique constraint logic manually or simplified lookup
	var existingGov models.Governorate
	if err := db.Where("country_id = ? AND name_en = ?", egypt.ID, "Damietta").First(&existingGov).Error; err == nil {
		damiettaGov = existingGov
	} else {
		if err := db.Create(&damiettaGov).Error; err != nil {
			log.Printf("Failed to seed governorate Damietta: %v", err)
			return
		}
	}

	// 3. Seed Cities for Damietta
	cities := []models.City{
		{NameEn: "Damietta", NameAr: "دمياط"},
		{NameEn: "New Damietta", NameAr: "دمياط الجديدة"},
		{NameEn: "Ras El Bar", NameAr: "رأس البر"},
		{NameEn: "Faraskour", NameAr: "فارسكور"},
		{NameEn: "Kafr Saad", NameAr: "كفر سعد"},
		{NameEn: "Kafr El-Battikh", NameAr: "كفر البطيخ"},
		{NameEn: "Ezbet El-Borg", NameAr: "عزبة البرج"},
		{NameEn: "El Zarqa", NameAr: "الزرقا"},
		{NameEn: "El Roda", NameAr: "الروضة"},
		{NameEn: "El Serw", NameAr: "السرو"},
		{NameEn: "Meet Abou Ghalib", NameAr: "ميت أبو غالب"},
	}

	for _, city := range cities {
		city.GovernorateID = damiettaGov.ID
		// Check existence by name and gov ID
		var existingCity models.City
		if err := db.Where("governorate_id = ? AND name_en = ?", damiettaGov.ID, city.NameEn).First(&existingCity).Error; err == nil {
			continue // Already exists
		}

		if err := db.Create(&city).Error; err != nil {
			log.Printf("Failed to seed city %s: %v", city.NameEn, err)
		}
	}
}
