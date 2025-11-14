package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type PartHandler struct {
	service services.PartService
}

func NewPartHandler(service services.PartService) *PartHandler {
	return &PartHandler{service: service}
}

func (h *PartHandler) GetParts(c *gin.Context) {
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	search := c.Query("search")
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	parts, err := h.service.GetParts(currentUser.OrganizationID, search, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch parts"})
		return
	}

	c.JSON(http.StatusOK, parts)
}

func (h *PartHandler) GetPart(c *gin.Context) {
	partID, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	part, err := h.service.GetPart(uint(partID), currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch part"})
		return
	}
	if part == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Part not found"})
		return
	}

	c.JSON(http.StatusOK, part)
}

func (h *PartHandler) CreatePart(c *gin.Context) {
	var partIn schemas.PartCreate
	if err := c.ShouldBindJSON(&partIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	createdPart, err := h.service.CreatePart(partIn, currentUser.OrganizationID, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create part"})
		return
	}

	c.JSON(http.StatusCreated, createdPart)
}

func (h *PartHandler) UpdatePart(c *gin.Context) {
	partID, _ := strconv.Atoi(c.Param("id"))
	var partIn schemas.PartUpdate
	if err := c.ShouldBindJSON(&partIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	updatedPart, err := h.service.UpdatePart(uint(partID), partIn, currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update part"})
		return
	}
	if updatedPart == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Part not found"})
		return
	}

	c.JSON(http.StatusOK, updatedPart)
}

func (h *PartHandler) DeletePart(c *gin.Context) {
	partID, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	err := h.service.DeletePart(uint(partID), currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete part"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *PartHandler) AddInventoryItems(c *gin.Context) {
	partID, _ := strconv.Atoi(c.Param("id"))
	var payload schemas.AddItemsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	err := h.service.AddInventoryItems(uint(partID), payload, currentUser.OrganizationID, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add items"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Items added successfully"})
}

func (h *PartHandler) SetInventoryItemStatus(c *gin.Context) {
	itemID, _ := strconv.Atoi(c.Param("item_id"))
	var payload schemas.SetItemStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	updatedItem, err := h.service.SetInventoryItemStatus(uint(itemID), payload, currentUser.OrganizationID, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set item status"})
		return
	}
	if updatedItem == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

func (h *PartHandler) GetItemsForPart(c *gin.Context) {
	partID, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	statusStr := c.Query("status")
	var status *models.InventoryItemStatus
	if statusStr != "" {
		s := models.InventoryItemStatus(statusStr)
		status = &s
	}

	items, err := h.service.GetItemsForPart(uint(partID), status, currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *PartHandler) GetPartHistory(c *gin.Context) {
	partID, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	history, err := h.service.GetPartHistory(uint(partID), currentUser.OrganizationID, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
