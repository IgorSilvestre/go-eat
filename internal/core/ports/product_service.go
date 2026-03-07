package ports

import (
	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
)

type ProductService interface {
	Create(name, description string, ingredients []uuid.UUID, price float64, image string) (*domain.Product, error)
	Get(id uuid.UUID) (*domain.Product, error)
	List() ([]domain.Product, error)
	Update(id uuid.UUID, name, description string, ingredients []uuid.UUID, price float64, image string) (*domain.Product, error)
	Delete(id uuid.UUID) error
}
