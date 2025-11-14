package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-api/internal/middleware"
	"go-api/internal/models"
	"go-api/internal/services"
	"go-api/internal/schemas"
)

// FineHandler gerencia as requisições HTTP para multas.
type FineHandler struct {
	service services.FineService
}

// NewFineHandler cria uma nova instância de FineHandler.
func NewFineHandler(service services.FineService) *FineHandler {
	return &FineHandler{service: service}
}

// CreateFine lida com a criação de uma nova multa.
func (h *FineHandler) CreateFine(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	var fineInput schemas.FineCreate
	if err := c.ShouldBindJSON(&fineInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Regra de negócio: motorista só pode criar multa para si mesmo
	if currentUser.Role == models.RoleDriver && fineInput.DriverID != nil && *fineInput.DriverID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Você só pode registrar multas para si mesmo."})
		return
	}

	fine := models.Fine{
		Description:    fineInput.Description,
		InfractionCode: fineInput.InfractionCode,
		Date:           fineInput.Date,
		Value:          fineInput.Value,
		Status:         models.FineStatusPending,
		VehicleID:      fineInput.VehicleID,
		DriverID:       fineInput.DriverID,
		OrganizationID: orgID,
	}

	createdFine, err := h.service.CreateFine(&fine, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar multa"})
		return
	}

	c.JSON(http.StatusCreated, createdFine)
}

// GetFines lida com a listagem de multas.
func (h *FineHandler) GetFines(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	var fines []models.Fine
	var err error

	if currentUser.Role == models.RoleClienteAtivo || currentUser.Role == models.RoleClienteDemo {
		fines, err = h.service.GetFinesByOrganization(orgID, skip, limit)
	} else if currentUser.Role == models.RoleDriver {
		fines, err = h.service.GetFinesByDriver(currentUser.ID, orgID, skip, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar multas"})
		return
	}

	c.JSON(http.StatusOK, fines)
}

// UpdateFine lida com a atualização de uma multa.
func (h *FineHandler) UpdateFine(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	fineID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de multa inválido"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedFine, err := h.service.UpdateFine(uint(fineID), orgID, payload)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedFine)
}

// DeleteFine lida com a exclusão de uma multa.
func (h *FineHandler) DeleteFine(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	fineID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de multa inválido"})
		return
	}

	if err := h.service.DeleteFine(uint(fineID), orgID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
