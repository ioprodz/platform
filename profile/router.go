package profile

import (
	"github.com/gorilla/mux"
)

func ConfigureRouter(router *mux.Router) *mux.Router {

	router.HandleFunc("/profile", GetHandler).Methods("GET")
	return router
}
