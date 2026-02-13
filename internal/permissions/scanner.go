package permissions

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

// PermissionInfo holds information about a discovered permission
type PermissionInfo struct {
	Name        string
	Description string
	Module      string
	Action      string
}

// Scanner scans routes and manages permissions
type Scanner struct {
	db          *gorm.DB
	permissions map[string]PermissionInfo
}

// NewScanner creates a new permission scanner
func NewScanner(db *gorm.DB) *Scanner {
	return &Scanner{
		db:          db,
		permissions: make(map[string]PermissionInfo),
	}
}

// ScanAndSync scans routes and syncs permissions to database
func (s *Scanner) ScanAndSync(engine *gin.Engine) error {
	log.Println("🔍 Scanning routes for permissions...")

	// Add predefined permissions that we know exist
	s.addPredefinedPermissions()

	// Scan routes for additional permissions
	routes := engine.Routes()
	for _, route := range routes {
		s.extractPermissionsFromRoute(route)
	}

	return s.syncPermissionsToDatabase()
}

// addPredefinedPermissions adds permissions that we know exist from the middleware usage
func (s *Scanner) addPredefinedPermissions() {
	predefined := []PermissionInfo{
		// Admin permissions (from routes)
		{Name: "admins.view", Description: "View and list admins", Module: "admins", Action: "view"},
		{Name: "admins.create", Description: "Create new admins", Module: "admins", Action: "create"},
		{Name: "admins.update", Description: "Update existing admins", Module: "admins", Action: "update"},
		{Name: "admins.delete", Description: "Delete existing admins", Module: "admins", Action: "delete"},

		// User permissions (from routes)
		{Name: "users.view", Description: "View and list users", Module: "users", Action: "view"},
		{Name: "users.create", Description: "Create new users", Module: "users", Action: "create"},
		{Name: "users.update", Description: "Update existing users", Module: "users", Action: "update"},
		{Name: "users.delete", Description: "Delete existing users", Module: "users", Action: "delete"},

		// Future module permissions (can be added as needed)
		{Name: "roles.view", Description: "View and list roles", Module: "roles", Action: "view"},
		{Name: "roles.create", Description: "Create new roles", Module: "roles", Action: "create"},
		{Name: "roles.update", Description: "Update existing roles", Module: "roles", Action: "update"},
		{Name: "roles.delete", Description: "Delete existing roles", Module: "roles", Action: "delete"},

		{Name: "permissions.view", Description: "View and list permissions", Module: "permissions", Action: "view"},
		{Name: "permissions.create", Description: "Create new permissions", Module: "permissions", Action: "create"},
		{Name: "permissions.update", Description: "Update existing permissions", Module: "permissions", Action: "update"},
		{Name: "permissions.delete", Description: "Delete existing permissions", Module: "permissions", Action: "delete"},

		// File permissions
		{Name: "files.view", Description: "View and list files", Module: "files", Action: "view"},
		{Name: "files.create", Description: "Upload new files", Module: "files", Action: "create"},
		{Name: "files.update", Description: "Update existing files", Module: "files", Action: "update"},
		{Name: "files.delete", Description: "Delete existing files", Module: "files", Action: "delete"},

		// Product permissions
		{Name: "products.view", Description: "View and list products", Module: "products", Action: "view"},
		{Name: "products.create", Description: "Create new products", Module: "products", Action: "create"},
		{Name: "products.update", Description: "Update existing products", Module: "products", Action: "update"},
		{Name: "products.delete", Description: "Delete existing products", Module: "products", Action: "delete"},

		// Phase 2: Inventory permissions
		{Name: "inventory.adjust", Description: "Adjust product inventory quantities", Module: "inventory", Action: "adjust"},
		{Name: "inventory.view", Description: "View inventory levels", Module: "inventory", Action: "view"},

		// Phase 2: StoreFront permissions
		{Name: "storefronts.manage", Description: "Manage store fronts (CRUD)", Module: "storefronts", Action: "manage"},

		// Phase 2: SEO permissions
		{Name: "seo.manage", Description: "Manage product SEO metadata", Module: "seo", Action: "manage"},
	}

	for _, perm := range predefined {
		s.permissions[perm.Name] = perm
	}
}

// extractPermissionsFromRoute extracts permissions from route analysis
func (s *Scanner) extractPermissionsFromRoute(route gin.RouteInfo) {
	path := route.Path
	method := route.Method

	// Skip non-API routes
	if !strings.HasPrefix(path, "/api/") {
		return
	}

	// Skip auth routes
	if strings.Contains(path, "/login") ||
		strings.Contains(path, "/register") ||
		strings.Contains(path, "/refresh") ||
		strings.Contains(path, "/health") {
		return
	}

	// Extract module from path
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "api" {
		return
	}

	// partsCount := len(pathParts)
	module := pathParts[len(pathParts)-1]

	// if partsCount > 4 && strings.HasPrefix(pathParts[partsCount-2], ":") {
	// 	module = pathParts[partsCount-3]
	// }
	deductParts := 0
	if strings.HasSuffix(path, "deleted") ||
		strings.HasSuffix(path, "archived") ||
		strings.HasSuffix(path, "recover") ||
		strings.HasSuffix(path, "restore") ||
		strings.HasSuffix(path, "activate") ||
		strings.HasSuffix(path, "deactivate") {
		deductParts = 2
		module = pathParts[len(pathParts)-deductParts]
	}

	// if last is a param, take the one before it
	if strings.HasPrefix(module, ":") && len(pathParts) >= 3 {
		if deductParts == 0 {
			deductParts = 2
		} else {
			deductParts++
		}
		module = pathParts[len(pathParts)-deductParts]
	}

	// Infer action from HTTP method and path
	var action string
	switch method {
	case "GET":
		action = "view"
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		return
	}

	permissionName := fmt.Sprintf("%s.%s", module, action)

	// Only add if not already exists
	if _, exists := s.permissions[permissionName]; !exists {
		s.permissions[permissionName] = PermissionInfo{
			Name:        permissionName,
			Description: s.generateDescription(module, action),
			Module:      module,
			Action:      action,
		}
	}
}

// generateDescription generates a human-readable description
func (s *Scanner) generateDescription(module, action string) string {
	actionMap := map[string]string{
		"view":   "View and read",
		"create": "Create new",
		"update": "Update existing",
		"delete": "Delete existing",
	}

	actionDesc := actionMap[action]
	if actionDesc == "" {
		actionDesc = strings.Title(action)
	}

	return fmt.Sprintf("%s %s", actionDesc, module)
}

// syncPermissionsToDatabase syncs discovered permissions to database
func (s *Scanner) syncPermissionsToDatabase() error {
	log.Printf("🔍 Found %d permissions to sync", len(s.permissions))

	// Get existing permissions
	var existingPermissions []models.Permission
	if err := s.db.Find(&existingPermissions).Error; err != nil {
		return fmt.Errorf("failed to fetch existing permissions: %v", err)
	}

	// Create map for quick lookup
	existingMap := make(map[string]bool)
	for _, perm := range existingPermissions {
		existingMap[perm.Name] = true
	}

	// Find new permissions
	var newPermissions []models.Permission
	for _, permInfo := range s.permissions {
		if !existingMap[permInfo.Name] {
			newPermissions = append(newPermissions, models.Permission{
				Name:        permInfo.Name,
				Description: permInfo.Description,
			})
			log.Printf("➕ New permission: %s - %s", permInfo.Name, permInfo.Description)
		}
	}

	// Insert new permissions
	if len(newPermissions) > 0 {
		if err := s.db.Create(&newPermissions).Error; err != nil {
			return fmt.Errorf("failed to create permissions: %v", err)
		}

		log.Printf("✅ Added %d new permissions", len(newPermissions))

		// Assign to Super Admin
		if err := s.assignToSuperAdmin(newPermissions); err != nil {
			log.Printf("⚠️  Failed to assign permissions to Super Admin: %v", err)
		}
	} else {
		log.Println("✅ No new permissions to add")
	}

	return nil
}

// assignToSuperAdmin assigns new permissions to Super Admin role
func (s *Scanner) assignToSuperAdmin(newPermissions []models.Permission) error {
	var superAdminRole models.Role
	if err := s.db.Where("name = ?", "Super Admin").First(&superAdminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("⚠️  Super Admin role not found")
			return nil
		}
		return err
	}

	// Assign new permissions
	if err := s.db.Model(&superAdminRole).Association("Permissions").Append(&newPermissions); err != nil {
		return err
	}

	log.Printf("🔐 Assigned %d permissions to Super Admin", len(newPermissions))
	return nil
}

// AddCustomPermission allows adding custom permissions
func (s *Scanner) AddCustomPermission(name, description string) {
	parts := strings.Split(name, ".")
	module := "system"
	action := name

	if len(parts) == 2 {
		module = parts[0]
		action = parts[1]
	}

	s.permissions[name] = PermissionInfo{
		Name:        name,
		Description: description,
		Module:      module,
		Action:      action,
	}
}
