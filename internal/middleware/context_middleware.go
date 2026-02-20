package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/utils"
)

// ContextMiddleware registers the gin.Context for the current goroutine
// and ensures cleanup after the request is processed
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Register the context for the current goroutine
		utils.SetContextForGoroutine(c)

		// Ensure cleanup after the request is done
		c.Next()

		// Cleanup the context for the goroutine
		utils.CleanupContextForGoroutine()
	}
}
