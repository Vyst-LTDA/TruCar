package services

import (
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type OrganizationService interface {
	GetOrganizations(skip, limit int, status *string) ([]models.Organization, error)
	UpdateOrganization(orgID uint, orgIn schemas.OrganizationUpdate) (*models.Organization, error)
}

type organizationService struct {
	repo repositories.OrganizationRepository
}

func NewOrganizationService(repo repositories.OrganizationRepository) OrganizationService {
	return &organizationService{repo: repo}
}

func (s *organizationService) GetOrganizations(skip, limit int, status *string) ([]models.Organization, error) {
	return s.repo.FindAll(skip, limit, status)
}

func (s *organizationService) UpdateOrganization(orgID uint, orgIn schemas.OrganizationUpdate) (*models.Organization, error) {
	org, err := s.repo.FindByID(orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, nil // Not found
	}

	if orgIn.Name != "" {
		org.Name = orgIn.Name
	}
	if orgIn.Sector != "" {
		org.Sector = models.Sector(orgIn.Sector)
	}

	return s.repo.Update(org)
}
