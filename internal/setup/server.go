package setup

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/logger"
)

// CreateAndStartServer creates a new server and starting listing
func CreateAndStartServer(addr string, app config.Application) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      getRoutes(&app),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Log.InfoLog.Infof("Server started at %s", addr)
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	go func() {
		<-serverCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Log.ErrorLog.DPanicf("error shutting down server: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.ErrorLog.DPanicf("error starting server: %v", err)
		serverStopCtx() // Stop the server
		return err
	}
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		logger.Log.InfoLog.Info("Stopping server...")
		serverStopCtx() // Call the serverStopCtx function to stop the server
	}()
	return nil
}
