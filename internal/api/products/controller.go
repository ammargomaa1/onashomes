package products

import (
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

	res := c.service.CreateProduct(ctx.Request.Context(), req)
	utils.WriteResource(ctx, res)
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

	res := c.service.UpdateProduct(ctx.Request.Context(), id, req)
	utils.WriteResource(ctx, res)
}

// AdminListProducts lists products for admin
func (c *Controller) AdminListProducts(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)

	res := c.service.ListProducts(ctx.Request.Context(), pagination)
	utils.WriteResource(ctx, res)
}

// AdminGetProduct returns a single product with variants and configuration
func (c *Controller) AdminGetProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	res := c.service.GetProductByID(ctx.Request.Context(), id)
	utils.WriteResource(ctx, res)
}

// PublicListProducts lists active products for public consumers
func (c *Controller) PublicListProducts(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)
	q := ctx.Query("q")

	res := c.service.PublicListProducts(ctx.Request.Context(), pagination, q)
	utils.WriteResource(ctx, res)
}

// PublicGetProduct returns a single active product with active variants for public consumers
func (c *Controller) PublicGetProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	res := c.service.PublicGetProduct(ctx.Request.Context(), id)
	utils.WriteResource(ctx, res)
}

// AdminDeleteProduct performs a soft delete
func (c *Controller) AdminDeleteProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	res := c.service.DeleteProduct(ctx.Request.Context(), id)
	utils.WriteResource(ctx, res)
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

	res := c.service.CreateVariant(ctx.Request.Context(), productID, req)
	utils.WriteResource(ctx, res)
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

	res := c.service.UpdateVariant(ctx.Request.Context(), productID, variantID, req)
	utils.WriteResource(ctx, res)
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

	res := c.service.ListVariantAddOns(ctx.Request.Context(), productID, variantID)
	utils.WriteResource(ctx, res)
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

	res := c.service.CreateVariantAddOn(ctx.Request.Context(), productID, variantID, req.AddOnProductID)
	utils.WriteResource(ctx, res)
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

	res := c.service.UpdateVariantAddOn(ctx.Request.Context(), productID, variantID, addOnID, req.AddOnProductID)
	utils.WriteResource(ctx, res)
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

	res := c.service.DeleteVariantAddOn(ctx.Request.Context(), productID, variantID, addOnID)
	utils.WriteResource(ctx, res)
}
