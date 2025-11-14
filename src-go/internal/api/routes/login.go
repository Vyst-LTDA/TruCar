package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterLoginRoutes(handler *api.AuthHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.POST("/login/access-token", handler.Login)
	}
}
