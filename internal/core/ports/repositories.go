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

// IngredientRepository is the interface for ingredient data operations.
type IngredientRepository interface {
	Create(ingredient *domain.Ingredient) error
	GetByID(id uuid.UUID) (*domain.Ingredient, error)
	GetAll() ([]domain.Ingredient, error)
	Update(ingredient *domain.Ingredient) error
	Delete(id uuid.UUID) error
	GetManyByIDs(ids []uuid.UUID) ([]domain.Ingredient, error)
}

// ProductRepository is the interface for product data operations.
type ProductRepository interface {
	Create(product *domain.Product) error
	GetByID(id uuid.UUID) (*domain.Product, error)
	GetAll() ([]domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uuid.UUID) error
}

// OrderRepository is the interface for order data operations.
type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id uuid.UUID) (*domain.Order, error)
	GetAll() ([]domain.Order, error)
	CreateOrderProduct(op *domain.OrderProduct) error
	GetOrderProductsByOrderID(orderID uuid.UUID) ([]domain.OrderProduct, error)
}
