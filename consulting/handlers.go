package consulting

import (
	"ioprodz/common/ui"
	"net/http"
)

func OverviewHandler(w http.ResponseWriter, r *http.Request) {
	ui.RenderPage(w, r, "consulting/overview", nil)
}

func ITStrategyHandler(w http.ResponseWriter, r *http.Request) {
	ui.RenderPage(w, r, "consulting/it-strategy", nil)
}

func CoachingHandler(w http.ResponseWriter, r *http.Request) {
	ui.RenderPage(w, r, "consulting/coaching", nil)
}
