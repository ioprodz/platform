package home

import (
	"ioprodz/common/ui"
	"net/http"
)

type PageData struct {
	Name string
}

func GetHandler(w http.ResponseWriter, r *http.Request) {

	data := PageData{Name: "John"}

	ui.RenderPage(w, "home/template", data)

}
