package routes

import (
	"github.com/gorilla/mux"
	"github.com/umerwaheed/backend_golang/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	//user routes
	router.HandleFunc("/users/{id}", controllers.GetById).Methods("GET")

	//address routes
	router.HandleFunc("/address", controllers.Create).Methods("POST")
	router.HandleFunc("/address/user/{id}", controllers.GetByUserId).Methods("GET")
	router.HandleFunc("/address/{id}", controllers.UpdateById).Methods("PATCH")
	router.HandleFunc("/address/{id}", controllers.DeleteById).Methods("DELETE")

	return router
}
