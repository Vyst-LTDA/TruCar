package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"

	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

type UserService interface {
	GetUsers(orgID uint, skip, limit int) ([]models.User, error)
	GetUser(userID, orgID uint) (*models.User, error)
	CreateUser(userIn schemas.UserCreate, orgID uint) (*models.User, error)
	UpdateUser(userID, orgID uint, userIn schemas.UserUpdate) (*models.User, error)
	DeleteUser(userID, orgID uint) error
	GetAllUsers(skip, limit int) ([]models.User, error)
	GetDemoUsers() ([]models.User, error)
	ActivateUser(userID uint) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUsers(orgID uint, skip, limit int) ([]models.User, error) {
	return s.repo.FindByOrganization(orgID, skip, limit)
}

func (s *userService) GetUser(userID, orgID uint) (*models.User, error) {
	return s.repo.FindByID(userID, orgID)
}

func (s *userService) CreateUser(userIn schemas.UserCreate, orgID uint) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userIn.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:          userIn.Email,
		FullName:       userIn.FullName,
		HashedPassword: string(hashedPassword),
		Role:           userIn.Role,
		OrganizationID: orgID,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(userID, orgID uint, userIn schemas.UserUpdate) (*models.User, error) {
	user, err := s.repo.FindByID(userID, orgID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Or return a not found error
	}

	if userIn.FullName != "" {
		user.FullName = userIn.FullName
	}
	if userIn.Email != "" {
		user.Email = userIn.Email
	}
	if userIn.IsActive != nil {
		user.IsActive = *userIn.IsActive
	}

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(userID, orgID uint) error {
	user, err := s.repo.FindByID(userID, orgID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil // Or return a not found error
	}

	return s.repo.Delete(user)
}

func (s *userService) GetAllUsers(skip, limit int) ([]models.User, error) {
	return s.repo.FindAll(skip, limit)
}

func (s *userService) GetDemoUsers() ([]models.User, error) {
	return s.repo.FindByRole(models.RoleClienteDemo)
}

func (s *userService) ActivateUser(userID uint) (*models.User, error) {
	user, err := s.repo.FindByIDUnscoped(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Not found
	}

	if user.Role != models.RoleClienteDemo {
		return nil, errors.New("user is not a demo client")
	}

	user.Role = models.RoleClienteAtivo
	err = s.repo.Update(user)
	return user, err
}
