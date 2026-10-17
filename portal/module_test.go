package portal

import "testing"

func TestTilesForListsGrantedServicesFirstAndLocksTheRest(t *testing.T) {
	services := []Service{
		{ProjectId: "2", Name: "Git"},
		{ProjectId: "1", Name: "Mail"},
		{ProjectId: "", Name: "Registry"},
		{ProjectId: "1", Name: "Webmail"},
	}

	tiles := tilesFor(services, []string{"1", "3"})

	got := []string{}
	for _, tile := range tiles {
		got = append(got, tile.Name)
	}
	want := []string{"Mail", "Webmail", "Git", "Registry"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected order %v, got %v", want, got)
		}
	}
	if !tiles[0].Granted || !tiles[1].Granted || tiles[2].Granted || tiles[3].Granted {
		t.Fatalf("expected only Mail and Webmail granted, got %+v", tiles)
	}
}

func TestServiceWithoutProjectIsLockedEvenForEmptyGrant(t *testing.T) {
	tiles := tilesFor([]Service{{Name: "Registry"}}, []string{""})
	if tiles[0].Granted {
		t.Fatal("a service without a project id must not open")
	}
}

func TestCatalogEntriesAreComplete(t *testing.T) {
	known := map[string]bool{}
	for _, c := range categories {
		known[c.Key] = true
	}
	for _, service := range catalog {
		if service.Name == "" || service.Url == "" || !known[service.Category] {
			t.Errorf("incomplete catalog entry %+v", service)
		}
	}
}
