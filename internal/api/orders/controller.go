package orders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/orders/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateOrder(ctx *gin.Context) {
	var req requests.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	// Get admin ID from auth context
	adminIDVal, exists := ctx.Get("entity_id")
	if !exists {
		utils.ErrorResponse(ctx, 401, "Unauthorized", nil)
		return
	}
	adminID := adminIDVal.(int64)
	res := c.service.CreateOrder(req, adminID)
	utils.WriteResource(ctx, res)
}

func (c *Controller) UpdateOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid order id")
		return
	}

	var req requests.UpdateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res := c.service.UpdateOrder(id, req)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ConfirmOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid order id")
		return
	}

	res := c.service.ConfirmOrder(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) CancelOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid order id")
		return
	}

	res := c.service.CancelOrder(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ListOrders(ctx *gin.Context) {
	var filter requests.OrderFilterRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	pagination := utils.ParsePaginationParams(ctx)
	res := c.service.ListOrders(filter, pagination)
	utils.WriteResource(ctx, res)
}

func (c *Controller) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid order id")
		return
	}

	res := c.service.GetOrder(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) GetOrderMeta(ctx *gin.Context) {
	res := c.service.GetOrderMeta()
	utils.WriteResource(ctx, res)
}
