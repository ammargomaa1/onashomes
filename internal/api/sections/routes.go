package sections

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	adminRoutes := router.Group("/admin/sections")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		adminRoutes.GET("", middleware.RequirePermission("sections.view"), controller.ListSections)
		adminRoutes.GET("/:id", middleware.RequirePermission("sections.view"), controller.GetSection)
		adminRoutes.POST("", middleware.RequirePermission("sections.create"), controller.CreateSection)
		adminRoutes.PUT("/:id", middleware.RequirePermission("sections.update"), controller.UpdateSection)
		adminRoutes.DELETE("/:id", middleware.RequirePermission("sections.delete"), controller.DeleteSection)
		adminRoutes.PUT("/:id/restore", middleware.RequirePermission("sections.update"), controller.RestoreSection)
		adminRoutes.GET("/deleted", middleware.RequirePermission("sections.view"), controller.ListDeletedSections)
	}

	// Public dropdown endpoint (no auth required)
	adminRoutes.GET("/dropdown", controller.ListSectionsForDropdown)
}
