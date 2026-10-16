package policies

type CurrentUser struct {
	Id        string
	SessionId string
	Name      string
	Email     string
	AvatarUrl string
	Provider  string
	// Zitadel project ids the user holds a role in (ioprodz accounts only)
	ProjectIds []string
}

type KeyType = string

const CurrentUserCtxKey KeyType = "sessionUser"

func (u *CurrentUser) IsAuthenticated() bool {
	return u.Id != ""
}

func (u *CurrentUser) HasIoprodzAccount() bool {
	return u.Provider == "ioprodz"
}
