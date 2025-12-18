package users

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	users := router.Group("/users")
	{
		// Public routes
		users.POST("/register", controller.Register)
		users.POST("/login", controller.Login)
		users.POST("/refresh", controller.RefreshToken)

		// Protected routes (require authentication)
		authenticated := users.Group("")
		authenticated.Use(middleware.AuthMiddleware(), middleware.UserAuthMiddleware())
		{
			authenticated.GET("/profile", controller.GetProfile)
			authenticated.PUT("/profile", controller.UpdateProfile)
		}

		// Admin only routes
		adminOnly := users.Group("")
		adminOnly.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
		{
			// List users - requires users.view permission
			adminOnly.GET("", middleware.RequirePermission("users.view"), controller.List)
		}
	}
}
