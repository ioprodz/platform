package auth_security

import (
	"errors"
	"html"
	auth_infra "ioprodz/auth/_infra"
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/config"
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/mileusna/useragent"
	"github.com/xeonx/timeago"
)

type SessionView struct {
	Id        string
	Title     string
	IsCurrent bool
	CreatedOn string
	LastUsed  string
}

func CreateSecurityPageHandler(sessionRepo auth_models.SessionRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		sessions := sessionRepo.GetByAccountId(user.Id)

		sessionViewList := []SessionView{}
		for _, session := range sessions {
			ua := useragent.Parse(session.UaString)

			createdOn, _ := time.Parse(time.RFC3339, session.FirstCreatedAt)
			lastUsed, _ := time.Parse(time.RFC3339, session.LastUsedAt)

			sessionViewList = append(sessionViewList, SessionView{
				Id:        session.Id,
				Title:     ua.Name + " on " + ua.Device + " " + ua.OS + "",
				CreatedOn: timeago.English.Format(createdOn),
				LastUsed:  timeago.English.Format(lastUsed),
				IsCurrent: session.Id == user.SessionId,
			})
		}

		page := SecurityPage{Sessions: sessionViewList}
		if user.HasIoprodzAccount() {
			page.Ioprodz = ioprodzView(auth_infra.GetIoprodzToken(r))
		}
		ui.RenderPage(w, r, "auth/security/settings", page)
	}
}

type SecurityPage struct {
	Sessions []SessionView
	// sign-in methods of an ioprodz account; nil for other providers
	Ioprodz *IoprodzView
}

type IoprodzView struct {
	// no usable token (signed in before it was kept, expired, or weaker than
	// the account now requires)
	SignInAgain bool
	Error       string
	Methods     auth_infra.IoprodzSignInMethods
	Links       AddLinks
}

// AddLinks open Zitadel's login pages (with our branding) that add a sign-in
// method; they work from the user's login session there.
type AddLinks struct {
	Passkey       string
	SecurityKey   string
	Authenticator string
	Password      string
}

// KeyGroup is one list of removable keys on the page (passkeys, security keys).
type KeyGroup struct {
	Title    string
	Hint     string
	Kind     string
	Keys     []auth_infra.IoprodzKey
	Add      string
	AddLabel string
}

func (v IoprodzView) PasskeyGroup() KeyGroup {
	return KeyGroup{Title: "Passkeys", Hint: "Sign in with your fingerprint, face or device PIN.", Kind: "passkey",
		Keys: v.Methods.Passkeys, Add: v.Links.Passkey, AddLabel: "Add a passkey"}
}

func (v IoprodzView) SecurityKeyGroup() KeyGroup {
	return KeyGroup{Title: "Security keys", Hint: "USB or NFC keys used as a second step.", Kind: "key",
		Keys: v.Methods.SecurityKeys, Add: v.Links.SecurityKey, AddLabel: "Add a key"}
}

func addLinks(issuer string, loginName string) AddLinks {
	query := "?" + url.Values{"loginName": {loginName}}.Encode()
	login := issuer + "/ui/v2/login"
	return AddLinks{
		Passkey:       login + "/passkey/set" + query,
		SecurityKey:   login + "/u2f/set" + query,
		Authenticator: login + "/otp/time-based/set" + query,
		Password:      login + "/password/change" + query,
	}
}

func ioprodzView(token auth_infra.IoprodzToken) *IoprodzView {
	if !token.Valid() {
		return &IoprodzView{SignInAgain: true}
	}
	issuer := config.Load().IOPRODZ_AUTH_ISSUER
	view := &IoprodzView{Links: addLinks(issuer, token.LoginName)}
	methods, err := auth_infra.FetchIoprodzSignInMethods(issuer, token)
	if errors.Is(err, auth_infra.ErrIoprodzSignInAgain) {
		return &IoprodzView{SignInAgain: true}
	}
	if err != nil {
		log.Printf("ioprodz sign-in methods for %s: %v", token.UserId, err)
		view.Error = "Could not read your sign-in methods: " + err.Error()
		return view
	}
	view.Methods = methods
	return view
}

// CreateRemoveSignInMethodHandler removes a passkey, security key or the
// authenticator app from the user's ioprodz account. It answers with the text
// that replaces the Remove button.
func CreateRemoveSignInMethodHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := auth_infra.GetIoprodzToken(r)
		if !token.Valid() {
			w.Write([]byte("Sign in again first"))
			return
		}
		vars := mux.Vars(r)
		err := auth_infra.RemoveIoprodzSignInMethod(config.Load().IOPRODZ_AUTH_ISSUER, token.AccessToken, vars["kind"], vars["id"])
		if errors.Is(err, auth_infra.ErrIoprodzSignInAgain) {
			w.Write([]byte("Sign in again first"))
			return
		}
		if err != nil {
			log.Printf("remove ioprodz %s for %s: %v", vars["kind"], token.UserId, err)
			w.Write([]byte("Not removed: " + html.EscapeString(err.Error())))
			return
		}
		w.Write([]byte("Removed"))
	}
}

func CreateRevokeSessionHandler(sessionRepo auth_models.SessionRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		vars := mux.Vars(r)
		sessionId := vars["id"]

		// only the user's own sessions; someone else's answers like a missing one
		session, err := sessionRepo.Get(sessionId)
		if err != nil || session.AccountId != user.Id {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		sessionRepo.Delete(sessionId)
		w.Write([]byte("Disconnected Successfully"))
	}
}
