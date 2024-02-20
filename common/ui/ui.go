package ui

import (
	"net/http"
	"text/template"
)

func RenderPage(w http.ResponseWriter, tmpl string, data interface{}) {
	tpl, err := template.ParseFiles("common/ui/layout.html", "common/ui/header.html", "common/ui/footer.html", tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Render404(w http.ResponseWriter) {
	RenderPage(w, "common/ui/not-found", nil)
}
