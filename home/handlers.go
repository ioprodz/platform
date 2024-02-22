package home

import (
	"ioprodz/common/ui"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {

	ui.RenderPage(w, r, "home/template", nil)

}
