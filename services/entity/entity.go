package entity

import (
	"context"
	"errors"
	"time"

	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/token"
	"github.com/imabg/sync/pkg/types"
	"github.com/imabg/sync/pkg/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type EntityServiceCtx struct {
	entityModel *models.EntityCtx
	sessionModel *models.SessionCtx
	config config.Application
	log config.Logger
}

func EntityServiceInit(app *config.Application) *EntityServiceCtx {
	return &EntityServiceCtx{
		entityModel: models.NewEntityModel(*app.MongoClient),
		sessionModel: models.NewSessionModel(*app.MongoClient),
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

func (e *EntityServiceCtx) Login(ctx context.Context, loginData types.LoginDTO) (types.LoginResp,error) {
	details, err := e.FindByEmail(ctx, loginData.Email)
	if err != nil {
		return types.LoginResp{}, err
	}

	isPwdMatch := e.entityModel.IsPwdCorrect(details.Password, loginData.Password)
	if !isPwdMatch {
		return types.LoginResp{}, errors.New("password not correct")
	}
	tokenCtx := token.New(e.config.Env.JwtSecretKey)
	t, err := tokenCtx.Generate(token.CustomClaimData{UserId: details.UserId, Email: details.Email}, 4 * time.Hour)
	if err != nil {
		return types.LoginResp{}, err 
	}

	err = e.sessionModel.Create(ctx, &models.Session{UserId: details.UserId, AccessToken: t.Token, ExpiredAt: t.ExpireAt, LastIP: loginData.IPAddr,IsExpired: false, LastUserAgent: loginData.UserAgent})

	if err != nil {
		return types.LoginResp{}, err
	}

	res := types.LoginResp{AccessToken: t.Token, ExpiresAt: t.ExpireAt}

	return res, nil
}