package attributes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/attributes/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

// Controller handles HTTP requests for attributes and attribute values.
type Controller struct {
	service Service
}

// NewController creates a new attributes controller.
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// CreateAttribute handles POST /api/admin/attributes
func (c *Controller) CreateAttribute(ctx *gin.Context) {
	var req requests.AttributeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.CreateAttribute(req)
	utils.WriteResource(ctx, res)
}

// UpdateAttribute handles PUT /api/admin/attributes/:id
func (c *Controller) UpdateAttribute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	var req requests.AttributeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.UpdateAttribute(id, req)
	utils.WriteResource(ctx, res)
}

// ListAttributes handles GET /api/admin/attributes
func (c *Controller) ListAttributes(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListAttributes(pagination)
	utils.WriteResource(ctx, res)
}

// ListDeletedAttributes handles GET /api/admin/attributes/deleted
func (c *Controller) ListDeletedAttributes(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListDeletedAttributes(pagination)
	utils.WriteResource(ctx, res)
}

// GetAttribute handles GET /api/admin/attributes/:id
func (c *Controller) GetAttribute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	res := c.service.GetAttributeByID(id)
	utils.WriteResource(ctx, res)
}

// DeleteAttribute handles DELETE /api/admin/attributes/:id
func (c *Controller) DeleteAttribute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	res := c.service.DeleteAttribute(id)
	utils.WriteResource(ctx, res)
}

// RecoverAttribute handles PUT /api/admin/attributes/:id/recover
func (c *Controller) RecoverAttribute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	res := c.service.RecoverAttribute(id)
	utils.WriteResource(ctx, res)
}

// ListAttributeValues handles GET /api/admin/attributes/:id/values
func (c *Controller) ListAttributeValues(ctx *gin.Context) {
	idParam := ctx.Param("id")
	attributeID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListAttributeValues(attributeID, pagination)
	utils.WriteResource(ctx, res)
}

// ListDeletedAttributeValues handles GET /api/admin/attributes/:id/values/deleted
func (c *Controller) ListDeletedAttributeValues(ctx *gin.Context) {
	idParam := ctx.Param("id")
	attributeID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListDeletedAttributeValues(attributeID, pagination)
	utils.WriteResource(ctx, res)
}

// RecoverAttributeValue handles PUT /api/admin/attributes/:id/values/:valueId/recover
func (c *Controller) RecoverAttributeValue(ctx *gin.Context) {
	idParam := ctx.Param("id")
	attributeID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	valueParam := ctx.Param("valueId")
	valueID, err := strconv.ParseInt(valueParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid value id")
		return
	}

	var req requests.AttributeValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.RecoverAttributeValue(attributeID, valueID)
	utils.WriteResource(ctx, res)
}

// DeleteAttributeValue handles DELETE /api/admin/attributes/:id/values/:valueId
func (c *Controller) DeleteAttributeValue(ctx *gin.Context) {
	idParam := ctx.Param("id")
	attributeID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	valueParam := ctx.Param("valueId")
	valueID, err := strconv.ParseInt(valueParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid value id")
		return
	}

	res := c.service.DeleteAttributeValue(attributeID, valueID)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ListAttributesForDropdown(ctx *gin.Context) {
	res := c.service.ListAttributesForDropdown()
	utils.WriteResource(ctx, res)
}
