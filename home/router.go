package home

import (
	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) *mux.Router {
	router.HandleFunc("/", GetHandler).Methods("GET")
	return router
}
