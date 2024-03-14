package members_models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type UserProject struct {
	Id          string
	UserId      string
	Link        string
	Title       string
	Description string
}

func (m UserProject) GetId() string {
	return m.Id
}

func (m *UserProject) SetTitle(title string) {
	m.Title = title
}

func (m *UserProject) SetDescription(description string) {
	m.Description = description
}

func (m *UserProject) SetLink(link string) {
	m.Link = link
}

func UserProjectFromJSON(jsonData []byte) UserProject {
	var up UserProject
	if err := json.Unmarshal(jsonData, &up); err != nil {
		panic("unable to parse User Project json")
	}
	return up
}

func NewUserProject(userId string, title string, description string, link string) UserProject {
	return UserProject{Id: uuid.NewString(), UserId: userId, Title: title, Description: description, Link: link}
}
