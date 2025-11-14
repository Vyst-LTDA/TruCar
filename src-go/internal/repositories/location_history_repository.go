package repositories

import (
	"go-api/internal/models"
	"gorm.io/gorm"
)

// LocationHistoryRepository define a interface para operações de banco de dados para LocationHistory.
type LocationHistoryRepository interface {
	Create(history *models.LocationHistory) error
}

type locationHistoryRepository struct {
	db *gorm.DB
}

// NewLocationHistoryRepository cria uma nova instância de LocationHistoryRepository.
func NewLocationHistoryRepository(db *gorm.DB) LocationHistoryRepository {
	return &locationHistoryRepository{db: db}
}

// Create cria um novo registro de histórico de localização.
func (r *locationHistoryRepository) Create(history *models.LocationHistory) error {
	return r.db.Create(history).Error
}
