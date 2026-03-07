package ports

import (
	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
)

// IngredientRepository is the interface for ingredient data operations.
type IngredientRepository interface {
	Create(ingredient *domain.Ingredient) error
	GetByID(id uuid.UUID) (*domain.Ingredient, error)
	GetAll() ([]domain.Ingredient, error)
	Update(ingredient *domain.Ingredient) error
	Delete(id uuid.UUID) error
	GetManyByIDs(ids []uuid.UUID) ([]domain.Ingredient, error)
}
