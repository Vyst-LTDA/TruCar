package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
	"go-api/internal/middleware"
	"go-api/internal/models"
)

// RegisterFineRoutes registra as rotas de multas.
func RegisterFineRoutes(handler *api.FineHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		finesGroup := router.Group("/fines")

		// Rotas acessíveis por Motoristas e Gestores
		finesGroup.Use(middleware.AuthorizationMiddleware(models.RoleDriver, models.RoleClienteAtivo, models.RoleClienteDemo))
		{
			finesGroup.POST("/", handler.CreateFine)
			finesGroup.GET("/", handler.GetFines)
		}

		// Rotas acessíveis apenas por Gestores
		managerRoutes := finesGroup.Group("/")
		managerRoutes.Use(middleware.AuthorizationMiddleware(models.RoleClienteAtivo, models.RoleClienteDemo))
		{
			managerRoutes.PUT("/:id", handler.UpdateFine)
			managerRoutes.DELETE("/:id", handler.DeleteFine)
		}
	}
}
