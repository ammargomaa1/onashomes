package locations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetCountries(ctx *gin.Context) {
	countries, err := c.service.GetCountries()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch countries", nil)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Countries fetched successfully", countries)
}

func (c *Controller) GetGovernorates(ctx *gin.Context) {
	countryIDStr := ctx.Param("country_id")
	countryID, err := strconv.ParseInt(countryIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid country ID", nil)
		return
	}

	governorates, err := c.service.GetGovernorates(countryID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch governorates", nil)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Governorates fetched successfully", governorates)
}

func (c *Controller) GetCities(ctx *gin.Context) {
	governorateIDStr := ctx.Param("governorate_id")
	governorateID, err := strconv.ParseInt(governorateIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid governorate ID", nil)
		return
	}

	cities, err := c.service.GetCities(governorateID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch cities", nil)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Cities fetched successfully", cities)
}
