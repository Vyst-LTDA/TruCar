package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type ImplementHandler struct {
	service services.ImplementService
}

func NewImplementHandler(service services.ImplementService) *ImplementHandler {
	return &ImplementHandler{service: service}
}

func (h *ImplementHandler) GetImplements(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	managementList := c.Query("management-list") == "true"

	implements, err := h.service.GetImplements(orgID, managementList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch implements"})
		return
	}

	var publicImplements []schemas.ImplementPublic
	for _, implement := range implements {
		publicImplements = append(publicImplements, schemas.ToImplementPublic(implement))
	}

	c.JSON(http.StatusOK, publicImplements)
}

func (h *ImplementHandler) CreateImplement(c *gin.Context) {
	var implementIn schemas.ImplementCreate
	if err := c.ShouldBindJSON(&implementIn); err != nil {
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

	createdImplement, err := h.service.CreateImplement(implementIn, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create implement"})
		return
	}

	c.JSON(http.StatusCreated, schemas.ToImplementPublic(*createdImplement))
}

func (h *ImplementHandler) GetImplement(c *gin.Context) {
	implementID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid implement ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	implement, err := h.service.GetImplement(uint(implementID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch implement"})
		return
	}
	if implement == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Implement not found"})
		return
	}

	c.JSON(http.StatusOK, schemas.ToImplementPublic(*implement))
}

func (h *ImplementHandler) UpdateImplement(c *gin.Context) {
	implementID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid implement ID"})
		return
	}

	var implementIn schemas.ImplementUpdate
	if err := c.ShouldBindJSON(&implementIn); err != nil {
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

	updatedImplement, err := h.service.UpdateImplement(uint(implementID), orgID, implementIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update implement"})
		return
	}
	if updatedImplement == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Implement not found"})
		return
	}

	c.JSON(http.StatusOK, schemas.ToImplementPublic(*updatedImplement))
}

func (h *ImplementHandler) DeleteImplement(c *gin.Context) {
	implementID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid implement ID"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)
	orgID := currentUser.OrganizationID

	err = h.service.DeleteImplement(uint(implementID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete implement"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
