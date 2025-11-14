package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type FuelLogHandler struct {
	service services.FuelLogService
}

func NewFuelLogHandler(service services.FuelLogService) *FuelLogHandler {
	return &FuelLogHandler{service: service}
}

func (h *FuelLogHandler) GetFuelLogs(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	fuelLogs, err := h.service.GetFuelLogs(currentUser, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fuel logs"})
		return
	}

	c.JSON(http.StatusOK, fuelLogs)
}

func (h *FuelLogHandler) CreateFuelLog(c *gin.Context) {
	var fuelLogIn schemas.FuelLogCreate
	if err := c.ShouldBindJSON(&fuelLogIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	createdFuelLog, err := h.service.CreateFuelLog(fuelLogIn, currentUser.(models.User))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create fuel log"})
		return
	}

	c.JSON(http.StatusCreated, createdFuelLog)
}

func (h *FuelLogHandler) GetFuelLog(c *gin.Context) {
	fuelLogID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fuel log ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	fuelLog, err := h.service.GetFuelLog(uint(fuelLogID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fuel log"})
		return
	}
	if fuelLog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fuel log not found"})
		return
	}

	c.JSON(http.StatusOK, fuelLog)
}

func (h *FuelLogHandler) UpdateFuelLog(c *gin.Context) {
	fuelLogID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fuel log ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	var fuelLogIn schemas.FuelLogUpdate
	if err := c.ShouldBindJSON(&fuelLogIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedFuelLog, err := h.service.UpdateFuelLog(uint(fuelLogID), orgID, fuelLogIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fuel log"})
		return
	}
	if updatedFuelLog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fuel log not found"})
		return
	}

	c.JSON(http.StatusOK, updatedFuelLog)
}

func (h *FuelLogHandler) DeleteFuelLog(c *gin.Context) {
	fuelLogID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fuel log ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	err = h.service.DeleteFuelLog(uint(fuelLogID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete fuel log"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
