package admins

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	admins := router.Group("/admin")
	{
		// Public routes
		admins.POST("/login", controller.Login)
		admins.POST("/refresh", controller.RefreshToken)

		// Protected routes (require admin authentication)
		authenticated := admins.Group("")
		authenticated.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
		{
			// List admins - requires admins.view permission
			authenticated.GET("", middleware.RequirePermission("admins.view"), controller.List)

			// Get admin by ID - requires admins.view permission
			authenticated.GET("/:id", middleware.RequirePermission("admins.view"), controller.GetByID)

			// Create admin - requires admins.create permission
			authenticated.POST("", middleware.RequirePermission("admins.create"), controller.Create)

			// Update admin - requires admins.update permission
			authenticated.PUT("/:id", middleware.RequirePermission("admins.update"), controller.Update)

			// Delete admin - requires admins.delete permission
			authenticated.DELETE("/:id", middleware.RequirePermission("admins.delete"), controller.Delete)
		}
	}
}
