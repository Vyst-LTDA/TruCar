package repositories

import (
	"gorm.io/gorm"
	"go-api/internal/models"
)

type OrganizationRepository interface {
	FindAll(skip, limit int, status *string) ([]models.Organization, error)
	Update(org *models.Organization) (*models.Organization, error)
	FindByID(orgID uint) (*models.Organization, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) FindAll(skip, limit int, status *string) ([]models.Organization, error) {
	var orgs []models.Organization
	query := r.db.Offset(skip).Limit(limit)
	if status != nil {
		// Esta lógica assume que o status de uma organização pode ser determinado
		// a partir de seus usuários. A implementação exata pode precisar ser ajustada
		// com base na lógica de negócios real.
	}
	if err := query.Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *organizationRepository) Update(org *models.Organization) (*models.Organization, error) {
	err := r.db.Save(org).Error
	return org, err
}

func (r *organizationRepository) FindByID(orgID uint) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.First(&org, orgID).Error; err != nil {
		return nil, err
	}
	return &org, nil
}
