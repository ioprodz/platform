package home

import (
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)

	meta := ui.PageMeta{
		Title:       "",
		Description: "ioprodz helps software teams meet business goals through strategic consulting, hands-on coaching, and production-grade AI platform components.",
		Path:        "/",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "home/template", map[string]interface{}{"IsAuthenticated": user.IsAuthenticated()}, meta)
}
