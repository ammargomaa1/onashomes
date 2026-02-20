package admins

import (
	"strconv"

	"github.com/gin-gonic/gin"
	adminreq "github.com/onas/ecommerce-api/internal/api/admins/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) Create(c *gin.Context) {
	var req adminreq.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.Create(req.Email, req.Password, req.FirstName, req.LastName, req.RoleID)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) Login(c *gin.Context) {
	var req adminreq.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.Login(req.Email, req.Password)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) RefreshToken(c *gin.Context) {
	var req adminreq.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.RefreshToken(req.RefreshToken)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid admin ID", nil)
		return
	}

	res := ctrl.service.GetByID(id)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid admin ID", nil)
		return
	}

	var req adminreq.UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.Update(id, req.FirstName, req.LastName, req.RoleID)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid admin ID", nil)
		return
	}

	res := ctrl.service.Delete(id)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) List(c *gin.Context) {
	pagination := utils.ParsePaginationParams(c)

	res := ctrl.service.List(pagination)
	utils.WriteResource(c, res)
}
