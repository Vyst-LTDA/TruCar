package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go-api/internal/models"
	"go-api/internal/schemas"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	router.GET("/users", GetUsers)
	router.POST("/users", CreateUserHandler)
	router.GET("/users/:id", GetUser)
	router.PUT("/users/:id", UpdateUserHandler)
	router.DELETE("/users/:id", DeleteUserHandler)
}

func GetUsers(c *gin.Context) {
	orgID, _ := strconv.Atoi(c.DefaultQuery("org_id", "1"))
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	users, err := GetUsersByOrganization(uint(orgID), skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func CreateUserHandler(c *gin.Context) {
	var userIn schemas.UserCreate
	if err := c.ShouldBindJSON(&userIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userIn.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:          userIn.Email,
		FullName:       userIn.FullName,
		HashedPassword: string(hashedPassword),
		Role:           userIn.Role,
		OrganizationID: 1, // Hardcoded for now
	}

	if err := CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	orgID, _ := strconv.Atoi(c.DefaultQuery("org_id", "1"))

	user, err := GetUserByID(uint(userID), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	orgID, _ := strconv.Atoi(c.DefaultQuery("org_id", "1"))

	user, err := GetUserByID(uint(userID), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var userIn schemas.UserUpdate
	if err := c.ShouldBindJSON(&userIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userIn.FullName != "" {
		user.FullName = userIn.FullName
	}

	if userIn.Email != "" {
		user.Email = userIn.Email
	}

	if userIn.IsActive != nil {
		user.IsActive = *userIn.IsActive
	}

	if err := UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeleteUserHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	orgID, _ := strconv.Atoi(c.DefaultQuery("org_id", "1"))

	user, err := GetUserByID(uint(userID), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := DeleteUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
