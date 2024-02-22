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

		ui.RenderPage(w, r, "blog/admin/list", repo.List())
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
		ui.RenderPage(w, r, "blog/admin/edit", blog)
	}
}

func CreateCreatePageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ui.RenderPage(w, r, "blog/admin/create", nil)
	}
}

func CreateCreateBlogHandler(repo blog_models.BlogRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		title := r.Form.Get("title")

		prompt := "you are going to write the content (more or less 3 short paragraphs) of a blog post that has the title '" + title + "' as well as 3 post titles that are related to subject"
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
			related[index] = blog_models.RelatedPosts{
				Id: "", Title: relatedTitle}
		}

		blog := blog_models.NewBlog(title, aiBlog.Content, related)
		blgStr, _ := json.Marshal(blog)
		fmt.Println("BLOOG", string(blgStr))

		repo.Create(*blog)
		w.Write([]byte(blog.Id))
	}
}
