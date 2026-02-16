package database

import (
	"fmt"
	"log"

	"github.com/onas/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

// SeedOrderStatuses seeds the initial order statuses and currencies
func SeedOrderStatuses(db *gorm.DB) error {
	fmt.Println("Seeding order statuses and currencies...")

	// 1. Order Statuses
	orderStatuses := []models.OrderStatus{
		{NameEn: "Draft", NameAr: "مسودة", Slug: "draft"},
		{NameEn: "Pending Payment", NameAr: "بانتظار الدفع", Slug: "pending_payment"},
		{NameEn: "Paid", NameAr: "مدفوع", Slug: "paid"},
		{NameEn: "Confirmed", NameAr: "مؤكد", Slug: "confirmed"},
		{NameEn: "Fulfilled", NameAr: "تم التجهيز", Slug: "fulfilled"},
		{NameEn: "Completed", NameAr: "مكتمل", Slug: "completed"},
		{NameEn: "Cancelled", NameAr: "ملغي", Slug: "cancelled"},
		{NameEn: "Returned", NameAr: "مرجع", Slug: "returned"},
		{NameEn: "Refunded", NameAr: "معاد المبلغ", Slug: "refunded"},
	}

	for _, status := range orderStatuses {
		if err := db.Where("slug = ?", status.Slug).FirstOrCreate(&status).Error; err != nil {
			return fmt.Errorf("failed to seed order status %s: %w", status.Slug, err)
		}
	}

	// 2. Payment Statuses
	paymentStatuses := []models.PaymentStatus{
		{NameEn: "Unpaid", NameAr: "غير مدفوع", Slug: "unpaid"},
		{NameEn: "Pending", NameAr: "قيد الانتظار", Slug: "pending"},
		{NameEn: "Paid", NameAr: "مدفوع", Slug: "paid"},
		{NameEn: "Refunded", NameAr: "معاد المبلغ", Slug: "refunded"},
		{NameEn: "Failed", NameAr: "فشل الدفع", Slug: "failed"},
	}

	for _, status := range paymentStatuses {
		if err := db.Where("slug = ?", status.Slug).FirstOrCreate(&status).Error; err != nil {
			return fmt.Errorf("failed to seed payment status %s: %w", status.Slug, err)
		}
	}

	// 3. Fulfillment Statuses
	fulfillmentStatuses := []models.FulfillmentStatus{
		{NameEn: "Unfulfilled", NameAr: "غير مجهز", Slug: "unfulfilled"},
		{NameEn: "Partially Fulfilled", NameAr: "مجهز جزئياً", Slug: "partially_fulfilled"},
		{NameEn: "Fulfilled", NameAr: "تم التجهيز", Slug: "fulfilled"},
	}

	for _, status := range fulfillmentStatuses {
		if err := db.Where("slug = ?", status.Slug).FirstOrCreate(&status).Error; err != nil {
			return fmt.Errorf("failed to seed fulfillment status %s: %w", status.Slug, err)
		}
	}

	// 4. Currencies
	currencies := []models.Currency{
		{NameEn: "Saudi Riyal", NameAr: "ريال سعودي", Code: "SAR", Symbol: "SAR"},
		{NameEn: "Egyptian Pound", NameAr: "جنيه مصري", Code: "EGP", Symbol: "EGP"},
		{NameEn: "US Dollar", NameAr: "دولار أمريكي", Code: "USD", Symbol: "$"},
	}

	for _, currency := range currencies {
		if err := db.Where("code = ?", currency.Code).FirstOrCreate(&currency).Error; err != nil {
			return fmt.Errorf("failed to seed currency %s: %w", currency.Code, err)
		}
	}

	log.Println("✓ Order statuses, payment statuses, fulfillment statuses, and currencies seeded successfully")
	return nil
}
