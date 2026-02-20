package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/config"
	"github.com/onas/ecommerce-api/internal/api/admins"
	"github.com/onas/ecommerce-api/internal/api/attributes"
	"github.com/onas/ecommerce-api/internal/api/brand"
	"github.com/onas/ecommerce-api/internal/api/categories"
	"github.com/onas/ecommerce-api/internal/api/customers"
	"github.com/onas/ecommerce-api/internal/api/dashboard"
	"github.com/onas/ecommerce-api/internal/api/files"
	"github.com/onas/ecommerce-api/internal/api/inventory"
	"github.com/onas/ecommerce-api/internal/api/locations"
	"github.com/onas/ecommerce-api/internal/api/orders"
	"github.com/onas/ecommerce-api/internal/api/products"
	"github.com/onas/ecommerce-api/internal/api/sections"
	"github.com/onas/ecommerce-api/internal/api/stats"
	"github.com/onas/ecommerce-api/internal/api/storefronts"
	"github.com/onas/ecommerce-api/internal/api/suppliers"
	"github.com/onas/ecommerce-api/internal/api/users"
	"github.com/onas/ecommerce-api/internal/database"
	"github.com/onas/ecommerce-api/internal/middleware"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/permissions"
	"github.com/onas/ecommerce-api/internal/services"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Check for migration-only flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--migrate-only":
			runMigrationsOnly()
			return
		case "--rollback":
			runRollback()
			return
		case "--migration-status":
			showMigrationStatus()
			return
		}
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database connection
	db := database.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm DB: %v", err)
	}

	// Run SQL migrations first
	migrator := database.NewMigrator(db, "migrations")
	if err := migrator.RunMigrations(); err != nil {
		log.Fatalf("Failed to run SQL migrations: %v", err)
	}

	// Run GORM auto migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run auto migrations: %v", err)
	}

	// Seed default data
	if err := database.SeedDefaultData(db); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ContextMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "E-commerce API is running",
		})
	})

	// API v1 routes
	api := router.Group("/api")
	{
		// User module
		userRepo := users.NewRepository(db)
		userService := users.NewService(userRepo)
		userController := users.NewController(userService)
		users.RegisterRoutes(api, userController)

		// Admin module
		adminRepo := admins.NewRepository(db)
		adminService := admins.NewService(adminRepo)
		adminController := admins.NewController(adminService)
		admins.RegisterRoutes(api, adminController)

		// File module
		fileService := services.NewFileService(db)
		fileController := files.NewController(fileService)
		files.RegisterRoutes(api, fileController)

		// Supplier module
		supplierRepo := suppliers.NewRepository(db)
		supplierService := suppliers.NewService(*supplierRepo)
		suppliers.RegisterRoutes(api, *supplierService)

		// Attributes module (admin)
		attributeRepo := attributes.NewRepository(db)
		attributeService := attributes.NewService(*attributeRepo)
		attributeController := attributes.NewController(*attributeService)
		attributes.RegisterRoutes(api, attributeController)

		// Products module (admin + public)
		productRepo := products.NewRepository(db)
		productService := products.NewService(sqlDB, productRepo)
		productController := products.NewController(productService)
		products.RegisterRoutes(api, productController)

		// Brand module (admin)
		brandRepo := brand.NewRepository()
		brandService := brand.NewService(brandRepo)
		brandController := brand.NewController(brandService)
		brand.RegisterRoutes(api, brandController)

		// Sections module (admin)
		sectionsRepo := sections.NewRepository()
		sectionsService := sections.NewService(sectionsRepo, db)
		sectionsController := sections.NewController(sectionsService)
		sections.RegisterRoutes(api, sectionsController)

		// Categories module (admin)
		categoriesRepo := categories.NewRepository()
		categoriesService := categories.NewService(categoriesRepo, db)
		categoriesController := categories.NewController(categoriesService)
		categories.RegisterRoutes(api, categoriesController)

		// StoreFront module (Phase 2)
		sfRepo := storefronts.NewRepository(db)
		sfService := storefronts.NewService(sfRepo)
		sfController := storefronts.NewController(sfService)
		storefronts.RegisterRoutes(api, sfController)

		// Inventory module (Phase 2)
		invRepo := inventory.NewRepository(db)
		invService := inventory.NewService(db, invRepo)
		invController := inventory.NewController(invService)
		inventory.RegisterRoutes(api, invController)

		// Products Phase 2 (V2)
		productV2Repo := products.NewV2Repository(db)
		productV2Service := products.NewServiceV2(db, productV2Repo, invRepo)
		productV2Controller := products.NewControllerV2(productV2Service)
		products.RegisterV2Routes(api, productV2Controller)

		// Customers module
		customerRepo := customers.NewRepository(db)
		customerService := customers.NewService(customerRepo)
		customerController := customers.NewController(customerService)
		customers.RegisterRoutes(api, customerController)

		// Location module
		locationRepo := locations.NewRepository(db)
		locationService := locations.NewService(locationRepo)
		locationController := locations.NewController(locationService)
		locations.RegisterRoutes(api, locationController)

		// Orders module (Priority 1)
		orderRepo := orders.NewRepository(db)
		orderService := orders.NewService(db, orderRepo, invService, customerService)
		orderController := orders.NewController(orderService)
		orders.RegisterRoutes(api, orderController)

		// Stats module
		statsHandler := stats.NewHandler(db)
		stats.RegisterRoutes(api, statsHandler)

		// Dashboard module
		dashboardHandler := dashboard.NewHandler(db)
		dashboard.RegisterRoutes(api, dashboardHandler)
	}

	// Scan routes and sync permissions to database
	log.Println("🔍 Scanning routes for permissions...")
	permissionScanner := permissions.NewScanner(db)
	if err := permissionScanner.ScanAndSync(router); err != nil {
		log.Printf("⚠️  Failed to scan and sync permissions: %v", err)
	}

	// Run permission fixes (dev only)
	// if err := fixPermissions(db); err != nil {
	// 	log.Printf("⚠️  Permission fix warning: %v", err)
	// }

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("\n🚀 Server starting on http://localhost%s\n", addr)
	fmt.Println("📚 API Documentation: http://localhost" + addr + "/api")
	fmt.Println("💚 Health Check: http://localhost" + addr + "/health")

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func fixPermissions(db *gorm.DB) error {
	log.Println("🔧 Running permission fix diagnostics...")

	// using models.Permission to avoid "model value required" and ensure correct table usage
	var perm models.Permission
	// Find or Create permission
	if err := db.Where(models.Permission{Name: "orders.edit"}).
		Attrs(models.Permission{Description: "Edit orders"}).
		FirstOrCreate(&perm).Error; err != nil {
		log.Printf("⚠️  Failed to fix permissions: %v", err)
		return nil // Don't crash app, just log
	}

	log.Printf("Permission orders.edit ready (ID: %d)", perm.ID)

	// Insert linkage for Role 1 and Role 2
	roles := []int64{1, 2}
	for _, roleID := range roles {
		var count int64
		db.Table("role_permissions").Where("role_id = ? AND permission_id = ?", roleID, perm.ID).Count(&count)
		if count == 0 {
			log.Printf("Assigning orders.edit (ID: %d) to Role %d...", perm.ID, roleID)
			if err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, perm.ID).Error; err != nil {
				log.Printf("Failed to assign permission to Role %d: %v", roleID, err)
			} else {
				log.Printf("✅ Assigned orders.edit to Role %d", roleID)
			}
		} else {
			log.Printf("✅ Role %d already has orders.edit (ID: %d)", roleID, perm.ID)
		}
	}

	return nil
}

// runMigrationsOnly runs only the SQL migrations without starting the server
func runMigrationsOnly() {
	log.Println("🔄 Running migrations only...")

	db := database.GetDB()
	migrator := database.NewMigrator(db, "migrations")

	if err := migrator.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Migrations completed successfully")
}

// runRollback rolls back the last migration
func runRollback() {
	log.Println("🔄 Rolling back last migration...")

	db := database.GetDB()
	migrator := database.NewMigrator(db, "migrations")

	if err := migrator.RollbackMigration(); err != nil {
		log.Fatalf("Failed to rollback migration: %v", err)
	}

	log.Println("✅ Rollback completed successfully")
}

// showMigrationStatus shows the current migration status
func showMigrationStatus() {
	log.Println("📋 Migration Status:")

	db := database.GetDB()
	migrator := database.NewMigrator(db, "migrations")

	// This is a simplified status - you could enhance this to show more details
	if err := migrator.RunMigrations(); err != nil {
		log.Printf("❌ Error checking migrations: %v", err)
		return
	}

	log.Println("✅ All migrations are up to date")
}
