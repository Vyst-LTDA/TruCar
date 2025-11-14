package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
	"go-api/internal/middleware"
	"go-api/internal/models"
)

// RegisterDashboardRoutes registra as rotas do dashboard.
func RegisterDashboardRoutes(handler *api.DashboardHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		// Grupo para rotas de gestor (CLIENTE_ATIVO, CLIENTE_DEMO)
		managerDashboardRoutes := router.Group("/dashboard")
		managerDashboardRoutes.Use(middleware.AuthorizationMiddleware(models.RoleClienteAtivo, models.RoleClienteDemo))
		{
			managerDashboardRoutes.GET("/manager", handler.GetManagerDashboard)
			managerDashboardRoutes.GET("/vehicles/positions", handler.GetVehiclePositions)
		}

		// Rota específica para CLIENTE_DEMO
		demoRoutes := router.Group("/dashboard")
		demoRoutes.Use(middleware.AuthorizationMiddleware(models.RoleClienteDemo))
		{
			demoRoutes.GET("/demo-stats", handler.GetDemoStats)
		}

		// Rota específica para motorista (DRIVER)
		driverDashboardRoutes := router.Group("/dashboard")
		driverDashboardRoutes.Use(middleware.AuthorizationMiddleware(models.RoleDriver))
		{
			driverDashboardRoutes.GET("/driver", handler.GetDriverDashboard)
		}
	}
}
