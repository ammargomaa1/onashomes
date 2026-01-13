package users

import (
	"github.com/gin-gonic/gin"
	userreq "github.com/onas/ecommerce-api/internal/api/users/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) Register(c *gin.Context) {
	var req userreq.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.Register(req.Email, req.Password, req.FirstName, req.LastName)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) Login(c *gin.Context) {
	var req userreq.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.Login(req.Email, req.Password)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) RefreshToken(c *gin.Context) {
	var req userreq.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.RefreshToken(req.RefreshToken)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) GetProfile(c *gin.Context) {
	entityID, _ := c.Get("entity_id")
	userID := entityID.(int64)

	res := ctrl.service.GetProfile(userID)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) UpdateProfile(c *gin.Context) {
	var req userreq.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	entityID, _ := c.Get("entity_id")
	userID := entityID.(int64)

	res := ctrl.service.UpdateProfile(userID, req.FirstName, req.LastName)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) List(c *gin.Context) {
	pagination := utils.ParsePaginationParams(c)

	res := ctrl.service.List(pagination)
	utils.WriteResource(c, res)
}
