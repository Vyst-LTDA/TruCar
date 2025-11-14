package repositories

import (
	"errors"

	"gorm.io/gorm"

	"go-api/internal/models"
)

type UserRepository interface {
	FindByID(userID, orgID uint) (*models.User, error)
	FindByIDUnscoped(userID uint) (*models.User, error) // Para Super Admin
	FindByEmail(email string) (*models.User, error)
	FindByOrganization(orgID uint, skip, limit int) ([]models.User, error)
	FindAll(skip, limit int) ([]models.User, error) // Para Super Admin
	FindByRole(role models.UserRole) ([]models.User, error) // Para Super Admin
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(userID, orgID uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ? AND organization_id = ?", userID, orgID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByOrganization(orgID uint, skip, limit int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("organization_id = ?", orgID).Offset(skip).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}

func (r *userRepository) FindByIDUnscoped(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(skip, limit int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Offset(skip).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
