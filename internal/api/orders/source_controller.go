package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/utils"
)

// GetOrderSources lists all available order sources
func (c *Controller) GetOrderSources(ctx *gin.Context) {
	sources, err := c.service.GetOrderSources()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch order sources", nil)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Order sources retrieved successfully", sources)
}
