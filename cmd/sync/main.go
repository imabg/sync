package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/imabg/sync/internal/config"
	"github.com/imabg/sync/internal/controller"
	"github.com/imabg/sync/internal/database"
	"github.com/imabg/sync/pkg/validate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)
var serverStopCtx context.CancelFunc
var serverCtx context.Context

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
	app.MongoClient = dbCtx.GetMongoDatabase(client)
  stopChan := make(chan os.Signal, 1)
    signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-stopChan
        app.InfoLog.Info("Stopping server...")
        serverStopCtx() // Call the serverStopCtx function to stop the server
    }()

	createAndStartServer(app.Env.ServerAddr, getRoutes(app), *app)
}

// CreateAndStartServer creates a new server and starting listing
func createAndStartServer(addr string, handlers http.Handler, app config.Application) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      handlers,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.InfoLog.Infof("Server started at %s", addr)
	serverCtx, serverStopCtx = context.WithCancel(context.Background())
	go func() {
		<-serverCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			app.ErrorLog.DPanicf("error shutting down server: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.ErrorLog.DPanicf("error starting server: %v", err)
		serverStopCtx() // Stop the server
		return err
	}
	return nil
}

func getRoutes(app *config.Application) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	v1Routes := r.PathPrefix("/api/v1").Subrouter()
	userCtrl := controller.NewUser(app)
	userRoutes := v1Routes.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("/create",userCtrl.CreateUser).Methods("POST")
	return r
}
