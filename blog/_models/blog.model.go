package blog_models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type RelatedPosts struct {
	Id    string
	Title string
}

type Blog struct {
	Id           string
	Title        string
	Body         string
	RelatedPosts []RelatedPosts
}

func BlogFromJSON(jsonData []byte) Blog {
	var blog Blog
	if err := json.Unmarshal(jsonData, &blog); err != nil {
		panic("unable to parse QNA json")

	}
	return blog
}

func NewBlog(title string, body string, related []RelatedPosts) *Blog {
	return &Blog{
		Id:           uuid.NewString(),
		Title:        title,
		Body:         body,
		RelatedPosts: related,
	}
}
