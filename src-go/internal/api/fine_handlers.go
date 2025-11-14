package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type FineHandler struct {
	service services.FineService
}

func NewFineHandler(service services.FineService) *FineHandler {
	return &FineHandler{service: service}
}

func (h *FineHandler) GetFines(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	fines, err := h.service.GetFines(currentUser, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fines"})
		return
	}

	c.JSON(http.StatusOK, fines)
}

func (h *FineHandler) CreateFine(c *gin.Context) {
	var fineIn schemas.FineCreate
	if err := c.ShouldBindJSON(&fineIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	createdFine, err := h.service.CreateFine(fineIn, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create fine"})
		return
	}

	c.JSON(http.StatusCreated, createdFine)
}

func (h *FineHandler) GetFine(c *gin.Context) {
	fineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fine ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	fine, err := h.service.GetFine(uint(fineID), currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fine"})
		return
	}
	if fine == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fine not found"})
		return
	}

	c.JSON(http.StatusOK, fine)
}

func (h *FineHandler) UpdateFine(c *gin.Context) {
	fineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fine ID"})
		return
	}

	var fineIn schemas.FineUpdate
	if err := c.ShouldBindJSON(&fineIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	updatedFine, err := h.service.UpdateFine(uint(fineID), fineIn, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fine"})
		return
	}
	if updatedFine == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fine not found"})
		return
	}

	c.JSON(http.StatusOK, updatedFine)
}

func (h *FineHandler) DeleteFine(c *gin.Context) {
	fineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fine ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	err = h.service.DeleteFine(uint(fineID), currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete fine"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
