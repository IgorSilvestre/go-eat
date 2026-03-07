package ports

import (
	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
)

// ProductRepository is the interface for product data operations.
type ProductRepository interface {
	Create(product *domain.Product) error
	GetByID(id uuid.UUID) (*domain.Product, error)
	GetAll() ([]domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uuid.UUID) error
}
