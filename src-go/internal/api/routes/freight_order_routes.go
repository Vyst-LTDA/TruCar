package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterFreightOrderRoutes(handler *api.FreightOrderHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/freight-orders", handler.GetFreightOrders)
		router.POST("/freight-orders", handler.CreateFreightOrder)
		router.GET("/freight-orders/open", handler.GetOpenFreightOrders)
		router.GET("/freight-orders/my-pending", handler.GetMyPendingFreightOrders)
		router.GET("/freight-orders/:id", handler.GetFreightOrderByID)
		router.PUT("/freight-orders/:id/claim", handler.ClaimFreightOrder)
		router.POST("/freight-orders/:order_id/start-leg/:stop_point_id", handler.StartJourneyForStop)
		router.PUT("/freight-orders/:order_id/complete-stop/:stop_point_id", handler.CompleteStopPoint)
	}
}
