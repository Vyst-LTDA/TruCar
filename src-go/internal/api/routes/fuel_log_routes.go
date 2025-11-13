package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterFuelLogRoutes(handler *api.FuelLogHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/fuel-logs", handler.GetFuelLogs)
		router.POST("/fuel-logs", handler.CreateFuelLog)
		router.GET("/fuel-logs/:id", handler.GetFuelLog)
		router.PUT("/fuel-logs/:id", handler.UpdateFuelLog)
		router.DELETE("/fuel-logs/:id", handler.DeleteFuelLog)
	}
}
