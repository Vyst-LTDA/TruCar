package services

import (
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type FuelLogService interface {
	GetFuelLogs(user models.User, skip, limit int) ([]models.FuelLog, error)
	GetFuelLog(fuelLogID, orgID uint) (*models.FuelLog, error)
	CreateFuelLog(fuelLogIn schemas.FuelLogCreate, currentUser models.User) (*models.FuelLog, error)
	UpdateFuelLog(fuelLogID, orgID uint, fuelLogIn schemas.FuelLogUpdate) (*models.FuelLog, error)
	DeleteFuelLog(fuelLogID, orgID uint) error
}

type fuelLogService struct {
	repo repositories.FuelLogRepository
}

func NewFuelLogService(repo repositories.FuelLogRepository) FuelLogService {
	return &fuelLogService{repo: repo}
}

func (s *fuelLogService) GetFuelLogs(user models.User, skip, limit int) ([]models.FuelLog, error) {
	if user.Role == models.RoleClienteAtivo || user.Role == models.RoleClienteDemo {
		return s.repo.FindByOrganization(user.OrganizationID, skip, limit)
	}
	return s.repo.FindByUser(user.ID, user.OrganizationID, skip, limit)
}

func (s *fuelLogService) GetFuelLog(fuelLogID, orgID uint) (*models.FuelLog, error) {
	return s.repo.FindByID(fuelLogID, orgID)
}

func (s *fuelLogService) CreateFuelLog(fuelLogIn schemas.FuelLogCreate, currentUser models.User) (*models.FuelLog, error) {
	userID := currentUser.ID
	if fuelLogIn.UserID != nil {
		userID = *fuelLogIn.UserID
	}

	fuelLog := &models.FuelLog{
		Odometer:        fuelLogIn.Odometer,
		Liters:          fuelLogIn.Liters,
		TotalCost:       fuelLogIn.TotalCost,
		VehicleID:       fuelLogIn.VehicleID,
		UserID:          userID,
		ReceiptPhotoURL: fuelLogIn.ReceiptPhotoURL,
		OrganizationID:  currentUser.OrganizationID,
	}

	err := s.repo.Create(fuelLog)
	return fuelLog, err
}

func (s *fuelLogService) UpdateFuelLog(fuelLogID, orgID uint, fuelLogIn schemas.FuelLogUpdate) (*models.FuelLog, error) {
	fuelLog, err := s.repo.FindByID(fuelLogID, orgID)
	if err != nil {
		return nil, err
	}
	if fuelLog == nil {
		return nil, nil // Or return a not found error
	}

	if fuelLogIn.Odometer != nil {
		fuelLog.Odometer = *fuelLogIn.Odometer
	}
	if fuelLogIn.Liters != nil {
		fuelLog.Liters = *fuelLogIn.Liters
	}
	if fuelLogIn.TotalCost != nil {
		fuelLog.TotalCost = *fuelLogIn.TotalCost
	}
	if fuelLogIn.VehicleID != nil {
		fuelLog.VehicleID = *fuelLogIn.VehicleID
	}
	if fuelLogIn.ReceiptPhotoURL != nil {
		fuelLog.ReceiptPhotoURL = fuelLogIn.ReceiptPhotoURL
	}

	err = s.repo.Update(fuelLog)
	return fuelLog, err
}

func (s *fuelLogService) DeleteFuelLog(fuelLogID, orgID uint) error {
	fuelLog, err := s.repo.FindByID(fuelLogID, orgID)
	if err != nil {
		return err
	}
	if fuelLog == nil {
		return nil // Or return a not found error
	}
	return s.repo.Delete(fuelLog)
}
