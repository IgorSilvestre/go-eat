package ports

import (
	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
)

// UserRepository is the interface for user data operations.
type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	GetAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
}
