package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"go-api/internal/models"
)

var ErrVehicleNotAvailable = errors.New("vehicle is not available for a new journey")

type JourneyRepository interface {
	FindByID(journeyID, orgID uint) (*models.Journey, error)
	FindByOrganization(orgID uint, skip, limit int, driverID, vehicleID *uint, dateFrom, dateTo *time.Time) ([]models.Journey, error)
	Create(journey *models.Journey) (*models.Journey, error)
	Update(journey *models.Journey) (*models.Journey, error)
	Delete(journey *models.Journey) error
	CheckVehicleAvailability(vehicleID uint) (bool, error)
	UpdateVehicleStatus(vehicleID uint, status models.VehicleStatus) error
	UpdateVehicleMileage(vehicleID uint, mileage int) error
}

type journeyRepository struct {
	db *gorm.DB
}

func NewJourneyRepository(db *gorm.DB) JourneyRepository {
	return &journeyRepository{db: db}
}

func (r *journeyRepository) FindByID(journeyID, orgID uint) (*models.Journey, error) {
	var journey models.Journey
	if err := r.db.Where("id = ? AND organization_id = ?", journeyID, orgID).First(&journey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &journey, nil
}

func (r *journeyRepository) FindByOrganization(orgID uint, skip, limit int, driverID, vehicleID *uint, dateFrom, dateTo *time.Time) ([]models.Journey, error) {
	var journeys []models.Journey
	query := r.db.Where("organization_id = ?", orgID)

	if driverID != nil {
		query = query.Where("driver_id = ?", *driverID)
	}
	if vehicleID != nil {
		query = query.Where("vehicle_id = ?", *vehicleID)
	}
	if dateFrom != nil {
		query = query.Where("start_time >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("end_time <= ?", *dateTo)
	}

	if err := query.Offset(skip).Limit(limit).Find(&journeys).Error; err != nil {
		return nil, err
	}
	return journeys, nil
}

func (r *journeyRepository) Create(journey *models.Journey) (*models.Journey, error) {
	err := r.db.Create(journey).Error
	return journey, err
}

func (r *journeyRepository) Update(journey *models.Journey) (*models.Journey, error) {
	err := r.db.Save(journey).Error
	return journey, err
}

func (r *journeyRepository) Delete(journey *models.Journey) error {
	return r.db.Delete(journey).Error
}

func (r *journeyRepository) CheckVehicleAvailability(vehicleID uint) (bool, error) {
	var vehicle models.Vehicle
	if err := r.db.First(&vehicle, vehicleID).Error; err != nil {
		return false, err
	}
	return vehicle.Status == models.StatusAvailable, nil
}

func (r *journeyRepository) UpdateVehicleStatus(vehicleID uint, status models.VehicleStatus) error {
	return r.db.Model(&models.Vehicle{}).Where("id = ?", vehicleID).Update("status", status).Error
}

func (r *journeyRepository) UpdateVehicleMileage(vehicleID uint, mileage int) error {
	return r.db.Model(&models.Vehicle{}).Where("id = ?", vehicleID).Update("current_km", mileage).Error
}
