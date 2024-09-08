package models

import (
	"context"
	"errors"
	"time"

	"github.com/imabg/sync/pkg/uuid"
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

//TODO: need to fix
func (userCtx *UserCtx) FindOne(ctx context.Context, filter interface{}) (error) {
	var user User
	_ = userCtx.col.FindOne(ctx, filter).Decode(&user)
	if user.UserId != "" {
		return errors.New("duplicate user exists")
	}
	return nil
}