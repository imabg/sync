package controller

import (
	"context"
	"net/http"

	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/response"
	"github.com/imabg/sync/pkg/validate"
	"go.mongodb.org/mongo-driver/bson"
)

type IUser interface {
	CreateUser(http.ResponseWriter, *http.Request)
}

type UserCtx struct {
		userCtx context.Context
		userModel *models.UserCtx
		config config.Application
}

func NewUser(app *config.Application) IUser {
	ctx := context.Background()
	userCol := models.NewUserModel(*app.MongoClient)
	return &UserCtx{
		userModel: userCol,
		config: *app,
		userCtx: ctx,
	} 
}

func(u *UserCtx) CreateUser(w http.ResponseWriter, r *http.Request) {
	// op := "user.CreateUser"
	var user models.User
	err := validate.GetPayload(r, &user)
	if err != nil {
		custErr := errors.NewCustomError(http.StatusBadRequest, err.Error(), "", "")	
		response.SendWithError(w, http.StatusBadRequest, *custErr)
		return
	}
	err = u.userModel.FindOne(u.userCtx, bson.M{"email": user.Email})
	if err != nil {
		custErr := errors.NewCustomError(http.StatusConflict, err.Error(), "", "")
		response.SendWithError(w, http.StatusBadRequest, *custErr)
		return
	}
	err = u.userModel.InsertOne(u.userCtx, &user)
	if err != nil {
		custErr := errors.NewCustomError(http.StatusBadRequest, err.Error(), "", "")	
		response.SendWithError(w, http.StatusInternalServerError, *custErr)
		return
	}
	response.Send(w, http.StatusCreated, user)
}