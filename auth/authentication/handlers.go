package auth_authentication

import (
	"fmt"
	"hash/fnv"
	auth_infra "ioprodz/auth/_infra"
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/config"
	"log"
	"net/http"
	"slices"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/mileusna/useragent"
)

func CreateOAuthCallbackHandler(accountRepo auth_models.AccountRepository, sessionRepo auth_models.SessionRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			log.Printf("auth callback failed: provider=%s err=%v ua=%q referer=%q",
				r.URL.Query().Get("state"), err, r.UserAgent(), r.Referer())
			http.Redirect(w, r, "/login?error=auth_failed", http.StatusTemporaryRedirect)
			return
		}

		ua := useragent.Parse(r.UserAgent())
		account := getAccount(accountRepo, user)
		session := getSession(account, user, ua, sessionRepo)

		auth_infra.SetAuthCookie(w, r, auth_infra.CookieData{Id: session.Id})
		keepIoprodzToken(w, r, user)

		w.Header().Set("Location", landingPage(w, r, user))
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}

// nextCookie remembers where an ioprodz sign-in should land (/security after
// "sign in again"); only the pages listed in landingPages are honoured.
var nextCookie = "ioprodz-next"

var landingPages = []string{"/security"}

func rememberNext(w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get("next")
	if !slices.Contains(landingPages, next) {
		return
	}
	http.SetCookie(w, &http.Cookie{Name: nextCookie, Value: next, Path: "/auth", MaxAge: 60 * 15, HttpOnly: true, SameSite: http.SameSiteLaxMode})
}

func landingPage(w http.ResponseWriter, r *http.Request, user goth.User) string {
	if user.Provider != "ioprodz" {
		return "/feed"
	}
	if next, err := r.Cookie(nextCookie); err == nil {
		http.SetCookie(w, &http.Cookie{Name: nextCookie, Path: "/auth", MaxAge: -1})
		if slices.Contains(landingPages, next.Value) {
			return next.Value
		}
	}
	return "/portal"
}

// keepIoprodzToken keeps the user's access token for the account's sign-in
// methods on /security.
func keepIoprodzToken(w http.ResponseWriter, r *http.Request, user goth.User) {
	if user.Provider != "ioprodz" {
		return
	}
	loginName, _ := user.RawData["preferred_username"].(string)
	if loginName == "" {
		loginName = user.Email
	}
	auth_infra.SetIoprodzToken(w, r, auth_infra.IoprodzToken{
		AccessToken: user.AccessToken,
		UserId:      user.UserID,
		LoginName:   loginName,
		ExpiresAt:   user.ExpiresAt,
	})
}

// projectIds lists the Zitadel projects an ioprodz account holds a role in.
// Other providers have none. On error the portal shows no services rather
// than failing the login.
func projectIds(user goth.User) []string {
	if user.Provider != "ioprodz" {
		return nil
	}
	ids, err := auth_infra.FetchIoprodzProjectIds(config.Load().IOPRODZ_AUTH_ISSUER, user.AccessToken)
	if err != nil {
		log.Printf("ioprodz account grants for %s: %v", user.UserID, err)
		return nil
	}
	return ids
}

func getSession(account auth_models.Account, user goth.User, ua useragent.UserAgent, sessionRepo auth_models.SessionRepository) auth_models.Session {
	sessionHash := getHash(account, ua)
	session, sessionNotFound := sessionRepo.GetByHash(sessionHash)

	if sessionNotFound != nil {
		session = auth_models.NewSession(account.Id, ua.String, sessionHash, user.AvatarURL, user.Name, user.Email, user.Provider)
		session.ProjectIds = projectIds(user)
		sessionRepo.Create(session)
	} else {
		session.SetLastUsedNow()
		session.Provider = user.Provider
		session.ProjectIds = projectIds(user)
		sessionRepo.Update(session)
	}
	return session
}

func getHash(account auth_models.Account, ua useragent.UserAgent) string {
	b := []byte(account.Id + ua.String)
	hash := fnv.New64a()
	hash.Write(b)
	return fmt.Sprint(hash.Sum64())
}

func getAccount(accountRepo auth_models.AccountRepository, user goth.User) auth_models.Account {
	account, accountNotFound := accountRepo.GetByProviderId(user.Provider, user.UserID)

	if accountNotFound != nil {
		account = auth_models.NewAccount(user.Email, user.Provider, user.UserID)
		accountRepo.Create(account)
	}
	return account
}

func CreateOAuthLoginHandler(accountRepo auth_models.AccountRepository, sessionRepo auth_models.SessionRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// try to get the user without re-authenticating, except when the
		// sign-in is asked for again (?next=): the last token was refused
		again := r.URL.Query().Get("next") != ""
		if user, err := gothic.CompleteUserAuth(w, r); err == nil && !again {
			ua := useragent.Parse(r.UserAgent())
			account := getAccount(accountRepo, user)
			session := getSession(account, user, ua, sessionRepo)
			auth_infra.SetAuthCookie(w, r, auth_infra.CookieData{Id: session.Id})
			keepIoprodzToken(w, r, user)
			w.Header().Set("Location", landingPage(w, r, user))
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			rememberNext(w, r)
			gothic.BeginAuthHandler(w, r)
		}

	}
}

func CreateLogoutHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		auth_infra.ClearAuthCookie(w, r)
		auth_infra.ClearIoprodzToken(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
