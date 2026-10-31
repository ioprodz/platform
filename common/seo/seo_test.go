package seo

import (
	blog_infra "ioprodz/blog/_infra"
	"net/http/httptest"
	"strings"
	"testing"
)

func sitemap(t *testing.T, repo *blog_infra.BlogMemoryRepository) string {
	t.Helper()
	t.Setenv("APP_ENV", "production") // config.Load reads the environment, no .env
	t.Setenv("BASE_URL", "https://ioprodz.com")
	w := httptest.NewRecorder()
	CreateSitemapHandler(repo)(w, httptest.NewRequest("GET", "/sitemap.xml", nil))
	return w.Body.String()
}

func TestSitemapLeavesOutAnEmptyBlog(t *testing.T) {
	body := sitemap(t, blog_infra.CreateMemoryBlogRepo()) // seeded posts are all unpublished

	if strings.Contains(body, "/blog</loc>") {
		t.Fatalf("expected no /blog without published posts, got\n%s", body)
	}
	if !strings.Contains(body, "/ar/solutions/collaborative-editing</loc>") {
		t.Fatalf("expected the Arabic solution pages, got\n%s", body)
	}
}

func TestSitemapGivesOnlyPostsALastmod(t *testing.T) {
	repo := blog_infra.CreateMemoryBlogRepo()
	posts, _ := repo.List()
	post := posts[0]
	post.PublishedAt = "2026-03-14T10:00:00Z"
	repo.Update(post)

	body := sitemap(t, repo)

	if !strings.Contains(body, "/blog</loc>") {
		t.Fatalf("expected /blog once a post is published, got\n%s", body)
	}
	if n := strings.Count(body, "<lastmod>"); n != 1 || !strings.Contains(body, "<lastmod>2026-03-14</lastmod>") {
		t.Fatalf("expected one lastmod, the post's publish date; got %d in\n%s", n, body)
	}
}
