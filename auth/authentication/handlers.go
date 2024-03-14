package auth_authentication

import (
	"fmt"
	"hash/fnv"
	auth_infra "ioprodz/auth/_infra"
	auth_models "ioprodz/auth/_models"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/mileusna/useragent"
)

func CreateOAuthCallbackHandler(accountRepo auth_models.AccountRepository, sessionRepo auth_models.SessionRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Println("Unauthorized: " + err.Error())
			return
		}

		ua := useragent.Parse(r.UserAgent())
		account := getAccount(accountRepo, user)
		session := getSession(account, user, ua, sessionRepo)

		auth_infra.SetAuthCookie(w, r, auth_infra.CookieData{Id: session.Id})

		w.Header().Set("Location", "/feed")
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}

func getSession(account auth_models.Account, user goth.User, ua useragent.UserAgent, sessionRepo auth_models.SessionRepository) auth_models.Session {
	sessionHash := getHash(account, ua)
	session, sessionNotFound := sessionRepo.GetByHash(sessionHash)

	if sessionNotFound != nil {
		session = auth_models.NewSession(account.Id, ua.String, sessionHash, user.AvatarURL, user.Name, user.Email)
		sessionRepo.Create(session)
	} else {
		session.SetLastUsedNow()
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
		// try to get the user without re-authenticating
		if user, err := gothic.CompleteUserAuth(w, r); err == nil {
			ua := useragent.Parse(r.UserAgent())
			account := getAccount(accountRepo, user)
			session := getSession(account, user, ua, sessionRepo)
			auth_infra.SetAuthCookie(w, r, auth_infra.CookieData{Id: session.Id})
			w.Header().Set("Location", "/feed")
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			gothic.BeginAuthHandler(w, r)
		}

	}
}

func CreateLogoutHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		auth_infra.ClearAuthCookie(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
