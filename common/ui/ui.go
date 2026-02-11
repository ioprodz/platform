package ui

import (
	"bytes"
	"ioprodz/common/config"
	"ioprodz/common/policies"
	"log"
	"net/http"
	"sync"
	"text/template"
)

// PageMeta holds SEO metadata for a page.
type PageMeta struct {
	Title                string
	Description          string
	Path                 string // e.g. "/consulting/coaching"
	OGType               string // "website" or "article"
	Keywords             string
	ArticlePublishedTime string // RFC3339, only for blog articles
}

// DefaultMeta returns fallback metadata for pages that haven't been migrated yet.
func DefaultMeta() PageMeta {
	return PageMeta{
		Title:       "",
		Description: "ioprodz helps software teams meet business goals through strategic consulting, hands-on coaching, and production-grade AI platform components.",
		Path:        "/",
		OGType:      "website",
	}
}

// FullTitle returns "PageTitle | ioprodz" or just the brand title for the homepage.
func (m PageMeta) FullTitle() string {
	if m.Title == "" {
		return "ioprodz — Engineering Leadership & AI-Powered Software"
	}
	return m.Title + " | ioprodz"
}

// CanonicalURL returns the full canonical URL for this page.
func (m PageMeta) CanonicalURL() string {
	return config.Load().BASE_URL + m.Path
}

// OGImageURL returns the absolute URL of the Open Graph image.
func (m PageMeta) OGImageURL() string {
	return config.Load().BASE_URL + "/static/img/og-image.png"
}

var (
	templateCache sync.Map
	commonFiles   = []string{"common/ui/layout.html", "common/ui/header.html", "common/ui/footer.html"}
)

func getTemplate(files ...string) (*template.Template, error) {
	key := files[len(files)-1] // use the page-specific template as cache key
	if cached, ok := templateCache.Load(key); ok {
		return cached.(*template.Template), nil
	}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	templateCache.Store(key, tpl)
	return tpl, nil
}

func RenderPageWithMeta(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}, meta PageMeta) {
	user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)

	tpl, err := getTemplate(append(commonFiles, tmpl+".html")...)
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, "layout", map[string]interface{}{
		"contentData":     data,
		"isAuthenticated": user.IsAuthenticated(),
		"user":            user,
		"meta":            meta,
		"baseURL":         config.Load().BASE_URL,
	})
	if err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	RenderPageWithMeta(w, r, tmpl, data, DefaultMeta())
}

func RenderAdminPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)

	adminFiles := append(commonFiles, "common/ui/admin-layout.html", tmpl+".html")
	tpl, err := getTemplate(adminFiles...)
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	meta := PageMeta{Title: "Admin", Description: "Admin panel", OGType: "website"}
	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, "layout", map[string]interface{}{
		"contentData":     data,
		"isAuthenticated": user.IsAuthenticated(),
		"layout":          "admin",
		"user":            user,
		"meta":            meta,
		"baseURL":         config.Load().BASE_URL,
	})
	if err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func Render404(w http.ResponseWriter, r *http.Request) {
	meta := PageMeta{Title: "Page Not Found", Description: "The requested page was not found.", OGType: "website"}
	RenderPageWithMeta(w, r, "common/ui/not-found", nil, meta)
}
