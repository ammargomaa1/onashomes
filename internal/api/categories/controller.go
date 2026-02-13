package categories

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/categories/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) ListCategories(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListCategories(pagination)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ListCategoriesBySection(ctx *gin.Context) {
	sectionIDParam := ctx.Param("section_id")
	sectionID, err := strconv.ParseInt(sectionIDParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid section id")
		return
	}

	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListCategoriesBySection(sectionID, pagination)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ListDeletedCategories(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListDeletedCategories(pagination)
	utils.WriteResource(ctx, res)
}

func (c *Controller) GetCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid category id")
		return
	}

	res := c.service.GetCategoryByID(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) CreateCategory(ctx *gin.Context) {
	var req requests.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.CreateCategory(req)
	utils.WriteResource(ctx, res)
}

func (c *Controller) UpdateCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid category id")
		return
	}

	var req requests.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := c.service.UpdateCategory(id, req)
	utils.WriteResource(ctx, res)
}

func (c *Controller) DeleteCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid category id")
		return
	}

	res := c.service.DeleteCategory(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) RestoreCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid category id")
		return
	}

	res := c.service.RestoreCategory(id)
	utils.WriteResource(ctx, res)
}

func (c *Controller) ListCategoriesForDropdown(ctx *gin.Context) {
	res := c.service.ListCategoriesForDropdown()
	utils.WriteResource(ctx, res)
}
