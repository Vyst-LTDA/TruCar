package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api/internal/middleware"
	"go-api/internal/models"
	"go-api/internal/services"
)

// DashboardHandler gerencia as requisições HTTP para o dashboard.
type DashboardHandler struct {
	service services.DashboardService
}

// NewDashboardHandler cria uma nova instância de DashboardHandler.
func NewDashboardHandler(service services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetManagerDashboard lida com a busca de dados para o dashboard do gestor.
func (h *DashboardHandler) GetManagerDashboard(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	period := c.DefaultQuery("period", "last_30_days")

	data, err := h.service.GetManagerDashboard(orgID, period, currentUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados do dashboard"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetDriverDashboard lida com a busca de dados para o dashboard do motorista.
func (h *DashboardHandler) GetDriverDashboard(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	data, err := h.service.GetDriverDashboard(currentUser.ID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados do dashboard do motorista"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetVehiclePositions lida com a busca das posições dos veículos.
func (h *DashboardHandler) GetVehiclePositions(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	data, err := h.service.GetVehiclePositions(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar posições dos veículos"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetDemoStats lida com a busca de estatísticas da conta demo.
func (h *DashboardHandler) GetDemoStats(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	data, err := h.service.GetDemoStats(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar estatísticas da conta demo"})
		return
	}

	c.JSON(http.StatusOK, data)
}
