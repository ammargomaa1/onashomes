package brand

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/brand/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) ListBrands(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListBrands(pagination)
	utils.WriteResource(ctx, res)
}

func (c *Controller) GetBrand(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid brand id")
		return
	}

	res := c.service.GetBrandByID(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) CreateBrand(ctx *gin.Context) {
	// Implementation for creating a brand
	var req requests.CreateBrandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.CreateBrand(req)
	utils.WriteResource(ctx, res)
}

func (c *Controller) UpdateBrand(ctx *gin.Context) {
	// Implementation for updating a brand
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid brand id")
		return
	}

	// Assuming similar request structure as CreateBrandRequest
	var req requests.CreateBrandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.UpdateBrand(id, req)
	utils.WriteResource(ctx, res)
}

func (c *Controller) DeleteBrand(ctx *gin.Context) {
	// Implementation for deleting a brand
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid brand id")
		return
	}

	res := c.service.DeleteBrand(id)
	utils.WriteResource(ctx, res)
}


func (c *Controller) RestoreBrand(ctx *gin.Context) {
	// Implementation for restoring a brand
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid brand id")
		return
	}

	res := c.service.RestoreBrand(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ListDeletedBrands(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListDeletedBrands(pagination)
	utils.WriteResource(ctx, res)
}