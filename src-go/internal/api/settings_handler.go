package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api/internal/middleware"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

// SettingsHandler gerencia as requisições HTTP para configurações.
type SettingsHandler struct {
	orgService services.OrganizationService
}

// NewSettingsHandler cria uma nova instância de SettingsHandler.
func NewSettingsHandler(orgService services.OrganizationService) *SettingsHandler {
	return &SettingsHandler{orgService: orgService}
}

// GetFuelIntegrationSettings lida com a busca das configurações de integração de combustível.
func (h *SettingsHandler) GetFuelIntegrationSettings(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	org, err := h.orgService.GetFuelIntegrationSettings(orgID)
	if err != nil || org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organização não encontrada"})
		return
	}

	response := schemas.OrganizationFuelIntegrationPublic{
		FuelProviderName: *org.FuelProviderName,
		IsAPIKeySet:      org.EncryptedFuelProviderAPIKey != nil,
		IsAPISecretSet:   org.EncryptedFuelProviderAPISecret != nil,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateFuelIntegrationSettings lida com a atualização das configurações de integração de combustível.
func (h *SettingsHandler) UpdateFuelIntegrationSettings(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	var settingsIn schemas.OrganizationFuelIntegrationUpdate
	if err := c.ShouldBindJSON(&settingsIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedOrg, err := h.orgService.UpdateFuelIntegrationSettings(orgID, settingsIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	response := schemas.OrganizationFuelIntegrationPublic{
		FuelProviderName: *updatedOrg.FuelProviderName,
		IsAPIKeySet:      updatedOrg.EncryptedFuelProviderAPIKey != nil,
		IsAPISecretSet:   updatedOrg.EncryptedFuelProviderAPISecret != nil,
	}

	c.JSON(http.StatusOK, response)
}
