package repositories

import (
	"errors"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type ImplementRepository interface {
	FindByID(implementID, orgID uint) (*models.Implement, error)
	FindByOrganization(orgID uint, managementList bool) ([]models.Implement, error)
	Create(implement *models.Implement) error
	Update(implement *models.Implement) error
	Delete(implement *models.Implement) error
}

type implementRepository struct {
	db *gorm.DB
}

func NewImplementRepository(db *gorm.DB) ImplementRepository {
	return &implementRepository{db: db}
}

func (r *implementRepository) FindByID(implementID, orgID uint) (*models.Implement, error) {
	var implement models.Implement
	if err := r.db.Where("id = ? AND organization_id = ?", implementID, orgID).First(&implement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &implement, nil
}

func (r *implementRepository) FindByOrganization(orgID uint, managementList bool) ([]models.Implement, error) {
	var implements []models.Implement
	query := r.db.Where("organization_id = ?", orgID)

	if !managementList {
		query = query.Where("status = ?", models.ImplementStatusAvailable)
	}

	if err := query.Find(&implements).Error; err != nil {
		return nil, err
	}
	return implements, nil
}

func (r *implementRepository) Create(implement *models.Implement) error {
	return r.db.Create(implement).Error
}

func (r *implementRepository) Update(implement *models.Implement) error {
	return r.db.Save(implement).Error
}

func (r *implementRepository) Delete(implement *models.Implement) error {
	return r.db.Delete(implement).Error
}
