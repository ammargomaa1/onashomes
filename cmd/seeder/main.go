package main

import (
	"log"

	"github.com/onas/ecommerce-api/config"
	"github.com/onas/ecommerce-api/internal/database"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
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
		&models.Role{},
		&models.Admin{},
		// Location Models
		&models.Country{},
		&models.Governorate{},
		&models.City{},
		&models.OrderAddress{},
	)

	log.Println("Seeding Database...")

	seedCurrencies(db)
	seedOrderStatuses(db)
	seedPaymentStatuses(db)
	seedFulfillmentStatuses(db)
	seedRoles(db)
	seedAdmin(db)
	seedLocations(db)

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

func seedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "Super Admin", Description: "Full access to all resources"},
		{Name: "Admin", Description: "Standard admin access"},
	}

	for _, role := range roles {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"description"}),
		}).Create(&role).Error; err != nil {
			log.Printf("Failed to seed role %s: %v", role.Name, err)
		}
	}
}

func seedAdmin(db *gorm.DB) {
	var role models.Role
	if err := db.Where("name = ?", "Super Admin").First(&role).Error; err != nil {
		log.Printf("Failed to find Super Admin role: %v", err)
		return
	}

	hp, _ := utils.HashPassword("admin123")
	admin := models.Admin{
		Email:     "admin@onashomes.com",
		Password:  hp,
		FirstName: "Super",
		LastName:  "Admin",
		RoleID:    role.ID,
		IsActive:  true,
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"password", "first_name", "last_name", "role_id", "is_active"}),
	}).Create(&admin).Error; err != nil {
		log.Printf("Failed to seed admin user: %v", err)
	}
}
