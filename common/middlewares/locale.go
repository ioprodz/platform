package middlewares

import (
	"context"
	"ioprodz/common/i18n"
	"net/http"
	"strings"
)

type localeCtxKeyType string

const LocaleCtxKey localeCtxKeyType = "locale"
const CanonicalPathCtxKey localeCtxKeyType = "canonicalPath"

// skipPrefixes are paths that must never be treated as locale-prefixed and
// don't participate in the i18n system at all.
var skipPrefixes = []string{
	"/static/",
	"/auth/",
	"/admin/",
	"/members/",
	"/qna/",
	"/logout",
	"/robots.txt",
	"/sitemap.xml",
	"/llms.txt",
	"/profile",
	"/feed",
}

// translatablePrefixes are the public paths that are in-scope for translation.
// A request under /fr/... or /ar/... that doesn't match one of these is
// redirected to the unprefixed (English) equivalent.
var translatablePrefixes = []string{
	"/about",
	"/login",
	"/consulting",
	"/solutions",
}

func shouldSkip(path string) bool {
	for _, p := range skipPrefixes {
		if path == p || strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func isTranslatable(path string) bool {
	return IsTranslatablePath(path)
}

// IsTranslatablePath reports whether a canonical (locale-stripped) path is
// in-scope for i18n. Used by the UI layer to decide whether to emit hreflang
// alternates for a page.
func IsTranslatablePath(path string) bool {
	if path == "/" {
		return true
	}
	for _, p := range translatablePrefixes {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// LocaleResolver detects the request locale from the URL path prefix, stashes
// it plus the canonical (de-prefixed) path in the request context, and rewrites
// r.URL.Path so downstream middlewares and the router see the canonical path.
func LocaleResolver(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		lang := i18n.DefaultLang

		if shouldSkip(path) {
			ctx := context.WithValue(r.Context(), LocaleCtxKey, lang)
			ctx = context.WithValue(ctx, CanonicalPathCtxKey, path)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		stripped := path
		for _, code := range i18n.PrefixedLangs {
			prefix := "/" + string(code)
			if path == prefix {
				lang = code
				stripped = "/"
				break
			}
			if strings.HasPrefix(path, prefix+"/") {
				lang = code
				stripped = path[len(prefix):]
				break
			}
		}

		if lang != i18n.DefaultLang && !isTranslatable(stripped) {
			target := stripped
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), LocaleCtxKey, lang)
		ctx = context.WithValue(ctx, CanonicalPathCtxKey, stripped)

		r2 := r.Clone(ctx)
		r2.URL.Path = stripped
		r2.URL.RawPath = ""
		next.ServeHTTP(w, r2)
	})
}
