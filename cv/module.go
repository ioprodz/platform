package cv

import (
	cv_infra "ioprodz/cv/_infra"
	cv_viewer "ioprodz/cv/viewer"

	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) {

	cvRepo := cv_infra.CreateMongoCVRepo()

	router.HandleFunc("/cv/{id}", cv_viewer.CreateGetCvHandler(cvRepo)).Methods("GET")
}
