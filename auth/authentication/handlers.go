package auth_authentication

import (
	"fmt"
	auth_infra "ioprodz/auth/_infra"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func CreateOAuthCallbackHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Println("Unauthorized: " + err.Error())
			return
		}

		auth_infra.SetUserSession(w, r, auth_infra.SessionData{Id: user.UserID, Email: user.Email, AvatarUrl: user.AvatarURL})
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}

func CreateOAuthLoginHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// try to get the user without re-authenticating
		if user, err := gothic.CompleteUserAuth(w, r); err == nil {
			auth_infra.SetUserSession(w, r, auth_infra.SessionData{Id: user.UserID, Email: user.Email, AvatarUrl: user.AvatarURL})
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			gothic.BeginAuthHandler(w, r)
		}

	}
}

func CreateLogoutHandler() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		auth_infra.ClearSessionHandler(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
