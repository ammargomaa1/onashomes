package database

import (
	"fmt"
	"log"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

// AutoMigrate runs automatic migrations for all models
func AutoMigrate(db *gorm.DB) error {
	fmt.Println("Running auto migrations...")

	err := db.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
		&models.Admin{},
		&models.File{},
		&models.Section{},
		&models.Category{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
		return err
	}

	fmt.Println("✓ Migrations completed successfully")
	return nil
}

// SeedDefaultData seeds initial data for roles and admin user
func SeedDefaultData(db *gorm.DB) error {
	fmt.Println("Seeding default data...")

	// Check if roles already exist
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)
	if roleCount > 0 {
		fmt.Println("✓ Roles already exist, skipping role creation...")
		return ensureDefaultAdmin(db)
	}

	// Create Super Admin role (permissions will be assigned by scanner)
	superAdminRole := models.Role{
		Name:        "Super Admin",
		Description: "Full system access - all permissions automatically assigned",
	}

	if err := db.Create(&superAdminRole).Error; err != nil {
		return fmt.Errorf("failed to create super admin role: %v", err)
	}

	// Create default Admin role (limited permissions)
	adminRole := models.Role{
		Name:        "Admin",
		Description: "Standard admin access with limited permissions",
	}

	if err := db.Create(&adminRole).Error; err != nil {
		return fmt.Errorf("failed to create admin role: %v", err)
	}

	// Create default admin user
	if err := createDefaultAdmin(db, adminRole.ID); err != nil {
		return err
	}

	fmt.Println("✓ Default roles and admin user created successfully")
	return nil
}

// ensureDefaultAdmin ensures the default admin user exists
func ensureDefaultAdmin(db *gorm.DB) error {
	var adminCount int64
	db.Model(&models.Admin{}).Where("email = ?", "admin@onas.com").Count(&adminCount)

	if adminCount > 0 {
		fmt.Println("✓ Default admin already exists")
		return nil
	}

	// Get Admin role
	var adminRole models.Role
	if err := db.Where("name = ?", "Admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("admin role not found: %v", err)
	}

	return createDefaultAdmin(db, adminRole.ID)
}

// createDefaultAdmin creates the default admin user
func createDefaultAdmin(db *gorm.DB, roleID int64) error {
	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	admin := models.Admin{
		Email:     "admin@onas.com",
		Password:  hashedPassword,
		FirstName: "Admin",
		LastName:  "Admin",
		RoleID:    roleID,
		IsActive:  true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin: %v", err)
	}

	fmt.Println("✓ Default admin user created: admin@onas.com / password")
	return nil
}
