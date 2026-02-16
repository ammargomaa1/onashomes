package orders

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	g := router.Group("/admin/orders")
	g.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())

	g.POST("", middleware.RequirePermission("orders.create"), controller.CreateOrder)
	g.GET("", middleware.RequirePermission("orders.view"), controller.ListOrders)
	g.GET("/meta", middleware.RequirePermission("orders.view"), controller.GetOrderMeta)
	g.GET("/:id", middleware.RequirePermission("orders.view"), controller.GetOrder)
	g.PUT("/:id", middleware.RequirePermission("orders.edit"), controller.UpdateOrder)
	g.POST("/:id/confirm", middleware.RequirePermission("orders.confirm"), controller.ConfirmOrder)
	g.POST("/:id/cancel", middleware.RequirePermission("orders.cancel"), controller.CancelOrder)
}
