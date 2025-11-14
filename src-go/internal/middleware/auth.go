package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-api/internal/core"
	"go-api/internal/models"
	"go-api/internal/services"
)

func AuthMiddleware(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		claims, err := core.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		user, err := userService.GetUser(claims.UserID, claims.OrganizationID)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("currentUser", *user)
		c.Next()
	}
}

// GetOrganizationID extrai o ID da organização do usuário autenticado no contexto.
func GetOrganizationID(c *gin.Context) (uint, bool) {
	user, exists := c.Get("currentUser")
	if !exists {
		return 0, false
	}

	// Faz o type assertion para o tipo User do pacote models.
	currentUser, ok := user.(models.User)
	if !ok {
		return 0, false
	}

	return currentUser.OrganizationID, true
}
