package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type JourneyHandler struct {
	service services.JourneyService
}

func NewJourneyHandler(service services.JourneyService) *JourneyHandler {
	return &JourneyHandler{service: service}
}

func (h *JourneyHandler) GetJourneys(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	var driverID, vehicleID *uint
	if val, err := strconv.Atoi(c.Query("driver_id")); err == nil {
		id := uint(val)
		driverID = &id
	}
	if val, err := strconv.Atoi(c.Query("vehicle_id")); err == nil {
		id := uint(val)
		vehicleID = &id
	}

	var dateFrom, dateTo *time.Time
	if val, err := time.Parse("2006-01-02", c.Query("date_from")); err == nil {
		dateFrom = &val
	}
	if val, err := time.Parse("2006-01-02", c.Query("date_to")); err == nil {
		dateTo = &val
	}

	journeys, err := h.service.GetJourneys(orgID, skip, limit, driverID, vehicleID, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch journeys"})
		return
	}

	c.JSON(http.StatusOK, journeys)
}

func (h *JourneyHandler) StartJourney(c *gin.Context) {
	var journeyIn schemas.JourneyCreate
	if err := c.ShouldBindJSON(&journeyIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	driverID := currentUser.ID
	orgID := currentUser.OrganizationID

	createdJourney, err := h.service.StartJourney(journeyIn, driverID, orgID)
	if err != nil {
		if err == repositories.ErrVehicleNotAvailable {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start journey"})
		return
	}

	c.JSON(http.StatusCreated, createdJourney)
}

func (h *JourneyHandler) EndJourney(c *gin.Context) {
	journeyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid journey ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	var journeyIn schemas.JourneyUpdate
	if err := c.ShouldBindJSON(&journeyIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endedJourney, updatedVehicle, err := h.service.EndJourney(uint(journeyID), orgID, journeyIn.EndMileage, journeyIn.EndEngineHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end journey"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"journey": endedJourney, "vehicle": updatedVehicle})
}

func (h *JourneyHandler) DeleteJourney(c *gin.Context) {
	journeyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid journey ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	err = h.service.DeleteJourney(uint(journeyID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete journey"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
