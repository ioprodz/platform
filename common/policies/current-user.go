package policies

type CurrentUser struct {
	Id        string
	SessionId string
	Name      string
	Email     string
	AvatarUrl string
}

type KeyType = string

const CurrentUserCtxKey KeyType = "sessionUser"

func (u *CurrentUser) IsAuthenticated() bool {
	return u.Id != ""
}
