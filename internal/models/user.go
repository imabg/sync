package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// User schema
type User struct {
	CreatedAt          time.Time            `json:"createdAt"`
	UpdatedAt          time.Time            `json:"updatedAt"`
	Email              *string              `json:"email" validate:"required,email"`
	UserId             string               `json:"userId"`
	Name               *string              `json:"name" validate:"required"`
	OnboardingTimeline []OnboardingTimeline `json:"onboardingTimeline"`
}

type OnboardingTimeline struct {
	Operation   string    `json:"operation"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completedAt"`
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

func (userCtx *UserCtx) CreateTimelineEvent(op string, done bool) OnboardingTimeline {
	if !done {
		return OnboardingTimeline{
			Operation: op,
			Completed: done,
		}
	}
	return OnboardingTimeline{
		Operation:   op,
		Completed:   done,
		CompletedAt: time.Now(),
	}
}

// InsertOne insert single user record
func (userCtx *UserCtx) InsertOne(ctx context.Context, data *User) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	_, err := userCtx.col.InsertOne(ctx, &data)
	// Need to create a transaction for adding onboarding-event
	return err
}

func (userCtx *UserCtx) FindOne(ctx context.Context, condition bson.M) error {
	err := userCtx.col.FindOne(ctx, &condition)
	// we get err if no document found
	if err.Err() == nil {
		return errors.New("User already exists")
	}
	return nil
}

func (userCtx *UserCtx) FindOneAndUpdate(ctx context.Context, findCondition bson.M, update *User) error {
	update.UpdatedAt = time.Now()
	res := userCtx.col.FindOneAndUpdate(ctx, findCondition, update)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
