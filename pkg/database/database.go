package database

import (
	"context"

	"github.com/imabg/sync/pkg/config"
)

type DatabaseCtx struct {
	context context.Context
	config  config.Application
}

func NewDB(ctx context.Context, config config.Application) *DatabaseCtx {
	return &DatabaseCtx{
		context: ctx,
		config:  config,
	}
}
