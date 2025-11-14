package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type FreightOrderHandler struct {
	service services.FreightOrderService
}

func NewFreightOrderHandler(service services.FreightOrderService) *FreightOrderHandler {
	return &FreightOrderHandler{service: service}
}

func (h *FreightOrderHandler) GetFreightOrders(c *gin.Context) {
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	orders, err := h.service.GetFreightOrders(currentUser, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch freight orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *FreightOrderHandler) CreateFreightOrder(c *gin.Context) {
	var orderIn schemas.FreightOrderCreate
	if err := c.ShouldBindJSON(&orderIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	createdOrder, err := h.service.CreateFreightOrder(orderIn, currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create freight order"})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

func (h *FreightOrderHandler) GetFreightOrderByID(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	order, err := h.service.GetFreightOrderByID(uint(orderID), currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch freight order"})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Freight order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *FreightOrderHandler) GetOpenFreightOrders(c *gin.Context) {
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	orders, err := h.service.GetOpenFreightOrders(currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch open freight orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *FreightOrderHandler) GetMyPendingFreightOrders(c *gin.Context) {
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	orders, err := h.service.GetMyPendingFreightOrders(currentUser.ID, currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending freight orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *FreightOrderHandler) ClaimFreightOrder(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("id"))
	var claimIn schemas.FreightOrderClaim
	if err := c.ShouldBindJSON(&claimIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	claimedOrder, err := h.service.ClaimFreightOrder(uint(orderID), claimIn, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to claim freight order"})
		return
	}

	c.JSON(http.StatusOK, claimedOrder)
}

func (h *FreightOrderHandler) StartJourneyForStop(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("order_id"))
	stopPointID, _ := strconv.Atoi(c.Param("stop_point_id"))
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	journey, err := h.service.StartJourneyForStop(uint(orderID), uint(stopPointID), currentUser.ID, currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start journey for stop"})
		return
	}

	c.JSON(http.StatusCreated, journey)
}

func (h *FreightOrderHandler) CompleteStopPoint(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("order_id"))
	stopPointID, _ := strconv.Atoi(c.Param("stop_point_id"))

	var data struct {
		JourneyID  uint `json:"journey_id"`
		EndMileage int  `json:"end_mileage"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	completedStop, err := h.service.CompleteStopPoint(uint(orderID), uint(stopPointID), currentUser.ID, currentUser.OrganizationID, data.JourneyID, data.EndMileage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete stop point"})
		return
	}

	c.JSON(http.StatusOK, completedStop)
}
