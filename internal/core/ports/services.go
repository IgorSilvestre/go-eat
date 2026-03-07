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

type IngredientService interface {
	Create(name string, price float64) (*domain.Ingredient, error)
	Get(id uuid.UUID) (*domain.Ingredient, error)
	List() ([]domain.Ingredient, error)
	Update(id uuid.UUID, name string, price float64) (*domain.Ingredient, error)
	Delete(id uuid.UUID) error
}

type ProductService interface {
	Create(name, description string, ingredients []uuid.UUID, price float64) (*domain.Product, error)
	Get(id uuid.UUID) (*domain.Product, error)
	List() ([]domain.Product, error)
	Update(id uuid.UUID, name, description string, ingredients []uuid.UUID, price float64) (*domain.Product, error)
	Delete(id uuid.UUID) error
}

// OrderItemInput represents the frontend payload for adding an item to an order.
type OrderItemInput struct {
	ProductID  uuid.UUID   `json:"product_id"`
	Adicionais []uuid.UUID `json:"adicionais"`
	Removed    []uuid.UUID `json:"removed"`
}

type OrderService interface {
	CreateOrder(userID uuid.UUID, items []OrderItemInput) (*domain.Order, error)
	Get(id uuid.UUID) (*domain.Order, error)
	List() ([]domain.Order, error)
}
