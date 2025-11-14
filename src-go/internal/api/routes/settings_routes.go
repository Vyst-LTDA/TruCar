package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
	"go-api/internal/middleware"
	"go-api/internal/models"
)

// RegisterSettingsRoutes registra as rotas de configurações.
func RegisterSettingsRoutes(handler *api.SettingsHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		settings := router.Group("/settings")
		settings.Use(middleware.AuthorizationMiddleware(models.RoleClienteAtivo, models.RoleClienteDemo))
		{
			settings.GET("/fuel-integration", handler.GetFuelIntegrationSettings)
			settings.PUT("/fuel-integration", handler.UpdateFuelIntegrationSettings)
		}
	}
}
