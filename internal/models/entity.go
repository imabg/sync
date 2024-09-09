package models

import (
	"context"
	"time"

	"github.com/imabg/sync/pkg/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Entity struct {
	UserId string `json:"userId"`
	Password string `json:"-" validate:"required"`
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
func NewEntity (client mongo.Database) *EntityCtx {
	return &EntityCtx{
		col: client.Collection(colName),
	}
}

func (e *EntityCtx) Insert(ctx context.Context, data *Entity) (string, error) {
	data.UserId = uuid.GenerateShortId(8)
	// Create password
	err := e.encryptPwd(data)
	if err != nil {
		return "", err
	}
	res, err := e.col.InsertOne(ctx, &data)
	if err != nil {
		return "", err
	}
	r := res.InsertedID.(string)
	return r,  nil
}

func (e *EntityCtx) encryptPwd(data *Entity) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return err
	}
	data.Password = string(bytes)
	return nil
}