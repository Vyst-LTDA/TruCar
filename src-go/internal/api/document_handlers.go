package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

type DocumentHandler struct {
	service services.DocumentService
}

func NewDocumentHandler(service services.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) GetDocuments(c *gin.Context) {
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	var expiringInDays *int
	if val, err := strconv.Atoi(c.Query("expiring_in_days")); err == nil {
		expiringInDays = &val
	}

	docs, err := h.service.GetDocuments(currentUser.OrganizationID, skip, limit, expiringInDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, docs)
}

func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var docIn schemas.DocumentCreate
	if err := c.ShouldBind(&docIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	createdDoc, err := h.service.CreateDocument(docIn, file, currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}

	c.JSON(http.StatusCreated, createdDoc)
}

func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	docID, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	err := h.service.DeleteDocument(uint(docID), currentUser.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
