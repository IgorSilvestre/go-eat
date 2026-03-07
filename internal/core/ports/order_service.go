package ports

import (
	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
)

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
