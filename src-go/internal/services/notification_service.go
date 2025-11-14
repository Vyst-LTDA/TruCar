package services

import (
	"go-api/internal/logging"
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go.uber.org/zap"
)

type NotificationService interface {
	CreateNotificationAsync(notification *models.Notification)
	GetNotificationsForUser(userID, orgID uint) ([]models.Notification, error)
	GetUnreadCountForUser(userID, orgID uint) (int64, error)
	MarkAsRead(notificationID, userID, orgID uint) (*models.Notification, error)
}

type notificationService struct {
	repo repositories.NotificationRepository
}

func NewNotificationService(repo repositories.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) CreateNotificationAsync(notification *models.Notification) {
	go func() {
		err := s.repo.Create(notification)
		if err != nil {
			logging.Logger.Error("Failed to create notification in background", zap.Error(err))
		}
	}()
}

func (s *notificationService) GetNotificationsForUser(userID, orgID uint) ([]models.Notification, error) {
	return s.repo.FindByUser(userID, orgID)
}

func (s *notificationService) GetUnreadCountForUser(userID, orgID uint) (int64, error) {
	return s.repo.CountUnreadByUser(userID, orgID)
}

func (s *notificationService) MarkAsRead(notificationID, userID, orgID uint) (*models.Notification, error) {
	notification, err := s.repo.FindByIDAndUser(notificationID, userID, orgID)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, nil // Not found
	}

	if !notification.IsRead {
		notification.IsRead = true
		err = s.repo.Update(notification)
		if err != nil {
			return nil, err
		}
	}

	return notification, nil
}
