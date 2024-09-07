package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/controller"
	"github.com/imabg/sync/internal/database"
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
	app := &config.Application{ErrorLog: sugar.Named("Error"), InfoLog: sugar.Named("Info"), Env: *env}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbCtx := database.NewDB(ctx, *app)
	client, err := dbCtx.CreateMongoConnection()
	if err != nil {
		app.ErrorLog.DPanicf("While creating DB connection %v", err)
	}
	defer func(ctx database.DatabaseCtx, client *mongo.Client) {
		if err = ctx.DiscountMongoConnection(client); err != nil {
			app.ErrorLog.DPanicf("While closing DB connection %v", err)
		}
	}(*dbCtx, client)
	app.Client = client
	createAndStartServer(ctx, app.Env.ServerAddr, getRoutes(app), *app)
}

// CreateAndStartServer creates a new server and starting listing
func createAndStartServer(ctx context.Context, addr string, handlers http.Handler, app config.Application) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      handlers,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.InfoLog.Infof("Server started at %s", addr)
	err := srv.ListenAndServe()
	if err != nil {
		go func(ctx context.Context, app config.Application) {
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			if err := srv.Shutdown(shutdownCtx); err != nil {
				app.ErrorLog.DPanicf("error shutting down server %v", err)
			}
		}(ctx, app)
	}
	return nil
}

func getRoutes(app *config.Application) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	userCtrl := controller.NewUser(app)	
	userRoutes := r.PathPrefix("/users")
	userRoutes.HandlerFunc(userCtrl.CreateUser).Methods("POST")
	return r
}
