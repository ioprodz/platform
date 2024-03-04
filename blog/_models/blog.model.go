package blog_models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/xeonx/timeago"
)

type RelatedPost struct {
	Id    string
	Title string
}

type Editor struct {
	Id        string
	Name      string
	AvatarUrl string
}

type Blog struct {
	Id string
	// blog
	Title    string
	Abstract string
	Body     string

	// search & discovery
	ReadingTime  int8
	Keywords     []string
	RelatedPosts []RelatedPost

	// timestamps
	CreatedAt   string
	PublishedAt string

	// review
	Reviewed bool
	Editors  []Editor
}

func (b Blog) GetId() string {
	return b.Id
}

func (b *Blog) SetContent(body string, related []RelatedPost) {
	b.Reviewed = false
	b.Body = body
	b.RelatedPosts = related
}

func (b *Blog) SetAsReviewed() {
	b.Reviewed = true
}

func (b *Blog) AddEditor(editor Editor) {
	if !b.hasEditor(editor.Id) {

		b.Editors = append(b.Editors, editor)
	}
}

func (b *Blog) hasEditor(editorId string) bool {
	for _, editor := range b.Editors {
		if editor.Id == editorId {
			return true
		}
	}
	return false
}

func (b Blog) IsReviewed() bool {
	return b.Reviewed
}

func (b *Blog) SetAsPublished() {
	b.PublishedAt = time.Now().Format(time.RFC3339)
}

func (b Blog) PublishedAtHumanReadable() string {
	date, err := time.Parse(time.RFC3339, b.PublishedAt)
	if err != nil {
		return "no publish date"
	}
	return timeago.English.Format(date)
}

func (b Blog) IsPublished() bool {
	return b.PublishedAt != ""
}

func BlogFromJSON(jsonData []byte) Blog {
	var blog Blog
	if err := json.Unmarshal(jsonData, &blog); err != nil {
		panic("unable to parse QNA json")

	}
	return blog
}

func NewBlog(title string, body string, related []RelatedPost) *Blog {
	return &Blog{
		Id:           uuid.NewString(),
		Title:        title,
		Body:         body,
		RelatedPosts: related,
		CreatedAt:    time.Now().Format(time.RFC3339),
		PublishedAt:  "",

		Abstract:    "",
		ReadingTime: 3,
		Keywords:    []string{},
	}
}
