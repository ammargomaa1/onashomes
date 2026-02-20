package storefronts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/storefronts/requests"
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

	if err := db.AutoMigrate(&models.StoreFront{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func setupRouter(controller *Controller) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("/api")

	// Simplified routes without auth middleware for testing
	api.POST("/admin/storefronts", controller.Create)
	api.GET("/admin/storefronts/:id", controller.GetByID)
	api.PUT("/admin/storefronts/:id", controller.Update)
	api.DELETE("/admin/storefronts/:id", controller.Delete)

	return r
}

func createTestController(db *gorm.DB) *Controller {
	repo := NewRepository(db)
	service := NewService(repo)
	return NewController(service)
}

func TestCreateStoreFront_Success(t *testing.T) {
	db := setupTestDB(t)
	controller := createTestController(db)
	router := setupRouter(controller)

	body := requests.CreateStoreFrontRequest{
		Name:            "Test Store",
		Slug:            "test-store",
		Domain:          "test.example.com",
		Currency:        "SAR",
		DefaultLanguage: "ar",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/admin/storefronts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	// Verify stored in DB
	var sf models.StoreFront
	if err := db.First(&sf).Error; err != nil {
		t.Fatalf("expected storefront in DB: %v", err)
	}
	if sf.Name != "Test Store" {
		t.Fatalf("expected name 'Test Store', got '%s'", sf.Name)
	}
	if sf.Slug != "test-store" {
		t.Fatalf("expected slug 'test-store', got '%s'", sf.Slug)
	}
}

func TestCreateStoreFront_DuplicateSlug(t *testing.T) {
	db := setupTestDB(t)
	controller := createTestController(db)
	router := setupRouter(controller)

	// Create first one
	db.Create(&models.StoreFront{
		Name: "Existing", Slug: "test-store", Domain: "existing.com",
		Currency: "SAR", DefaultLanguage: "ar", IsActive: true,
	})

	body := requests.CreateStoreFrontRequest{
		Name:            "Another Store",
		Slug:            "test-store",
		Domain:          "another.com",
		Currency:        "SAR",
		DefaultLanguage: "ar",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/admin/storefronts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for duplicate slug, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestCreateStoreFront_DuplicateDomain(t *testing.T) {
	db := setupTestDB(t)
	controller := createTestController(db)
	router := setupRouter(controller)

	db.Create(&models.StoreFront{
		Name: "Existing", Slug: "existing", Domain: "test.example.com",
		Currency: "SAR", DefaultLanguage: "ar", IsActive: true,
	})

	body := requests.CreateStoreFrontRequest{
		Name:            "Another Store",
		Slug:            "another",
		Domain:          "test.example.com",
		Currency:        "SAR",
		DefaultLanguage: "ar",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/admin/storefronts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for duplicate domain, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestGetStoreFront_Success(t *testing.T) {
	db := setupTestDB(t)
	controller := createTestController(db)
	router := setupRouter(controller)

	sf := models.StoreFront{
		Name: "My Store", Slug: "my-store", Domain: "mystore.com",
		Currency: "SAR", DefaultLanguage: "ar", IsActive: true,
	}
	db.Create(&sf)

	req, _ := http.NewRequest("GET", "/api/admin/storefronts/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestGetStoreFront_NotFound(t *testing.T) {
	db := setupTestDB(t)
	controller := createTestController(db)
	router := setupRouter(controller)

	req, _ := http.NewRequest("GET", "/api/admin/storefronts/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateStoreFront_Success(t *testing.T) {
	db := setupTestDB(t)
	controller := createTestController(db)
	router := setupRouter(controller)

	sf := models.StoreFront{
		Name: "Old Name", Slug: "old-name", Domain: "old.com",
		Currency: "SAR", DefaultLanguage: "ar", IsActive: true,
	}
	db.Create(&sf)

	body := requests.UpdateStoreFrontRequest{
		Name:            "New Name",
		Domain:          "new.com",
		Currency:        "USD",
		DefaultLanguage: "en",
		IsActive:        true,
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", "/api/admin/storefronts/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify update in DB
	var updated models.StoreFront
	db.First(&updated, 1)
	if updated.Name != "New Name" {
		t.Fatalf("expected name 'New Name', got '%s'", updated.Name)
	}
	if updated.Domain != "new.com" {
		t.Fatalf("expected domain 'new.com', got '%s'", updated.Domain)
	}
}

func TestDeleteStoreFront_Success(t *testing.T) {
	db := setupTestDB(t)
	controller := createTestController(db)
	router := setupRouter(controller)

	sf := models.StoreFront{
		Name: "To Delete", Slug: "to-delete", Domain: "del.com",
		Currency: "SAR", DefaultLanguage: "ar", IsActive: true,
	}
	db.Create(&sf)

	req, _ := http.NewRequest("DELETE", "/api/admin/storefronts/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, w.Code, w.Body.String())
	}

	// Verify deleted
	var count int64
	db.Model(&models.StoreFront{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 storefronts, got %d", count)
	}
}
