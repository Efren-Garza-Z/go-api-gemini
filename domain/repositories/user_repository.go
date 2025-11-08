package repositories

import (
	"errors"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.UserDB) error
	FindAll() ([]models.UserDB, error)
	FindByID(id uint) (*models.UserDB, error)
	Update(user *models.UserDB) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.UserDB) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll() ([]models.UserDB, error) {
	var users []models.UserDB
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(id uint) (*models.UserDB, error) {
	var user models.UserDB
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.UserDB) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserDB{}, id).Error
}
