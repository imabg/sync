package controller

import (
	"context"
	"net/http"

	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/response"
	"github.com/imabg/sync/pkg/types"
	"github.com/imabg/sync/pkg/validate"
	"github.com/imabg/sync/services/entity"
)


type IEntity interface {
	SingUp(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

type EntityCtx struct {
	entityCtx context.Context
	service entity.EntityServiceCtx
	config config.Application
}

func NewEntity(app *config.Application) IEntity {
	ctx := context.Background()
	return &EntityCtx{
		entityCtx: ctx,
		service: *entity.EntityServiceInit(app),
		config: *app,
	}
}



func (e *EntityCtx) SingUp(w http.ResponseWriter, r *http.Request) {
	var entity models.Entity
	err := validate.GetPayload(r, &entity)
	if err != nil {
		response.SendWithError(w, http.StatusBadRequest, *errors.BadRequestError(err.Error()))
		return
	} 
	d, err := e.service.FindByEmail(e.entityCtx, entity.Email)
	if err == nil && d.Email != "" {
		response.SendWithError(w, http.StatusConflict, *errors.ConflictError("user already exists"))
		return
	}
	err = e.service.SignUp(e.entityCtx, &entity)
	if err != nil {
		response.SendWithError(w, http.StatusBadRequest, *errors.BadRequestError(err.Error()))
		return
	}
	res := types.EntityResp{UserId: entity.UserId, Email: entity.Email, UpdatedAt: entity.UpdatedAt}
	response.Send(w, http.StatusCreated, res)
}

func(e *EntityCtx) Login(w http.ResponseWriter, r *http.Request) {
	var login types.LoginDTO
	err := validate.GetPayload(r, &login)
	if err != nil {
		response.SendWithError(w, http.StatusBadRequest, *errors.BadRequestError(err.Error()))
		return
	}
	eDetail, err := e.service.Login(e.entityCtx, login)
	if err != nil {
		response.SendWithError(w, http.StatusBadRequest, *errors.BadRequestError(err.Error()))
		return
	}
	response.Send(w, http.StatusOK, eDetail)

}