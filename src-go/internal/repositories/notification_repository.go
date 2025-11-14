package repositories

import (
	"errors"
	"gorm.io/gorm"
	"go-api/internal/models"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByUser(userID, orgID uint) ([]models.Notification, error)
	CountUnreadByUser(userID, orgID uint) (int64, error)
	FindByIDAndUser(notificationID, userID, orgID uint) (*models.Notification, error)
	Update(notification *models.Notification) error
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

func (r *notificationRepository) FindByUser(userID, orgID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ? AND organization_id = ?", userID, orgID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) CountUnreadByUser(userID, orgID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).Where("user_id = ? AND organization_id = ? AND is_read = ?", userID, orgID, false).Count(&count).Error
	return count, err
}

func (r *notificationRepository) FindByIDAndUser(notificationID, userID, orgID uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.Where("id = ? AND user_id = ? AND organization_id = ?", notificationID, userID, orgID).First(&notification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) Update(notification *models.Notification) error {
	return r.db.Save(notification).Error
}
