package products

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

// RegisterRoutes sets up the product routes
func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	// Public product routes
	publicRoutes := router.Group("/products")
	{
		publicRoutes.GET("", controller.PublicListProducts)
		publicRoutes.GET("/:id", controller.PublicGetProduct)
	}

	// Admin routes (protected by admin auth middleware)
	adminRoutes := router.Group("/admin/products")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		// Product CRUD
		adminRoutes.GET("", middleware.RequirePermission("products.view"), controller.AdminListProducts)
		adminRoutes.GET("/:id", middleware.RequirePermission("products.view"), controller.AdminGetProduct)
		adminRoutes.POST("", middleware.RequirePermission("products.create"), controller.AdminCreateProduct)
		adminRoutes.PUT("/:id", middleware.RequirePermission("products.update"), controller.AdminUpdateProduct)
		adminRoutes.DELETE("/:id", middleware.RequirePermission("products.delete"), controller.AdminDeleteProduct)

		// Variants
		adminRoutes.POST("/:id/variants", middleware.RequirePermission("products.update"), controller.AdminCreateVariant)
		adminRoutes.PUT("/:id/variants/:variantId", middleware.RequirePermission("products.update"), controller.AdminUpdateVariant)
	}
}
