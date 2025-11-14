package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-api/internal/middleware"
	"go-api/internal/services"
)

// ReportHandler gerencia as requisições HTTP para relatórios.
type ReportHandler struct {
	service services.ReportService
}

// NewReportHandler cria uma nova instância de ReportHandler.
func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetDashboardSummary lida com a busca do resumo de dados do dashboard.
func (h *ReportHandler) GetDashboardSummary(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	summary, err := h.service.GetDashboardSummary(orgID, thirtyDaysAgo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
