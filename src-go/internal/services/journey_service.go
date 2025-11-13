package services

import (
	"time"

	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type JourneyService interface {
	GetJourneys(orgID uint, skip, limit int, driverID, vehicleID *uint, dateFrom, dateTo *time.Time) ([]models.Journey, error)
	StartJourney(journeyIn schemas.JourneyCreate, driverID, orgID uint) (*models.Journey, error)
	EndJourney(journeyID, orgID uint, endMileage *int, endEngineHours *float64) (*models.Journey, *models.Vehicle, error)
	DeleteJourney(journeyID, orgID uint) error
}

type journeyService struct {
	journeyRepo repositories.JourneyRepository
	vehicleRepo repositories.VehicleRepository
}

func NewJourneyService(journeyRepo repositories.JourneyRepository, vehicleRepo repositories.VehicleRepository) JourneyService {
	return &journeyService{journeyRepo: journeyRepo, vehicleRepo: vehicleRepo}
}

func (s *journeyService) GetJourneys(orgID uint, skip, limit int, driverID, vehicleID *uint, dateFrom, dateTo *time.Time) ([]models.Journey, error) {
	return s.journeyRepo.FindByOrganization(orgID, skip, limit, driverID, vehicleID, dateFrom, dateTo)
}

func (s *journeyService) StartJourney(journeyIn schemas.JourneyCreate, driverID, orgID uint) (*models.Journey, error) {
	vehicle, err := s.vehicleRepo.FindByID(journeyIn.VehicleID, orgID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil || vehicle.Status != models.StatusAvailable {
		return nil, repositories.ErrVehicleNotAvailable
	}

	journey := &models.Journey{
		VehicleID:               journeyIn.VehicleID,
		TripType:                journeyIn.TripType,
		DestinationAddress:      journeyIn.DestinationAddress,
		TripDescription:         journeyIn.TripDescription,
		ImplementID:             journeyIn.ImplementID,
		DestinationStreet:       journeyIn.DestinationStreet,
		DestinationNeighborhood: journeyIn.DestinationNeighborhood,
		DestinationCity:         journeyIn.DestinationCity,
		DestinationState:        journeyIn.DestinationState,
		DestinationCEP:          journeyIn.DestinationCEP,
		DriverID:                driverID,
		OrganizationID:          orgID,
		StartTime: 				 time.Now(),
	}

	createdJourney, err := s.journeyRepo.Create(journey)
	if err != nil {
		return nil, err
	}

	vehicle.Status = models.StatusInUse
	err = s.vehicleRepo.Update(vehicle)
	return createdJourney, err
}

func (s *journeyService) EndJourney(journeyID, orgID uint, endMileage *int, endEngineHours *float64) (*models.Journey, *models.Vehicle, error) {
	journey, err := s.journeyRepo.FindByID(journeyID, orgID)
	if err != nil {
		return nil, nil, err
	}
	if journey == nil {
		return nil, nil, nil // Or return a not found error
	}

	now := time.Now()
	journey.EndTime = &now
	journey.EndMileage = endMileage
	journey.EndEngineHours = endEngineHours
	journey.IsActive = false

	updatedJourney, err := s.journeyRepo.Update(journey)
	if err != nil {
		return nil, nil, err
	}

	vehicle, err := s.vehicleRepo.FindByID(journey.VehicleID, orgID)
	if err != nil {
		return nil, nil, err
	}
	if vehicle == nil {
		return nil, nil, nil // Or return a not found error
	}

	vehicle.Status = models.StatusAvailable
	if endMileage != nil {
		vehicle.CurrentKM = *endMileage
	}

	err = s.vehicleRepo.Update(vehicle)
	return updatedJourney, vehicle, err
}

func (s *journeyService) DeleteJourney(journeyID, orgID uint) error {
	journey, err := s.journeyRepo.FindByID(journeyID, orgID)
	if err != nil {
		return err
	}
	if journey == nil {
		return nil // Or return a not found error
	}
	return s.journeyRepo.Delete(journey)
}
