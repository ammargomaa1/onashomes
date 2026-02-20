package products

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

// RegisterV2Routes sets up V2 product routes (Phase 2)
func RegisterV2Routes(router *gin.RouterGroup, controller *ControllerV2) {
	// V2 Admin routes
	adminRoutes := router.Group("/admin/products/v2")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		// V2 Product CRUD
		adminRoutes.GET("", middleware.RequirePermission("products.view"), controller.AdminListProductsV2)
		adminRoutes.GET("/:id", middleware.RequirePermission("products.view"), controller.AdminGetProductV2)
		adminRoutes.POST("", middleware.RequirePermission("products.create"), controller.AdminCreateProductV2)
		adminRoutes.PUT("/:id", middleware.RequirePermission("products.update"), controller.AdminUpdateProductV2)

		// Status management
		adminRoutes.PATCH("/:id/status", middleware.RequirePermission("products.update"), controller.AdminUpdateProductStatus)

		// SEO management
		adminRoutes.PUT("/:id/seo", middleware.RequirePermission("seo.manage"), controller.AdminUpdateProductSEO)
		adminRoutes.GET("/:id/seo/suggest", middleware.RequirePermission("seo.manage"), controller.AdminSuggestSEO)

		// V2 Variants (with pricing)
		adminRoutes.POST("/:id/variants", middleware.RequirePermission("products.update"), controller.AdminCreateVariantV2)
		adminRoutes.PUT("/:id/variants/:variantId", middleware.RequirePermission("products.update"), controller.AdminUpdateVariantV2)

		// Product images
		adminRoutes.POST("/:id/images", middleware.RequirePermission("products.update"), controller.AdminAddProductImages)
		adminRoutes.DELETE("/:id/images/:imageId", middleware.RequirePermission("products.update"), controller.AdminRemoveProductImage)
		adminRoutes.PUT("/:id/images/cover", middleware.RequirePermission("products.update"), controller.AdminSetCoverImage)
	}

	// Storefront public routes (domain-resolved)
	sfRoutes := router.Group("/storefront")
	sfRoutes.Use(middleware.StoreFrontResolver())
	{
		sfRoutes.GET("/products", controller.StorefrontListProducts)
		sfRoutes.GET("/products/:slug", controller.StorefrontGetProduct)
		sfRoutes.GET("/products/:slug/structured-data", controller.StorefrontGetStructuredData)
	}
}
