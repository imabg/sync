package setup

import (
	"github.com/gorilla/mux"
	"github.com/imabg/sync/internal/controller"
	"github.com/imabg/sync/pkg/config"
)

func getRoutes(app *config.Application) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	v1Routes := r.PathPrefix("/api/v1").Subrouter()
	v1Routes.Use(AuthMiddleware)
	publicRoutes := r.PathPrefix("/api/v1").Subrouter()
	userCtrl := controller.NewUser(app)
	entityCtrl := controller.NewEntity(app)
	notebookCtrl := controller.NewNotebook(app)
	// open-routes
	publicRoutes.HandleFunc("/signup", entityCtrl.SingUp).Methods("POST")
	publicRoutes.HandleFunc("/login", entityCtrl.Login).Methods("POST")
	// user
	v1Routes.HandleFunc("/users/create", userCtrl.CreateUser).Methods("POST")
	v1Routes.HandleFunc("/users/get", userCtrl.Get).Methods("GET")
	// notebook
	v1Routes.HandleFunc("/notebook/create", notebookCtrl.CreateNoteBook).Methods("POST")

	return r
}