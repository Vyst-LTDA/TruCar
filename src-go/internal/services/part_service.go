package services

import (
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type PartService interface {
	GetParts(orgID uint, search string, skip, limit int) ([]models.Part, error)
	GetPart(partID, orgID uint) (*models.Part, error)
	CreatePart(partIn schemas.PartCreate, orgID uint, userID uint) (*models.Part, error)
	UpdatePart(partID uint, partIn schemas.PartUpdate, orgID uint) (*models.Part, error)
	DeletePart(partID, orgID uint) error
	AddInventoryItems(partID uint, payload schemas.AddItemsPayload, orgID uint, userID uint) error
	SetInventoryItemStatus(itemID uint, payload schemas.SetItemStatusPayload, orgID uint, userID uint) (*models.InventoryItem, error)
	GetItemsForPart(partID uint, status *models.InventoryItemStatus, orgID uint) ([]models.InventoryItem, error)
	GetPartHistory(partID, orgID uint, skip, limit int) ([]models.InventoryTransaction, error)
}

type partService struct {
	partRepo         repositories.PartRepository
	transactionRepo  repositories.InventoryTransactionRepository
	notificationService NotificationService
}

func NewPartService(partRepo repositories.PartRepository, transactionRepo repositories.InventoryTransactionRepository, notificationService NotificationService) PartService {
	return &partService{partRepo: partRepo, transactionRepo: transactionRepo, notificationService: notificationService}
}

func (s *partService) GetParts(orgID uint, search string, skip, limit int) ([]models.Part, error) {
	return s.partRepo.FindByOrganization(orgID, search, skip, limit)
}

func (s *partService) GetPart(partID, orgID uint) (*models.Part, error) {
	return s.partRepo.FindByID(partID, orgID)
}

func (s *partService) CreatePart(partIn schemas.PartCreate, orgID uint, userID uint) (*models.Part, error) {
	part := &models.Part{
		Name:           partIn.Name,
		Category:       models.PartCategory(partIn.Category),
		Value:          partIn.Value,
		SerialNumber:   partIn.SerialNumber,
		MinimumStock:   partIn.MinimumStock,
		PartNumber:     partIn.PartNumber,
		Brand:          partIn.Brand,
		Location:       partIn.Location,
		Notes:          partIn.Notes,
		LifespanKM:     partIn.LifespanKM,
		OrganizationID: orgID,
	}
	createdPart, err := s.partRepo.Create(part)
	if err != nil {
		return nil, err
	}

	if partIn.InitialQuantity > 0 {
		addItemsPayload := schemas.AddItemsPayload{
			Quantity: partIn.InitialQuantity,
			Notes:    "Ajuste inicial",
		}
		err = s.AddInventoryItems(createdPart.ID, addItemsPayload, orgID, userID)
	}

	return createdPart, err
}

func (s *partService) UpdatePart(partID uint, partIn schemas.PartUpdate, orgID uint) (*models.Part, error) {
	part, err := s.partRepo.FindByID(partID, orgID)
	if err != nil {
		return nil, err
	}
	if part == nil {
		return nil, nil // Not found
	}

	// Update fields

	return s.partRepo.Update(part)
}

func (s *partService) DeletePart(partID, orgID uint) error {
	part, err := s.partRepo.FindByID(partID, orgID)
	if err != nil {
		return err
	}
	if part == nil {
		return nil // Not found
	}
	return s.partRepo.Delete(part)
}

func (s *partService) AddInventoryItems(partID uint, payload schemas.AddItemsPayload, orgID uint, userID uint) error {
	for i := 0; i < payload.Quantity; i++ {
		item := &models.InventoryItem{
			PartID:         partID,
			OrganizationID: orgID,
		}
		_, err := s.partRepo.CreateItem(item)
		if err != nil {
			return err
		}

		transaction := &models.InventoryTransaction{
			ItemID:          item.ID,
			PartID:          &partID,
			UserID:          &userID,
			TransactionType: models.TransactionTypeEntrada,
			Notes:           &payload.Notes,
		}
		err = s.transactionRepo.Create(transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *partService) SetInventoryItemStatus(itemID uint, payload schemas.SetItemStatusPayload, orgID uint, userID uint) (*models.InventoryItem, error) {
	item, err := s.partRepo.FindItemByID(itemID, orgID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil // Not found
	}

	item.Status = models.InventoryItemStatus(payload.NewStatus)
	item.InstalledOnVehicleID = payload.RelatedVehicleID

	updatedItem, err := s.partRepo.UpdateItem(item)
	if err != nil {
		return nil, err
	}

	transaction := &models.InventoryTransaction{
		ItemID:           itemID,
		PartID:           &item.PartID,
		UserID:           &userID,
		TransactionType:  models.TransactionType(payload.NewStatus),
		Notes:            payload.Notes,
		RelatedVehicleID: payload.RelatedVehicleID,
	}
	err = s.transactionRepo.Create(transaction)

	// Check stock and send notification

	return updatedItem, err
}

func (s *partService) GetItemsForPart(partID uint, status *models.InventoryItemStatus, orgID uint) ([]models.InventoryItem, error) {
	return s.partRepo.FindItemsByPartID(partID, status)
}

func (s *partService) GetPartHistory(partID, orgID uint, skip, limit int) ([]models.InventoryTransaction, error) {
	return s.transactionRepo.FindByPartID(partID, skip, limit)
}
