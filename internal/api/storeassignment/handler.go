package storeassignment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/middleware"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

// AssignStoresRequest is the shared request body for all store assignment endpoints
type AssignStoresRequest struct {
	StoreFrontIDs []int64 `json:"store_front_ids" binding:"required"`
}

// Handler provides generic store assignment endpoints for any entity type.
// It operates directly on the pivot tables and doesn't need entity-specific logic.
type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// AssignStores replaces the store assignments for an entity.
// pivotTable: e.g. "brand_storefront"
// entityCol: e.g. "brand_id"
func (h *Handler) AssignStores(pivotTable, entityCol string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		entityID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			utils.ValidationErrorResponse(ctx, "invalid id")
			return
		}

		var req AssignStoresRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.ValidationErrorResponse(ctx, err.Error())
			return
		}

		// Validate that all store front IDs exist
		var count int64
		h.db.Model(&models.StoreFront{}).Where("id IN ? AND is_active = true", req.StoreFrontIDs).Count(&count)
		if count != int64(len(req.StoreFrontIDs)) {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "One or more store front IDs are invalid or inactive", nil)
			return
		}

		tx := h.db.Begin()

		// Delete existing assignments
		if err := tx.Exec("DELETE FROM "+pivotTable+" WHERE "+entityCol+" = ?", entityID).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update store assignments", err)
			return
		}

		// Insert new assignments
		if len(req.StoreFrontIDs) > 0 {
			for _, sfID := range req.StoreFrontIDs {
				if err := tx.Exec(
					"INSERT INTO "+pivotTable+" ("+entityCol+", store_front_id) VALUES (?, ?)",
					entityID, sfID,
				).Error; err != nil {
					tx.Rollback()
					utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to assign stores", err)
					return
				}
			}
		}

		if err := tx.Commit().Error; err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit store assignments", err)
			return
		}

		// Return updated assignments
		var storeIDs []int64
		h.db.Table(pivotTable).Where(entityCol+" = ?", entityID).Pluck("store_front_id", &storeIDs)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Store assignments updated",
			"data":    gin.H{"store_front_ids": storeIDs},
		})
	}
}

// GetStores returns the store assignments for a specific entity
func (h *Handler) GetStores(pivotTable, entityCol string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		entityID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			utils.ValidationErrorResponse(ctx, "invalid id")
			return
		}

		var stores []models.StoreFront
		h.db.Table("store_fronts").
			Joins("INNER JOIN "+pivotTable+" ON store_fronts.id = "+pivotTable+".store_front_id").
			Where(pivotTable+"."+entityCol+" = ?", entityID).
			Find(&stores)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Store assignments retrieved",
			"data":    stores,
		})
	}
}

// ListByStore returns entities assigned to the resolved storefront.
// entityTable: e.g. "brands"
// entityCol: e.g. "brand_id"
func (h *Handler) ListByStore(pivotTable, entityTable, entityCol string, dest interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sfID, exists := ctx.Get("store_front_id")
		if !exists {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Store not resolved", nil)
			return
		}

		h.db.Table(entityTable).
			Joins("INNER JOIN "+pivotTable+" ON "+entityTable+".id = "+pivotTable+"."+entityCol).
			Where(pivotTable+".store_front_id = ?", sfID).
			Find(dest)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": entityTable + " retrieved for store",
			"data":    dest,
		})
	}
}

// RegisterEntityStoreRoutes registers store assignment routes for a specific entity.
// This is a convenience method that sets up admin + storefront routes.
func (h *Handler) RegisterEntityStoreRoutes(
	adminGroup *gin.RouterGroup,
	storefrontGroup *gin.RouterGroup,
	entityPath string, // e.g. "brands"
	pivotTable string, // e.g. "brand_storefront"
	entityCol string, // e.g. "brand_id"
	permission string, // e.g. "brands.update"
) {
	// Admin: assign stores to entity
	adminGroup.PUT("/admin/"+entityPath+"/:id/stores",
		middleware.AuthMiddleware(),
		middleware.AdminAuthMiddleware(),
		middleware.RequirePermission(permission),
		h.AssignStores(pivotTable, entityCol),
	)

	// Admin: get entity store assignments
	adminGroup.GET("/admin/"+entityPath+"/:id/stores",
		middleware.AuthMiddleware(),
		middleware.AdminAuthMiddleware(),
		middleware.RequirePermission(permission),
		h.GetStores(pivotTable, entityCol),
	)
}
