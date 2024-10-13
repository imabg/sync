package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ISettingEntity interface {
	Create(context.Context, *Setting) error
	FindOneAndUpdate(context.Context, bson.M, *Setting) error
	FindByUserId(context.Context, *Setting) error
}

type Setting struct {
	UserId               string    `json:"userId" validate:"required" bson:"userId"`
	SettingId            string    `json:"settingId" validate:"required" bson:"settingId"`
	IsOptForDailyDigest  bool      `json:"isOptForDailyDigest" validate:"required" bson:"isOptForDailyDigest"`
	Timezone             string    `json:"timezone" validate:"required" bson:"timezone"`
	OptedTime            time.Time `json:"optedTime" bson:"optedTime"`
	MaximumNotePerEmail  int64     `json:"maximumNotePerEmail" bson:"maximumNotePerEmail"`
	OptForWeekendSpecial bool      `json:"optForWeekendSpecial" bson:"optForWeekendSpecial"`
	CreatedAt            time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt" bson:"updatedAt"`
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

func (s *SettingCtx) FindByUserId(ctx context.Context, data *Setting) error {
	return s.col.FindOne(ctx, bson.M{"userId": &data.UserId}).Decode(&data)
}