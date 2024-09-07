package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// User schema
type User struct {
	Base
	Email     *string            `json:"email" validate:"required"`
	UserId    *string             `json:"userId"`
	Name      *string            `json:"name" validate:"required"`
}

const COLLECTION_NAME = "users"

func NewUserModel(client mongo.Client, dbName string) *UserCtx {
	return &UserCtx{
		col: client.Database(dbName).Collection(COLLECTION_NAME),
	}
}

type UserCtx struct {
	col *mongo.Collection
}

// InsertOne insert single record
func(userCtx *UserCtx) InsertOne(ctx context.Context, data User) error {
	_, err := userCtx.col.InsertOne(ctx, data)
	return err
}