package user

import (
	"context"
	"time"

	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
)

type UserServiceCtx struct {
	userModel models.IUserEntity
	config    config.Application
}

func UserServiceInit(app *config.Application) *UserServiceCtx {
	return &UserServiceCtx{
		userModel: models.NewUserModel(*app.MongoClient),
		config:    *app,
	}
}

func (u *UserServiceCtx) CreateNewUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return u.userModel.InsertOne(ctx, user)
}

func (u *UserServiceCtx) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var usr models.User
	err := u.userModel.FindOne(ctx, bson.M{"email": email}, &usr)
	if err != nil {
		return usr, err
	}
	return usr, nil
}
