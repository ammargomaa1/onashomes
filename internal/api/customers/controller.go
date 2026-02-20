package customers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

// List handles GET /admin/customers
func (c *Controller) List(ctx *gin.Context) {
	pagination := utils.ParsePaginationParams(ctx)
	filter := ctx.Query("search")
	var storeFrontID *int64
	if idStr := ctx.Query("store_front_id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			storeFrontID = &id
		}
	}

	customers, meta, err := c.service.ListCustomers(pagination, filter, storeFrontID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch customers", err.Error())
		return
	}

	response := utils.Response{
		Success: true,
		Message: "Customers retrieved successfully",
		Data:    customers,
		Meta:    meta,
	}
	ctx.JSON(http.StatusOK, response)
}

// Search handles GET /admin/customers/search?q=...
func (c *Controller) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		utils.SuccessResponse(ctx, http.StatusOK, "Empty query", []models.Customer{})
		return
	}

	var storeFrontID *int64
	if idStr := ctx.Query("store_front_id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			storeFrontID = &id
		}
	}

	customers, err := c.service.SearchCustomers(query, storeFrontID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to search customers", err.Error())
		return
	}

	// Map to minimal DTO if needed, but for now returning a list of model.Customer is sufficient
	utils.SuccessResponse(ctx, http.StatusOK, "Customers retrieved", customers)
}

// Create handles POST /admin/customers
func (c *Controller) Create(ctx *gin.Context) {
	var input struct {
		FirstName    string `json:"first_name" binding:"required"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		Phone        string `json:"phone" binding:"required"`
		StoreFrontID int64  `json:"store_front_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	customer := models.Customer{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		Phone:        input.Phone,
		StoreFrontID: input.StoreFrontID,
	}

	if err := c.service.CreateCustomer(&customer); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create customer", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Customer created successfully", customer)
}

// Get handles GET /admin/customers/:id
func (c *Controller) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid customer ID", err.Error())
		return
	}

	customer, err := c.service.GetCustomer(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Customer not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Customer retrieved successfully", customer)
}

// Update handles PUT /admin/customers/:id
func (c *Controller) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid customer ID", err.Error())
		return
	}

	var input models.Customer
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	if err := c.service.UpdateCustomer(id, &input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update customer", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Customer updated successfully", nil)
}
