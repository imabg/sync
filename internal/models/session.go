package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ISessionEntity interface {
	Create(context.Context, *Session) error
	FindOne(context.Context, bson.D, *Session) error
	FindOneAndUpdate(context.Context, bson.M, bson.D) error
}

// TODO: support multiple sessions

type Session struct {
	UserId        string    `json:"user_id" validate:"required" bson:"userId"`
	AccessToken   string    `json:"accessToken" validate:"required" bson:"accessToken"`
	ExpiredAt     time.Time `json:"expiredAt" validate:"required" bson:"expiredAt"`
	LastIP        string    `json:"lastIp" validate:"required" bson:"lastIp"`
	LastUserAgent string    `json:"userAgent" bson:"userAgent"`
	CreatedAt     time.Time `json:"createdAt" validate:"required" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" validate:"required" bson:"updatedAt"`
	IsExpired     bool      `json:"isExpired" bson:"isExpired"`
	IsActive      bool      `json:"isActive" validate:"required" bson:"isActive"`
}

type SessionCtx struct {
	col *mongo.Collection
}

const sessionCollName = "sessions"

func NewSessionModel(client mongo.Database) *SessionCtx {
	return &SessionCtx{
		col: client.Collection(sessionCollName),
	}
}

// TODO: Only create single session against single userId

func (s *SessionCtx) Create(ctx context.Context, data *Session) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	_, err := s.col.InsertOne(ctx, &data)
	return err
}

func (s *SessionCtx) FindOne(ctx context.Context, findCondition bson.D, details *Session) error {
	return s.col.FindOne(ctx, &findCondition).Decode(details)
}

func (s *SessionCtx) FindOneAndUpdate(ctx context.Context, findCondition bson.M, update bson.D) error {
	r := s.col.FindOneAndUpdate(ctx, findCondition, update)
	return r.Err()
}
