package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterMaintenanceRoutes(handler *api.MaintenanceHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/maintenance", handler.GetMaintenanceRequests)
		router.POST("/maintenance", handler.CreateMaintenanceRequest)
		router.GET("/maintenance/:id", handler.GetMaintenanceRequest)
		router.PUT("/maintenance/:id/status", handler.UpdateMaintenanceRequestStatus)
		router.DELETE("/maintenance/:id", handler.DeleteMaintenanceRequest)
		router.GET("/maintenance/:id/comments", handler.GetMaintenanceComments)
		router.POST("/maintenance/:id/comments", handler.CreateMaintenanceComment)
	}
}
