package services

import (
	"time"

	"go-api/internal/models"
	"go-api/internal/repositories"
)

// VehicleCostService define a interface para a lógica de negócios relacionada a custos de veículos.
type VehicleCostService interface {
	GetCostsByOrganization(orgID uint, startDate, endDate *time.Time) ([]models.VehicleCost, error)
}

type vehicleCostService struct {
	repo repositories.VehicleCostRepository
}

// NewVehicleCostService cria uma nova instância de VehicleCostService.
func NewVehicleCostService(repo repositories.VehicleCostRepository) VehicleCostService {
	return &vehicleCostService{repo: repo}
}

// GetCostsByOrganization busca os custos de uma organização, aplicando a lógica de negócios necessária.
func (s *vehicleCostService) GetCostsByOrganization(orgID uint, startDate, endDate *time.Time) ([]models.VehicleCost, error) {
	// Aqui, a lógica de negócios pode ser adicionada no futuro (ex: validações, enriquecimento de dados).
	return s.repo.FindByOrganizationID(orgID, startDate, endDate)
}
