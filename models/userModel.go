package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	FirstName string             `json:"first_name,omitempty" validate:"required,min=2,max=30" bson:"first_name"`
	LastName  string             `json:"last_name,omitempty" validate:"required,min=2,max=30" bson:"last_name"`
	Password  string             `json:"password,omitempty" validate:"required,min=6,max=30" bson:"password"`
	Email     string             `json:"email,omitempty" validate:"email,required" bson:"email"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}
