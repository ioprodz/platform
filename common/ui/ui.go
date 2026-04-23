package ui

import (
	"bytes"
	"ioprodz/common/config"
	"ioprodz/common/i18n"
	"ioprodz/common/middlewares"
	"ioprodz/common/policies"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// PageMeta holds SEO metadata for a page.
//
// Two authoring modes are supported and can be mixed per-page:
//   - Literal: set Title/Description directly (used by blog, CV, admin).
//   - Translated: set TitleKey/DescKey to "file:dotted.path" translation keys.
//     The Resolved* methods look them up against the current request locale,
//     with English fallback.
type PageMeta struct {
	Title                string
	Description          string
	TitleKey             string
	DescKey              string
	Path                 string // canonical path without locale prefix, e.g. "/consulting/coaching"
	OGType               string // "website" or "article"
	Keywords             string
	ArticlePublishedTime string // RFC3339, only for blog articles
	NoIndex              bool   // true → emits noindex,nofollow robots meta
}

// DefaultMeta returns fallback metadata for pages that haven't been migrated yet.
func DefaultMeta() PageMeta {
	return PageMeta{
		Title:       "",
		Description: "I turn vibe-coded software into real products in 90 days — on sovereign AI infrastructure (OpenAI, Claude, Gemini, Grok, or local), shipped daily, with the boring engineering every AI tutorial skips: evals, rollback, observability, domain modeling. Independent engineering consultancy serving Europe and North Africa.",
		Path:        "/",
		OGType:      "website",
	}
}

// ResolvedTitle returns the localized title if TitleKey is set, otherwise the literal Title.
func (m PageMeta) ResolvedTitle(lang i18n.Lang) string {
	if m.TitleKey != "" {
		return i18n.T(lang, m.TitleKey)
	}
	return m.Title
}

// ResolvedDescription returns the localized description if DescKey is set, otherwise the literal Description.
func (m PageMeta) ResolvedDescription(lang i18n.Lang) string {
	if m.DescKey != "" {
		return i18n.T(lang, m.DescKey)
	}
	return m.Description
}

// FullTitle returns "PageTitle | ioprodz" or just the brand title for the homepage.
func (m PageMeta) FullTitle(lang i18n.Lang) string {
	t := m.ResolvedTitle(lang)
	if t == "" {
		if v := i18n.T(lang, "common:meta.siteTitleHome"); v != "" && v != "common:meta.siteTitleHome" {
			return v
		}
		return "ioprodz — From vibe coding to real product, in 90 days"
	}
	return t + " | ioprodz"
}

// CanonicalURL returns the full canonical URL for this page in the given language.
func (m PageMeta) CanonicalURL(lang i18n.Lang) string {
	return config.Load().BASE_URL + i18n.MetaFor(lang).URLPrefix + m.Path
}

// OGImageURL returns the absolute URL of the Open Graph image.
func (m PageMeta) OGImageURL() string {
	return config.Load().BASE_URL + "/static/img/og-image.png"
}

// LangFrom extracts the resolved locale from the request context.
func LangFrom(r *http.Request) i18n.Lang {
	if v, ok := r.Context().Value(middlewares.LocaleCtxKey).(i18n.Lang); ok {
		return v
	}
	return i18n.DefaultLang
}

// CanonicalPathFrom returns the request's canonical path (without locale prefix).
func CanonicalPathFrom(r *http.Request) string {
	if v, ok := r.Context().Value(middlewares.CanonicalPathCtxKey).(string); ok {
		return v
	}
	return "/"
}

// stubFuncs are registered at parse time so that templates referencing `t` or
// `localizedPath` parse cleanly. The real per-request implementations replace
// these via tpl.Clone().Funcs(...) in RenderPageWithMeta.
var stubFuncs = template.FuncMap{
	"t":             func(string) string { return "" },
	"localizedPath": func(i18n.Lang, string) string { return "" },
	"lp":            func(string) string { return "" },
}

var (
	templateCache sync.Map
	commonFiles   = []string{"common/ui/layout.html", "common/ui/header.html", "common/ui/footer.html"}
)

func getTemplate(files ...string) (*template.Template, error) {
	key := files[len(files)-1]
	if cached, ok := templateCache.Load(key); ok {
		return cached.(*template.Template), nil
	}
	tpl, err := template.New(filepath.Base(files[0])).Funcs(stubFuncs).ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	templateCache.Store(key, tpl)
	return tpl, nil
}

func perRequestFuncs(lang i18n.Lang) template.FuncMap {
	buildPath := func(l i18n.Lang, path string) string {
		m := i18n.MetaFor(l)
		if path == "" {
			path = "/"
		}
		if m.URLPrefix == "" {
			return path
		}
		if path == "/" {
			return m.URLPrefix + "/"
		}
		return m.URLPrefix + path
	}
	return template.FuncMap{
		"t":             func(key string) string { return i18n.T(lang, key) },
		"localizedPath": buildPath,
		"lp":            func(path string) string { return buildPath(lang, path) },
	}
}

func renderWithLayout(w http.ResponseWriter, r *http.Request, files []string, layoutData map[string]interface{}) {
	lang := LangFrom(r)

	base, err := getTemplate(files...)
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tpl, err := base.Clone()
	if err != nil {
		log.Printf("template clone error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tpl.Funcs(perRequestFuncs(lang))

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "layout", layoutData); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func RenderPageWithMeta(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}, meta PageMeta) {
	user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
	lang := LangFrom(r)
	langMeta := i18n.MetaFor(lang)
	canonicalPath := CanonicalPathFrom(r)
	if meta.Path == "" {
		meta.Path = canonicalPath
	}

	layoutData := map[string]interface{}{
		"contentData":       data,
		"isAuthenticated":   user.IsAuthenticated(),
		"user":              user,
		"meta":              meta,
		"baseURL":           config.Load().BASE_URL,
		"lang":              string(lang),
		"langMeta":          langMeta,
		"canonicalPath":     canonicalPath,
		"allLangs":          i18n.AllMetas(),
		"pageTitle":         meta.FullTitle(lang),
		"pageDescription":   meta.ResolvedDescription(lang),
		"pageCanonicalURL":  meta.CanonicalURL(lang),
		"translatable":     middlewares.IsTranslatablePath(meta.Path),
	}

	renderWithLayout(w, r, append(commonFiles, tmpl+".html"), layoutData)
}

func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	RenderPageWithMeta(w, r, tmpl, data, DefaultMeta())
}

func RenderAdminPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
	lang := LangFrom(r)
	langMeta := i18n.MetaFor(lang)
	canonicalPath := CanonicalPathFrom(r)

	meta := PageMeta{Title: "Admin", Description: "Admin panel", OGType: "website", Path: canonicalPath}

	layoutData := map[string]interface{}{
		"contentData":      data,
		"isAuthenticated":  user.IsAuthenticated(),
		"layout":           "admin",
		"user":             user,
		"meta":             meta,
		"baseURL":          config.Load().BASE_URL,
		"lang":             string(lang),
		"langMeta":         langMeta,
		"canonicalPath":    canonicalPath,
		"allLangs":         i18n.AllMetas(),
		"pageTitle":        meta.FullTitle(lang),
		"pageDescription":  meta.ResolvedDescription(lang),
		"pageCanonicalURL": meta.CanonicalURL(lang),
		"translatable":    false,
	}

	adminFiles := append(commonFiles, "common/ui/admin-layout.html", tmpl+".html")
	renderWithLayout(w, r, adminFiles, layoutData)
}

func Render404(w http.ResponseWriter, r *http.Request) {
	meta := PageMeta{Title: "Page Not Found", Description: "The requested page was not found.", OGType: "website"}
	RenderPageWithMeta(w, r, "common/ui/not-found", nil, meta)
}
