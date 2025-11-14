package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

// RegisterNotificationRoutes registra as rotas de notificações.
func RegisterNotificationRoutes(handler *api.NotificationHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		notifications := router.Group("/notifications")
		{
			notifications.GET("/", handler.GetNotifications)
			notifications.GET("/unread-count", handler.GetUnreadCount)
			notifications.POST("/:id/read", handler.MarkAsRead)
		}
	}
}
