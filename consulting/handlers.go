package consulting

import (
	"ioprodz/common/ui"
	"net/http"
)

func OverviewHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		TitleKey: "consulting:meta.overview.title",
		DescKey:  "consulting:meta.overview.description",
		Path:     "/consulting",
		OGType:   "website",
	}
	ui.RenderPageWithMeta(w, r, "consulting/overview", nil, meta)
}

func ITStrategyHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		TitleKey: "consulting:meta.itStrategy.title",
		DescKey:  "consulting:meta.itStrategy.description",
		Path:     "/consulting/it-strategy",
		OGType:   "website",
	}
	ui.RenderPageWithMeta(w, r, "consulting/it-strategy", nil, meta)
}

func CoachingHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		TitleKey: "consulting:meta.coaching.title",
		DescKey:  "consulting:meta.coaching.description",
		Path:     "/consulting/coaching",
		OGType:   "website",
	}
	ui.RenderPageWithMeta(w, r, "consulting/coaching", nil, meta)
}
