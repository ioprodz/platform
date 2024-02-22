package auth

import (
	"net/http"
)

var store = NewAppCookieStore()

var sessionName = "ioprodz-session"

type SessionError struct {
	msg string
}

func (e *SessionError) Error() string { return e.msg }

func SetUserSession(w http.ResponseWriter, r *http.Request, userId string) error {
	session, err := store.Get(r, sessionName)
	session.Values["userId"] = userId
	session.Save(r, w)
	return err
}

func GetUserSession(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := store.Get(r, sessionName)

	if err != nil {
		return "", &SessionError{msg: "could not load session"}
	}
	userId, ok := session.Values["userId"].(string)

	if !ok {
		return "", &SessionError{msg: "could not find userId on session"}
	}

	return userId, nil
}

func ClearSessionHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	for key := range session.Values {
		delete(session.Values, key)
	}
	session.Save(r, w)
}
