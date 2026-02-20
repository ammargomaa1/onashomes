package inventory

import (
	"testing"

	"github.com/onas/ecommerce-api/internal/api/inventory/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.StoreFront{},
		&models.ProductVariant{},
		&models.VariantInventory{},
		&models.InventoryAdjustment{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestAdjustInventory_Success(t *testing.T) {
	db := setupTestDB(t)

	// Seed a store and variant
	sf := models.StoreFront{Name: "Test", Slug: "test", Domain: "test.com", Currency: "SAR", DefaultLanguage: "ar", IsActive: true}
	db.Create(&sf)

	variant := models.ProductVariant{ProductID: 1, SKU: "SKU-001", IsActive: true}
	db.Create(&variant)

	repo := NewRepository(db)
	service := NewService(db, repo)

	req := requests.AdjustInventoryRequest{
		ProductVariantID: variant.ID,
		StoreFrontID:     sf.ID,
		Adjustment:       10,
		Reason:           "restock",
		Notes:            "Initial stock",
	}

	res := service.AdjustInventory(req, 1)
	if res.GetStatusCode() != 200 {
		t.Fatalf("expected status 200, got %d: %s", res.GetStatusCode(), res.GetMessage())
	}

	// Verify inventory record
	var inv models.VariantInventory
	db.Where("product_variant_id = ? AND store_front_id = ?", variant.ID, sf.ID).First(&inv)
	if inv.Quantity != 10 {
		t.Fatalf("expected quantity 10, got %d", inv.Quantity)
	}

	// Verify adjustment record
	var adj models.InventoryAdjustment
	db.Where("variant_inventory_id = ?", inv.ID).First(&adj)
	if adj.AdjustmentAmount != 10 {
		t.Fatalf("expected adjustment +10, got %d", adj.AdjustmentAmount)
	}
	if adj.PreviousQuantity != 0 {
		t.Fatalf("expected previous_quantity 0, got %d", adj.PreviousQuantity)
	}
	if adj.NewQuantity != 10 {
		t.Fatalf("expected new_quantity 10, got %d", adj.NewQuantity)
	}
}

func TestAdjustInventory_NegativeStockPrevented(t *testing.T) {
	db := setupTestDB(t)

	sf := models.StoreFront{Name: "Test", Slug: "test", Domain: "test.com", Currency: "SAR", DefaultLanguage: "ar", IsActive: true}
	db.Create(&sf)

	variant := models.ProductVariant{ProductID: 1, SKU: "SKU-002", IsActive: true}
	db.Create(&variant)

	// Pre-set inventory to 5
	db.Create(&models.VariantInventory{
		ProductVariantID: variant.ID, StoreFrontID: sf.ID,
		Quantity: 5, ReservedQuantity: 0, LowStockThreshold: 3,
	})

	repo := NewRepository(db)
	service := NewService(db, repo)

	// Try to remove 10 from stock of 5 → should fail
	req := requests.AdjustInventoryRequest{
		ProductVariantID: variant.ID,
		StoreFrontID:     sf.ID,
		Adjustment:       -10,
		Reason:           "adjustment",
		Notes:            "Should fail",
	}

	res := service.AdjustInventory(req, 1)
	if res.GetStatusCode() != 400 {
		t.Fatalf("expected status 400 for negative stock, got %d: %s", res.GetStatusCode(), res.GetMessage())
	}

	// Verify stock unchanged
	var inv models.VariantInventory
	db.Where("product_variant_id = ? AND store_front_id = ?", variant.ID, sf.ID).First(&inv)
	if inv.Quantity != 5 {
		t.Fatalf("expected quantity unchanged at 5, got %d", inv.Quantity)
	}
}

func TestAdjustInventory_SuccessiveAdjustments(t *testing.T) {
	db := setupTestDB(t)

	sf := models.StoreFront{Name: "Test", Slug: "test", Domain: "test.com", Currency: "SAR", DefaultLanguage: "ar", IsActive: true}
	db.Create(&sf)

	variant := models.ProductVariant{ProductID: 1, SKU: "SKU-003", IsActive: true}
	db.Create(&variant)

	repo := NewRepository(db)
	service := NewService(db, repo)

	// Add 10
	res := service.AdjustInventory(requests.AdjustInventoryRequest{
		ProductVariantID: variant.ID, StoreFrontID: sf.ID,
		Adjustment: 10, Reason: "restock",
	}, 1)
	if res.GetStatusCode() != 200 {
		t.Fatalf("first adjust failed: %s", res.GetMessage())
	}

	// Remove 3
	res = service.AdjustInventory(requests.AdjustInventoryRequest{
		ProductVariantID: variant.ID, StoreFrontID: sf.ID,
		Adjustment: -3, Reason: "adjustment",
	}, 1)
	if res.GetStatusCode() != 200 {
		t.Fatalf("second adjust failed: %s", res.GetMessage())
	}

	// Verify final quantity = 7
	var inv models.VariantInventory
	db.Where("product_variant_id = ? AND store_front_id = ?", variant.ID, sf.ID).First(&inv)
	if inv.Quantity != 7 {
		t.Fatalf("expected quantity 7, got %d", inv.Quantity)
	}

	// Verify 2 adjustment records
	var adjCount int64
	db.Model(&models.InventoryAdjustment{}).Where("variant_inventory_id = ?", inv.ID).Count(&adjCount)
	if adjCount != 2 {
		t.Fatalf("expected 2 adjustment records, got %d", adjCount)
	}
}

func TestVariantInventory_AvailableQuantity(t *testing.T) {
	inv := models.VariantInventory{Quantity: 20, ReservedQuantity: 5, LowStockThreshold: 10}

	if inv.AvailableQuantity() != 15 {
		t.Fatalf("expected available 15, got %d", inv.AvailableQuantity())
	}
}

func TestVariantInventory_IsLowStock(t *testing.T) {
	tests := []struct {
		name     string
		inv      models.VariantInventory
		expected bool
	}{
		{"above threshold", models.VariantInventory{Quantity: 20, ReservedQuantity: 0, LowStockThreshold: 5}, false},
		{"at threshold", models.VariantInventory{Quantity: 5, ReservedQuantity: 0, LowStockThreshold: 5}, true},
		{"below threshold", models.VariantInventory{Quantity: 3, ReservedQuantity: 0, LowStockThreshold: 5}, true},
		{"with reserved", models.VariantInventory{Quantity: 10, ReservedQuantity: 7, LowStockThreshold: 5}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inv.IsLowStock() != tt.expected {
				t.Fatalf("expected IsLowStock=%v for %+v", tt.expected, tt.inv)
			}
		})
	}
}
