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
		member, _ := repo.Get("member-id")
		ui.RenderPage(w, r, "members/studio/profile", member)
	}
}

func CreateSaveProfileHandler(repo members_models.MembersRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		member, _ := repo.Get("member-id")

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		name := r.Form.Get("name")
		bio := r.Form.Get("bio")

		names := r.Form["links_names[]"]
		urls := r.Form["links_urls[]"]

		links := []members_models.Link{}
		for index, name := range names {
			links = append(links, members_models.Link{Name: name, Url: urls[index]})

		}

		member.Name = name
		member.Bio = bio
		member.Links = links

		repo.Update(member)
		w.Write([]byte("ok"))
	}
}
