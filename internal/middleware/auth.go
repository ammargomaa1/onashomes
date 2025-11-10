package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/utils"
)

// AuthMiddleware validates JWT token and sets user/admin info in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, 401, "Authorization header required", nil)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, 401, "Invalid authorization header format", nil)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token, utils.AccessToken)
		if err != nil {
			utils.ErrorResponse(c, 401, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Set claims in context for use in handlers
		c.Set("entity_id", claims.EntityID)
		c.Set("entity_type", claims.EntityType)
		if claims.RoleID != nil {
			c.Set("role_id", *claims.RoleID)
		}

		c.Next()
	}
}

// AdminAuthMiddleware ensures the authenticated entity is an admin
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		entityType, exists := c.Get("entity_type")
		if !exists || entityType != utils.EntityAdmin {
			utils.ErrorResponse(c, 403, "Admin access required", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserAuthMiddleware ensures the authenticated entity is a user
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		entityType, exists := c.Get("entity_type")
		if !exists || entityType != utils.EntityUser {
			utils.ErrorResponse(c, 403, "User access required", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
