package members_studio

import (
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	members_models "ioprodz/members/_models"
	"net/http"
)

type PageData struct {
	Name string
}

func CreateGetHandler(repo members_models.MembersRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		member, memberNotFound := repo.Get(user.Id)
		if memberNotFound != nil {
			member = members_models.Member{
				Id:        user.Id,
				Name:      user.Name,
				Email:     user.Email,
				AvatarUrl: user.AvatarUrl,
			}
			repo.Create(member)
		}
		ui.RenderPage(w, r, "members/studio/profile", member)
	}
}

func CreateSaveProfileHandler(repo members_models.MembersRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		member, memberNotFound := repo.Get(user.Id)
		if memberNotFound != nil {
			w.Write([]byte("nok"))
		}

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
