package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/imabg/sync/internal/setup"
	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/database"
	"github.com/imabg/sync/pkg/logger"
	"github.com/imabg/sync/pkg/validate"
)

func main() {
	env := config.NewEnv()
	err := logger.Setup()
	validate.SetupValidation()
	app := &config.Application{Env: *env}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbCtx := database.NewDB(ctx, *app)
	client, err := dbCtx.CreatePostgresConnection()
	if err != nil {
		logger.Log.ErrorLog.DPanicf("While creating DB connection %v", err)
	}
	defer func(ctx database.DatabaseCtx, client *sql.DB) {
		if err = ctx.DisconnectPostgresConnection(client); err != nil {
			logger.Log.ErrorLog.DPanicf("While closing DB connection %v", err)
		}
	}(*dbCtx, client)
	// app.MongoClient = dbCtx.GetMongoDatabase(client)
	app.PostgresClient = client

	err = setup.CreateAndStartServer(env.ServerAddr, *app)
	if err != nil {
		panic(err)
	}
}
