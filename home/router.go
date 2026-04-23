package home

import (
	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) *mux.Router {
	router.HandleFunc("/", GetHandler).Methods("GET")
	router.HandleFunc("/about", GetAboutHandler).Methods("GET")
	router.HandleFunc("/login", GetLoginHandler).Methods("GET")
	return router
}
