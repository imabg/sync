package controller

import (
	"context"
	"net/http"

	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/response"
	"github.com/imabg/sync/pkg/token"

	"github.com/imabg/sync/pkg/validate"
	"github.com/imabg/sync/services/user"
)

type IUser interface {
	CreateUser(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
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

func(u *UserCtx) Get(w http.ResponseWriter, r *http.Request) {
	claim:= r.Context().Value("claims").(token.CustomClaimData)

	usr, err := u.service.FindByEmail(r.Context(), claim.Email)
	
	if err != nil {
		response.SendWithError(w, http.StatusNotFound, *errors.NotFound(err.Error()))
		return
	}

	response.Send(w, http.StatusOK, usr)
}