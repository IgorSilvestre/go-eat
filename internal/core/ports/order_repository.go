package ports

import (
	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
)

// OrderRepository is the interface for order data operations.
type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id uuid.UUID) (*domain.Order, error)
	GetAll() ([]domain.Order, error)
	CreateOrderProduct(op *domain.OrderProduct) error
	GetOrderProductsByOrderID(orderID uuid.UUID) ([]domain.OrderProduct, error)
}
