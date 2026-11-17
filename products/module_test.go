package products

import (
	"context"
	"ioprodz/common/i18n"
	"ioprodz/common/middlewares"
	"ioprodz/common/policies"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// Never on ioprodz.com (claudeshell docs/ioprodz-inventory.md): client work,
// proofs of concept, do-not-mention demos and the home lab.
var neverNamed = []string{
	"Trocaro", "Liberio", "TOntoN", "Tonton", "Zadini", "GSM Prod", "gsmprod", "Radiant",
	"missionshams", "SGS", "bison", "devexpertise",
	"ts-classdiagram", "smsd", "cvadz", "cards-against",
	"claudeshell", "himarr", "catacomb", "catcam", "lasvegas",
}

func setup(t *testing.T) *mux.Router {
	t.Helper()
	wd, _ := os.Getwd()
	if err := os.Chdir(".."); err != nil { // templates, locales and llms.txt are read from the repo root
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(wd) })
	t.Setenv("APP_ENV", "production") // config.Load reads the environment, no .env
	t.Setenv("BASE_URL", "https://ioprodz.com")
	if err := i18n.Load(); err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	ConfigureModule(router)
	return router
}

func get(t *testing.T, router *mux.Router, path string) string {
	t.Helper()
	req := httptest.NewRequest("GET", path, nil)
	req = req.WithContext(context.WithValue(req.Context(), policies.CurrentUserCtxKey, policies.CurrentUser{}))
	w := httptest.NewRecorder()
	middlewares.LocaleResolver(router).ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET %s: %d", path, w.Code)
	}
	return w.Body.String()
}

func TestProductsShowXdrivenFirstAsTheFlagshipInEachLanguage(t *testing.T) {
	router := setup(t)
	for path, want := range map[string][]string{
		"/products":    {"Our software factory.", "Flagship"},
		"/fr/products": {"Notre usine logicielle.", "Produit phare"},
		"/ar/products": {"مصنعنا البرمجي.", "المنتج الرئيسي"},
	} {
		body := get(t, router, path)
		for _, s := range append(want, `href="https://www.xdriven.dev"`) {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
		noRawKey(t, path, body)
	}
	if Products[0].Key != "xdriven" {
		t.Errorf("xdriven must come first, got %s", Products[0].Key)
	}
}

func TestProductsListXdrivenChatInEachLanguage(t *testing.T) {
	router := setup(t)
	for path, want := range map[string]string{
		"/products":    "Your mail as chats:",
		"/fr/products": "Vos e-mails en conversations :",
		"/ar/products": "بريدك الإلكتروني على شكل محادثات:",
	} {
		body := get(t, router, path)
		for _, s := range []string{want, "xDriven chat", `href="https://chat.ioprodz.com"`} {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
	}
}

func TestProductsListMuzioInEachLanguage(t *testing.T) {
	router := setup(t)
	for path, want := range map[string]string{
		"/products":    "A social music app:",
		"/fr/products": "Une application musicale sociale :",
		"/ar/products": "تطبيق موسيقى اجتماعي:",
	} {
		body := get(t, router, path)
		for _, s := range []string{want, "Muzio", `href="https://muzio.ioprodz.com"`} {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
	}
}

func TestProductsListMissionFitInEachLanguage(t *testing.T) {
	router := setup(t)
	for path, want := range map[string]string{
		"/products":    "IT freelance missions, driven from your AI:",
		"/fr/products": "Les missions freelance IT, pilotées depuis votre IA :",
		"/ar/products": "مهام العمل الحر في تقنية المعلومات، من مساعدك الذكي:",
	} {
		body := get(t, router, path)
		for _, s := range []string{want, "MissionFit", `href="https://missionfit.app"`} {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
	}
}

func TestProductsListPolyseeInEachLanguage(t *testing.T) {
	router := setup(t)
	for path, want := range map[string]string{
		"/products":    "One workspace for software delivery:",
		"/fr/products": "Un espace de travail unique pour la livraison logicielle :",
		"/ar/products": "مساحة عمل واحدة لتسليم البرمجيات:",
	} {
		body := get(t, router, path)
		for _, s := range []string{want, "Polysee", `href="https://polysee.dev"`} {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
	}
}

func TestProductsListBledimmoInEachLanguage(t *testing.T) {
	router := setup(t)
	for path, want := range map[string]string{
		"/products":    "The Algerian real-estate network:",
		"/fr/products": "immobilier algérien : les prix réels du marché",
		"/ar/products": "شبكة العقار الجزائري:",
	} {
		body := get(t, router, path)
		for _, s := range []string{want, "Bledimmo", `href="https://bled.immo"`} {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
	}
}

func TestProductsListDoctoCabInEachLanguage(t *testing.T) {
	router := setup(t)
	for path, want := range map[string]string{
		"/products":    "Patient management for medical practices:",
		"/fr/products": "La gestion des patients pour les cabinets médicaux :",
		"/ar/products": "إدارة المرضى في العيادات الطبية:",
	} {
		body := get(t, router, path)
		for _, s := range []string{want, "DoctoCab", `href="https://doctocab.ioprodz.com"`} {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
	}
}

func TestProductsListMonaraInEachLanguage(t *testing.T) {
	router := setup(t)
	for path, want := range map[string]string{
		"/products":    "An AI cofounder for founders:",
		"/fr/products": "Un cofondateur IA pour les fondateurs :",
		"/ar/products": "شريك مؤسس بالذكاء الاصطناعي لروّاد الأعمال:",
	} {
		body := get(t, router, path)
		for _, s := range []string{want, "Monara", `href="https://monara.ioprodz.com"`} {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
	}
}

func TestLabListsOnlyDemosAndLabInEachLanguage(t *testing.T) {
	router := setup(t)
	for _, path := range []string{"/lab", "/fr/lab", "/ar/lab"} {
		body := get(t, router, path)
		for _, s := range []string{"SpecJest", "cloviz", "localvoice", "https://github.com/ioprodz/specjest", "https://github.com/ioprodz/cloviz"} {
			if !strings.Contains(body, s) {
				t.Errorf("%s: expected %q", path, s)
			}
		}
		noRawKey(t, path, body)
	}
}

var rawKey = regexp.MustCompile(`(products|common):[a-z][A-Za-z.]*`)

// A missing translation renders its key ("products:items.x.body") instead of the copy.
func noRawKey(t *testing.T, path, body string) {
	t.Helper()
	if key := rawKey.FindString(body); key != "" {
		t.Errorf("%s: a raw translation key shows: %q", path, key)
	}
}

func TestNoClientPocOrHomeLabNameOnThePagesOrInLlmsTxt(t *testing.T) {
	router := setup(t)
	pages := map[string]string{}
	for _, path := range []string{"/products", "/fr/products", "/ar/products", "/lab", "/fr/lab", "/ar/lab"} {
		pages[path] = get(t, router, path)
	}
	llms, err := os.ReadFile("llms.txt")
	if err != nil {
		t.Fatal(err)
	}
	pages["llms.txt"] = string(llms)
	for where, body := range pages {
		for _, name := range neverNamed {
			if strings.Contains(body, name) {
				t.Errorf("%s names %q", where, name)
			}
		}
	}
}
