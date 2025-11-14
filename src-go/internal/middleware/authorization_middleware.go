package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api/internal/models"
)

func AuthorizationMiddleware(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("currentUser")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		currentUser := user.(models.User)

		isAllowed := false
		for _, role := range allowedRoles {
			if currentUser.Role == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			return
		}

		c.Next()
	}
}
