package blog_models

import (
	"encoding/json"
	"time"

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
	CreatedAt    string
	PublishedAt  string
	Reviewed     bool
	RelatedPosts []RelatedPosts
}

func (b Blog) GetId() string {
	return b.Id
}

func (b *Blog) SetContent(body string, related []RelatedPosts) {
	b.Reviewed = false
	b.Body = body
	b.RelatedPosts = related
}

func (b *Blog) SetAsReviewed() {
	b.Reviewed = true
}

func (b *Blog) IsReviewed() bool {
	return b.Reviewed
}

func (b *Blog) SetAsPublished() {
	b.PublishedAt = time.Now().Format(time.RFC3339)
}

func (b *Blog) IsPublished() bool {
	return b.PublishedAt != ""
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
		CreatedAt:    time.Now().Format("2006-01-02T15:04:05Z"),
		PublishedAt:  "",
	}
}
