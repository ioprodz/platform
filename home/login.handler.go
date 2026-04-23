package home

import (
	"ioprodz/common/ui"
	"net/http"
)

func GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		TitleKey: "home:meta.login.title",
		DescKey:  "home:meta.login.description",
		Path:     "/login",
		OGType:   "website",
		NoIndex:  true,
	}
	ui.RenderPageWithMeta(w, r, "home/login.template", nil, meta)
}
