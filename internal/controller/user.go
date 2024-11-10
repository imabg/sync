package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/logger"
	"github.com/imabg/sync/pkg/response"
	"github.com/imabg/sync/pkg/token"

	"github.com/imabg/sync/pkg/validate"
	"github.com/imabg/sync/services/email"
	"github.com/imabg/sync/services/setting"
	"github.com/imabg/sync/services/user"
)

type IUser interface {
	CreateUser(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}

type UserCtx struct {
	userCtx        context.Context
	service        user.UserServiceCtx
	settingService setting.SettingServiceCtx
	config         config.Application
}

func NewUser(app *config.Application) IUser {
	ctx := context.Background()
	return &UserCtx{
		config:         *app,
		service:        *user.UserServiceInit(app),
		settingService: *setting.SettingServiceInit(app),
		userCtx:        ctx,
	}
}

func (u *UserCtx) CreateUser(w http.ResponseWriter, r *http.Request) {
	op := "user.CreateUser"
	logger.Log.InfoLog.Infof("%s is called", op)
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
	// Create default setting
	go u.settingService.DefaultSetting(u.userCtx, user.UserId)

	portNum, _ := strconv.Atoi(u.config.Env.MailerPort)
	// Send welcome mail
	e := email.SetupMailer(u.config.Env.MailerHost, portNum, u.config.Env.MailerUsername, u.config.Env.MailerPassword)
	msg, err := e.CreateMessage(u.config.Env.MailerSender, *user.Email, []string{}, email.WELCOME_EMAIL, email.TemplateConfig{
		TemplateDir:  "templates",
		TemplateFile: "welcome.html",
	}, email.WelcomeEmail{FullName: *user.Name})
	if err != nil {
		response.SendWithError(w, http.StatusInternalServerError, *errors.BadRequestError(err.Error()))
		return
	}
	if err = e.Send(*msg); err != nil {
		response.SendWithError(w, http.StatusInternalServerError, *errors.BadRequestError(err.Error()))
		return
	}
	response.Send(w, http.StatusCreated, user)
}

func (u *UserCtx) Get(w http.ResponseWriter, r *http.Request) {
	claim := r.Context().Value("claims").(token.CustomClaimData)

	usr, err := u.service.FindByEmail(r.Context(), claim.Email)
	if err != nil {
		response.SendWithError(w, http.StatusNotFound, *errors.NotFound(err.Error()))
		return
	}

	response.Send(w, http.StatusOK, usr)
}
