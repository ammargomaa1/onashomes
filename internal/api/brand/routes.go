package brand

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	adminRoutes := router.Group("/admin/brands")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	adminRoutes.Use() // Add necessary middleware here
	{
		adminRoutes.GET("", middleware.RequirePermission("brands.view"), controller.ListBrands)
		adminRoutes.GET("/:id", middleware.RequirePermission("brands.view"), controller.GetBrand)
		adminRoutes.POST("", middleware.RequirePermission("brands.create"), controller.CreateBrand)
		adminRoutes.PUT("/:id", middleware.RequirePermission("brands.update"), controller.UpdateBrand)
		adminRoutes.DELETE("/:id", middleware.RequirePermission("brands.delete"), controller.DeleteBrand)
		adminRoutes.PUT("/:id/restore", middleware.RequirePermission("brands.update"), controller.RestoreBrand)
		adminRoutes.GET("/deleted", middleware.RequirePermission("brands.view"), controller.ListDeletedBrands)
	}

	// Dropdown endpoint (admin auth, no permission check)
	adminRoutes.GET("/dropdown", controller.ListBrandsForDropdown)
}
