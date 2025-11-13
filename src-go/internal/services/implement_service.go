package services

import (
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type ImplementService interface {
	GetImplements(orgID uint, managementList bool) ([]models.Implement, error)
	GetImplement(implementID, orgID uint) (*models.Implement, error)
	CreateImplement(implementIn schemas.ImplementCreate, orgID uint) (*models.Implement, error)
	UpdateImplement(implementID, orgID uint, implementIn schemas.ImplementUpdate) (*models.Implement, error)
	DeleteImplement(implementID, orgID uint) error
}

type implementService struct {
	repo repositories.ImplementRepository
}

func NewImplementService(repo repositories.ImplementRepository) ImplementService {
	return &implementService{repo: repo}
}

func (s *implementService) GetImplements(orgID uint, managementList bool) ([]models.Implement, error) {
	return s.repo.FindByOrganization(orgID, managementList)
}

func (s *implementService) GetImplement(implementID, orgID uint) (*models.Implement, error) {
	return s.repo.FindByID(implementID, orgID)
}

func (s *implementService) CreateImplement(implementIn schemas.ImplementCreate, orgID uint) (*models.Implement, error) {
	implement := &models.Implement{
		Name:           implementIn.Name,
		Brand:          implementIn.Brand,
		VehicleModel:   implementIn.VehicleModel,
		Year:           implementIn.Year,
		Identifier:     implementIn.Identifier,
		Type:           implementIn.Type,
		OrganizationID: orgID,
	}
	err := s.repo.Create(implement)
	return implement, err
}

func (s *implementService) UpdateImplement(implementID, orgID uint, implementIn schemas.ImplementUpdate) (*models.Implement, error) {
	implement, err := s.repo.FindByID(implementID, orgID)
	if err != nil {
		return nil, err
	}
	if implement == nil {
		return nil, nil // Or return a not found error
	}

	if implementIn.Name != "" {
		implement.Name = implementIn.Name
	}
	if implementIn.Brand != "" {
		implement.Brand = implementIn.Brand
	}
	if implementIn.VehicleModel != "" {
		implement.VehicleModel = implementIn.VehicleModel
	}
	if implementIn.Year != 0 {
		implement.Year = implementIn.Year
	}
	if implementIn.Identifier != "" {
		implement.Identifier = implementIn.Identifier
	}
	if implementIn.Type != "" {
		implement.Type = implementIn.Type
	}
	if implementIn.Status != "" {
		implement.Status = models.ImplementStatus(implementIn.Status)
	}

	err = s.repo.Update(implement)
	return implement, err
}

func (s *implementService) DeleteImplement(implementID, orgID uint) error {
	implement, err := s.repo.FindByID(implementID, orgID)
	if err != nil {
		return err
	}
	if implement == nil {
		return nil // Or return a not found error
	}
	return s.repo.Delete(implement)
}
