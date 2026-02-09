package consulting

import (
	"ioprodz/common/ui"
	"net/http"
)

func OverviewHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		Title:       "Consulting Services",
		Description: "Strategic technology leadership, IT strategy, and hands-on team coaching for software organizations.",
		Path:        "/consulting",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "consulting/overview", nil, meta)
}

func ITStrategyHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		Title:       "IT Strategy Consulting",
		Description: "Technology roadmaps, architecture decisions, and hands-on engineering leadership that align technology to business outcomes.",
		Path:        "/consulting/it-strategy",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "consulting/it-strategy", nil, meta)
}

func CoachingHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		Title:       "Team & Individual Coaching",
		Description: "Build a learning culture rooted in Extreme Programming and Lean principles. Coaching on code reviews, TDD, pair programming, and engineering practices.",
		Path:        "/consulting/coaching",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "consulting/coaching", nil, meta)
}
