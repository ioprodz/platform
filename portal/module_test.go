package portal

import "testing"

func TestServicesForShowsOnlyGrantedProjects(t *testing.T) {
	catalog := loadCatalog(`[
		{"projectId":"1","name":"Mail","url":"https://mail.example"},
		{"projectId":"1","name":"Webmail","url":"https://webmail.example"},
		{"projectId":"2","name":"Git","url":"https://git.example"}
	]`)

	visible := servicesFor(catalog, []string{"1", "3"})

	if len(visible) != 2 || visible[0].Name != "Mail" || visible[1].Name != "Webmail" {
		t.Fatalf("expected Mail and Webmail, got %+v", visible)
	}
}

func TestLoadCatalogToleratesMissingOrInvalidConfig(t *testing.T) {
	if len(loadCatalog("")) != 0 || len(loadCatalog("not json")) != 0 {
		t.Fatal("expected an empty catalog")
	}
}
