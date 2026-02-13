package products

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/products/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

// ControllerV2 handles Phase 2 product endpoints
type ControllerV2 struct {
	service *ServiceV2
}

func NewControllerV2(service *ServiceV2) *ControllerV2 {
	return &ControllerV2{service: service}
}

// AdminCreateProductV2 creates a product with V2 fields
func (ctrl *ControllerV2) AdminCreateProductV2(ctx *gin.Context) {
	var req requests.CreateProductV2Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := ctrl.service.CreateProductV2(req)
	utils.WriteResource(ctx, res)
}

// AdminUpdateProductV2 updates a product with V2 fields
func (ctrl *ControllerV2) AdminUpdateProductV2(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	var req requests.UpdateProductV2Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := ctrl.service.UpdateProductV2(id, req)
	utils.WriteResource(ctx, res)
}

// AdminUpdateProductStatus changes product status
func (ctrl *ControllerV2) AdminUpdateProductStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	var req requests.UpdateProductStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := ctrl.service.UpdateProductStatus(id, req)
	utils.WriteResource(ctx, res)
}

// AdminUpdateProductSEO updates product SEO metadata
func (ctrl *ControllerV2) AdminUpdateProductSEO(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	var req requests.ProductSEORequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := ctrl.service.UpdateProductSEO(id, req)
	utils.WriteResource(ctx, res)
}

// AdminSuggestSEO auto-generates SEO suggestion
func (ctrl *ControllerV2) AdminSuggestSEO(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	res := ctrl.service.SuggestSEO(id)
	utils.WriteResource(ctx, res)
}

// AdminListProductsV2 lists products with V2 filtering
func (ctrl *ControllerV2) AdminListProductsV2(ctx *gin.Context) {
	var filter requests.ProductFilterRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	pagination := utils.ParsePaginationParams(ctx)
	res := ctrl.service.ListAdminProductsV2(filter, pagination)
	utils.WriteResource(ctx, res)
}

// AdminGetProductV2 gets a product with V2 detail
func (ctrl *ControllerV2) AdminGetProductV2(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	res := ctrl.service.GetProductV2ByID(id)
	utils.WriteResource(ctx, res)
}

// AdminCreateVariantV2 creates a variant with V2 fields
func (ctrl *ControllerV2) AdminCreateVariantV2(ctx *gin.Context) {
	productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	var req requests.CreateVariantV2Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := ctrl.service.CreateVariantV2(productID, req)
	utils.WriteResource(ctx, res)
}

// AdminUpdateVariantV2 updates a variant with V2 fields
func (ctrl *ControllerV2) AdminUpdateVariantV2(ctx *gin.Context) {
	productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	variantID, err := strconv.ParseInt(ctx.Param("variantId"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid variant id")
		return
	}

	var req requests.CreateVariantV2Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := ctrl.service.UpdateVariantV2(productID, variantID, req)
	utils.WriteResource(ctx, res)
}

// AdminAddProductImages adds images to a product
func (ctrl *ControllerV2) AdminAddProductImages(ctx *gin.Context) {
	productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	var req requests.AddProductImagesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := ctrl.service.AddProductImages(productID, req)
	utils.WriteResource(ctx, res)
}

// AdminRemoveProductImage removes an image from a product
func (ctrl *ControllerV2) AdminRemoveProductImage(ctx *gin.Context) {
	productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	imageID, err := strconv.ParseInt(ctx.Param("imageId"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid image id")
		return
	}

	res := ctrl.service.RemoveProductImage(productID, imageID)
	utils.WriteResource(ctx, res)
}

// AdminSetCoverImage sets the cover image for a product
func (ctrl *ControllerV2) AdminSetCoverImage(ctx *gin.Context) {
	productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(ctx, "invalid product id")
		return
	}

	var req requests.SetCoverImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	res := ctrl.service.SetCoverImage(productID, req)
	utils.WriteResource(ctx, res)
}

// StorefrontListProducts lists products for a resolved store
func (ctrl *ControllerV2) StorefrontListProducts(ctx *gin.Context) {
	sfID, exists := ctx.Get("store_front_id")
	if !exists {
		utils.ErrorResponse(ctx, 400, "Store not resolved", nil)
		return
	}

	var filter requests.ProductFilterRequest
	_ = ctx.ShouldBindQuery(&filter)

	pagination := utils.ParsePaginationParams(ctx)
	res := ctrl.service.StorefrontListProducts(sfID.(int64), filter, pagination)
	utils.WriteResource(ctx, res)
}

// StorefrontGetProduct gets a product by slug for a resolved store
func (ctrl *ControllerV2) StorefrontGetProduct(ctx *gin.Context) {
	sfID, exists := ctx.Get("store_front_id")
	if !exists {
		utils.ErrorResponse(ctx, 400, "Store not resolved", nil)
		return
	}

	slug := ctx.Param("slug")
	if slug == "" {
		utils.ValidationErrorResponse(ctx, "slug is required")
		return
	}

	res := ctrl.service.StorefrontGetProduct(sfID.(int64), slug)
	utils.WriteResource(ctx, res)
}

// StorefrontGetStructuredData returns JSON-LD structured data
func (ctrl *ControllerV2) StorefrontGetStructuredData(ctx *gin.Context) {
	sfID, exists := ctx.Get("store_front_id")
	if !exists {
		utils.ErrorResponse(ctx, 400, "Store not resolved", nil)
		return
	}

	slug := ctx.Param("slug")
	if slug == "" {
		utils.ValidationErrorResponse(ctx, "slug is required")
		return
	}

	res := ctrl.service.GetProductStructuredData(sfID.(int64), slug)
	utils.WriteResource(ctx, res)
}
