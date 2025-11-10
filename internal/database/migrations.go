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
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
		return err
	}

	fmt.Println("✓ Migrations completed successfully")
	return nil
}

// SeedDefaultData seeds initial data for roles and permissions
func SeedDefaultData(db *gorm.DB) error {
	fmt.Println("Seeding default data...")

	// Check if data already exists
	var count int64
	db.Model(&models.Permission{}).Count(&count)
	if count > 0 {
		fmt.Println("✓ Data already seeded, skipping...")
		return nil
	}

	// Create default permissions
	permissions := []models.Permission{
		{Name: "users.view", Description: "View users"},
		{Name: "users.create", Description: "Create users"},
		{Name: "users.update", Description: "Update users"},
		{Name: "users.delete", Description: "Delete users"},
		{Name: "admins.view", Description: "View admins"},
		{Name: "admins.create", Description: "Create admins"},
		{Name: "admins.update", Description: "Update admins"},
		{Name: "admins.delete", Description: "Delete admins"},
		{Name: "roles.view", Description: "View roles"},
		{Name: "roles.create", Description: "Create roles"},
		{Name: "roles.update", Description: "Update roles"},
		{Name: "roles.delete", Description: "Delete roles"},
	}

	for _, perm := range permissions {
		if err := db.Create(&perm).Error; err != nil {
			return fmt.Errorf("failed to create permission %s: %v", perm.Name, err)
		}
	}

	// Create super admin role with all permissions
	var allPermissions []models.Permission
	db.Find(&allPermissions)

	superAdminRole := models.Role{
		Name:        "Super Admin",
		Description: "Full system access",
		Permissions: allPermissions,
	}

	if err := db.Create(&superAdminRole).Error; err != nil {
		return fmt.Errorf("failed to create super admin role: %v", err)
	}

	// Create default admin role
	adminPermissions := []models.Permission{}
	db.Where("name LIKE ?", "users.%").Find(&adminPermissions)

	adminRole := models.Role{
		Name:        "Admin",
		Description: "Standard admin access",
		Permissions: adminPermissions,
	}

	if err := db.Create(&adminRole).Error; err != nil {
		return fmt.Errorf("failed to create admin role: %v", err)
	}

	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	admin := models.Admin{
		Email:     "admin@onas.com",
		Password:  hashedPassword,
		FirstName: "Admin",
		LastName:  "Admin",
		RoleID:    adminRole.ID,
		IsActive:  true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin: %v", err)
	}

	fmt.Println("✓ Default data seeded successfully")
	return nil
}
