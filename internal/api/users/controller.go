package users

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (ctrl *Controller) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := ctrl.service.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, 201, "User registered successfully", user)
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

func (ctrl *Controller) GetProfile(c *gin.Context) {
	entityID, _ := c.Get("entity_id")
	userID := entityID.(uuid.UUID)

	user, err := ctrl.service.GetProfile(userID)
	if err != nil {
		utils.ErrorResponse(c, 404, "User not found", nil)
		return
	}

	utils.SuccessResponse(c, 200, "Profile retrieved successfully", user)
}

func (ctrl *Controller) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	entityID, _ := c.Get("entity_id")
	userID := entityID.(uuid.UUID)

	user, err := ctrl.service.UpdateProfile(userID, req.FirstName, req.LastName)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, 200, "Profile updated successfully", user)
}

func (ctrl *Controller) List(c *gin.Context) {
	pagination := utils.ParsePaginationParams(c)

	users, total, err := ctrl.service.List(pagination)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to retrieve users", err.Error())
		return
	}

	pagination.SetTotal(total)
	utils.SuccessResponseWithMeta(c, 200, "Users retrieved successfully", users, pagination.GetMeta())
}
