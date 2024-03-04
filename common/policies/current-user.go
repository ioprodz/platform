package policies

type CurrentUser struct {
	Id        string
	Name      string
	Email     string
	AvatarUrl string
}

type KeyType = int

const CurrentUserCtxKey KeyType = iota

func (u *CurrentUser) IsAuthenticated() bool {
	return u.Id != ""
}
