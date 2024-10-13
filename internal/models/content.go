package models

import "time"

type Content struct {
	ID        string    `json:"id" bson:"id"`
	Text      string    `json:"text" bson:"text" validate:"required"`
	Tags      []string  `json:"tags" bson:"tags"`
	IsLike    bool      `json:"isLike" bson:"isLike"`
	Source    string    `json:"source" bson:"source"`
	User      User      `json:"user" bson:"user"`
	Notebook  Notebook  `json:"notebook" bson:"notebook"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
