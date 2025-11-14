package services

import (
	"go-api/internal/logging"
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
	"time"
	"go.uber.org/zap"
)

// GPSService define a interface para a lógica de negócios de GPS.
type GPSService interface {
	ProcessPing(ping schemas.LocationCreate, orgID uint)
}

type gpsService struct {
	vehicleRepo repositories.VehicleRepository
	historyRepo repositories.LocationHistoryRepository
}

// NewGPSService cria uma nova instância de GPSService.
func NewGPSService(
	vehicleRepo repositories.VehicleRepository,
	historyRepo repositories.LocationHistoryRepository,
) GPSService {
	return &gpsService{
		vehicleRepo: vehicleRepo,
		historyRepo: historyRepo,
	}
}

// ProcessPing processa um novo ping de localização de forma assíncrona.
func (s *gpsService) ProcessPing(ping schemas.LocationCreate, orgID uint) {
	go func() {
		// 1. Atualizar a localização atual do veículo
		vehicle, err := s.vehicleRepo.FindByID(ping.VehicleID, orgID)
		if err != nil || vehicle == nil {
			logging.Logger.Error("Failed to find vehicle for GPS ping", zap.Error(err), zap.Uint("vehicleID", ping.VehicleID), zap.Uint("orgID", orgID))
			return
		}

		now := time.Now()
		vehicle.LastLatitude = &ping.Latitude
		vehicle.LastLongitude = &ping.Longitude
		vehicle.LastLocationUpdate = &now

		if err := s.vehicleRepo.Update(vehicle); err != nil {
			logging.Logger.Error("Failed to update vehicle location", zap.Error(err), zap.Uint("vehicleID", vehicle.ID))
			return
		}

		// 2. Criar um registro no histórico de localização
		history := &models.LocationHistory{
			VehicleID: ping.VehicleID,
			Latitude:  ping.Latitude,
			Longitude: ping.Longitude,
			Timestamp: now,
		}

		if err := s.historyRepo.Create(history); err != nil {
			logging.Logger.Error("Failed to create location history", zap.Error(err), zap.Uint("vehicleID", vehicle.ID))
		}
	}()
}
