package blog_reader

import (
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/ui"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func CreateListBlogs(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		list, _ := repo.ListPublished()
		meta := ui.PageMeta{
			Title:       "Blog",
			Description: "Insights and tutorials on software engineering, architecture, TDD, and technology leadership.",
			Path:        "/blog",
			OGType:      "website",
		}
		ui.RenderPageWithMeta(w, r, "blog/reader/list", list, meta)
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

		description := blog.Abstract
		if description == "" {
			description = blog.Title
		}
		keywords := ""
		if len(blog.Keywords) > 0 {
			keywords = strings.Join(blog.Keywords, ", ")
		}

		meta := ui.PageMeta{
			Title:                blog.Title,
			Description:          description,
			Path:                 "/blog/" + blog.Id,
			OGType:               "article",
			Keywords:             keywords,
			ArticlePublishedTime: blog.PublishedAt,
		}
		ui.RenderPageWithMeta(w, r, "blog/reader/view", blog, meta)
	}
}
