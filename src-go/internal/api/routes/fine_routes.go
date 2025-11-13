package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterFineRoutes(handler *api.FineHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/fines", handler.GetFines)
		router.POST("/fines", handler.CreateFine)
		router.GET("/fines/:id", handler.GetFine)
		router.PUT("/fines/:id", handler.UpdateFine)
		router.DELETE("/fines/:id", handler.DeleteFine)
	}
}
