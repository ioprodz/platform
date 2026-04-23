package home

import (
	"ioprodz/common/ui"
	"net/http"
)

func GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		Title:       "Sign in",
		Description: "Owner sign-in page.",
		Path:        "/login",
		OGType:      "website",
		NoIndex:     true,
	}
	ui.RenderPageWithMeta(w, r, "home/login.template", nil, meta)
}
