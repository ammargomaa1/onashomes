package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/config"
	"github.com/onas/ecommerce-api/internal/api/admins"
	"github.com/onas/ecommerce-api/internal/api/users"
	"github.com/onas/ecommerce-api/internal/database"
	"github.com/onas/ecommerce-api/internal/middleware"
	"github.com/onas/ecommerce-api/internal/permissions"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database connection
	db := database.GetDB()

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed default data
	if err := database.SeedDefaultData(db); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
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
	}

	// Scan routes and sync permissions to database
	log.Println("🔍 Scanning routes for permissions...")
	permissionScanner := permissions.NewScanner(db)
	if err := permissionScanner.ScanAndSync(router); err != nil {
		log.Printf("⚠️  Failed to scan and sync permissions: %v", err)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("\n🚀 Server starting on http://localhost%s\n", addr)
	fmt.Println("📚 API Documentation: http://localhost" + addr + "/api")
	fmt.Println("💚 Health Check: http://localhost" + addr + "/health")
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
