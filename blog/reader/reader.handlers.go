package blog_reader

import (
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/ui"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateListBlogs(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		list, _ := repo.List()
		ui.RenderPage(w, r, "blog/reader/list", list)
	}
}

func CreateViewBlog(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		blogId := vars["id"]
		blog, err := repo.Get(blogId)
		if err != nil {
			ui.Render404(w, r)
			return
		}

		ui.RenderPage(w, r, "blog/reader/view", blog)
	}
}
