package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"go-api/internal/models"
)

// VehicleCostRepository define a interface para operações de banco de dados para VehicleCost.
type VehicleCostRepository interface {
	WithTx(tx *gorm.DB) VehicleCostRepository
	FindByOrganizationID(orgID uint, startDate, endDate *time.Time) ([]models.VehicleCost, error)
	Create(cost *models.VehicleCost) error
	FindByFineID(fineID uint) (*models.VehicleCost, error)
	Update(cost *models.VehicleCost) error
	Delete(cost *models.VehicleCost) error
}

type vehicleCostRepository struct {
	db *gorm.DB
}

// NewVehicleCostRepository cria uma nova instância de VehicleCostRepository.
func NewVehicleCostRepository(db *gorm.DB) VehicleCostRepository {
	return &vehicleCostRepository{db: db}
}

func (r *vehicleCostRepository) WithTx(tx *gorm.DB) VehicleCostRepository {
	return &vehicleCostRepository{db: tx}
}

// FindByOrganizationID busca custos de uma organização, com filtro opcional por data.
func (r *vehicleCostRepository) FindByOrganizationID(orgID uint, startDate, endDate *time.Time) ([]models.VehicleCost, error) {
	var costs []models.VehicleCost
	query := r.db.Where("organization_id = ?", orgID)

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	err := query.Order("date desc").Find(&costs).Error
	return costs, err
}

// Create cria um novo registro de custo.
func (r *vehicleCostRepository) Create(cost *models.VehicleCost) error {
	return r.db.Create(cost).Error
}

// FindByFineID busca um custo associado a uma multa.
func (r *vehicleCostRepository) FindByFineID(fineID uint) (*models.VehicleCost, error) {
	var cost models.VehicleCost
	if err := r.db.Where("fine_id = ?", fineID).First(&cost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cost, nil
}

// Update atualiza um registro de custo.
func (r *vehicleCostRepository) Update(cost *models.VehicleCost) error {
	return r.db.Save(cost).Error
}

// Delete remove um registro de custo.
func (r *vehicleCostRepository) Delete(cost *models.VehicleCost) error {
	return r.db.Delete(cost).Error
}
