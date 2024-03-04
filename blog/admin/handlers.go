package blog_admin

import (
	"encoding/json"
	"fmt"
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/clients/openaiClient"
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	"net/http"

	"github.com/gorilla/mux"
)

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

func CreateReviewHandler(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blogId := vars["id"]
		blog, err := repo.Get(blogId)
		if err != nil {
			ui.Render404(w, r)
			return
		}

		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		blog.AddEditor(blog_models.Editor{
			Id:        user.Id,
			Name:      user.Name,
			AvatarUrl: user.AvatarUrl,
		})
		blog.SetAsReviewed()
		repo.Update(blog)
		w.Write([]byte("ok"))
	}
}

func CreatePublishHandler(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blogId := vars["id"]
		blog, err := repo.Get(blogId)
		if err != nil {
			ui.Render404(w, r)
			return
		}

		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		blog.AddEditor(blog_models.Editor{
			Id:        user.Id,
			Name:      user.Name,
			AvatarUrl: user.AvatarUrl,
		})

		blog.SetAsPublished()

		repo.Update(blog)
		w.Write([]byte("ok"))
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

type AIBlogPost struct {
	Content          string
	Abstract         string
	RelatedBlogPosts []string
	Keywords         []string
	ReadingTime      int8
}

func CreateCreateBlogHandler(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			blog = *blog_models.NewBlog(title, "", []blog_models.RelatedPost{})
		}

		// directives
		paragraphCount := r.Form.Get("paragraphCount")
		guidelines := "# Guidelines\n"
		guidelines += "- be practical and cite concrete examples\n"
		guidelines += "- don't repeat the title in the content\n"
		guidelines += "- " + paragraphCount + " short paragraphs\n"
		if r.Form.Get("useEmojis") == "active" {
			guidelines += "- use emojis\n"
		}
		if r.Form.Get("useMarkdown") == "active" {
			guidelines += "- use markdown\n"
		}
		if r.Form.Get("useMermaid") == "active" {
			guidelines += "- use mermaid\n"
		}

		guidelines += ", use mermaid"
		command := "Command: \n"
		command += "- write the content of a blog post that has the title '" + title + "'\n"
		command += "- an abstract that is apealing to open the post \n"
		command += "- 3 related post titles \n"
		command += "- 5 keywords\n"
		command += "- approximative time to read it (minutes)\n"
		prompt := command + "(" + guidelines + ")"
		outputFormat := "{content: string, abstract:string, relatedBlogPosts: string[], keywords:string[], readingTime:number}"

		aiResponse, err := openaiClient.JsonPrompt(prompt, outputFormat)
		if err != nil {
			fmt.Println("Error getting prompt from ai response", err)
		}
		var aiBlog *AIBlogPost
		if err := json.Unmarshal([]byte(aiResponse), &aiBlog); err != nil {
			fmt.Println("Error parsing json from ai response", err)
		}

		related := make([]blog_models.RelatedPost, len(aiBlog.RelatedBlogPosts))
		for index, relatedTitle := range aiBlog.RelatedBlogPosts {
			relatedBlogPost := blog_models.NewBlog(relatedTitle, "", []blog_models.RelatedPost{})
			repo.Create(*relatedBlogPost)
			related[index] = blog_models.RelatedPost{
				Id: relatedBlogPost.Id, Title: relatedBlogPost.Title}
		}

		blog.SetContent(aiBlog.Content, related)
		blog.Keywords = aiBlog.Keywords
		blog.ReadingTime = aiBlog.ReadingTime
		blog.Abstract = aiBlog.Abstract

		if existingPostId != "" {
			repo.Update(blog)
		} else {
			repo.Create(blog)
		}
		w.Write([]byte(blog.Id))
	}
}
