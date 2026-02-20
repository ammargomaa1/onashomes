package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

// RegisterRoutes registers stats routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	g := router.Group("/admin/stats")
	g.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())

	g.GET("/item-order-counts", middleware.RequirePermission("orders.view"), handler.GetItemOrderCounts)
}
