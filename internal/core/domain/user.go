package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Email       string    `bson:"email" json:"email"`
	PhoneNumber string    `bson:"phone_number" json:"phone_number"`
	ClerkID     string    `bson:"clerk_id,omitempty" json:"clerk_id,omitzero"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
