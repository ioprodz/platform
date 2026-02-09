package seo

import (
	"encoding/xml"
	"fmt"
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/config"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func ConfigureRoutes(router *mux.Router, blogRepo blog_models.BlogRepository) {
	router.HandleFunc("/robots.txt", RobotsHandler).Methods("GET")
	router.HandleFunc("/sitemap.xml", CreateSitemapHandler(blogRepo)).Methods("GET")
	router.HandleFunc("/llms.txt", LLMSTxtHandler).Methods("GET")
}

func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	baseURL := config.Load().BASE_URL
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, `User-agent: *
Allow: /

Sitemap: %s/sitemap.xml

# LLM-friendly site summary
# See %s/llms.txt for a structured overview of this site
`, baseURL, baseURL)
}

func LLMSTxtHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("llms.txt")
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(data)
}

type urlSet struct {
	XMLName xml.Name  `xml:"urlset"`
	XMLNS   string    `xml:"xmlns,attr"`
	URLs    []siteURL `xml:"url"`
}

type siteURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

func CreateSitemapHandler(blogRepo blog_models.BlogRepository) http.HandlerFunc {
	type staticPage struct {
		path       string
		changefreq string
		priority   string
	}

	pages := []staticPage{
		{"/", "weekly", "1.0"},
		{"/consulting", "monthly", "0.8"},
		{"/consulting/it-strategy", "monthly", "0.7"},
		{"/consulting/coaching", "monthly", "0.7"},
		{"/solutions", "monthly", "0.8"},
		{"/solutions/ai-engine", "monthly", "0.7"},
		{"/solutions/chat-collaboration", "monthly", "0.7"},
		{"/solutions/collaborative-editing", "monthly", "0.7"},
		{"/solutions/search-rag", "monthly", "0.7"},
		{"/blog", "weekly", "0.8"},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := config.Load().BASE_URL
		now := time.Now().Format("2006-01-02")

		urls := make([]siteURL, 0, len(pages)+10)
		for _, p := range pages {
			urls = append(urls, siteURL{
				Loc:        baseURL + p.path,
				LastMod:    now,
				ChangeFreq: p.changefreq,
				Priority:   p.priority,
			})
		}

		posts, err := blogRepo.ListPublished()
		if err == nil {
			for _, post := range posts {
				lastMod := now
				if post.PublishedAt != "" {
					if t, err := time.Parse(time.RFC3339, post.PublishedAt); err == nil {
						lastMod = t.Format("2006-01-02")
					}
				}
				urls = append(urls, siteURL{
					Loc:        baseURL + "/blog/" + post.Id,
					LastMod:    lastMod,
					ChangeFreq: "monthly",
					Priority:   "0.6",
				})
			}
		}

		sitemap := urlSet{
			XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
			URLs:  urls,
		}

		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Write([]byte(xml.Header))
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		enc.Encode(sitemap)
	}
}
