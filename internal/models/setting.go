package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Setting struct {
	UserId string `json:"userId" validate:"required"`
	IsOptForDailyDigest bool `json:"isOptForDailyDigest" validate:"required"`
	Timezone string `json:"timezone" validate:"required"`
	OptedTime time.Time `json:"optedTime"`
	MaximumNote int64 `json:"maximumNote"`
	OptForWeekendSpecial bool `json:"optForWeekendSpecial"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

const settingColName = "settings"

type SettingCtx struct {
	col *mongo.Collection
}

func NewSettingModel(client mongo.Database) *SettingCtx {
	return &SettingCtx{
		col: client.Collection(settingColName),
	}
}

func (s *SettingCtx) Create(ctx context.Context, data *Setting) error {
	_, err := s.col.InsertOne(ctx, &data)
	return err
}

func (s *SettingCtx) FindOneAndUpdate(ctx context.Context, findCondition bson.M, update *Setting) error {
	update.UpdatedAt = time.Now()
	r := s.col.FindOneAndUpdate(ctx, findCondition, update)
	return r.Err()
}