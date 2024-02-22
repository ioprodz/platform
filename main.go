package main

import (
	"fmt"
	"ioprodz/auth"
	"ioprodz/common/config"
	"ioprodz/common/middlewares"
	"ioprodz/home"
	"ioprodz/profile"
	"ioprodz/qna"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	configuration := config.Load()
	//openaiClient.Prompt("you are going to ask me 5 questions about 'CI/CD' to assess my knowledge", "{questions:string[]}")

	router := mux.NewRouter()
	// configure module routers
	auth.ConfigureModule(router)
	home.ConfigureModule(router)
	profile.ConfigureModule(router)
	qna.ConfigureModule(router)

	router.Use(middlewares.RequestLogger)
	router.Use(auth.AuthorizeRequest)
	http.Handle("/", router)

	// Start the HTTP server
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":"+configuration.PORT, nil)
}
