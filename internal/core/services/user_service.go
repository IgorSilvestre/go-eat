package services

import (
	"restaurant-api/internal/core/domain"
	"restaurant-api/internal/core/ports"

	"github.com/google/uuid"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(name, email, phone string) (*domain.User, error) {
	user := &domain.User{
		Name:        name,
		Email:       email,
		PhoneNumber: phone,
	}
	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Get(id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) List() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *userService) Update(id uuid.UUID, name, email, phone string) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.Email = email
	user.PhoneNumber = phone

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
