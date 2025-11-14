package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

// RegisterReportGeneratorRoutes registra as rotas de geração de relatórios.
func RegisterReportGeneratorRoutes(handler *api.ReportGeneratorHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		reports := router.Group("/reports")
		{
			reports.POST("/generate-pdf", handler.GeneratePDF)
		}
	}
}
