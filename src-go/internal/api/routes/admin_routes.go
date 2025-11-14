package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterAdminRoutes(handler *api.AdminHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/organizations", handler.GetOrganizations)
		router.PUT("/organizations/:id", handler.UpdateOrganization)
		router.GET("/users/all", handler.GetAllUsers)
		router.GET("/users/demo", handler.GetDemoUsers)
		router.POST("/users/:id/activate", handler.ActivateUser)
		router.POST("/users/:id/impersonate", handler.ImpersonateUser)
	}
}
