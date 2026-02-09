package main

import (
	"fmt"
	"log"
	"ioprodz/auth"
	"ioprodz/blog"
	"ioprodz/common/config"
	"ioprodz/common/middlewares"
	"ioprodz/consulting"
	"ioprodz/cv"
	"ioprodz/home"
	"ioprodz/members"
	"ioprodz/qna"
	"ioprodz/solutions"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	configuration := config.Load()
	router := mux.NewRouter()
	//db.NewMongoConnection()

	// Create static subrouter with strict matching
	staticRouter := router.PathPrefix("/static").Subrouter()
	staticRouter.PathPrefix("/favicon/").Handler(http.StripPrefix("/static/favicon/", http.FileServer(http.Dir("common/ui/favicon/"))))
	staticRouter.PathPrefix("/cv-osmane-kalache/").Handler(http.StripPrefix("/static/cv-osmane-kalache/", http.FileServer(http.Dir("cv_osm/"))))

	// Hook global middlewares
	router.Use(middlewares.RequestLogger)

	// Configure module routers (specific routes first)
	auth.ConfigureModule(router)
	members.ConfigureModule(router)
	qna.ConfigureModule(router)
	blog.ConfigureModule(router)
	cv.ConfigureModule(router)

	// Configure consulting module
	consulting.ConfigureModule(router)

	// Configure solutions module
	solutions.ConfigureModule(router)

	// Configure home module last (has catch-all "/" route)
	home.ConfigureModule(router)

	// Mount routes to the HTTP server
	http.Handle("/", router)

	// Start the HTTP server
	fmt.Println("Server listening on port " + configuration.PORT)
	log.Fatal(http.ListenAndServe(":"+configuration.PORT, nil))
}
