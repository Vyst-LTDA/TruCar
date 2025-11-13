package services

import (
	"errors"

	"go-api/internal/core"
	"go-api/internal/repositories"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if !core.CheckPasswordHash(password, user.HashedPassword) {
		return "", ErrInvalidCredentials
	}

	token, err := core.GenerateJWT(user.ID, user.OrganizationID)
	if err != nil {
		return "", err
	}

	return token, nil
}
