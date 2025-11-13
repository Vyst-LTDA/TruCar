package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type MaintenanceHandler struct {
	service services.MaintenanceService
}

func NewMaintenanceHandler(service services.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{service: service}
}

func (h *MaintenanceHandler) GetMaintenanceRequests(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	search := c.DefaultQuery("search", "")

	requests, err := h.service.GetMaintenanceRequests(orgID, skip, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance requests"})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *MaintenanceHandler) CreateMaintenanceRequest(c *gin.Context) {
	var reqIn schemas.MaintenanceRequestCreate
	if err := c.ShouldBindJSON(&reqIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	createdReq, err := h.service.CreateMaintenanceRequest(reqIn, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create maintenance request"})
		return
	}

	c.JSON(http.StatusCreated, createdReq)
}

func (h *MaintenanceHandler) GetMaintenanceRequest(c *gin.Context) {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	req, err := h.service.GetMaintenanceRequest(uint(reqID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance request"})
		return
	}
	if req == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance request not found"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *MaintenanceHandler) UpdateMaintenanceRequestStatus(c *gin.Context) {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var reqIn schemas.MaintenanceRequestUpdate
	if err := c.ShouldBindJSON(&reqIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	updatedReq, err := h.service.UpdateMaintenanceRequestStatus(uint(reqID), reqIn, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance request"})
		return
	}
	if updatedReq == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance request not found"})
		return
	}

	c.JSON(http.StatusOK, updatedReq)
}

func (h *MaintenanceHandler) DeleteMaintenanceRequest(c *gin.Context) {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	err = h.service.DeleteMaintenanceRequest(uint(reqID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete maintenance request"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *MaintenanceHandler) GetMaintenanceComments(c *gin.Context) {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	comments, err := h.service.GetMaintenanceComments(uint(reqID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *MaintenanceHandler) CreateMaintenanceComment(c *gin.Context) {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var commentIn schemas.MaintenanceCommentCreate
	if err := c.ShouldBindJSON(&commentIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	createdComment, err := h.service.CreateMaintenanceComment(commentIn, uint(reqID), currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}
