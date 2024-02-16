package home

import (
	"net/http"
	"text/template"
)

type PageData struct {
	Name string
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// Define data to be passed to the template
	data := PageData{Name: "John"}

	// Parse the template file
	tmpl, err := template.ParseFiles("home/template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the provided data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
