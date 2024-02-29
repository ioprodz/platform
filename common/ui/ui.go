package ui

import (
	"ioprodz/auth"
	"net/http"
	"text/template"
)

func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {

	user, err := auth.GetUserSession(w, r)
	isAuthenticated := err == nil
	tpl, err := template.ParseFiles("common/ui/layout.html", "common/ui/header.html", "common/ui/footer.html", tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.ExecuteTemplate(w, "layout", map[string]interface{}{"contentData": data, "isAuthenticated": isAuthenticated, "user": user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderAdminPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {

	user, err := auth.GetUserSession(w, r)
	isAuthenticated := err == nil
	tpl, err := template.ParseFiles("common/ui/layout.html", "common/ui/header.html", "common/ui/footer.html", "common/ui/admin-layout.html", tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "layout", map[string]interface{}{"contentData": data, "isAuthenticated": isAuthenticated, "layout": "admin", "user": user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Render404(w http.ResponseWriter, r *http.Request) {
	RenderPage(w, r, "common/ui/not-found", nil)
}
