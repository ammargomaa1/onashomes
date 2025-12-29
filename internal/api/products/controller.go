package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

type AdminProductAttributeRequest struct {
	AttributeID     int64   `json:"attribute_id" binding:"required"`
	SortOrder       int     `json:"sort_order"`
	AllowedValueIDs []int64 `json:"allowed_value_ids"`
}

type AdminProductRequest struct {
	Name        string                         `json:"name" binding:"required"`
	Description string                         `json:"description"`
	Price       float64                        `json:"price" binding:"required"`
	IsActive    bool                           `json:"is_active"`
	Attributes  []AdminProductAttributeRequest `json:"attributes"`
}

// AdminCreateProduct handles product creation
func (c *Controller) AdminCreateProduct(ctx *gin.Context) {
	var req AdminProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	product, err := c.service.CreateProduct(ctx.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Product created successfully", product)
}

// AdminUpdateProduct handles product update
func (c *Controller) AdminUpdateProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	var req AdminProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	product, err := c.service.UpdateProduct(ctx.Request.Context(), id, req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Product updated successfully", product)
}

// AdminListProducts lists products for admin
func (c *Controller) AdminListProducts(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	products, total, err := c.service.ListProducts(ctx.Request.Context(), pagination)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve products", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(ctx, http.StatusOK, "Products retrieved successfully", products, pagination.GetMeta())
}

// AdminGetProduct returns a single product with variants and configuration
func (c *Controller) AdminGetProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	product, err := c.service.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Product not found", nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Product retrieved successfully", product)
}

// PublicListProducts lists active products for public consumers
func (c *Controller) PublicListProducts(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)
	q := ctx.Query("q")

	products, total, err := c.service.PublicListProducts(ctx.Request.Context(), pagination, q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve products", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(ctx, http.StatusOK, "Products retrieved successfully", products, pagination.GetMeta())
}

// PublicGetProduct returns a single active product with active variants for public consumers
func (c *Controller) PublicGetProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	product, err := c.service.PublicGetProduct(ctx.Request.Context(), id)
	if err != nil {
		// For public API, treat any error as not found to avoid leaking internals
		utils.ErrorResponse(ctx, http.StatusNotFound, "Product not found", nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Product retrieved successfully", product)
}

// AdminDeleteProduct performs a soft delete
func (c *Controller) AdminDeleteProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	if err := c.service.DeleteProduct(ctx.Request.Context(), id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// AdminCreateVariant creates a variant for a product
func (c *Controller) AdminCreateVariant(ctx *gin.Context) {
	idParam := ctx.Param("id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	var req AdminVariantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	variant, err := c.service.CreateVariant(ctx.Request.Context(), productID, req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Variant created successfully", variant)
}

// AdminUpdateVariant updates an existing variant
func (c *Controller) AdminUpdateVariant(ctx *gin.Context) {
	idParam := ctx.Param("id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	variantParam := ctx.Param("variantId")
	variantID, err := strconv.ParseInt(variantParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid variant id")
		return
	}

	var req AdminVariantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	variant, err := c.service.UpdateVariant(ctx.Request.Context(), productID, variantID, req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Variant updated successfully", variant)
}
