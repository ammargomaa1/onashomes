package products

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/products/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// AdminCreateProduct handles product creation
func (c *Controller) AdminCreateProduct(ctx *gin.Context) {
	var req requests.AdminProductRequest
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

	var req requests.AdminProductRequest
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

	var req requests.AdminVariantRequest
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

	var req requests.AdminVariantRequest
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

// AdminListVariantAddOns lists add-ons for a specific variant
func (c *Controller) AdminListVariantAddOns(ctx *gin.Context) {
	productParam := ctx.Param("id")
	productID, err := strconv.ParseInt(productParam, 10, 64)
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

	addOns, err := c.service.ListVariantAddOns(ctx.Request.Context(), productID, variantID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Product or variant not found", nil)
		return
	}
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Variant add-ons retrieved successfully", addOns)
}

// AdminCreateVariantAddOn creates a new add-on mapping for a variant
func (c *Controller) AdminCreateVariantAddOn(ctx *gin.Context) {
	productParam := ctx.Param("id")
	productID, err := strconv.ParseInt(productParam, 10, 64)
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

	var req requests.VariantAddOnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	addOn, err := c.service.CreateVariantAddOn(ctx.Request.Context(), productID, variantID, req.AddOnProductID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Product, variant or add-on product not found", nil)
		return
	}
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Variant add-on created successfully", addOn)
}

// AdminUpdateVariantAddOn updates an existing add-on mapping for a variant
func (c *Controller) AdminUpdateVariantAddOn(ctx *gin.Context) {
	productParam := ctx.Param("id")
	productID, err := strconv.ParseInt(productParam, 10, 64)
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

	addOnParam := ctx.Param("addOnId")
	addOnID, err := strconv.ParseInt(addOnParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid add-on id")
		return
	}

	var req requests.VariantAddOnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	addOn, err := c.service.UpdateVariantAddOn(ctx.Request.Context(), productID, variantID, addOnID, req.AddOnProductID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Product, variant or add-on not found", nil)
		return
	}
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Variant add-on updated successfully", addOn)
}

// AdminDeleteVariantAddOn deletes an add-on mapping from a variant
func (c *Controller) AdminDeleteVariantAddOn(ctx *gin.Context) {
	productParam := ctx.Param("id")
	productID, err := strconv.ParseInt(productParam, 10, 64)
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

	addOnParam := ctx.Param("addOnId")
	addOnID, err := strconv.ParseInt(addOnParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid add-on id")
		return
	}

	if err := c.service.DeleteVariantAddOn(ctx.Request.Context(), productID, variantID, addOnID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Product, variant or add-on not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.Status(http.StatusNoContent)
}
