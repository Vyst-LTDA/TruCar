package repositories

import (
	"gorm.io/gorm"

	"go-api/internal/models"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	// Outras funções de busca e atualização podem ser adicionadas aqui no futuro
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}
