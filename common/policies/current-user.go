package policies

type CurrentUser struct {
	Id        string
	Name      string
	AvatarUrl string
}

type KeyType = string

const CurrentUserCtxKey KeyType = "currentUser"

func (u *CurrentUser) IsAuthenticated() bool {
	return u.Id != ""
}
