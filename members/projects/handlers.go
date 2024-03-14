package members_projects

import (
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	members_models "ioprodz/members/_models"
	"net/http"

	"github.com/gorilla/mux"
)

type PageData struct {
	Name string
}

func CreateCreateHandler(repo members_models.UserProjectRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		title := r.Form.Get("title")
		description := r.Form.Get("description")
		link := r.Form.Get("link")

		project := members_models.NewUserProject(user.Id, title, description, link)

		projectCreationError := repo.Create(project)
		if projectCreationError != nil {
			w.Write([]byte("nok"))
		}

		w.Write([]byte("ok"))
	}
}

func CreateUpdateHandler(repo members_models.UserProjectRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		project, _ := repo.Get(id)
		if project.UserId != user.Id {
			w.Write([]byte("not your project"))
			return
		}

		project.SetTitle(r.Form.Get("title"))
		project.SetDescription(r.Form.Get("description"))

		repo.Create(project)
		w.Write([]byte("ok"))
	}
}

func CreateGetHandler(repo members_models.UserProjectRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		projects, err := repo.GetByUserId(user.Id)
		if err != nil {
			w.Write([]byte("error fetching list"))
		}
		ui.RenderPage(w, r, "members/projects/list", projects)
	}
}

func CreateDeleteHandler(repo members_models.UserProjectRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		project, _ := repo.Get(id)

		if project.UserId != user.Id {
			w.Write([]byte("not your project"))
			return
		}
		err := repo.Delete(project.Id)
		if err != nil {
			w.Write([]byte("error deleting project"))
			return
		}
		w.Write([]byte("ok"))
	}
}
