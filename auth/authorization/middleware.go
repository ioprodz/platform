package auth_authorization

import (
	"context"
	auth_infra "ioprodz/auth/_infra"
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/policies"
	"net/http"
	"strings"
)

type Paths []string

var public Paths = Paths{
	"/",
	"/auth/github",
	"/auth/google",
	"/blog",
}

var authCallback Paths = Paths{
	"/auth/github/callback",
	"/auth/google/callback",
}

func (paths *Paths) matchPath(path string) bool {
	for _, s := range *paths {
		if s == "/" {
			if s == path {
				return true
			}
		} else {
			if strings.HasPrefix(path, s) {
				return true
			}
		}
	}
	return false
}

func CreateRequestAuthorization(sessionRepo auth_models.SessionRepository) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if authCallback.matchPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			cookie, _ := auth_infra.GetAuthCookie(w, r)
			session, sessionError := sessionRepo.Get(cookie.Id)
			ctx := context.WithValue(r.Context(), policies.CurrentUserCtxKey, policies.CurrentUser{
				Id:        session.AccountId,
				SessionId: session.Id,
				Name:      session.Name,
				Email:     session.Email,
				AvatarUrl: session.AvatarUrl,
			})
			autnenticated := sessionError == nil
			isPublic := public.matchPath(r.URL.Path)

			if !autnenticated && !isPublic {
				auth_infra.ClearAuthCookie(w, r)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
