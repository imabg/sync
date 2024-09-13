package controller

import (
	"context"
	"net/http"

	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/response"
	"github.com/imabg/sync/pkg/validate"
	"github.com/imabg/sync/services/user"
)

type IUser interface {
	CreateUser(http.ResponseWriter, *http.Request)
}

type UserCtx struct {
		userCtx context.Context
		service user.UserServiceCtx
		config config.Application
		log config.Logger
}

func NewUser(app *config.Application) IUser {
	ctx := context.Background()
	return &UserCtx{
		config: *app,
		service: *user.UserServiceInit(app),
		userCtx: ctx,
		log: app.Log,
	} 
}

func(u *UserCtx) CreateUser(w http.ResponseWriter, r *http.Request) {
	op := "user.CreateUser"
	u.log.InfoLog.Infof("[UserController]: started %s", op)
	defer u.log.InfoLog.Infof("[UserController]: finished %s", op)
	var user models.User
	err := validate.GetPayload(r, &user)
	if err != nil {
		response.SendWithError(w, http.StatusBadRequest, *errors.BadRequestError(err.Error()))
		return
	}
	err = u.service.CreateNewUser(u.userCtx, &user)
	if err != nil {
		response.SendWithError(w, http.StatusInternalServerError, *errors.BadRequestError(err.Error()))
		return
	}
	response.Send(w, http.StatusCreated, user)
}