package auth

import (
	auth_infra "ioprodz/auth/_infra"
	auth_authentication "ioprodz/auth/authentication"
	auth_security "ioprodz/auth/security"
	"ioprodz/common/config"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"

	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func NewOAuthCookieStore() *sessions.CookieStore {
	var conf = config.Load()
	store := sessions.NewCookieStore([]byte(conf.AUTH_OAUTH_COOKIE_SECRET))
	store.Options.Path = "/"
	store.Options.Domain = conf.APP_DOMAIN
	store.Options.HttpOnly = true
	store.Options.Secure = conf.IS_PRODUCTION
	store.Options.MaxAge = 86400 * 7

	return store
}

func ConfigureModule(router *mux.Router) {

	// Configure Goth/Gothic lib
	baseUrl := config.Load().BASE_URL

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), baseUrl+"/auth/google/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), baseUrl+"/auth/github/callback"),
	)
	gothic.Store = NewOAuthCookieStore()

	// Configure routes

	accountRepo := auth_infra.CreateMemoryAccountRepo()
	sessionRepo := auth_infra.CreateMemorySessionRepo()

	router.HandleFunc("/auth/{provider}/callback", auth_authentication.CreateOAuthCallbackHandler(accountRepo, sessionRepo)).Methods("GET")
	router.HandleFunc("/auth/{provider}", auth_authentication.CreateOAuthLoginHandler()).Methods("GET")
	router.HandleFunc("/logout", auth_authentication.CreateLogoutHandler()).Methods("GET")

	router.HandleFunc("/security", auth_security.CreateSecurityPageHandler(sessionRepo)).Methods("GET")

}
