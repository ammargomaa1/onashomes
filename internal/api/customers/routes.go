package customers

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

// RegisterRoutes registers the customer routes
func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware()) // Ensure admin authentication
	{
		customers := admin.Group("/customers")
		{
			customers.GET("/search", controller.Search)
			customers.GET("", controller.List)
			customers.GET("/:id", controller.Get)
			customers.POST("", controller.Create)
			customers.PUT("/:id", controller.Update)
		}
	}
}
