package repositories

import (
	"errors"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type PartRepository interface {
	FindByID(partID, orgID uint) (*models.Part, error)
	FindByOrganization(orgID uint, search string, skip, limit int) ([]models.Part, error)
	Create(part *models.Part) (*models.Part, error)
	Update(part *models.Part) (*models.Part, error)
	Delete(part *models.Part) error
	FindItemByID(itemID, orgID uint) (*models.InventoryItem, error)
	CreateItem(item *models.InventoryItem) (*models.InventoryItem, error)
	UpdateItem(item *models.InventoryItem) (*models.InventoryItem, error)
	FindItemsByPartID(partID uint, status *models.InventoryItemStatus) ([]models.InventoryItem, error)
	CountByOrganization(orgID uint) (int64, error)
}

type partRepository struct {
	db *gorm.DB
}

func NewPartRepository(db *gorm.DB) PartRepository {
	return &partRepository{db: db}
}

func (r *partRepository) FindByID(partID, orgID uint) (*models.Part, error) {
	var part models.Part
	if err := r.db.Where("id = ? AND organization_id = ?", partID, orgID).First(&part).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &part, nil
}

func (r *partRepository) FindByOrganization(orgID uint, search string, skip, limit int) ([]models.Part, error) {
	var parts []models.Part
	query := r.db.Where("organization_id = ?", orgID)
	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("name LIKE ? OR part_number LIKE ? OR brand LIKE ?", searchQuery, searchQuery, searchQuery)
	}
	if err := query.Offset(skip).Limit(limit).Find(&parts).Error; err != nil {
		return nil, err
	}
	return parts, nil
}

func (r *partRepository) Create(part *models.Part) (*models.Part, error) {
	err := r.db.Create(part).Error
	return part, err
}

func (r *partRepository) Update(part *models.Part) (*models.Part, error) {
	err := r.db.Save(part).Error
	return part, err
}

func (r *partRepository) Delete(part *models.Part) error {
	return r.db.Delete(part).Error
}

func (r *partRepository) FindItemByID(itemID, orgID uint) (*models.InventoryItem, error) {
	var item models.InventoryItem
	if err := r.db.Where("id = ? AND organization_id = ?", itemID, orgID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *partRepository) CreateItem(item *models.InventoryItem) (*models.InventoryItem, error) {
	err := r.db.Create(item).Error
	return item, err
}

func (r *partRepository) UpdateItem(item *models.InventoryItem) (*models.InventoryItem, error) {
	err := r.db.Save(item).Error
	return item, err
}

func (r *partRepository) FindItemsByPartID(partID uint, status *models.InventoryItemStatus) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	query := r.db.Where("part_id = ?", partID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *partRepository) CountByOrganization(orgID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Part{}).Where("organization_id = ?", orgID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
