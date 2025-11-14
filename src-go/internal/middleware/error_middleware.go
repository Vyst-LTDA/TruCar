package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			log.Printf("Error: %v", err.Err)

			// Aqui você pode adicionar lógica para mapear diferentes tipos de erro
			// para diferentes códigos de status e mensagens.
			appErr := &AppError{
				Code:    http.StatusInternalServerError,
				Message: "An internal server error occurred",
			}

			c.JSON(appErr.Code, appErr)
			return
		}
	}
}
