package user

import (
	"context"
	"time"

	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
)

type UserServiceCtx struct {
	userModel *models.UserCtx
	config config.Application
	log config.Logger
}

func UserServiceInit(app *config.Application) *UserServiceCtx {
	return &UserServiceCtx{
		userModel: models.NewUserModel(*app.MongoClient),
		config: *app,
		log: app.Log,
	}
}

func (u *UserServiceCtx) CreateNewUser(ctx context.Context, user *models.User) error {
	op := "userService.CreateNewUser"
	u.log.InfoLog.Infof("[UserService]: started %s function", op)
	defer u.log.InfoLog.Infof("[UserService]: %s function", op)
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