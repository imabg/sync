package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/response"
	"github.com/imabg/sync/pkg/validate"
	"go.mongodb.org/mongo-driver/bson"
)


type IEntity interface {
	SingUp(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

type EntityCtx struct {
	entityCtx context.Context
	entityModel *models.EntityCtx
	config config.Application
}

func NewEntity(app *config.Application) IEntity {
	ctx := context.Background()
	entityCol := models.NewEntityModel(*app.MongoClient)
	return &EntityCtx{
		entityCtx: ctx,
		entityModel: entityCol,
		config: *app,
	}
}

type EntityResp struct {
	UserId string `json:"userId"`
	Email string `json:"email"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (e *EntityCtx) SingUp(w http.ResponseWriter, r *http.Request) {
	var entity models.Entity
	err := validate.GetPayload(r, &entity)
	if err != nil {
		custErr := errors.NewCustomError(http.StatusBadRequest, err.Error(), "", "")
		response.SendWithError(w, http.StatusBadRequest, *custErr)
		return
	} 
	err = e.entityModel.FindOne(e.entityCtx, bson.M{"email": entity.Email}, &entity)
	if err != nil {
		custErr := errors.NewCustomError(http.StatusConflict, err.Error(), "", "")
		response.SendWithError(w, http.StatusConflict, *custErr)
		return
	}
	err = e.entityModel.Insert(e.entityCtx, &entity)
	if err != nil {
		custErr := errors.NewCustomError(http.StatusBadRequest, err.Error(), "", "")
		response.SendWithError(w, http.StatusBadRequest, *custErr)
		return
	}
	res := EntityResp{UserId: entity.UserId, Email: entity.Email, UpdatedAt: entity.UpdatedAt}
	response.Send(w, http.StatusCreated, res)
}

func(e *EntityCtx) Login(w http.ResponseWriter, r *http.Request) {
	var login LoginDTO	
	var entity models.Entity
	err := validate.GetPayload(r, &login)
	if err != nil {
		custErr := errors.NewCustomError(http.StatusBadRequest, err.Error(), "", "")
		response.SendWithError(w, http.StatusBadRequest, *custErr)
		return
	}
	err = e.entityModel.FindOne(e.entityCtx, bson.M{"email": login.Email}, &entity )
	if err == nil {
		custErr := errors.NewCustomError(http.StatusNotFound, "user not found", "", "")
		response.SendWithError(w, http.StatusNotFound, *custErr)
		return
	}

	isPwdMatch := e.entityModel.IsPwdCorrect(entity.Password, login.Password)
	if !isPwdMatch {
		custErr := errors.NewCustomError(http.StatusBadRequest, "password not correct", "", "")
		response.SendWithError(w, http.StatusBadRequest, *custErr)
		return
	}

	response.Send(w, http.StatusOK, "login")

}