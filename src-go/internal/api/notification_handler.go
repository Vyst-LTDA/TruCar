package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-api/internal/middleware"
	"go-api/internal/models"
	"go-api/internal/services"
)

// NotificationHandler gerencia as requisições HTTP para notificações.
type NotificationHandler struct {
	service services.NotificationService
}

// NewNotificationHandler cria uma nova instância de NotificationHandler.
func NewNotificationHandler(service services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// GetNotifications lida com a busca de notificações para o usuário logado.
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	notifications, err := h.service.GetNotificationsForUser(currentUser.ID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetUnreadCount lida com a contagem de notificações não lidas.
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	count, err := h.service.GetUnreadCountForUser(currentUser.ID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count unread notifications"})
		return
	}

	c.JSON(http.StatusOK, count)
}

// MarkAsRead lida com a marcação de uma notificação como lida.
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notificação inválido"})
		return
	}

	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	notification, err := h.service.MarkAsRead(uint(notificationID), currentUser.ID, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	if notification == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, notification)
}
