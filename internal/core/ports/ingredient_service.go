package ports

import (
	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
)

type IngredientService interface {
	Create(name string, price float64) (*domain.Ingredient, error)
	Get(id uuid.UUID) (*domain.Ingredient, error)
	List() ([]domain.Ingredient, error)
	Update(id uuid.UUID, name string, price float64) (*domain.Ingredient, error)
	Delete(id uuid.UUID) error
}
