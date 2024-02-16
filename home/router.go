package home

import (
	"github.com/gorilla/mux"
)

func ConfigureRouter(router *mux.Router) *mux.Router {
	router.HandleFunc("/", GetHandler).Methods("GET")
	return router
}
