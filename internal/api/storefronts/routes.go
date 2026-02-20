package storefronts

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	adminRoutes := router.Group("/admin/storefronts")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		adminRoutes.GET("", middleware.RequirePermission("storefronts.manage"), controller.List)
		adminRoutes.GET("/:id", middleware.RequirePermission("storefronts.manage"), controller.GetByID)
		adminRoutes.POST("", middleware.RequirePermission("storefronts.manage"), controller.Create)
		adminRoutes.PUT("/:id", middleware.RequirePermission("storefronts.manage"), controller.Update)
		adminRoutes.DELETE("/:id", middleware.RequirePermission("storefronts.manage"), controller.Delete)
	}
}
