package repositories

import (
	"errors"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type FreightOrderRepository interface {
	FindByID(orderID, orgID uint) (*models.FreightOrder, error)
	FindByOrganization(orgID uint, skip, limit int) ([]models.FreightOrder, error)
	FindByStatus(orgID uint, status models.FreightStatus) ([]models.FreightOrder, error)
	FindPendingByDriver(driverID, orgID uint) ([]models.FreightOrder, error)
	CreateWithStops(order *models.FreightOrder) (*models.FreightOrder, error)
	Update(order *models.FreightOrder) (*models.FreightOrder, error)
	UpdateStopPoint(stopPoint *models.StopPoint) (*models.StopPoint, error)
}

type freightOrderRepository struct {
	db *gorm.DB
}

func NewFreightOrderRepository(db *gorm.DB) FreightOrderRepository {
	return &freightOrderRepository{db: db}
}

func (r *freightOrderRepository) FindByID(orderID, orgID uint) (*models.FreightOrder, error) {
	var order models.FreightOrder
	if err := r.db.Preload("StopPoints").Preload("Vehicle").Preload("Driver").Where("id = ? AND organization_id = ?", orderID, orgID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *freightOrderRepository) FindByOrganization(orgID uint, skip, limit int) ([]models.FreightOrder, error) {
	var orders []models.FreightOrder
	if err := r.db.Preload("StopPoints").Where("organization_id = ?", orgID).Offset(skip).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *freightOrderRepository) FindByStatus(orgID uint, status models.FreightStatus) ([]models.FreightOrder, error) {
	var orders []models.FreightOrder
	if err := r.db.Preload("StopPoints").Where("organization_id = ? AND status = ?", orgID, status).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *freightOrderRepository) FindPendingByDriver(driverID, orgID uint) ([]models.FreightOrder, error) {
	var orders []models.FreightOrder
	statuses := []models.FreightStatus{models.FreightStatusClaimed, models.FreightStatusInTransit}
	if err := r.db.Preload("StopPoints").Where("driver_id = ? AND organization_id = ? AND status IN (?)", driverID, orgID, statuses).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *freightOrderRepository) CreateWithStops(order *models.FreightOrder) (*models.FreightOrder, error) {
	err := r.db.Create(order).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(order.ID, order.OrganizationID)
}

func (r *freightOrderRepository) Update(order *models.FreightOrder) (*models.FreightOrder, error) {
	err := r.db.Save(order).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(order.ID, order.OrganizationID)
}

func (r *freightOrderRepository) UpdateStopPoint(stopPoint *models.StopPoint) (*models.StopPoint, error) {
	err := r.db.Save(stopPoint).Error
	return stopPoint, err
}
