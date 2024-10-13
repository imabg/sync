package models

import (
	"context"
	"time"

	"github.com/imabg/sync/pkg/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Notebook struct {
	ID         string    `json:"id"`
	UserId     string    `json:"userId" bson:"userId"`
	Name       string    `json:"name" bson:"name" validate:"required"`
	IsArchived bool      `json:"isArchived" bson:"isArchived"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
}

type NotebookCtx struct {
	col *mongo.Collection
}

const notebookCollName = "notebooks"

func NewNotebookModel(client mongo.Database) *NotebookCtx {
	return &NotebookCtx{
		col: client.Collection(notebookCollName),
	}
}

func (n *NotebookCtx) Create(ctx context.Context, data *Notebook) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.ID = uuid.GenerateUUID()
	_, err := n.col.InsertOne(ctx, &data)
	return err
}
