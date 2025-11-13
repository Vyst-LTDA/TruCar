package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterJourneyRoutes(handler *api.JourneyHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/journeys", handler.GetJourneys)
		router.POST("/journeys/start", handler.StartJourney)
		router.PUT("/journeys/:id/end", handler.EndJourney)
		router.DELETE("/journeys/:id", handler.DeleteJourney)
	}
}
