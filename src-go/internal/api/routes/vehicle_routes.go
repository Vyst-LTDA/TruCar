package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterVehicleRoutes(handler *api.VehicleHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/vehicles", handler.GetVehicles)
		router.POST("/vehicles", handler.CreateVehicle)
		router.GET("/vehicles/:id", handler.GetVehicle)
		router.PUT("/vehicles/:id", handler.UpdateVehicle)
		router.DELETE("/vehicles/:id", handler.DeleteVehicle)
	}
}
