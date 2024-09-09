package models

import (
	"context"
	"errors"
	"time"

	"github.com/imabg/sync/pkg/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// User schema
type User struct {
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Email     *string            `json:"email" validate:"required,email"`
	UserId    string             `json:"userId"`
	Name      *string            `json:"name" validate:"required"`
}

const COLLECTION_NAME = "users"

func NewUserModel(client mongo.Database) *UserCtx {
	return &UserCtx{
		col: client.Collection(COLLECTION_NAME),
	}
}

type UserCtx struct {
	col *mongo.Collection
}

// InsertOne insert single user record
func(userCtx *UserCtx) InsertOne(ctx context.Context, data *User) error {
	data.UserId = uuid.GenerateShortId(0)
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	_, err := userCtx.col.InsertOne(ctx, &data)
	return err
}

func (userCtx *UserCtx) FindOne(ctx context.Context, condition bson.M) (error) {
	err := userCtx.col.FindOne(ctx, &condition)
	// we get err if no document found 
	if err.Err() == nil {
	return errors.New("User already exists")
	}
	return nil
}