package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: support multiple sessions

type Session struct {
	UserId string `json:"user_id" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
	ExpiredAt time.Time `json:"expired_at" validate:"required"`
	LastIP string `json:"last_ip" validate:"required"`
	LastUserAgent string `json:"user_agent"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	IsExpired bool `json:"is_expired"`
	IsActive bool `json:"is_active" validate:"required"`
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

func(s *SessionCtx) Create(ctx context.Context, data *Session) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	_, err := s.col.InsertOne(ctx, &data)
	return err
}

func(s *SessionCtx) FindOne(ctx context.Context, findCondition bson.D, details *Session) error {
	return s.col.FindOne(ctx, &findCondition).Decode(details)
}

func(s *SessionCtx) FindOneAndUpdate(ctx context.Context, findCondition bson.M, update bson.D) error {
	r := s.col.FindOneAndUpdate(ctx, findCondition, update)
	return r.Err()
}