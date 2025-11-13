package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-api/internal/models"
	"go-api/internal/schemas"
)

func CreateImplement(c *gin.Context) {
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

	implement := models.Implement{
		Name:           implementIn.Name,
		Brand:          implementIn.Brand,
		VehicleModel:          implementIn.VehicleModel,
		Year:           implementIn.Year,
		Identifier:     implementIn.Identifier,
		Type:           implementIn.Type,
		OrganizationID: currentUser.OrganizationID,
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&implement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create implement"})
		return
	}

	c.JSON(http.StatusCreated, schemas.ToImplementPublic(implement))
}

func GetImplements(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	var implements []models.Implement
	db := c.MustGet("db").(*gorm.DB)

	query := db.Where("organization_id = ?", currentUser.OrganizationID)

	// Filtro de status "disponível" do endpoint original
	if c.Query("management-list") != "true" {
		query = query.Where("status = ?", models.ImplementStatusAvailable)
	}

	if err := query.Find(&implements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve implements"})
		return
	}

	var publicImplements []schemas.ImplementPublic
	for _, implement := range implements {
		publicImplements = append(publicImplements, schemas.ToImplementPublic(implement))
	}

	c.JSON(http.StatusOK, publicImplements)
}

func GetImplementByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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

	var implement models.Implement
	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&implement, "id = ? AND organization_id = ?", id, currentUser.OrganizationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Implement not found"})
		return
	}

	c.JSON(http.StatusOK, schemas.ToImplementPublic(implement))
}

func UpdateImplement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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

	var implement models.Implement
	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&implement, "id = ? AND organization_id = ?", id, currentUser.OrganizationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Implement not found"})
		return
	}

	var implementIn schemas.ImplementUpdate
	if err := c.ShouldBindJSON(&implementIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&implement).Updates(implementIn)
	db.First(&implement, implement.ID)

	c.JSON(http.StatusOK, schemas.ToImplementPublic(implement))
}

func DeleteImplement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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

	var implement models.Implement
	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&implement, "id = ? AND organization_id = ?", id, currentUser.OrganizationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Implement not found"})
		return
	}

	db.Delete(&implement)

	c.Status(http.StatusNoContent)
}
