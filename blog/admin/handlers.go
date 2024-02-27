package blog_admin

import (
	"encoding/json"
	"fmt"
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/clients/openaiClient"
	"ioprodz/common/ui"
	"net/http"

	"github.com/gorilla/mux"
)

type AIBlogPost struct {
	Content          string
	RelatedBlogPosts []string
}

func CreateListPageHandler(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		list, _ := repo.List()
		ui.RenderAdminPage(w, r, "blog/admin/list", list)
	}
}

func CreateEditPageHandler(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blogId := vars["id"]
		blog, err := repo.Get(blogId)
		if err != nil {
			ui.Render404(w, r)
			return
		}
		ui.RenderAdminPage(w, r, "blog/admin/edit", blog)
	}
}

func CreateCreatePageHandler(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		existingPostId := r.URL.Query().Get("postId")

		if existingPostId != "" {
			existingBlog, err := repo.Get(existingPostId)
			if err != nil {
				http.Error(w, "this post does not exist", http.StatusBadRequest)
			}
			ui.RenderAdminPage(w, r, "blog/admin/create", existingBlog)
		} else {
			ui.RenderAdminPage(w, r, "blog/admin/create", nil)
		}
	}
}

func CreateCreateBlogHandler(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// ensure a blog post in db
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		title := r.Form.Get("title")
		existingPostId := r.Form.Get("postId")
		var blog blog_models.Blog
		if existingPostId != "" {
			existingBlog, err := repo.Get(existingPostId)
			if err != nil {
				http.Error(w, "this post does not exist", http.StatusBadRequest)
			}
			existingBlog.Title = title
			blog = existingBlog

		} else {
			blog = *blog_models.NewBlog(title, "", []blog_models.RelatedPosts{})
		}

		// directives
		paragraphCount := r.Form.Get("paragraphCount")

		guidelines := paragraphCount + " short paragraphs"
		if r.Form.Get("useEmojis") == "active" {
			guidelines += ", use emojis"
		}
		if r.Form.Get("useMarkdown") == "active" {
			guidelines += ", use markdown"
		}
		if r.Form.Get("useMermaid") == "active" {
			guidelines += ", use mermaid"
		}

		prompt := "you are going to write the content (" + guidelines + ") of a blog post that has the title '" + title + "' as well as 3 post titles that are related to subject"
		outputFormat := "{content: string,relatedBlogPosts: string[]}"

		aiResponse, err := openaiClient.JsonPrompt(prompt, outputFormat)
		if err != nil {
			fmt.Println("Error getting prompt from ai response", err)
		}
		var aiBlog *AIBlogPost
		if err := json.Unmarshal([]byte(aiResponse), &aiBlog); err != nil {
			fmt.Println("Error parsing json from ai response", err)
		}

		related := make([]blog_models.RelatedPosts, len(aiBlog.RelatedBlogPosts))
		for index, relatedTitle := range aiBlog.RelatedBlogPosts {
			relatedBlogPost := blog_models.NewBlog(relatedTitle, "", []blog_models.RelatedPosts{})
			repo.Create(*relatedBlogPost)
			related[index] = blog_models.RelatedPosts{
				Id: relatedBlogPost.Id, Title: relatedBlogPost.Title}
		}

		blog.SetContent(aiBlog.Content, related)

		blgStr, _ := json.Marshal(blog)
		fmt.Println("BLOOG", string(blgStr))

		if existingPostId != "" {
			repo.Update(blog)
		} else {
			repo.Create(blog)
		}
		w.Write([]byte(blog.Id))
	}
}
