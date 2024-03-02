package members_explore

import (
	"ioprodz/common/ui"
	"net/http"
)

func CreateGetHandler() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		ui.RenderPage(w, r, "members/explore/index", nil)
	}
}
