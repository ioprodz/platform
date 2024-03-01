package members_models

import "encoding/json"

type Account struct {
	Provider string
	Id       string
}

type Link struct {
	Name string
	Url  string
}

type Member struct {
	Id        string
	Email     string
	Name      string
	Bio       string
	AvatarUrl string
	Accounts  []Account
	Links     []Link
}

func (m Member) GetId() string {
	return m.Id
}

func MemberFromJSON(jsonData []byte) Member {
	var member Member
	if err := json.Unmarshal(jsonData, &member); err != nil {
		panic("unable to parse member json")
	}
	return member
}
