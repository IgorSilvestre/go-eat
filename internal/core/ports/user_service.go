package ports

import (
	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
)

type UserService interface {
	Create(name, email, phone string) (*domain.User, error)
	Get(id uuid.UUID) (*domain.User, error)
	List() ([]domain.User, error)
	Update(id uuid.UUID, name, email, phone string) (*domain.User, error)
	Delete(id uuid.UUID) error
}
