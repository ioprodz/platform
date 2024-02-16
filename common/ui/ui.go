package ui

import (
	"net/http"
	"text/template"
)

func RenderPage(w http.ResponseWriter, tmpl string, data interface{}) {
	tpl, err := template.ParseFiles("common/ui/layout.html", tmpl+".html")
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
