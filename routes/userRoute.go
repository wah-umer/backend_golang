package user_routes

import (
	"github.com/gorilla/mux"
	usercontrollers "github.com/umerwaheed/backend_golang/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", usercontrollers.GetById).Methods("GET")
	return router
}
