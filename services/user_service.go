package services

import (
	"errors"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/domain/repositories"
)

type UserService interface {
	CreateUser(input models.CreateUserInput) (*models.UserDB, error)
	GetAllUsers() ([]models.UserDB, error)
	GetUserByID(id uint) (*models.UserDB, error)
	UpdateUser(id uint, input models.CreateUserInput) (*models.UserDB, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(input models.CreateUserInput) (*models.UserDB, error) {
	user := &models.UserDB{
		FullName: input.FullName,
		Email:    input.Email,
		Password: input.Password, // ideal: hash aquí
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]models.UserDB, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*models.UserDB, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("usuario no encontrado")
	}
	return u, nil
}

func (s *userService) UpdateUser(id uint, input models.CreateUserInput) (*models.UserDB, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("usuario no encontrado")
	}
	u.FullName = input.FullName
	u.Email = input.Email
	// opcional: actualizar password si se envía
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
