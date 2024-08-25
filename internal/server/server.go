package server

import (
	"net/http"
	"time"
)

// TODO: need to update 
// CreateAndStartServer creates a new server and starting listing
func CreateAndStartServer(addr string, handlers http.Handler) error {
	srv := &http.Server{
		Addr: addr,
		Handler: handlers,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}