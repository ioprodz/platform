package members

import (
	member_infra "ioprodz/members/_infra"
	members_explore "ioprodz/members/explore"
	members_studio "ioprodz/members/studio"

	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) *mux.Router {

	repo := member_infra.CreateMemoryMemberRepo()

	router.HandleFunc("/profile", members_studio.CreateSaveProfileHandler(repo)).Methods("POST")
	router.HandleFunc("/profile", members_studio.CreateGetHandler(repo)).Methods("GET")

	router.HandleFunc("/explore", members_explore.CreateGetHandler())
	return router
}
