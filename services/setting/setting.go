package setting

import (
	"context"
	"time"

	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/uuid"
)

type SettingServiceCtx struct {
	settingModel models.ISettingEntity
	config config.Application
	log config.Logger
}

func SettingServiceInit(app *config.Application) *SettingServiceCtx {
	return &SettingServiceCtx{
		settingModel: models.NewSettingModel(*app.MongoClient),
		config: *app,
		log: app.Log,
	}
}

func (s *SettingServiceCtx) DefaultSetting(ctx context.Context, userId string) error {
	var data models.Setting
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.SettingId = uuid.GenerateUUID()
	data.IsOptForDailyDigest = true
	data.UserId = userId
	data.MaximumNotePerEmail = 5
	// TODO: add OptedTime
	return s.settingModel.Create(ctx, &data)
}

func (s *SettingServiceCtx) FindById(ctx context.Context, id string, setting *models.Setting) {

}