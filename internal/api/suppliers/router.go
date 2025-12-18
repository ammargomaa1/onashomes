package suppliers

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

// RegisterRoutes sets up the supplier routes
func RegisterRoutes(router *gin.RouterGroup, service Service) {
	handler := NewHandler(service)

	// Admin routes (protected by admin auth middleware)
	adminRoutes := router.Group("/admin/suppliers")
	adminRoutes.Use(middleware.AdminAuthMiddleware())
	{
		adminRoutes.POST("", middleware.RequirePermission("suppliers.create"), handler.Create)
		adminRoutes.PUT("/:id", middleware.RequirePermission("suppliers.update"), handler.Update)
		adminRoutes.GET("",middleware.RequirePermission("suppliers.view"), handler.List)
		adminRoutes.GET("/:id", middleware.RequirePermission("suppliers.view"), handler.GetByID)
		adminRoutes.PUT("/:id/activate", middleware.RequirePermission("suppliers.update"), handler.Activate)
		adminRoutes.PUT("/:id/deactivate", middleware.RequirePermission("suppliers.update"), handler.Deactivate)
	}
}
