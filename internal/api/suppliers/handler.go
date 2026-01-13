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
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	adminID, exists := c.Get("entity_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	adminIDInt64, ok := adminID.(int64)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "invalid admin id", nil)
		return
	}

	res := h.service.Create(req.CompanyName, req.ContactPersonName, req.ContactNumber, adminIDInt64)
	utils.WriteResource(c, res)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid supplier id")
		return
	}

	var req suppreq.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, exists := c.Get("entity_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	adminIDInt64, ok := adminID.(int64)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "invalid admin id", nil)
		return
	}

	res := h.service.Update(id, req.CompanyName, req.ContactPersonName, req.ContactNumber, adminIDInt64)
	utils.WriteResource(c, res)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid supplier id")
		return
	}

	res := h.service.GetByID(id)
	utils.WriteResource(c, res)
}

func (h *Handler) List(c *gin.Context) {
	pagination := utils.ParsePaginationParams(c)

	res := h.service.List(pagination)
	utils.WriteResource(c, res)
}

func (h *Handler) Activate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid supplier id")
		return
	}

	res := h.service.Activate(id)
	utils.WriteResource(c, res)
}

func (h *Handler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid supplier id")
		return
	}

	res := h.service.Deactivate(id)
	utils.WriteResource(c, res)
}
