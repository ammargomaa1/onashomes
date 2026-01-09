package attributes

import (
	"database/sql"
	"errors"
	"net/http"
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

	attr, err := c.service.CreateAttribute(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Attribute created successfully", attr)
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

	attr, err := c.service.UpdateAttribute(id, req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Attribute updated successfully", attr)
}

// ListAttributes handles GET /api/admin/attributes
func (c *Controller) ListAttributes(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	attrs, total, err := c.service.ListAttributes(pagination)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve attributes", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(ctx, http.StatusOK, "Attributes retrieved successfully", attrs, pagination.GetMeta())
}

// ListDeletedAttributes handles GET /api/admin/attributes/deleted
func (c *Controller) ListDeletedAttributes(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	attrs, total, err := c.service.ListDeletedAttributes(pagination)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted attributes", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(ctx, http.StatusOK, "Deleted attributes retrieved successfully", attrs, pagination.GetMeta())
}

// GetAttribute handles GET /api/admin/attributes/:id
func (c *Controller) GetAttribute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	attr, err := c.service.GetAttributeByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Attribute not found", nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Attribute retrieved successfully", attr)
}

// DeleteAttribute handles DELETE /api/admin/attributes/:id
func (c *Controller) DeleteAttribute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	if err := c.service.DeleteAttribute(id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// RecoverAttribute handles PUT /api/admin/attributes/:id/recover
func (c *Controller) RecoverAttribute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	attr, err := c.service.RecoverAttribute(id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Attribute not found", nil)
		return
	}
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Attribute recovered successfully", attr)
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

	values, total, err := c.service.ListAttributeValues(attributeID, pagination)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve attribute values", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(ctx, http.StatusOK, "Attribute values retrieved successfully", values, pagination.GetMeta())
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

	values, total, err := c.service.ListDeletedAttributeValues(attributeID, pagination)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted attribute values", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(ctx, http.StatusOK, "Deleted attribute values retrieved successfully", values, pagination.GetMeta())
}

// CreateAttributeValue handles POST /api/admin/attributes/:id/values
func (c *Controller) CreateAttributeValue(ctx *gin.Context) {
	idParam := ctx.Param("id")
	attributeID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid attribute id")
		return
	}

	var req requests.AttributeValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	value, err := c.service.CreateAttributeValue(attributeID, req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Attribute value created successfully", value)
}

// UpdateAttributeValue handles PUT /api/admin/attributes/:id/values/:valueId
func (c *Controller) UpdateAttributeValue(ctx *gin.Context) {
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

	value, err := c.service.UpdateAttributeValue(attributeID, valueID, req)
	if errors.Is(err, sql.ErrNoRows) {
		utils.ErrorResponse(ctx, 404, "some entities can not be found", nil)
		return
	}

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Attribute value updated successfully", value)
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

	value, err := c.service.RecoverAttributeValue(attributeID, valueID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.ErrorResponse(ctx, http.StatusNotFound, "some entities can not be found", nil)
		return
	}
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Attribute value recovered successfully", value)
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

	if err := c.service.DeleteAttributeValue(attributeID, valueID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.Status(http.StatusNoContent)
}
