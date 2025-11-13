package repositories

import (
	"errors"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type FineRepository interface {
	FindByID(fineID, orgID uint) (*models.Fine, error)
	FindByOrganization(orgID uint, skip, limit int) ([]models.Fine, error)
	FindByDriver(driverID, orgID uint, skip, limit int) ([]models.Fine, error)
	Create(fine *models.Fine) (*models.Fine, error)
	Update(fine *models.Fine) (*models.Fine, error)
	Delete(fine *models.Fine) error
}

type fineRepository struct {
	db *gorm.DB
}

func NewFineRepository(db *gorm.DB) FineRepository {
	return &fineRepository{db: db}
}

func (r *fineRepository) FindByID(fineID, orgID uint) (*models.Fine, error) {
	var fine models.Fine
	if err := r.db.Preload("Vehicle").Preload("Driver").Where("id = ? AND organization_id = ?", fineID, orgID).First(&fine).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &fine, nil
}

func (r *fineRepository) FindByOrganization(orgID uint, skip, limit int) ([]models.Fine, error) {
	var fines []models.Fine
	if err := r.db.Preload("Vehicle").Preload("Driver").Where("organization_id = ?", orgID).Offset(skip).Limit(limit).Find(&fines).Error; err != nil {
		return nil, err
	}
	return fines, nil
}

func (r *fineRepository) FindByDriver(driverID, orgID uint, skip, limit int) ([]models.Fine, error) {
	var fines []models.Fine
	if err := r.db.Preload("Vehicle").Preload("Driver").Where("driver_id = ? AND organization_id = ?", driverID, orgID).Offset(skip).Limit(limit).Find(&fines).Error; err != nil {
		return nil, err
	}
	return fines, nil
}

func (r *fineRepository) Create(fine *models.Fine) (*models.Fine, error) {
	err := r.db.Create(fine).Error
	if err != nil {
		return nil, err
	}
	// Recarregar para obter os dados pré-carregados
	return r.FindByID(fine.ID, fine.OrganizationID)
}

func (r *fineRepository) Update(fine *models.Fine) (*models.Fine, error) {
	err := r.db.Save(fine).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(fine.ID, fine.OrganizationID)
}

func (r *fineRepository) Delete(fine *models.Fine) error {
	return r.db.Delete(fine).Error
}
