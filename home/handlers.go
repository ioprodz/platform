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
		Description: "I turn vibe-coded software into real products in 90 days — on sovereign AI infrastructure (OpenAI, Claude, Gemini, Grok, or local), shipped daily, with the boring engineering every AI tutorial skips: evals, rollback, observability, domain modeling. Independent engineering consultancy serving Europe and North Africa.",
		Path:        "/",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "home/template", map[string]interface{}{"IsAuthenticated": user.IsAuthenticated()}, meta)
}
