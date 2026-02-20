package locations

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	routes := router.Group("/locations")
	{
		routes.GET("/countries", controller.GetCountries)
		routes.GET("/governorates/:country_id", controller.GetGovernorates)
		routes.GET("/cities/:governorate_id", controller.GetCities)
	}
}
