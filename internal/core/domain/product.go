package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID   `bson:"_id" json:"id"`
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description" json:"description"`
	Ingredients []uuid.UUID `bson:"ingredients" json:"ingredients"` // array of ingredient IDs
	Price       float64     `bson:"price" json:"price"`
	Image       string      `bson:"image" json:"image"`
	CreatedAt   time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `bson:"updated_at" json:"updated_at"`
}
