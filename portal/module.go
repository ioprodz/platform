package portal

import (
	"encoding/json"
	"ioprodz/common/config"
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	"log"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

// Service is one entry of the portal catalog. A user sees it when they hold
// a role in its Zitadel project. The catalog comes from PORTAL_SERVICES (JSON
// list) so service addresses stay out of this public repo.
type Service struct {
	ProjectId   string `json:"projectId"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Note        string `json:"note"`
}

func loadCatalog(raw string) []Service {
	services := []Service{}
	if raw == "" {
		return services
	}
	if err := json.Unmarshal([]byte(raw), &services); err != nil {
		log.Printf("PORTAL_SERVICES is not valid JSON: %v", err)
	}
	return services
}

func servicesFor(catalog []Service, projectIds []string) []Service {
	visible := []Service{}
	for _, service := range catalog {
		if slices.Contains(projectIds, service.ProjectId) {
			visible = append(visible, service)
		}
	}
	return visible
}

func ConfigureModule(router *mux.Router) *mux.Router {
	catalog := loadCatalog(config.Load().PORTAL_SERVICES)

	router.HandleFunc("/portal", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		meta := ui.PageMeta{Title: "Portal", Path: "/portal", OGType: "website", NoIndex: true}
		ui.RenderPageWithMeta(w, r, "portal/index", map[string]interface{}{
			"User":              user,
			"HasIoprodzAccount": user.HasIoprodzAccount(),
			"Services":          servicesFor(catalog, user.ProjectIds),
		}, meta)
	}).Methods("GET")
	return router
}
