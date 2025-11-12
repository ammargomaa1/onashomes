package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/database"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

// RequirePermission checks if the admin has the required permission
func RequirePermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role_id from context (set by AuthMiddleware)
		roleIDValue, exists := c.Get("role_id")
		if !exists {
			utils.ErrorResponse(c, 403, "No role assigned", nil)
			c.Abort()
			return
		}

		roleID, ok := roleIDValue.(int64)
		if !ok {
			utils.ErrorResponse(c, 403, "Invalid role ID", nil)
			c.Abort()
			return
		}

		// Check if role has the required permission
		db := database.GetDB()
		var role models.Role
		err := db.Preload("Permissions").Where("id = ?", roleID).First(&role).Error
		if err != nil {
			utils.ErrorResponse(c, 403, "Role not found", nil)
			c.Abort()
			return
		}

		// Check if permission exists in role's permissions
		hasPermission := false
		for _, perm := range role.Permissions {
			if perm.Name == permissionName {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.ErrorResponse(c, 403, "Insufficient permissions", map[string]string{
				"required": permissionName,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
