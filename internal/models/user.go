package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserEntity interface {
	InsertOne(context.Context, *User) error
	FindOne(context.Context, bson.M, *User) error
	FindOneAndUpdate(context.Context, bson.M, *User) error
}

// User schema
type User struct {
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
	Email      *string   `json:"email" validate:"required,email" bson:"email"`
	UserId     string    `json:"userId" validate:"required" bson:"userId"`
	Name       *string   `json:"name" validate:"required" bson:"name"`
	IsArchived bool      `json:"is_archived"`
}

const col_name = "users"

func NewUserModel(client mongo.Database) *UserCtx {
	return &UserCtx{
		col: client.Collection(col_name),
	}
}

type UserCtx struct {
	col *mongo.Collection
}

// InsertOne insert single user record
func (userCtx *UserCtx) InsertOne(ctx context.Context, data *User) error {
	_, err := userCtx.col.InsertOne(ctx, &data)
	//TODO:  Need to create a transaction for adding onboarding-event
	return err
}

func (userCtx *UserCtx) FindOne(ctx context.Context, condition bson.M, data *User) error {
	return userCtx.col.FindOne(ctx, &condition).Decode(data)
}

func (userCtx *UserCtx) FindOneAndUpdate(ctx context.Context, findCondition bson.M, update *User) error {
	update.UpdatedAt = time.Now()
	res := userCtx.col.FindOneAndUpdate(ctx, findCondition, update)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
