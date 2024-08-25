package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Base have common fields of all documents
type Base struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}