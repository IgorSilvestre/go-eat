package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderProduct struct {
	ID          uuid.UUID   `bson:"_id" json:"id"`
	OrderID     uuid.UUID   `bson:"order_id" json:"order_id"`
	ProductID   uuid.UUID   `bson:"product_id" json:"product_id"`
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description" json:"description"`
	Ingredients []uuid.UUID `bson:"ingredients" json:"ingredients"` // array of ingredient IDs included
	Adicionais  []uuid.UUID `bson:"adicionais" json:"adicionais"`   // array of extra ingredient IDs
	Removed     []uuid.UUID `bson:"removed" json:"removed"`         // array of removed ingredient IDs
	Price       float64     `bson:"price" json:"price"`             // calculated price (base product + adicionais)
	CreatedAt   time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `bson:"updated_at" json:"updated_at"`
}

type Order struct {
	ID          uuid.UUID      `bson:"_id" json:"id"`
	UserID      uuid.UUID      `bson:"user_id" json:"user_id"`
	PhoneNumber string         `bson:"phone_number" json:"phone_number"` // Stored alongside user_id as per diagram
	Items       []OrderProduct `bson:"-" json:"items"`                   // Stored separately in MongoDB, populated in memory
	TotalPrice  float64        `bson:"total_price" json:"total_price"`
	CreatedAt   time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at" json:"updated_at"`
}
