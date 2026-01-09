package suppliers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	suppreq "github.com/onas/ecommerce-api/internal/api/suppliers/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req suppreq.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, exists := c.Get("entity_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	adminIDInt64, ok := adminID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid admin id"})
		return
	}
	supplier, err := h.service.Create(req.CompanyName, req.ContactPersonName, req.ContactNumber, adminIDInt64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		return
	}

	var req suppreq.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, exists := c.Get("entity_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	adminIDInt64, ok := adminID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid admin id"})
		return
	}
	supplier, err := h.service.Update(id, req.CompanyName, req.ContactPersonName, req.ContactNumber, adminIDInt64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		return
	}

	supplier, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

func (h *Handler) List(c *gin.Context) {
	pagination := utils.ParsePaginationParams(c)

	suppliers, total, err := h.service.List(pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve suppliers", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(c, http.StatusOK, "Suppliers retrieved successfully", suppliers, pagination.GetMeta())
}

func (h *Handler) Activate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		return
	}

	if err := h.service.Activate(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		return
	}


	if err := h.service.Deactivate(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
