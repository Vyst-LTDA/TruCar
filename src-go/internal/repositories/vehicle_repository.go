package repositories

import (
	"errors"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type VehicleRepository interface {
	FindByID(vehicleID, orgID uint) (*models.Vehicle, error)
	FindByOrganization(orgID uint, skip, limit int, search string) ([]models.Vehicle, error)
	CountByOrganization(orgID uint, search string) (int64, error)
	Create(vehicle *models.Vehicle) error
	Update(vehicle *models.Vehicle) error
	Delete(vehicle *models.Vehicle) error
}

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) FindByID(vehicleID, orgID uint) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	if err := r.db.Where("id = ? AND organization_id = ?", vehicleID, orgID).First(&vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &vehicle, nil
}

func (r *vehicleRepository) FindByOrganization(orgID uint, skip, limit int, search string) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	query := r.db.Where("organization_id = ?", orgID)

	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("brand LIKE ? OR model LIKE ? OR license_plate LIKE ? OR identifier LIKE ?", searchQuery, searchQuery, searchQuery, searchQuery)
	}

	if err := query.Offset(skip).Limit(limit).Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (r *vehicleRepository) CountByOrganization(orgID uint, search string) (int64, error) {
	var count int64
	query := r.db.Model(&models.Vehicle{}).Where("organization_id = ?", orgID)

	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("brand LIKE ? OR model LIKE ? OR license_plate LIKE ? OR identifier LIKE ?", searchQuery, searchQuery, searchQuery, searchQuery)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *vehicleRepository) Create(vehicle *models.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *vehicleRepository) Update(vehicle *models.Vehicle) error {
	return r.db.Save(vehicle).Error
}

func (r *vehicleRepository) Delete(vehicle *models.Vehicle) error {
	return r.db.Delete(vehicle).Error
}
