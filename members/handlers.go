package members

import (
	"ioprodz/common/ui"
	"net/http"
)

type PageData struct {
	Name string
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// Define data to be passed to the template
	data := PageData{Name: "Smith"}

	// Parse the template file
	ui.RenderPage(w, r, "profile/template", data)
}
