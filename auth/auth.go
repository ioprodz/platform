package auth

import (
	auth_infra "ioprodz/auth/_infra"
	auth_authentication "ioprodz/auth/authentication"
	auth_authorization "ioprodz/auth/authorization"
	auth_security "ioprodz/auth/security"
	"ioprodz/common/config"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"

	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/openidConnect"
)

func NewOAuthCookieStore() *sessions.CookieStore {
	var conf = config.Load()
	store := sessions.NewCookieStore([]byte(conf.AUTH_OAUTH_COOKIE_SECRET))
	store.Options.Path = "/"
	// Intentionally no Domain: host-only cookie, bulletproof round-trip on the
	// exact host that served the /auth/{provider} response. Avoids cookie loss
	// when APP_DOMAIN is misconfigured or doesn't exactly match the request host.
	store.Options.HttpOnly = true
	store.Options.Secure = conf.IS_PRODUCTION
	// OAuth state round-trip takes seconds. 15 min buffer, not 7 days.
	store.Options.MaxAge = 60 * 15

	return store
}

func ConfigureModule(router *mux.Router) {

	// Configure Goth/Gothic lib
	conf := config.Load()
	baseUrl := conf.BASE_URL

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), baseUrl+"/auth/google/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), baseUrl+"/auth/github/callback"),
	)
	configureIoprodzAuth(conf.IOPRODZ_AUTH_ISSUER, conf.IOPRODZ_AUTH_CLIENT_ID, conf.IOPRODZ_AUTH_CLIENT_SECRET, baseUrl)
	gothic.Store = NewOAuthCookieStore()

	// Configure routes

	accountRepo := auth_infra.CreateAccountRepository()
	sessionRepo := auth_infra.CreateSessionRepository()

	router.HandleFunc("/auth/{provider}/callback", auth_authentication.CreateOAuthCallbackHandler(accountRepo, sessionRepo)).Methods("GET")
	router.HandleFunc("/auth/{provider}", auth_authentication.CreateOAuthLoginHandler(accountRepo, sessionRepo)).Methods("GET")
	router.HandleFunc("/logout", auth_authentication.CreateLogoutHandler()).Methods("GET")

	router.HandleFunc("/security", auth_security.CreateSecurityPageHandler(sessionRepo)).Methods("GET")
	router.HandleFunc("/security/sessions/{id}", auth_security.CreateRevokeSessionHandler(sessionRepo)).Methods("DELETE")

	router.Use(auth_authorization.CreateRequestAuthorization(sessionRepo))

}

// configureIoprodzAuth registers the ioprodz account (Zitadel OIDC) as goth provider "ioprodz".
// Skipped when no client id is configured.
func configureIoprodzAuth(issuer string, clientId string, clientSecret string, baseUrl string) {
	if clientId == "" {
		return
	}
	provider, err := openidConnect.NewNamed("ioprodz", clientId, clientSecret,
		baseUrl+"/auth/ioprodz/callback",
		issuer+"/.well-known/openid-configuration",
		"openid", "profile", "email",
		// lets the user's token read their own grants (portal)
		"urn:zitadel:iam:org:project:id:zitadel:aud",
	)
	if err != nil {
		log.Printf("ioprodz login disabled: %v", err)
		return
	}
	goth.UseProviders(provider)
}
