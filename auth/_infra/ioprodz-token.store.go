package auth_infra

import (
	"crypto/sha256"
	"ioprodz/common/config"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
)

// IoprodzToken is the user's own Zitadel access token, kept after an ioprodz
// sign-in so /security can list and remove the account's sign-in methods.
type IoprodzToken struct {
	AccessToken string
	UserId      string
	LoginName   string
	ExpiresAt   time.Time
}

func (t IoprodzToken) Valid() bool {
	return t.AccessToken != "" && time.Now().Before(t.ExpiresAt)
}

// The token is encrypted (not only signed), sent only to /security and
// dropped when it expires. Keys derive from the app cookie secret.
func newIoprodzTokenStore(secret string, secure bool) *sessions.CookieStore {
	hashKey := sha256.Sum256([]byte("ioprodz-token-hash:" + secret))
	blockKey := sha256.Sum256([]byte("ioprodz-token-block:" + secret))
	store := sessions.NewCookieStore(hashKey[:], blockKey[:])
	store.Options.Path = "/security"
	store.Options.HttpOnly = true
	store.Options.Secure = secure
	return store
}

// built on first use, so importing the package doesn't load the config
var ioprodzTokenStore = sync.OnceValue(func() *sessions.CookieStore {
	return newIoprodzTokenStore(config.Load().AUTH_APP_COOKIE_SECRET, config.Load().IS_PRODUCTION)
})

var ioprodzTokenCookieName = "ioprodz-token"

func SetIoprodzToken(w http.ResponseWriter, r *http.Request, token IoprodzToken) {
	cookie, _ := ioprodzTokenStore().New(r, ioprodzTokenCookieName)
	cookie.Values["token"] = token.AccessToken
	cookie.Values["user"] = token.UserId
	cookie.Values["login"] = token.LoginName
	cookie.Values["expires"] = token.ExpiresAt.Unix()
	cookie.Options.MaxAge = int(time.Until(token.ExpiresAt).Seconds())
	if cookie.Options.MaxAge <= 0 {
		return
	}
	if err := cookie.Save(r, w); err != nil {
		log.Printf("ioprodz token cookie: %v", err)
	}
}

// GetIoprodzToken returns the stored token, or an invalid one when there is
// none or it has expired.
func GetIoprodzToken(r *http.Request) IoprodzToken {
	cookie, err := ioprodzTokenStore().Get(r, ioprodzTokenCookieName)
	if err != nil {
		return IoprodzToken{}
	}
	accessToken, _ := cookie.Values["token"].(string)
	userId, _ := cookie.Values["user"].(string)
	loginName, _ := cookie.Values["login"].(string)
	expires, _ := cookie.Values["expires"].(int64)
	return IoprodzToken{
		AccessToken: accessToken,
		UserId:      userId,
		LoginName:   loginName,
		ExpiresAt:   time.Unix(expires, 0),
	}
}

func ClearIoprodzToken(w http.ResponseWriter, r *http.Request) {
	cookie, _ := ioprodzTokenStore().New(r, ioprodzTokenCookieName)
	cookie.Options.MaxAge = -1
	cookie.Save(r, w)
}
