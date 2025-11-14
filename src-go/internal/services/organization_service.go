package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go-api/internal/config"
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
	"io"
)

type OrganizationService interface {
	GetOrganizations(skip, limit int, status *string) ([]models.Organization, error)
	UpdateOrganization(orgID uint, orgIn schemas.OrganizationUpdate) (*models.Organization, error)
	GetFuelIntegrationSettings(orgID uint) (*models.Organization, error)
	UpdateFuelIntegrationSettings(orgID uint, settingsIn schemas.OrganizationFuelIntegrationUpdate) (*models.Organization, error)
}

type organizationService struct {
	repo repositories.OrganizationRepository
}

func NewOrganizationService(repo repositories.OrganizationRepository) OrganizationService {
	return &organizationService{repo: repo}
}

func (s *organizationService) GetOrganizations(skip, limit int, status *string) ([]models.Organization, error) {
	return s.repo.FindAll(skip, limit, status)
}

func (s *organizationService) UpdateOrganization(orgID uint, orgIn schemas.OrganizationUpdate) (*models.Organization, error) {
	org, err := s.repo.FindByID(orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, nil // Not found
	}

	if orgIn.Name != "" {
		org.Name = orgIn.Name
	}
	if orgIn.Sector != "" {
		org.Sector = models.Sector(orgIn.Sector)
	}

	return s.repo.Update(org)
}

func (s *organizationService) GetFuelIntegrationSettings(orgID uint) (*models.Organization, error) {
	return s.repo.FindByID(orgID)
}

func (s *organizationService) UpdateFuelIntegrationSettings(orgID uint, settingsIn schemas.OrganizationFuelIntegrationUpdate) (*models.Organization, error) {
	org, err := s.repo.FindByID(orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, nil // Not found
	}

	org.FuelProviderName = &settingsIn.FuelProviderName

	if settingsIn.APIKey != nil {
		encryptedKey, err := encrypt(*settingsIn.APIKey, config.AppConfig.SECRET_KEY)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt API key: %w", err)
		}
		org.EncryptedFuelProviderAPIKey = &encryptedKey
	}

	if settingsIn.APISecret != nil {
		encryptedSecret, err := encrypt(*settingsIn.APISecret, config.AppConfig.SECRET_KEY)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt API secret: %w", err)
		}
		org.EncryptedFuelProviderAPISecret = &encryptedSecret
	}

	return s.repo.Update(org)
}

// encrypt encrypts text with a key
func encrypt(text, keyString string) (string, error) {
	key := []byte(keyString)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// decrypt decrypts text with a key
func decrypt(encryptedText, keyString string) (string, error) {
	key := []byte(keyString)
	ciphertext, err := hex.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
