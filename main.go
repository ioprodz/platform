package main

import (
	"fmt"
	"ioprodz/auth"
	"ioprodz/blog"
	"ioprodz/common/config"
	"ioprodz/common/middlewares"
	"ioprodz/cv"
	"ioprodz/home"
	"ioprodz/members"
	"ioprodz/qna"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	configuration := config.Load()
	router := mux.NewRouter()
	//db.NewMongoConnection()

	// Configure module routers
	auth.ConfigureModule(router)
	home.ConfigureModule(router)
	members.ConfigureModule(router)
	qna.ConfigureModule(router)
	blog.ConfigureModule(router)
	cv.ConfigureModule(router)

	// Hook global middlewares
	router.Use(middlewares.RequestLogger)
	router.Use(auth.AuthorizeRequest)

	// Mount routes to the HTTP server
	http.Handle("/", router)

	// Start the HTTP server
	fmt.Println("Server listening on port " + configuration.PORT)
	http.ListenAndServe(":"+configuration.PORT, nil)
}
