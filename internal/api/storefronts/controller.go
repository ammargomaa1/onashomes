package storefronts

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/storefronts/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) Create(c *gin.Context) {
	var req requests.CreateStoreFrontRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.Create(req)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid store front id")
		return
	}

	var req requests.UpdateStoreFrontRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	res := ctrl.service.Update(id, req)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid store front id")
		return
	}

	res := ctrl.service.Delete(id)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationErrorResponse(c, "invalid store front id")
		return
	}

	res := ctrl.service.GetByID(id)
	utils.WriteResource(c, res)
}

func (ctrl *Controller) List(c *gin.Context) {
	pagination := utils.ParsePaginationParams(c)
	res := ctrl.service.List(pagination)
	utils.WriteResource(c, res)
}
