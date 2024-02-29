package members_studio

import (
	"ioprodz/common/ui"
	members_models "ioprodz/members/_models"
	"net/http"
)

type PageData struct {
	Name string
}

func CreateGetHandler(repo members_models.MembersRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// Define data to be passed to the template
		member, _ := repo.Get("member-id")

		// Parse the template file
		ui.RenderPage(w, r, "members/studio/profile", member)
	}
}

func CreateSaveProfileHandler(repo members_models.MembersRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// Define data to be passed to the template
		member, _ := repo.Get("member-id")

		// Parse the template file
		ui.RenderPage(w, r, "members/studio/profile", member)
	}
}
