package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-api/internal/services"
	"go-api/internal/middleware" // Importar para acessar GetOrganizationID
)

// VehicleCostHandler gerencia as requisições HTTP para custos de veículos.
type VehicleCostHandler struct {
	service services.VehicleCostService
}

// NewVehicleCostHandler cria uma nova instância de VehicleCostHandler.
func NewVehicleCostHandler(service services.VehicleCostService) *VehicleCostHandler {
	return &VehicleCostHandler{service: service}
}

// GetCosts lida com a busca de custos da organização.
func (h *VehicleCostHandler) GetCosts(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	// Parsing das datas de filtro da query string
	var startDate, endDate *time.Time
	if dateStr := c.Query("start_date"); dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de start_date inválido. Use AAAA-MM-DD."})
			return
		}
		startDate = &parsedDate
	}
	if dateStr := c.Query("end_date"); dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de end_date inválido. Use AAAA-MM-DD."})
			return
		}
		endDate = &parsedDate
	}

	costs, err := h.service.GetCostsByOrganization(orgID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar custos"})
		return
	}

	c.JSON(http.StatusOK, costs)
}
