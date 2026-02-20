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
	g.GET("/sources", middleware.RequirePermission("orders.view"), controller.GetOrderSources) // New Enpoint
	g.GET("/meta", middleware.RequirePermission("orders.view"), controller.GetOrderMeta)
	g.GET("/:id", middleware.RequirePermission("orders.view"), controller.GetOrder)
	g.PUT("/:id", middleware.RequirePermission("orders.edit"), controller.UpdateOrder)
	g.POST("/:id/confirm", middleware.RequirePermission("orders.confirm"), controller.ConfirmOrder)
	g.POST("/:id/cancel", middleware.RequirePermission("orders.cancel"), controller.CancelOrder)
	g.POST("/:id/out-for-delivery", middleware.RequirePermission("orders.edit"), controller.MarkOutForDelivery)
	g.POST("/:id/complete", middleware.RequirePermission("orders.edit"), controller.CompleteOrder)
	// Assuming logic handles permission check or reuse "orders.edit" if "orders.complete" doesn't exist.
	// But best practice is specific permission.
	// We'll see if we need to add permission to seed. For now let's assume reuse "orders.edit" or check seed.
	// Wait, `ConfirmOrder` uses "orders.confirm". `CancelOrder` uses "orders.cancel".
	// Let's us "orders.edit" for now to avoid permission errors if "orders.complete" is missing from DB seeds.
	// OR check `internal/permissions/scanner.go` - it scans constants?
	// `RequirePermission` takes a string. The scanner scans these strings. So if I add it here, the scanner will pick it up on next restart!
	// So "orders.complete" is fine.
}
