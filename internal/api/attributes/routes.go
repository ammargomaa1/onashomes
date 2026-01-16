package attributes

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

// RegisterRoutes sets up the attribute routes.
func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	// Admin routes (protected by admin auth middleware)
	adminRoutes := router.Group("/admin/attributes")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		// Attributes CRUD
		adminRoutes.GET("", middleware.RequirePermission("attributes.view"), controller.ListAttributes)
		adminRoutes.GET("/deleted", middleware.RequirePermission("attributes.view"), controller.ListDeletedAttributes)
		adminRoutes.GET("/:id", middleware.RequirePermission("attributes.view"), controller.GetAttribute)
		adminRoutes.POST("", middleware.RequirePermission("attributes.create"), controller.CreateAttribute)
		adminRoutes.PUT("/:id", middleware.RequirePermission("attributes.update"), controller.UpdateAttribute)
		adminRoutes.DELETE("/:id", middleware.RequirePermission("attributes.delete"), controller.DeleteAttribute)
		adminRoutes.PUT("/:id/recover", middleware.RequirePermission("attributes.update"), controller.RecoverAttribute)

		// Attribute values
		adminRoutes.GET("/:id/values", middleware.RequirePermission("attributes.view"), controller.ListAttributeValues)
		adminRoutes.GET("/:id/values/deleted", middleware.RequirePermission("attributes.view"), controller.ListDeletedAttributeValues)
		adminRoutes.PUT("/:id/values/:valueId/recover", middleware.RequirePermission("attributes.update"), controller.RecoverAttributeValue)
		adminRoutes.DELETE("/:id/values/:valueId", middleware.RequirePermission("attributes.update"), controller.DeleteAttributeValue)
	}
}
