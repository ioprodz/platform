package blog

import (
	blog_infra "ioprodz/blog/_infra"
	blog_models "ioprodz/blog/_models"
	blog_admin "ioprodz/blog/admin"
	blog_reader "ioprodz/blog/reader"

	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) blog_models.BlogRepository {

	blogRepo := blog_infra.CreateBlogRepository()

	// admin
	router.HandleFunc("/admin/blog", blog_admin.CreateListPageHandler(blogRepo)).Methods("GET")
	router.HandleFunc("/admin/blog/create", blog_admin.CreateCreatePageHandler(blogRepo)).Methods("GET")
	router.HandleFunc("/admin/blog/{id}", blog_admin.CreateEditPageHandler(blogRepo)).Methods("GET")

	// reader
	router.HandleFunc("/blog", blog_reader.CreateListBlogs(blogRepo)).Methods("GET")
	router.HandleFunc("/blog/{id}", blog_reader.CreateViewBlog(blogRepo)).Methods("GET")

	// api
	router.HandleFunc("/api/admin/blog", blog_admin.CreateCreateBlogHandler(blogRepo)).Methods("POST")
	router.HandleFunc("/api/admin/blog/{id}/review", blog_admin.CreateReviewHandler(blogRepo)).Methods("PUT")
	router.HandleFunc("/api/admin/blog/{id}/publish", blog_admin.CreatePublishHandler(blogRepo)).Methods("PUT")

	return blogRepo
}
