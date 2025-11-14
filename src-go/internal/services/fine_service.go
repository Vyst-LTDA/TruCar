package services

import (
	"fmt"
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type FineService interface {
	GetFines(user models.User, skip, limit int) ([]models.Fine, error)
	GetFine(fineID, orgID uint) (*models.Fine, error)
	CreateFine(fineIn schemas.FineCreate, user models.User) (*models.Fine, error)
	UpdateFine(fineID uint, fineIn schemas.FineUpdate, user models.User) (*models.Fine, error)
	DeleteFine(fineID, orgID uint) error
}

type fineService struct {
	fineRepo         repositories.FineRepository
	notificationService NotificationService
}

func NewFineService(fineRepo repositories.FineRepository, notificationService NotificationService) FineService {
	return &fineService{fineRepo: fineRepo, notificationService: notificationService}
}

func (s *fineService) GetFines(user models.User, skip, limit int) ([]models.Fine, error) {
	if user.Role == models.RoleClienteAtivo || user.Role == models.RoleClienteDemo {
		return s.fineRepo.FindByOrganization(user.OrganizationID, skip, limit)
	}
	return s.fineRepo.FindByDriver(user.ID, user.OrganizationID, skip, limit)
}

func (s *fineService) GetFine(fineID, orgID uint) (*models.Fine, error) {
	return s.fineRepo.FindByID(fineID, orgID)
}

func (s *fineService) CreateFine(fineIn schemas.FineCreate, user models.User) (*models.Fine, error) {
	fine := &models.Fine{
		Description:    fineIn.Description,
		InfractionCode: fineIn.InfractionCode,
		Date:           fineIn.Date,
		Value:          fineIn.Value,
		Status:         models.FineStatus(fineIn.Status),
		VehicleID:      fineIn.VehicleID,
		DriverID:       fineIn.DriverID,
		OrganizationID: user.OrganizationID,
	}

	createdFine, err := s.fineRepo.Create(fine)
	if err != nil {
		return nil, err
	}

	// Disparar notificação em segundo plano
	go func() {
		notification := &models.Notification{
			OrganizationID:   user.OrganizationID,
			UserID:           user.ID, // Temporário - idealmente, notificar todos os gestores
			Message:          fmt.Sprintf("Nova multa de R$%.2f registrada para o veículo.", createdFine.Value),
			NotificationType: models.NotificationTypeNewFineRegistered,
			RelatedEntityType: "fine",
			RelatedEntityID:  &createdFine.ID,
			RelatedVehicleID: &createdFine.VehicleID,
		}
		s.notificationService.CreateNotificationAsync(notification)
	}()

	return createdFine, nil
}

func (s *fineService) UpdateFine(fineID uint, fineIn schemas.FineUpdate, user models.User) (*models.Fine, error) {
	fine, err := s.fineRepo.FindByID(fineID, user.OrganizationID)
	if err != nil {
		return nil, err
	}
	if fine == nil {
		return nil, nil // Not found
	}

	if fineIn.Description != "" {
		fine.Description = fineIn.Description
	}
	if fineIn.InfractionCode != "" {
		fine.InfractionCode = fineIn.InfractionCode
	}
	if !fineIn.Date.IsZero() {
		fine.Date = fineIn.Date
	}
	if fineIn.Value > 0 {
		fine.Value = fineIn.Value
	}
	if fineIn.Status != "" {
		fine.Status = models.FineStatus(fineIn.Status)
	}
	if fineIn.VehicleID > 0 {
		fine.VehicleID = fineIn.VehicleID
	}
	if fineIn.DriverID != nil {
		fine.DriverID = fineIn.DriverID
	}

	return s.fineRepo.Update(fine)
}

func (s *fineService) DeleteFine(fineID, orgID uint) error {
	fine, err := s.fineRepo.FindByID(fineID, orgID)
	if err != nil {
		return err
	}
	if fine == nil {
		return nil // Not found
	}
	return s.fineRepo.Delete(fine)
}
