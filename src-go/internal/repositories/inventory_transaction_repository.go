package repositories

import (
	"gorm.io/gorm"

	"go-api/internal/models"
)

type InventoryTransactionRepository interface {
	Create(transaction *models.InventoryTransaction) error
	FindByPartID(partID uint, skip, limit int) ([]models.InventoryTransaction, error)
}

type inventoryTransactionRepository struct {
	db *gorm.DB
}

func NewInventoryTransactionRepository(db *gorm.DB) InventoryTransactionRepository {
	return &inventoryTransactionRepository{db: db}
}

func (r *inventoryTransactionRepository) Create(transaction *models.InventoryTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *inventoryTransactionRepository) FindByPartID(partID uint, skip, limit int) ([]models.InventoryTransaction, error) {
	var transactions []models.InventoryTransaction
	if err := r.db.Where("part_id = ?", partID).Offset(skip).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
