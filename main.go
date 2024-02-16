package main

import (
	"fmt"
	"ioprodz/home"
	"ioprodz/profile"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Mount module routers

	router := mux.NewRouter()
	router.HandleFunc("/profile", profile.GetHandler)
	router.HandleFunc("/", home.GetHandler)

	// Start the HTTP server
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", router)
}
