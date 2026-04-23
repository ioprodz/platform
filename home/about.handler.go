package home

import (
	"ioprodz/common/ui"
	"net/http"
)

func GetAboutHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		Title:       "Manifesto — The counter-voice to AI slop",
		Description: "Why I started ioprodz in late 2025. The counter-voice to AI slop: boring engineering is the competitive advantage, and a small team with the right practices outlasts a big team stacking tokens.",
		Path:        "/about",
		OGType:      "article",
	}
	ui.RenderPageWithMeta(w, r, "home/about.template", nil, meta)
}
