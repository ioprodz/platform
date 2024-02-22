package auth

import "github.com/gorilla/sessions"

const IsProd = false

func NewAuthCookieStore() *sessions.CookieStore {
	const (
		key    = "secret"
		MaxAge = 86400 * 7
	)
	store := sessions.NewCookieStore([]byte(key))
	store.Options.Path = "/"
	store.Options.Domain = "localhost"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	return store
}

func NewAppCookieStore() *sessions.CookieStore {
	const (
		key    = "secret2"
		MaxAge = 86400 * 1
	)
	store := sessions.NewCookieStore([]byte(key))
	store.Options.Path = "/"
	store.Options.Domain = "localhost"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	return store
}
