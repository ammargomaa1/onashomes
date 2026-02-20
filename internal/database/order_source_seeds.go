package database

import (
	"fmt"

	"github.com/onas/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

func SeedOrderSources(db *gorm.DB) error {
	fmt.Println("Seeding order sources...")

	sources := []models.OrderSource{
		{ID: 1, NameEn: "Phone Call", NameAr: "مكالمة هاتفية"},
		{ID: 2, NameEn: "Facebook", NameAr: "فيسبوك"},
		{ID: 3, NameEn: "Instagram", NameAr: "انستجرام"},
		{ID: 4, NameEn: "Whatsapp", NameAr: "واتساب"},
		{ID: 5, NameEn: "Nada", NameAr: "ندى"},
	}

	for _, source := range sources {
		if err := db.FirstOrCreate(&source, models.OrderSource{ID: source.ID}).Error; err != nil {
			return fmt.Errorf("failed to seed order source %s: %v", source.NameEn, err)
		}
	}

	fmt.Println("✓ Order sources seeded successfully")
	return nil
}
