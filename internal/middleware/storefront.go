package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/database"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

// StoreFrontResolver resolves the storefront from the Host header or query parameter
func StoreFrontResolver() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try query param first (for development), then Host header
		domain := c.Query("store_domain")
		if domain == "" {
			domain = c.Request.Host
			// Strip port if present
			if idx := strings.Index(domain, ":"); idx != -1 {
				domain = domain[:idx]
			}
		}

		if domain == "" {
			utils.ErrorResponse(c, 400, "Store domain is required", nil)
			c.Abort()
			return
		}

		db := database.GetDB()
		var storeFront models.StoreFront
		if err := db.Where("domain = ? AND is_active = true", domain).First(&storeFront).Error; err != nil {
			utils.ErrorResponse(c, 404, "Store not found for domain: "+domain, nil)
			c.Abort()
			return
		}

		c.Set("store_front_id", storeFront.ID)
		c.Set("store_front", storeFront)
		c.Set("store_currency", storeFront.Currency)
		c.Set("store_language", storeFront.DefaultLanguage)

		c.Next()
	}
}
