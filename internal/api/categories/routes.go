package categories

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	adminRoutes := router.Group("/admin/categories")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		adminRoutes.GET("", middleware.RequirePermission("categories.view"), controller.ListCategories)
		adminRoutes.GET("/:id", middleware.RequirePermission("categories.view"), controller.GetCategory)
		adminRoutes.POST("", middleware.RequirePermission("categories.create"), controller.CreateCategory)
		adminRoutes.PUT("/:id", middleware.RequirePermission("categories.update"), controller.UpdateCategory)
		adminRoutes.DELETE("/:id", middleware.RequirePermission("categories.delete"), controller.DeleteCategory)
		adminRoutes.PUT("/:id/restore", middleware.RequirePermission("categories.update"), controller.RestoreCategory)
		adminRoutes.GET("/deleted", middleware.RequirePermission("categories.view"), controller.ListDeletedCategories)
		adminRoutes.GET("/by-section/:section_id", middleware.RequirePermission("categories.view"), controller.ListCategoriesBySection)
	}

	// Dropdown endpoint (admin auth, no permission check)
	adminRoutes.GET("/dropdown", controller.ListCategoriesForDropdown)
}
