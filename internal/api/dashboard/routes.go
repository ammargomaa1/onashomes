package dashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

// RegisterRoutes registers dashboard routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	g := router.Group("/admin/dashboard")
	g.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())

	g.GET("/overview", middleware.RequirePermission("orders.view"), handler.GetOverview)
}
