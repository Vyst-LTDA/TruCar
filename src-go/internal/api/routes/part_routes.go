package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

func RegisterPartRoutes(handler *api.PartHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.GET("/parts", handler.GetParts)
		router.POST("/parts", handler.CreatePart)
		router.GET("/parts/:id", handler.GetPart)
		router.PUT("/parts/:id", handler.UpdatePart)
		router.DELETE("/parts/:id", handler.DeletePart)
		router.POST("/parts/:id/add-items", handler.AddInventoryItems)
		router.PUT("/items/:item_id/set-status", handler.SetInventoryItemStatus)
		router.GET("/parts/:id/items", handler.GetItemsForPart)
		router.GET("/parts/:id/history", handler.GetPartHistory)
	}
}
