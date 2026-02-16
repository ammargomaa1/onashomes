package inventory

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/inventory/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) AdjustInventory(c *gin.Context) {
	var req requests.AdjustInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// Get admin ID from auth context
	adminID, exists := c.Get("entity_id")
	if !exists {
		utils.ErrorResponse(c, 401, "Unauthorized", nil)
		return
	}

	res := ctrl.service.AdjustInventory(req, adminID.(int64))
	utils.WriteResource(c, res)
}

func (ctrl *Controller) BulkUpdateInventory(c *gin.Context) {
	var req requests.BulkInventoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// Get admin ID from auth context
	adminID, exists := c.Get("entity_id")
	if !exists {
		utils.ErrorResponse(c, 401, "Unauthorized", nil)
		return
	}

	res := ctrl.service.BulkAdjustInventory(req, adminID.(int64))
	utils.WriteResource(c, res)
}

func (ctrl *Controller) ListInventoryByStore(c *gin.Context) {
	storeFrontID, err := strconv.ParseInt(c.Param("storeFrontId"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid store front id")
		return
	}

	var req requests.InventoryFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, "invalid filter parameters")
		return
	}
	req.StoreFrontID = storeFrontID

	pagination := utils.ParsePaginationParams(c)
	res := ctrl.service.ListInventory(req, pagination)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) GetVariantInventory(c *gin.Context) {
	variantID, err := strconv.ParseInt(c.Param("variantId"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid variant id")
		return
	}

	storeFrontID, err := strconv.ParseInt(c.Param("storeFrontId"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid store front id")
		return
	}

	res := ctrl.service.GetVariantInventory(variantID, storeFrontID)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) GetLowStockAlerts(c *gin.Context) {
	storeFrontID, err := strconv.ParseInt(c.Param("storeFrontId"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid store front id")
		return
	}

	pagination := utils.ParsePaginationParams(c)
	low := true
	req := requests.InventoryFilterRequest{
		StoreFrontID: storeFrontID,
		LowStockOnly: &low,
	}
	res := ctrl.service.ListInventory(req, pagination)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) GetAdjustmentHistory(c *gin.Context) {
	inventoryID, err := strconv.ParseInt(c.Param("inventoryId"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid inventory id")
		return
	}

	pagination := utils.ParsePaginationParams(c)
	res := ctrl.service.GetAdjustmentHistory(inventoryID, pagination)
	utils.WriteResource(c, res)
}
