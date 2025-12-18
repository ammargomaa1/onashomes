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
		var permission models.Permission
		query := "SELECT permissions.id, permissions.name FROM role_permissions JOIN permissions ON role_permissions.permission_id = permissions.id WHERE role_permissions.role_id = ? AND permissions.name = ?"
		err := db.Raw(query, roleID, permissionName).Scan(&permission).Error
		if err != nil {
			utils.ErrorResponse(c, 403, "Permission not found", nil)
			c.Abort()
			return
		}

		// Check if permission exists in role's permissions
		hasPermission := false
		if permissionName == permission.Name {
			hasPermission = true
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
