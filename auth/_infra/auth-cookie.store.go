package auth_infra

import (
	"ioprodz/common/config"
	"net/http"

	"github.com/gorilla/sessions"
)

func NewAuthCookieStore() *sessions.CookieStore {
	var conf = config.Load()
	store := sessions.NewCookieStore([]byte(conf.AUTH_APP_COOKIE_SECRET))
	store.Options.Path = "/"
	store.Options.Domain = conf.APP_DOMAIN
	store.Options.HttpOnly = true
	store.Options.Secure = conf.IS_PRODUCTION
	store.Options.MaxAge = 86400

	return store
}

var cookieStore = NewAuthCookieStore()

var cookieName = "ioprodz-session"

type AuthCookieError struct {
	msg string
}
type CookieData struct {
	Id string
}

func (e *AuthCookieError) Error() string { return e.msg }

func SetAuthCookie(w http.ResponseWriter, r *http.Request, data CookieData) error {
	cookie, err := cookieStore.Get(r, cookieName)
	cookie.Values["id"] = data.Id
	cookie.Save(r, w)
	return err
}

func GetAuthCookie(w http.ResponseWriter, r *http.Request) (CookieData, error) {
	cookie, err := cookieStore.Get(r, cookieName)
	if err != nil {
		return CookieData{}, &AuthCookieError{msg: "could not load session"}
	}
	id, ok := cookie.Values["id"].(string)
	if !ok {
		return CookieData{}, &AuthCookieError{msg: "could not find userId on session"}
	}
	return CookieData{
		Id: id,
	}, nil
}

func ClearAuthCookie(w http.ResponseWriter, r *http.Request) {
	cookie, _ := cookieStore.Get(r, cookieName)
	for key := range cookie.Values {
		delete(cookie.Values, key)
	}
	cookie.Save(r, w)
}
