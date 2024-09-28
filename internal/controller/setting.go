package controller

import (
	"context"
	// "net/http"

	"github.com/imabg/sync/pkg/config"
)

// type ISetting interface {
// 	UpdateSetting(http.ResponseWriter, *http.Request)
// }

type SettingCtx struct {
	settingCtx context.Context
	config config.Application
	log config.Logger
}

// func NewSetting(app *config.Application) ISetting {
	
// }