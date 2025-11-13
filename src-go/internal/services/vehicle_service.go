package services

import (
	"context"
	"fmt"
	"time"

	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type VehicleService interface {
	GetVehicles(orgID uint, skip, limit int, search string) ([]models.Vehicle, int64, error)
	GetVehicle(vehicleID, orgID uint) (*models.Vehicle, error)
	CreateVehicle(vehicleIn schemas.VehicleCreate, orgID uint) (*models.Vehicle, error)
	UpdateVehicle(vehicleID, orgID uint, vehicleIn schemas.VehicleUpdate) (*models.Vehicle, error)
	DeleteVehicle(vehicleID, orgID uint) error
}

type vehicleService struct {
	repo repositories.VehicleRepository
	cache repositories.CacheRepository
}

func NewVehicleService(repo repositories.VehicleRepository, cache repositories.CacheRepository) VehicleService {
	return &vehicleService{repo: repo, cache: cache}
}

func (s *vehicleService) GetVehicles(orgID uint, skip, limit int, search string) ([]models.Vehicle, int64, error) {
	// Caching para listas é mais complexo e pode ser implementado depois.
	// Por enquanto, buscamos diretamente do banco.
	vehicles, err := s.repo.FindByOrganization(orgID, skip, limit, search)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountByOrganization(orgID, search)
	if err != nil {
		return nil, 0, err
	}
	return vehicles, total, nil
}

func (s *vehicleService) GetVehicle(vehicleID, orgID uint) (*models.Vehicle, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("vehicle:%d", vehicleID)

	var vehicle models.Vehicle
	err := s.cache.Get(ctx, cacheKey, &vehicle)
	if err == nil {
		// Cache hit
		return &vehicle, nil
	}

	// Cache miss
	dbVehicle, err := s.repo.FindByID(vehicleID, orgID)
	if err != nil {
		return nil, err
	}
	if dbVehicle == nil {
		return nil, nil // Not found
	}

	// Set cache for future requests
	s.cache.Set(ctx, cacheKey, dbVehicle, 10*time.Minute)

	return dbVehicle, nil
}

func (s *vehicleService) CreateVehicle(vehicleIn schemas.VehicleCreate, orgID uint) (*models.Vehicle, error) {
	vehicle := &models.Vehicle{
		Brand:              vehicleIn.Brand,
		Model:              vehicleIn.Model,
		Year:               vehicleIn.Year,
		LicensePlate:       vehicleIn.LicensePlate,
		Identifier:         vehicleIn.Identifier,
		PhotoURL:           vehicleIn.PhotoURL,
		CurrentKM:          vehicleIn.CurrentKM,
		CurrentEngineHours: vehicleIn.CurrentEngineHours,
		NextMaintenanceDate: vehicleIn.NextMaintenanceDate,
		NextMaintenanceKM:  vehicleIn.NextMaintenanceKM,
		MaintenanceNotes:   vehicleIn.MaintenanceNotes,
		TelemetryDeviceID:  vehicleIn.TelemetryDeviceID,
		OrganizationID:     orgID,
	}
	err := s.repo.Create(vehicle)
	return vehicle, err
}

func (s *vehicleService) UpdateVehicle(vehicleID, orgID uint, vehicleIn schemas.VehicleUpdate) (*models.Vehicle, error) {
	vehicle, err := s.repo.FindByID(vehicleID, orgID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, nil
	}

	if vehicleIn.Brand != nil {
		vehicle.Brand = *vehicleIn.Brand
	}
	// ... (outras atualizações)

	err = s.repo.Update(vehicle)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("vehicle:%d", vehicleID)
	s.cache.Delete(context.Background(), cacheKey)

	return vehicle, nil
}

func (s *vehicleService) DeleteVehicle(vehicleID, orgID uint) error {
	vehicle, err := s.repo.FindByID(vehicleID, orgID)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return nil
	}

	err = s.repo.Delete(vehicle)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("vehicle:%d", vehicleID)
	s.cache.Delete(context.Background(), cacheKey)

	return nil
}
