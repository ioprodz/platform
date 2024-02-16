package main

import (
	"fmt"
	"ioprodz/common/middlewares"
	"ioprodz/home"
	"ioprodz/profile"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	// configure module routers
	home.ConfigureRouter(router)
	profile.ConfigureRouter(router)

	router.Use(middlewares.RequestLogger)
	http.Handle("/", router)

	// Start the HTTP server
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
