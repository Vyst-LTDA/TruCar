package services

import (
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type MaintenanceService interface {
	GetMaintenanceRequests(orgID uint, skip, limit int, search string) ([]models.MaintenanceRequest, error)
	GetMaintenanceRequest(reqID, orgID uint) (*models.MaintenanceRequest, error)
	CreateMaintenanceRequest(reqIn schemas.MaintenanceRequestCreate, user models.User) (*models.MaintenanceRequest, error)
	UpdateMaintenanceRequestStatus(reqID uint, reqIn schemas.MaintenanceRequestUpdate, user models.User) (*models.MaintenanceRequest, error)
	DeleteMaintenanceRequest(reqID, orgID uint) error
	GetMaintenanceComments(reqID, orgID uint) ([]models.MaintenanceComment, error)
	CreateMaintenanceComment(commentIn schemas.MaintenanceCommentCreate, reqID uint, user models.User) (*models.MaintenanceComment, error)
}

type maintenanceService struct {
	repo repositories.MaintenanceRepository
}

func NewMaintenanceService(repo repositories.MaintenanceRepository) MaintenanceService {
	return &maintenanceService{repo: repo}
}

func (s *maintenanceService) GetMaintenanceRequests(orgID uint, skip, limit int, search string) ([]models.MaintenanceRequest, error) {
	return s.repo.FindByOrganization(orgID, skip, limit, search)
}

func (s *maintenanceService) GetMaintenanceRequest(reqID, orgID uint) (*models.MaintenanceRequest, error) {
	return s.repo.FindByID(reqID, orgID)
}

func (s *maintenanceService) CreateMaintenanceRequest(reqIn schemas.MaintenanceRequestCreate, user models.User) (*models.MaintenanceRequest, error) {
	req := &models.MaintenanceRequest{
		ProblemDescription: reqIn.ProblemDescription,
		VehicleID:          reqIn.VehicleID,
		Category:           reqIn.Category,
		ReportedByID:       &user.ID,
		OrganizationID:     user.OrganizationID,
	}
	err := s.repo.Create(req)
	return req, err
}

func (s *maintenanceService) UpdateMaintenanceRequestStatus(reqID uint, reqIn schemas.MaintenanceRequestUpdate, user models.User) (*models.MaintenanceRequest, error) {
	req, err := s.repo.FindByID(reqID, user.OrganizationID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, nil // Or return a not found error
	}

	req.Status = reqIn.Status
	req.ManagerNotes = reqIn.ManagerNotes
	req.ApprovedByID = &user.ID

	err = s.repo.Update(req)
	return req, err
}

func (s *maintenanceService) DeleteMaintenanceRequest(reqID, orgID uint) error {
	req, err := s.repo.FindByID(reqID, orgID)
	if err != nil {
		return err
	}
	if req == nil {
		return nil // Or return a not found error
	}
	return s.repo.Delete(req)
}

func (s *maintenanceService) GetMaintenanceComments(reqID, orgID uint) ([]models.MaintenanceComment, error) {
	return s.repo.FindCommentsByRequestID(reqID, orgID)
}

func (s *maintenanceService) CreateMaintenanceComment(commentIn schemas.MaintenanceCommentCreate, reqID uint, user models.User) (*models.MaintenanceComment, error) {
	comment := &models.MaintenanceComment{
		CommentText:    commentIn.CommentText,
		FileURL:        commentIn.FileURL,
		RequestID:      reqID,
		UserID:         &user.ID,
		OrganizationID: user.OrganizationID,
	}
	err := s.repo.CreateComment(comment)
	return comment, err
}
