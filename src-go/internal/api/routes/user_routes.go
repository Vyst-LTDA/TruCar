package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterUserRoutes(handler *api.UserHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/users", handler.GetUsers)
		router.POST("/users", handler.CreateUser)
		router.GET("/users/:id", handler.GetUser)
		router.PUT("/users/:id", handler.UpdateUser)
		router.DELETE("/users/:id", handler.DeleteUser)
	}
}
