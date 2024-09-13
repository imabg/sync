package entity

import (
	"context"
	"errors"
	"time"

	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/types"
	"github.com/imabg/sync/pkg/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type EntityServiceCtx struct {
	entityModel *models.EntityCtx
	config config.Application
	log config.Logger
}

func EntityServiceInit(app *config.Application) *EntityServiceCtx {
	return &EntityServiceCtx{
		entityModel: models.NewEntityModel(*app.MongoClient),
		config: *app,
		log: app.Log,
	}
}

func (e *EntityServiceCtx) SignUp(ctx context.Context, data *models.Entity) error {
	op := "entityService.SignUp"
	e.log.InfoLog.Infof("[EntityService]: started %s function", op)
	defer e.log.InfoLog.Infof("[EntityService]: %s function", op)
	data.UserId = uuid.GenerateShortId(8)
	err := e.entityModel.EncryptPwd(data)
	if err != nil {
		return  err
	}
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	return e.entityModel.Insert(ctx, data)
}

func (e *EntityServiceCtx) FindByEmail (ctx context.Context, email string) (models.Entity, error) {
	var entityData models.Entity
	op := "entityService.FindByEmail"
	e.log.InfoLog.Infof("[EntityService]: started %s function", op)
	defer e.log.InfoLog.Infof("[EntityService]: %s function", op)
	err := e.entityModel.FindOne(ctx, bson.M{"email": email}, &entityData)
	if err != nil {
		return models.Entity{}, err
	}
	return entityData, nil
}

func (e *EntityServiceCtx) Login(ctx context.Context, loginData types.LoginDTO) (types.EntityResp,error) {
	details, err := e.FindByEmail(ctx, loginData.Email)
	if err != nil {
		return types.EntityResp{}, err
	}

	isPwdMatch := e.entityModel.IsPwdCorrect(details.Password, loginData.Password)
	if !isPwdMatch {
		return types.EntityResp{}, errors.New("password not correct")
	}
	res := types.EntityResp{UserId: details.UserId,Email: details.Email, UpdatedAt: details.UpdatedAt}
	return res, nil
}