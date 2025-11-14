package services

import (
	"errors"
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

var ErrFreightOrderNotFound = errors.New("freight order not found")
var ErrVehicleNotFound = errors.New("vehicle not found")
var ErrStopPointNotFound = errors.New("stop point not found")
var ErrJourneyNotFound = errors.New("journey not found")
var ErrFreightNotAssigned = errors.New("freight not assigned to this driver")

type FreightOrderService interface {
	GetFreightOrders(user models.User, skip, limit int) ([]models.FreightOrder, error)
	GetOpenFreightOrders(orgID uint) ([]models.FreightOrder, error)
	GetMyPendingFreightOrders(driverID, orgID uint) ([]models.FreightOrder, error)
	GetFreightOrderByID(orderID, orgID uint) (*models.FreightOrder, error)
	CreateFreightOrder(orderIn schemas.FreightOrderCreate, orgID uint) (*models.FreightOrder, error)
	ClaimFreightOrder(orderID uint, claimIn schemas.FreightOrderClaim, driver models.User) (*models.FreightOrder, error)
	StartJourneyForStop(orderID, stopPointID, driverID, orgID uint) (*models.Journey, error)
	CompleteStopPoint(orderID, stopPointID, driverID, orgID uint, journeyID uint, endMileage int) (*models.StopPoint, error)
}

type freightOrderService struct {
	freightOrderRepo repositories.FreightOrderRepository
	vehicleRepo      repositories.VehicleRepository
	journeyService   JourneyService
}

func NewFreightOrderService(freightOrderRepo repositories.FreightOrderRepository, vehicleRepo repositories.VehicleRepository, journeyService JourneyService) FreightOrderService {
	return &freightOrderService{freightOrderRepo: freightOrderRepo, vehicleRepo: vehicleRepo, journeyService: journeyService}
}

func (s *freightOrderService) GetFreightOrders(user models.User, skip, limit int) ([]models.FreightOrder, error) {
	if user.Role == models.RoleClienteAtivo || user.Role == models.RoleClienteDemo {
		return s.freightOrderRepo.FindByOrganization(user.OrganizationID, skip, limit)
	}
	return []models.FreightOrder{}, nil // Drivers don't see all orders
}

func (s *freightOrderService) GetOpenFreightOrders(orgID uint) ([]models.FreightOrder, error) {
	return s.freightOrderRepo.FindByStatus(orgID, models.FreightStatusOpen)
}

func (s *freightOrderService) GetMyPendingFreightOrders(driverID, orgID uint) ([]models.FreightOrder, error) {
	return s.freightOrderRepo.FindPendingByDriver(driverID, orgID)
}

func (s *freightOrderService) GetFreightOrderByID(orderID, orgID uint) (*models.FreightOrder, error) {
	return s.freightOrderRepo.FindByID(orderID, orgID)
}

func (s *freightOrderService) CreateFreightOrder(orderIn schemas.FreightOrderCreate, orgID uint) (*models.FreightOrder, error) {
	var stopPoints []models.StopPoint
	for _, sp := range orderIn.StopPoints {
		stopPoints = append(stopPoints, models.StopPoint{
			SequenceOrder:    sp.SequenceOrder,
			Type:             models.StopPointType(sp.Type),
			Address:          sp.Address,
			CargoDescription: sp.CargoDescription,
			ScheduledTime:    sp.ScheduledTime,
		})
	}

	order := &models.FreightOrder{
		Description:        orderIn.Description,
		ScheduledStartTime: orderIn.ScheduledStartTime,
		ScheduledEndTime:   orderIn.ScheduledEndTime,
		ClientID:           orderIn.ClientID,
		OrganizationID:     orgID,
		StopPoints:         stopPoints,
	}
	return s.freightOrderRepo.CreateWithStops(order)
}

func (s *freightOrderService) ClaimFreightOrder(orderID uint, claimIn schemas.FreightOrderClaim, driver models.User) (*models.FreightOrder, error) {
	order, err := s.freightOrderRepo.FindByID(orderID, driver.OrganizationID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrFreightOrderNotFound
	}

	vehicle, err := s.vehicleRepo.FindByID(claimIn.VehicleID, driver.OrganizationID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, ErrVehicleNotFound
	}

	order.Status = models.FreightStatusClaimed
	order.DriverID = &driver.ID
	order.VehicleID = &vehicle.ID

	return s.freightOrderRepo.Update(order)
}

func (s *freightOrderService) StartJourneyForStop(orderID, stopPointID, driverID, orgID uint) (*models.Journey, error) {
	order, err := s.freightOrderRepo.FindByID(orderID, orgID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrFreightOrderNotFound
	}
	if order.DriverID == nil || *order.DriverID != driverID {
		return nil, ErrFreightNotAssigned
	}

	// A lógica para encontrar o ponto de parada e iniciar a jornada será implementada no futuro.
	// Por enquanto, vamos retornar uma jornada vazia.
	return &models.Journey{}, nil
}

func (s *freightOrderService) CompleteStopPoint(orderID, stopPointID, driverID, orgID uint, journeyID uint, endMileage int) (*models.StopPoint, error) {
	order, err := s.freightOrderRepo.FindByID(orderID, orgID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrFreightOrderNotFound
	}
	if order.DriverID == nil || *order.DriverID != driverID {
		return nil, ErrFreightNotAssigned
	}

	// A lógica para encontrar o ponto de parada e completar a jornada será implementada no futuro.
	// Por enquanto, vamos retornar um ponto de parada vazio.
	return &models.StopPoint{}, nil
}
