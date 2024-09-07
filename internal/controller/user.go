package controller

import (
	"net/http"

	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/validate"
)

type IUser interface {
	CreateUser(http.ResponseWriter, *http.Request)
}

type UserCtx struct {
		userModel *models.UserCtx
		config config.Application
}

func NewUser(app *config.Application) IUser {
	userCol := models.NewUserModel(*app.Client, app.Env.MongoURI)
	return &UserCtx{
		userModel: userCol,
		config: *app,
	} 
}

func(ctx *UserCtx) CreateUser(w http.ResponseWriter, r *http.Request) {
	// op := "user.CreateUser"
	var user models.User
	err := validate.ValidateRequestPayload(r, &user)
	if err != nil {
		return
	}
}