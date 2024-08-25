package main

import (
	"context"
	"time"

	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)



func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	app := &config.Application{ErrorLog: sugar.Named("Error"), InfoLog: sugar.Named("Info")}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbCtx := database.NewDB(ctx, *app, "")
	client, err := dbCtx.CreateMongoConnection()
	if err != nil {
		app.ErrorLog.DPanicf("While creating DB connection %v", err)
	}
	defer func (ctx database.DatabaseCtx, client *mongo.Client)  {
		if err = ctx.DiscountMongoConnection(client); err != nil	{
			app.ErrorLog.DPanicf("While closing DB connection %v", err)
		}
	}(*dbCtx, client)
	

}
