package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/internal/services"
	"go-api/internal/schemas"
)

type AdminHandler struct {
	orgService  services.OrganizationService
	userService services.UserService
	authService services.AuthService
}

func NewAdminHandler(orgService services.OrganizationService, userService services.UserService, authService services.AuthService) *AdminHandler {
	return &AdminHandler{orgService: orgService, userService: userService, authService: authService}
}

func (h *AdminHandler) GetOrganizations(c *gin.Context) {
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	status := c.Query("status")

	orgs, err := h.orgService.GetOrganizations(skip, limit, &status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func (h *AdminHandler) UpdateOrganization(c *gin.Context) {
	orgID, _ := strconv.Atoi(c.Param("id"))
	var orgIn schemas.OrganizationUpdate
	if err := c.ShouldBindJSON(&orgIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedOrg, err := h.orgService.UpdateOrganization(uint(orgID), orgIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}
	c.JSON(http.StatusOK, updatedOrg)
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	users, err := h.userService.GetAllUsers(skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) GetDemoUsers(c *gin.Context) {
	users, err := h.userService.GetDemoUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch demo users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) ActivateUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	activatedUser, err := h.userService.ActivateUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}
	c.JSON(http.StatusOK, activatedUser)
}

func (h *AdminHandler) ImpersonateUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	token, err := h.authService.Impersonate(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to impersonate user"})
		return
	}
	c.JSON(http.StatusOK, schemas.Token{AccessToken: token, TokenType: "bearer"})
}
