package auth

import (
	"net/http"
)

var store = NewAppCookieStore()

var sessionName = "ioprodz-session"

type SessionError struct {
	msg string
}

type SessionData struct {
	Id        string
	Email     string
	AvatarUrl string
}

func (e *SessionError) Error() string { return e.msg }

func SetUserSession(w http.ResponseWriter, r *http.Request, sessionData SessionData) error {
	session, err := store.Get(r, sessionName)
	session.Values["id"] = sessionData.Id
	session.Values["email"] = sessionData.Email
	session.Values["avatarUrl"] = sessionData.AvatarUrl
	session.Save(r, w)
	return err
}

func GetUserSession(w http.ResponseWriter, r *http.Request) (SessionData, error) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return SessionData{}, &SessionError{msg: "could not load session"}
	}
	id, ok := session.Values["id"].(string)
	if !ok {
		return SessionData{}, &SessionError{msg: "could not find userId on session"}
	}
	email, ok := session.Values["email"].(string)
	if !ok {
		return SessionData{}, &SessionError{msg: "could not find userId on session"}
	}
	avatarUrl, ok := session.Values["avatarUrl"].(string)
	if !ok {
		return SessionData{}, &SessionError{msg: "could not find userId on session"}
	}
	return SessionData{
		Id:        id,
		Email:     email,
		AvatarUrl: avatarUrl,
	}, nil
}

func ClearSessionHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	for key := range session.Values {
		delete(session.Values, key)
	}
	session.Save(r, w)
}
