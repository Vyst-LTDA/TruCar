package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api/internal/middleware"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

// GPSHandler gerencia as requisições HTTP para GPS.
type GPSHandler struct {
	service services.GPSService
}

// NewGPSHandler cria uma nova instância de GPSHandler.
func NewGPSHandler(service services.GPSService) *GPSHandler {
	return &GPSHandler{service: service}
}

// Ping lida com o recebimento de um ping de GPS.
func (h *GPSHandler) Ping(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	var ping schemas.LocationCreate
	if err := c.ShouldBindJSON(&ping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.service.ProcessPing(ping, orgID)

	c.Status(http.StatusNoContent)
}
