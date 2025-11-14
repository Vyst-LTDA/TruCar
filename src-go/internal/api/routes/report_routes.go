package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

// RegisterReportRoutes registra as rotas de relatórios.
func RegisterReportRoutes(handler *api.ReportHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		reports := router.Group("/reports")
		{
			reports.GET("/dashboard-summary", handler.GetDashboardSummary)
			// Adicionar outras rotas de relatório aqui
		}
	}
}
