package blog

import (
	blog_infra "ioprodz/blog/_infra"
	blog_admin "ioprodz/blog/admin"
	blog_reader "ioprodz/blog/reader"

	"github.com/gorilla/mux"
)

func ConfigureModule(router *mux.Router) {

	blogRepo := blog_infra.CreateMemoryBlogRepo()

	// admin
	router.HandleFunc("/admin/blog", blog_admin.CreateListPageHandler(blogRepo)).Methods("GET")
	router.HandleFunc("/admin/blog/create", blog_admin.CreateCreatePageHandler(blogRepo)).Methods("GET")
	router.HandleFunc("/admin/blog/{id}", blog_admin.CreateEditPageHandler(blogRepo)).Methods("GET")

	// reader
	router.HandleFunc("/blog", blog_reader.CreateListBlogs(blogRepo)).Methods("GET")
	router.HandleFunc("/blog/{id}", blog_reader.CreateViewBlog(blogRepo)).Methods("GET")

	// api
	router.HandleFunc("/api/admin/blog", blog_admin.CreateCreateBlogHandler(blogRepo)).Methods("POST")
}
