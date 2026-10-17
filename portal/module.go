package portal

import (
	"ioprodz/common/policies"
	"ioprodz/common/ui"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

// Tile is a catalog service as one user sees it.
type Tile struct {
	Service
	CategoryName string
	Initial      string
	Granted      bool
}

type CategoryCount struct {
	Category
	Count int
}

// tilesFor lists the whole catalog, the user's services first, each group in
// catalog order.
func tilesFor(services []Service, projectIds []string) []Tile {
	names := map[string]string{}
	for _, c := range categories {
		names[c.Key] = c.Name
	}
	tiles := []Tile{}
	for _, service := range services {
		tiles = append(tiles, Tile{
			Service:      service,
			CategoryName: names[service.Category],
			Initial:      string([]rune(service.Name)[:1]),
			Granted:      service.ProjectId != "" && slices.Contains(projectIds, service.ProjectId),
		})
	}
	slices.SortStableFunc(tiles, func(a, b Tile) int {
		if a.Granted == b.Granted {
			return 0
		}
		if a.Granted {
			return -1
		}
		return 1
	})
	return tiles
}

// categoryCounts lists the categories that have at least one service.
func categoryCounts(services []Service) []CategoryCount {
	counts := []CategoryCount{}
	for _, c := range categories {
		n := 0
		for _, service := range services {
			if service.Category == c.Key {
				n++
			}
		}
		if n > 0 {
			counts = append(counts, CategoryCount{Category: c, Count: n})
		}
	}
	return counts
}

func ConfigureModule(router *mux.Router) *mux.Router {
	router.HandleFunc("/portal", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(policies.CurrentUserCtxKey).(policies.CurrentUser)
		tiles := tilesFor(catalog, user.ProjectIds)
		granted := 0
		for _, tile := range tiles {
			if tile.Granted {
				granted++
			}
		}
		meta := ui.PageMeta{Title: "Portal", Path: "/portal", OGType: "website", NoIndex: true}
		ui.RenderPageWithMeta(w, r, "portal/index", map[string]interface{}{
			"User":              user,
			"HasIoprodzAccount": user.HasIoprodzAccount(),
			"Tiles":             tiles,
			"Categories":        categoryCounts(catalog),
			"AllFilter":         CategoryCount{Category: Category{Key: "all", Name: "All apps"}, Count: len(tiles)},
			"MineFilter":        CategoryCount{Category: Category{Key: "mine", Name: "My apps"}, Count: granted},
			"Total":             len(tiles),
			"GrantedCount":      granted,
		}, meta)
	}).Methods("GET")
	return router
}
