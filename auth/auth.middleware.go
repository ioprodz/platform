package auth

import (
	"net/http"
)

type Paths []string

var public Paths = Paths{
	"/",
	"/auth/github",
	"/auth/google",
}

var authCallback Paths = Paths{
	"/auth/github/callback",
	"/auth/google/callback",
}

func (paths *Paths) matchPath(path string) bool {
	found := false
	for _, s := range *paths {
		if s == path {
			found = true
			break
		}
	}
	return found
}

func AuthorizeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if authCallback.matchPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		isPublic := public.matchPath(r.URL.Path)
		_, err := GetUserSession(w, r)
		autnenticated := err == nil
		if !autnenticated {

			if isPublic {
				next.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			if isPublic {
				http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}
