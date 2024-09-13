package models

import (
	"context"
	"errors"
	"time"

	"github.com/imabg/sync/pkg/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Entity struct {
	UserId string `json:"userId"`
	Password string `json:"password" validate:"required"`
	OptForPasswordLess bool `json:"optForPasswordLess"`
	Email string `json:"email" validate:"required,email"`
	Source string `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type EntityCtx struct {
	col *mongo.Collection
}

var colName = "entities"
func NewEntityModel (client mongo.Database) *EntityCtx {
	return &EntityCtx{
		col: client.Collection(colName),
	}
}

func (e *EntityCtx) Insert(ctx context.Context, data *Entity) error {
	data.UserId = uuid.GenerateShortId(8)
	// Create password
	err := e.encryptPwd(data)
	if err != nil {
		return  err
	}
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	_, err = e.col.InsertOne(ctx, &data)
	if err != nil {
		return err
	}
	return nil
}

func (e *EntityCtx) encryptPwd(data *Entity) error {
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
	err := e.col.FindOne(ctx, &data).Decode(entity)
	if err == nil {
		return errors.New("User already exists")
	}
	return nil
}