// Package products serves ioprodz.com's Products and Lab pages (#1287), built from
// claudeshell's docs/ioprodz-inventory.md: internal products only, never client
// work, proofs of concept or the home lab. Each entry's copy is approved on its
// own ticket and lives in locales/<lang>/products.json.
package products

import (
	"ioprodz/common/ui"
	"net/http"

	"github.com/gorilla/mux"
)

// Entry is one product or lab item; Key names its copy under products:items.<Key>.
type Entry struct {
	Key  string
	Name string
	URL  string // empty: no public link
	Tag  string // a translation key for its badge, e.g. products:tags.flagship
}

// Products, flagship first. The others are added by their own tickets once
// their en/fr/ar copy is approved.
var Products = []Entry{
	{Key: "xdriven", Name: "xdriven", URL: "https://www.xdriven.dev", Tag: "products:tags.flagship"},
	{Key: "chat", Name: "xDriven chat", URL: "https://chat.ioprodz.com"},
	{Key: "muzio", Name: "Muzio", URL: "https://muzio.ioprodz.com"},
	{Key: "missionfit", Name: "MissionFit", URL: "https://missionfit.app"},
	{Key: "polysee", Name: "Polysee", URL: "https://polysee.dev"},
	{Key: "bledimmo", Name: "Bledimmo", URL: "https://bled.immo"},
}

// Lab: open source and experiments, never live products.
var Lab = []Entry{
	{Key: "specjest", Name: "SpecJest", URL: "https://github.com/ioprodz/specjest", Tag: "products:tags.openSource"},
	{Key: "cloviz", Name: "cloviz", URL: "https://github.com/ioprodz/cloviz", Tag: "products:tags.openSource"},
	{Key: "localvoice", Name: "localvoice", Tag: "products:tags.lab"},
}

func ConfigureModule(router *mux.Router) {
	router.HandleFunc("/products", page("products", "/products", Products)).Methods("GET")
	router.HandleFunc("/lab", page("lab", "/lab", Lab)).Methods("GET")
}

func page(name, path string, entries []Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		meta := ui.PageMeta{
			TitleKey: "products:meta." + name + ".title",
			DescKey:  "products:meta." + name + ".description",
			Path:     path,
			OGType:   "website",
		}
		ui.RenderPageWithMeta(w, r, "products/list", map[string]interface{}{
			"Page":    name,
			"Entries": entries,
		}, meta)
	}
}
