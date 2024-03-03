package members

import (
	members_infra "ioprodz/members/_infra"
	members_explore "ioprodz/members/explore"
	members_studio "ioprodz/members/studio"

	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) *mux.Router {

	repo := members_infra.CreateMembersRepository()

	router.HandleFunc("/profile", members_studio.CreateSaveProfileHandler(repo)).Methods("POST")
	router.HandleFunc("/profile", members_studio.CreateGetHandler(repo)).Methods("GET")

	router.HandleFunc("/explore", members_explore.CreateGetHandler())
	return router
}
