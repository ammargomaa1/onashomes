package main

import (
	"log"

	"github.com/onas/ecommerce-api/config"
	"github.com/onas/ecommerce-api/internal/database"
	"github.com/onas/ecommerce-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	config.LoadConfig()

	db := database.GetDB()

	// Auto Migrate
	db.AutoMigrate(
		&models.Currency{},
		&models.OrderStatus{},
		&models.PaymentStatus{},
		&models.FulfillmentStatus{},
		&models.Order{}, // Also migrate Order to update columns
	)

	log.Println("Seeding Database...")

	seedCurrencies(db)
	seedOrderStatuses(db)
	seedPaymentStatuses(db)
	seedFulfillmentStatuses(db)

	log.Println("Database Seeded Successfully!")
}

func seedCurrencies(db *gorm.DB) {
	currencies := []models.Currency{
		{
			Code:   "SAR",
			NameEn: "Saudi Riyal",
			NameAr: "ريال سعودي",
			Symbol: "SAR",
		},
		{
			Code:   "EGP",
			NameEn: "Egyptian Pound",
			NameAr: "جنيه مصري",
			Symbol: "EGP",
		},
	}

	for _, currency := range currencies {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name_en", "name_ar", "symbol"}),
		}).Create(&currency).Error; err != nil {
			log.Printf("Failed to seed currency %s: %v", currency.Code, err)
		}
	}
}

func seedOrderStatuses(db *gorm.DB) {
	statuses := []models.OrderStatus{
		{Slug: "draft", NameEn: "Draft", NameAr: "مسودة"},
		{Slug: "pending_payment", NameEn: "Pending Payment", NameAr: "بانتظار الدفع"},
		{Slug: "paid", NameEn: "Paid", NameAr: "مدفوع"},
		{Slug: "confirmed", NameEn: "Confirmed", NameAr: "مؤكد"},
		{Slug: "fulfilled", NameEn: "Fulfilled", NameAr: "تم التوصيل"},
		{Slug: "completed", NameEn: "Completed", NameAr: "مكتمل"},
		{Slug: "cancelled", NameEn: "Cancelled", NameAr: "ملغي"},
	}

	for _, status := range statuses {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "slug"}},
			DoUpdates: clause.AssignmentColumns([]string{"name_en", "name_ar"}),
		}).Create(&status).Error; err != nil {
			log.Printf("Failed to seed order status %s: %v", status.Slug, err)
		}
	}
}

func seedPaymentStatuses(db *gorm.DB) {
	statuses := []models.PaymentStatus{
		{Slug: "unpaid", NameEn: "Unpaid", NameAr: "غير مدفوع"},
		{Slug: "pending", NameEn: "Pending", NameAr: "قيد الانتظار"},
		{Slug: "paid", NameEn: "Paid", NameAr: "مدفوع"},
		{Slug: "refunded", NameEn: "Refunded", NameAr: "مرتجع"},
		{Slug: "failed", NameEn: "Failed", NameAr: "فشل الدفع"},
	}

	for _, status := range statuses {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "slug"}},
			DoUpdates: clause.AssignmentColumns([]string{"name_en", "name_ar"}),
		}).Create(&status).Error; err != nil {
			log.Printf("Failed to seed payment status %s: %v", status.Slug, err)
		}
	}
}

func seedFulfillmentStatuses(db *gorm.DB) {
	statuses := []models.FulfillmentStatus{
		{Slug: "unfulfilled", NameEn: "Unfulfilled", NameAr: "غير منفذ"},
		{Slug: "partially_fulfilled", NameEn: "Partially Fulfilled", NameAr: "منفذ جزئياً"},
		{Slug: "fulfilled", NameEn: "Fulfilled", NameAr: "منفذ"},
	}

	for _, status := range statuses {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "slug"}},
			DoUpdates: clause.AssignmentColumns([]string{"name_en", "name_ar"}),
		}).Create(&status).Error; err != nil {
			log.Printf("Failed to seed fulfillment status %s: %v", status.Slug, err)
		}
	}
}
