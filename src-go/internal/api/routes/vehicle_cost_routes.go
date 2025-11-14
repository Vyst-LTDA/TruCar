package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

// RegisterVehicleCostRoutes registra as rotas de custos de veículos.
func RegisterVehicleCostRoutes(handler *api.VehicleCostHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/costs", handler.GetCosts)
	}
}
