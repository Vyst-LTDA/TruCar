package services

import (
	"mime/multipart"

	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
	"go-api/internal/storage"
)

type DocumentService interface {
	GetDocuments(orgID uint, skip, limit int, expiringInDays *int) ([]models.Document, error)
	CreateDocument(docIn schemas.DocumentCreate, file *multipart.FileHeader, orgID uint) (*models.Document, error)
	DeleteDocument(docID, orgID uint) (bool, error)
}

type documentService struct {
	repo         repositories.DocumentRepository
	storageService storage.FileStorageService
}

func NewDocumentService(repo repositories.DocumentRepository, storageService storage.FileStorageService) DocumentService {
	return &documentService{repo: repo, storageService: storageService}
}

func (s *documentService) GetDocuments(orgID uint, skip, limit int, expiringInDays *int) ([]models.Document, error) {
	return s.repo.FindByOrganization(orgID, skip, limit, expiringInDays)
}

func (s *documentService) CreateDocument(docIn schemas.DocumentCreate, file *multipart.FileHeader, orgID uint) (*models.Document, error) {
	fileURL, err := s.storageService.Save(file, "documents")
	if err != nil {
		return nil, err
	}

	doc := &models.Document{
		DocumentType:   models.DocumentType(docIn.DocumentType),
		ExpiryDate:     docIn.ExpiryDate,
		Notes:          docIn.Notes,
		VehicleID:      docIn.VehicleID,
		DriverID:       docIn.DriverID,
		FileURL:        fileURL,
		OrganizationID: orgID,
	}

	return s.repo.Create(doc)
}

func (s *documentService) DeleteDocument(docID, orgID uint) (bool, error) {
	doc, err := s.repo.FindByID(docID, orgID)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil // Not found
	}

	err = s.storageService.Delete(doc.FileURL)
	if err != nil {
		// Log the error but don't block the deletion of the DB record
		// logging.Logger.Error("Failed to delete document file", zap.Error(err), zap.String("fileURL", doc.FileURL))
	}

	err = s.repo.Delete(doc)
	if err != nil {
		return false, err
	}

	return true, nil
}
