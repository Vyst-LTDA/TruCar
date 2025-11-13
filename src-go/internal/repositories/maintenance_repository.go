package repositories

import (
	"errors"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type MaintenanceRepository interface {
	FindByID(reqID, orgID uint) (*models.MaintenanceRequest, error)
	FindByOrganization(orgID uint, skip, limit int, search string) ([]models.MaintenanceRequest, error)
	Create(req *models.MaintenanceRequest) error
	Update(req *models.MaintenanceRequest) error
	Delete(req *models.MaintenanceRequest) error
	FindCommentsByRequestID(reqID, orgID uint) ([]models.MaintenanceComment, error)
	CreateComment(comment *models.MaintenanceComment) error
}

type maintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) MaintenanceRepository {
	return &maintenanceRepository{db: db}
}

func (r *maintenanceRepository) FindByID(reqID, orgID uint) (*models.MaintenanceRequest, error) {
	var req models.MaintenanceRequest
	if err := r.db.Where("id = ? AND organization_id = ?", reqID, orgID).First(&req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *maintenanceRepository) FindByOrganization(orgID uint, skip, limit int, search string) ([]models.MaintenanceRequest, error) {
	var reqs []models.MaintenanceRequest
	query := r.db.Where("organization_id = ?", orgID)

	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("problem_description LIKE ?", searchQuery)
	}

	if err := query.Offset(skip).Limit(limit).Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *maintenanceRepository) Create(req *models.MaintenanceRequest) error {
	return r.db.Create(req).Error
}

func (r *maintenanceRepository) Update(req *models.MaintenanceRequest) error {
	return r.db.Save(req).Error
}

func (r *maintenanceRepository) Delete(req *models.MaintenanceRequest) error {
	return r.db.Delete(req).Error
}

func (r *maintenanceRepository) FindCommentsByRequestID(reqID, orgID uint) ([]models.MaintenanceComment, error) {
	var comments []models.MaintenanceComment
	if err := r.db.Where("request_id = ? AND organization_id = ?", reqID, orgID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *maintenanceRepository) CreateComment(comment *models.MaintenanceComment) error {
	return r.db.Create(comment).Error
}
