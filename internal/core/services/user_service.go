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

func (s *userService) Create(name, email, phone, clerkID string) (*domain.User, error) {
	// Check for existing email
	existingByEmail, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingByEmail != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Check for existing clerk_id
	existingByClerkID, err := s.repo.GetByClerkID(clerkID)
	if err != nil {
		return nil, err
	}
	if existingByClerkID != nil {
		return nil, domain.ErrClerkIDAlreadyExists
	}

	user := &domain.User{
		Name:        name,
		Email:       email,
		PhoneNumber: phone,
		ClerkID:     clerkID,
	}
	err = s.repo.Create(user)
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

func (s *userService) Update(id uuid.UUID, name, email, phone, clerkID string) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check for existing email if it's being changed
	if email != user.Email {
		existingByEmail, err := s.repo.GetByEmail(email)
		if err != nil {
			return nil, err
		}
		if existingByEmail != nil {
			return nil, domain.ErrEmailAlreadyExists
		}
	}

	// Check for existing clerk_id if it's being changed
	if clerkID != user.ClerkID {
		existingByClerkID, err := s.repo.GetByClerkID(clerkID)
		if err != nil {
			return nil, err
		}
		if existingByClerkID != nil {
			return nil, domain.ErrClerkIDAlreadyExists
		}
	}

	user.Name = name
	user.Email = email
	user.PhoneNumber = phone
	user.ClerkID = clerkID

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
