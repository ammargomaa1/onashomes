package files

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
)

// RegisterRoutes registers file-related routes
func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	files := router.Group("/files")
	{
		// All file operations require admin authentication
		authenticated := files.Group("")
		authenticated.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
		{
			// Upload file - requires files.create permission
			authenticated.POST("/upload", middleware.RequirePermission("files.create"), controller.Upload)
			
			// Get upload configuration - requires files.view permission
			authenticated.GET("/config", middleware.RequirePermission("files.view"), controller.GetUploadConfig)
			
			// List files - requires files.view permission
			authenticated.GET("", middleware.RequirePermission("files.view"), controller.ListFiles)
			
			// Get file by ID - requires files.view permission
			authenticated.GET("/:id", middleware.RequirePermission("files.view"), controller.GetFile)
			
			// Delete file - requires files.delete permission
			authenticated.DELETE("/:id", middleware.RequirePermission("files.delete"), controller.DeleteFile)
		}
	}
}
