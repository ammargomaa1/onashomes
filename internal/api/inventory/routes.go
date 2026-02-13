package inventory

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	adminRoutes := router.Group("/admin/inventory")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		adminRoutes.POST("/adjust", middleware.RequirePermission("inventory.adjust"), controller.AdjustInventory)
		adminRoutes.GET("/store/:storeFrontId", middleware.RequirePermission("inventory.adjust"), controller.ListInventoryByStore)
		adminRoutes.GET("/variant/:variantId/store/:storeFrontId", middleware.RequirePermission("inventory.adjust"), controller.GetVariantInventory)
		adminRoutes.GET("/low-stock/:storeFrontId", middleware.RequirePermission("inventory.adjust"), controller.GetLowStockAlerts)
		adminRoutes.GET("/:inventoryId/history", middleware.RequirePermission("inventory.adjust"), controller.GetAdjustmentHistory)
	}
}
