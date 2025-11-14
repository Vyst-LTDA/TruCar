package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type FuelLogRepository interface {
	FindByID(fuelLogID, orgID uint) (*models.FuelLog, error)
	FindByOrganization(orgID uint, skip, limit int) ([]models.FuelLog, error)
	FindByUser(userID, orgID uint, skip, limit int) ([]models.FuelLog, error)
	Create(fuelLog *models.FuelLog) error
	Update(fuelLog *models.FuelLog) error
	Delete(fuelLog *models.FuelLog) error
	SumFuelByOrganization(orgID uint, startDate time.Time) (float64, error)
}

type fuelLogRepository struct {
	db *gorm.DB
}

func NewFuelLogRepository(db *gorm.DB) FuelLogRepository {
	return &fuelLogRepository{db: db}
}

func (r *fuelLogRepository) FindByID(fuelLogID, orgID uint) (*models.FuelLog, error) {
	var fuelLog models.FuelLog
	if err := r.db.Where("id = ? AND organization_id = ?", fuelLogID, orgID).First(&fuelLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &fuelLog, nil
}

func (r *fuelLogRepository) FindByOrganization(orgID uint, skip, limit int) ([]models.FuelLog, error) {
	var fuelLogs []models.FuelLog
	if err := r.db.Where("organization_id = ?", orgID).Offset(skip).Limit(limit).Find(&fuelLogs).Error; err != nil {
		return nil, err
	}
	return fuelLogs, nil
}

func (r *fuelLogRepository) FindByUser(userID, orgID uint, skip, limit int) ([]models.FuelLog, error) {
	var fuelLogs []models.FuelLog
	if err := r.db.Where("user_id = ? AND organization_id = ?", userID, orgID).Offset(skip).Limit(limit).Find(&fuelLogs).Error; err != nil {
		return nil, err
	}
	return fuelLogs, nil
}

func (r *fuelLogRepository) Create(fuelLog *models.FuelLog) error {
	return r.db.Create(fuelLog).Error
}

func (r *fuelLogRepository) Update(fuelLog *models.FuelLog) error {
	return r.db.Save(fuelLog).Error
}

func (r *fuelLogRepository) Delete(fuelLog *models.FuelLog) error {
	return r.db.Delete(fuelLog).Error
}

func (r *fuelLogRepository) SumFuelByOrganization(orgID uint, startDate time.Time) (float64, error) {
	var totalFuel float64
	err := r.db.Model(&models.FuelLog{}).
		Where("organization_id = ? AND date >= ?", orgID, startDate).
		Select("coalesce(sum(liters), 0)").
		Row().
		Scan(&totalFuel)
	return totalFuel, err
}
