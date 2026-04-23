package home

import (
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)

	meta := ui.PageMeta{
		DescKey: "home:meta.home.description",
		Path:    "/",
		OGType:  "website",
	}
	ui.RenderPageWithMeta(w, r, "home/template", map[string]interface{}{"IsAuthenticated": user.IsAuthenticated()}, meta)
}
