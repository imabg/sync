package main

import (
	"context"
	"time"

	"github.com/imabg/sync/internal/database"
	"github.com/imabg/sync/internal/setup"
	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/validate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)


func main() {
	env := config.NewEnv()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	validate.SetupValidation()
	log := &config.Logger{ErrorLog: sugar.Named("Error"), InfoLog: sugar.Named("Info")}
	app := &config.Application{Env: *env, Log: *log}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbCtx := database.NewDB(ctx, *app)
	client, err := dbCtx.CreateMongoConnection()
	if err != nil {
		app.Log.ErrorLog.DPanicf("While creating DB connection %v", err)
	}
	defer func(ctx database.DatabaseCtx, client *mongo.Client) {
		if err = ctx.DiscountMongoConnection(client); err != nil {
			app.Log.ErrorLog.DPanicf("While closing DB connection %v", err)
		}
	}(*dbCtx, client)
	app.MongoClient = dbCtx.GetMongoDatabase(client)

	err = setup.CreateAndStartServer(env.ServerAddr, *app)
	if err != nil {
		panic(err)
	}
}
