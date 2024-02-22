package auth

import (
	"ioprodz/common/config"

	"github.com/gorilla/sessions"
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

func NewAppCookieStore() *sessions.CookieStore {
	var conf = config.Load()
	store := sessions.NewCookieStore([]byte(conf.AUTH_APP_COOKIE_SECRET))
	store.Options.Path = "/"
	store.Options.Domain = conf.APP_DOMAIN
	store.Options.HttpOnly = true
	store.Options.Secure = conf.IS_PRODUCTION
	store.Options.MaxAge = 86400

	return store
}
