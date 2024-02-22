package auth

import (
	"fmt"
	"ioprodz/common/config"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/markbates/goth"

	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func ConfigureModule(router *mux.Router) {
	baseUrl := config.Load().BASE_URL
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), baseUrl+"/auth/google/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), baseUrl+"/auth/github/callback"),
	)

	gothic.Store = NewOAuthCookieStore()

	router.HandleFunc("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Println("Unauthorized: " + err.Error())
			return
		}

		SetUserSession(w, r, user.Email)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}).Methods("GET")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		ClearSessionHandler(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}).Methods("GET")

	router.HandleFunc("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		// try to get the user without re-authenticating
		if user, err := gothic.CompleteUserAuth(w, r); err == nil {
			SetUserSession(w, r, user.Email)
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	}).Methods("GET")

}
