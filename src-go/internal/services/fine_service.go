package services

import (
	"fmt"
	"go-api/internal/models"
	"go-api/internal/repositories"
	"gorm.io/gorm"
)

// FineService define a interface para a lógica de negócios de multas.
type FineService interface {
	GetFineByID(fineID, orgID uint) (*models.Fine, error)
	GetFinesByOrganization(orgID uint, skip, limit int) ([]models.Fine, error)
	GetFinesByDriver(driverID, orgID uint, skip, limit int) ([]models.Fine, error)
	CreateFine(fine *models.Fine, currentUser models.User) (*models.Fine, error)
	UpdateFine(fineID, orgID uint, payload map[string]interface{}) (*models.Fine, error)
	DeleteFine(fineID, orgID uint) error
}

type fineService struct {
	db              *gorm.DB
	fineRepo        repositories.FineRepository
	costRepo        repositories.VehicleCostRepository
	notificationSvc NotificationService
}

// NewFineService cria uma nova instância de FineService.
func NewFineService(
	db *gorm.DB,
	fineRepo repositories.FineRepository,
	costRepo repositories.VehicleCostRepository,
	notificationSvc NotificationService,
) FineService {
	return &fineService{
		db:              db,
		fineRepo:        fineRepo,
		costRepo:        costRepo,
		notificationSvc: notificationSvc,
	}
}

func (s *fineService) GetFineByID(fineID, orgID uint) (*models.Fine, error) {
	return s.fineRepo.FindByID(fineID, orgID)
}

func (s *fineService) GetFinesByOrganization(orgID uint, skip, limit int) ([]models.Fine, error) {
	return s.fineRepo.FindByOrganization(orgID, skip, limit)
}

func (s *fineService) GetFinesByDriver(driverID, orgID uint, skip, limit int) ([]models.Fine, error) {
	return s.fineRepo.FindByDriver(driverID, orgID, skip, limit)
}

func (s *fineService) CreateFine(fine *models.Fine, currentUser models.User) (*models.Fine, error) {
	var createdFine *models.Fine

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.fineRepo.WithTx(tx)
		costTxRepo := s.costRepo.WithTx(tx)

		// 1. Criar a multa
		var err error
		createdFine, err = txRepo.Create(fine)
		if err != nil {
			return err
		}

		// 2. Criar o custo associado
		cost := &models.VehicleCost{
			Description:    fmt.Sprintf("Custo da multa: %s", createdFine.Description),
			Amount:         createdFine.Value,
			Date:           createdFine.Date,
			CostType:       models.CostTypeMulta,
			VehicleID:      createdFine.VehicleID,
			OrganizationID: createdFine.OrganizationID,
			FineID:         &createdFine.ID,
		}
		if err := costTxRepo.Create(cost); err != nil {
			return fmt.Errorf("falha ao criar custo associado: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 3. Disparar notificação em background (goroutine)
	go func() {
		message := fmt.Sprintf("Nova multa de R$%.2f registrada para o veículo.", createdFine.Value)
		notification := &models.Notification{
			OrganizationID:    createdFine.OrganizationID,
			UserID:            0,
			Message:           message,
			IsRead:            false,
			NotificationType:  models.NotificationTypeNewFineRegistered,
			RelatedEntityType: "fine",
			RelatedEntityID:   &createdFine.ID,
			RelatedVehicleID:  &createdFine.VehicleID,
		}
		s.notificationSvc.CreateNotificationAsync(notification)
	}()

	return createdFine, nil
}

func (s *fineService) UpdateFine(fineID, orgID uint, payload map[string]interface{}) (*models.Fine, error) {
	var updatedFine *models.Fine

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.fineRepo.WithTx(tx)
		costTxRepo := s.costRepo.WithTx(tx)

		dbFine, err := txRepo.FindByID(fineID, orgID)
		if err != nil || dbFine == nil {
			return fmt.Errorf("multa não encontrada")
		}

		// Mapear o payload
		if val, ok := payload["description"].(string); ok {
			dbFine.Description = val
		}
		if val, ok := payload["value"].(float64); ok {
			dbFine.Value = val
		}
		// ... outros campos ...

		// Atualizar o custo associado
		cost, err := costTxRepo.FindByFineID(fineID)
		if err != nil {
			fmt.Printf("Erro ao buscar custo associado para a multa %d: %v\n", fineID, err)
		}
		if cost != nil {
			costUpdated := false
			if val, ok := payload["value"].(float64); ok {
				cost.Amount = val
				costUpdated = true
			}
			if val, ok := payload["description"].(string); ok {
				cost.Description = fmt.Sprintf("Custo da multa: %s", val)
				costUpdated = true
			}
			if costUpdated {
				if err := costTxRepo.Update(cost); err != nil {
					return fmt.Errorf("falha ao atualizar custo associado: %w", err)
				}
			}
		}

		updatedFine, err = txRepo.Update(dbFine)
		return err
	})

	return updatedFine, err
}

func (s *fineService) DeleteFine(fineID, orgID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.fineRepo.WithTx(tx)
		costTxRepo := s.costRepo.WithTx(tx)

		dbFine, err := txRepo.FindByID(fineID, orgID)
		if err != nil || dbFine == nil {
			return fmt.Errorf("multa não encontrada")
		}

		// Deletar o custo associado
		cost, err := costTxRepo.FindByFineID(fineID)
		if err != nil {
			fmt.Printf("Erro ao buscar custo associado para a multa %d: %v\n", fineID, err)
		}
		if cost != nil {
			if err := costTxRepo.Delete(cost); err != nil {
				return fmt.Errorf("falha ao deletar custo associado: %w", err)
			}
		}

		return txRepo.Delete(dbFine)
	})
}
