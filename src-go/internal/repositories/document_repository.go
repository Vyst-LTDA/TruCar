package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type DocumentRepository interface {
	FindByID(docID, orgID uint) (*models.Document, error)
	FindByOrganization(orgID uint, skip, limit int, expiringInDays *int) ([]models.Document, error)
	Create(doc *models.Document) (*models.Document, error)
	Delete(doc *models.Document) error
	CountByOrganization(orgID uint) (int64, error)
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) FindByID(docID, orgID uint) (*models.Document, error) {
	var doc models.Document
	if err := r.db.Where("id = ? AND organization_id = ?", docID, orgID).First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) FindByOrganization(orgID uint, skip, limit int, expiringInDays *int) ([]models.Document, error) {
	var docs []models.Document
	query := r.db.Where("organization_id = ?", orgID)

	if expiringInDays != nil {
		expiryDate := time.Now().AddDate(0, 0, *expiringInDays)
		query = query.Where("expiry_date <= ?", expiryDate)
	}

	if err := query.Offset(skip).Limit(limit).Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *documentRepository) Create(doc *models.Document) (*models.Document, error) {
	err := r.db.Create(doc).Error
	return doc, err
}

func (r *documentRepository) Delete(doc *models.Document) error {
	return r.db.Delete(doc).Error
}

func (r *documentRepository) CountByOrganization(orgID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Document{}).Where("organization_id = ?", orgID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
