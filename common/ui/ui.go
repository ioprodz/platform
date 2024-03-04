package ui

import (
	"ioprodz/common/policies"
	"net/http"
	"text/template"
)

func RenderPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {

	user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)

	tpl, err := template.ParseFiles("common/ui/layout.html", "common/ui/header.html", "common/ui/footer.html", tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.ExecuteTemplate(w, "layout", map[string]interface{}{"contentData": data, "isAuthenticated": user.IsAuthenticated(), "user": user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderAdminPage(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {

	user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)

	tpl, err := template.ParseFiles("common/ui/layout.html", "common/ui/header.html", "common/ui/footer.html", "common/ui/admin-layout.html", tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "layout", map[string]interface{}{"contentData": data, "isAuthenticated": user.IsAuthenticated(), "layout": "admin", "user": user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Render404(w http.ResponseWriter, r *http.Request) {
	RenderPage(w, r, "common/ui/not-found", nil)
}
