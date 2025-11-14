package services

import (
	"go-api/internal/logging"
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go.uber.org/zap"
)

type NotificationService interface {
	CreateNotificationAsync(notification *models.Notification)
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
