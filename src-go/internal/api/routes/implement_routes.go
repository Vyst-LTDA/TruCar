package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterImplementRoutes(handler *api.ImplementHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.POST("/implements", handler.CreateImplement)
		router.GET("/implements", handler.GetImplements)
		router.GET("/implements/:id", handler.GetImplement)
		router.PUT("/implements/:id", handler.UpdateImplement)
		router.DELETE("/implements/:id", handler.DeleteImplement)
	}
}
