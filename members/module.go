package members

import (
	members_infra "ioprodz/members/_infra"
	members_feed "ioprodz/members/feed"
	members_projects "ioprodz/members/projects"
	members_studio "ioprodz/members/studio"

	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) *mux.Router {

	membersRepo := members_infra.CreateMembersRepository()
	userProjectsRepo := members_infra.CreateUserProjectRepository()

	router.HandleFunc("/profile", members_studio.CreateSaveProfileHandler(membersRepo)).Methods("POST")
	router.HandleFunc("/profile", members_studio.CreateGetHandler(membersRepo)).Methods("GET")

	router.HandleFunc("/feed", members_feed.CreateGetHandler(membersRepo)).Methods("GET")

	router.HandleFunc("/projects", members_projects.CreateGetHandler(userProjectsRepo)).Methods("GET")
	router.HandleFunc("/projects", members_projects.CreateCreateHandler(userProjectsRepo)).Methods("POST")
	router.HandleFunc("/projects/{id}", members_projects.CreateUpdateHandler(userProjectsRepo)).Methods("PUT")
	router.HandleFunc("/projects/{id}", members_projects.CreateDeleteHandler(userProjectsRepo)).Methods("DELETE")
	return router
}
