package sections

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/sections/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) ListSections(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListSections(pagination)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ListDeletedSections(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListDeletedSections(pagination)
	utils.WriteResource(ctx, res)
}

func (c *Controller) GetSection(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid section id")
		return
	}

	res := c.service.GetSectionByID(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) CreateSection(ctx *gin.Context) {
	var req requests.CreateSectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.CreateSection(req)
	utils.WriteResource(ctx, res)
}

func (c *Controller) UpdateSection(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid section id")
		return
	}

	var req requests.CreateSectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.UpdateSection(id, req)
	utils.WriteResource(ctx, res)
}

func (c *Controller) DeleteSection(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid section id")
		return
	}

	res := c.service.DeleteSection(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) RestoreSection(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid section id")
		return
	}

	res := c.service.RestoreSection(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ListSectionsForDropdown(ctx *gin.Context) {
	res := c.service.ListSectionsForDropdown()
	utils.WriteResource(ctx, res)
}
