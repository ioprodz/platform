package home

import (
	"ioprodz/common/ui"
	"net/http"
)

func GetAboutHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		TitleKey: "home:meta.about.title",
		DescKey:  "home:meta.about.description",
		Path:     "/about",
		OGType:   "article",
	}
	ui.RenderPageWithMeta(w, r, "home/about.template", nil, meta)
}
