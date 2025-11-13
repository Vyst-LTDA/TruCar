package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type VehicleHandler struct {
	service services.VehicleService
}

func NewVehicleHandler(service services.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

func (h *VehicleHandler) GetVehicles(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	rowsPerPage, _ := strconv.Atoi(c.DefaultQuery("rowsPerPage", "8"))
	search := c.DefaultQuery("search", "")

	skip := (page - 1) * rowsPerPage

	vehicles, total, err := h.service.GetVehicles(orgID, skip, rowsPerPage, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vehicles": vehicles, "total_items": total})
}

func (h *VehicleHandler) CreateVehicle(c *gin.Context) {
	var vehicleIn schemas.VehicleCreate
	if err := c.ShouldBindJSON(&vehicleIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	createdVehicle, err := h.service.CreateVehicle(vehicleIn, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vehicle"})
		return
	}

	c.JSON(http.StatusCreated, createdVehicle)
}

func (h *VehicleHandler) GetVehicle(c *gin.Context) {
	vehicleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	vehicle, err := h.service.GetVehicle(uint(vehicleID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicle"})
		return
	}
	if vehicle == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

func (h *VehicleHandler) UpdateVehicle(c *gin.Context) {
	vehicleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	var vehicleIn schemas.VehicleUpdate
	if err := c.ShouldBindJSON(&vehicleIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	updatedVehicle, err := h.service.UpdateVehicle(uint(vehicleID), orgID, vehicleIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle"})
		return
	}
	if updatedVehicle == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	c.JSON(http.StatusOK, updatedVehicle)
}

func (h *VehicleHandler) DeleteVehicle(c *gin.Context) {
	vehicleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	err = h.service.DeleteVehicle(uint(vehicleID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
