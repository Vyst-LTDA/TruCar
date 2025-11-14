package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterDocumentRoutes(handler *api.DocumentHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/documents", handler.GetDocuments)
		router.POST("/documents", handler.CreateDocument)
		router.DELETE("/documents/:id", handler.DeleteDocument)
	}
}
