package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type IEntityEntity interface {
	Insert(context.Context, *Entity) error
	EncryptPwd(*Entity) error
	IsPwdCorrect(string, string) bool
	FindOne(context.Context, bson.M, *Entity) error
}

type Entity struct {
	UserId             string    `json:"userId" bson:"userId"`
	Password           string    `json:"password" bson:"password" validate:"required"`
	OptForPasswordLess bool      `json:"optForPasswordLess" bson:"optForPasswordLess"`
	Email              string    `json:"email" bson:"email" validate:"required,email"`
	IsEmailVerified    bool      `json:"isEmailVerified" bson:"isEmailVerified"`
	Source             string    `json:"source" bson:"source"`
	CreatedAt          time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" bson:"updatedAt"`
}

type EntityCtx struct {
	col *mongo.Collection
}

var colName = "entities"

func NewEntityModel(client mongo.Database) *EntityCtx {
	return &EntityCtx{
		col: client.Collection(colName),
	}
}

func (e *EntityCtx) Insert(ctx context.Context, data *Entity) error {
	_, err := e.col.InsertOne(ctx, &data)
	return err
}

func (e *EntityCtx) EncryptPwd(data *Entity) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return err
	}
	data.Password = string(bytes)
	return nil
}

func (e *EntityCtx) IsPwdCorrect(hashPwd string, currentPlainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(currentPlainPwd))
	return err == nil
}

func (e *EntityCtx) FindOne(ctx context.Context, data bson.M, entity *Entity) error {
	return e.col.FindOne(ctx, &data).Decode(entity)
}
