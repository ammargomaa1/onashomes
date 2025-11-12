package admins

import (
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

type CreateAdminRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	RoleID    int64  `json:"role_id" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateAdminRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	RoleID    int64  `json:"role_id" binding:"required"`
}

func (ctrl *Controller) Create(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	admin, err := ctrl.service.Create(req.Email, req.Password, req.FirstName, req.LastName, req.RoleID)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, 201, "Admin created successfully", admin)
}

func (ctrl *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	accessToken, refreshToken, err := ctrl.service.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, 401, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, 200, "Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (ctrl *Controller) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	accessToken, err := ctrl.service.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, 401, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, 200, "Token refreshed successfully", gin.H{
		"access_token": accessToken,
	})
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid admin ID", nil)
		return
	}

	admin, err := ctrl.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, 404, "Admin not found", nil)
		return
	}

	utils.SuccessResponse(c, 200, "Admin retrieved successfully", admin)
}

func (ctrl *Controller) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid admin ID", nil)
		return
	}

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	admin, err := ctrl.service.Update(id, req.FirstName, req.LastName, req.RoleID)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, 200, "Admin updated successfully", admin)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid admin ID", nil)
		return
	}

	err = ctrl.service.Delete(id)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, 200, "Admin deleted successfully", nil)
}

func (ctrl *Controller) List(c *gin.Context) {
	pagination := utils.ParsePaginationParams(c)

	admins, total, err := ctrl.service.List(pagination)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to retrieve admins", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(c, 200, "Admins retrieved successfully", admins, pagination.GetMeta())
}
