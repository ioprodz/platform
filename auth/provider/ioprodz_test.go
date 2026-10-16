package auth_provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/markbates/goth"
)

func TestIoprodzProviderIsRegisteredUnderItsRouteName(t *testing.T) {
	issuer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"issuer":                 "https://auth.example",
			"authorization_endpoint": "https://auth.example/oauth/v2/authorize",
			"token_endpoint":         "https://auth.example/oauth/v2/token",
			"userinfo_endpoint":      "https://auth.example/oidc/v1/userinfo",
		})
	}))
	defer issuer.Close()

	ConfigureIoprodz(issuer.URL, "client-id", "client-secret", "https://site.example")

	if _, err := goth.GetProvider("ioprodz"); err != nil {
		t.Fatalf("provider for /auth/ioprodz not registered: %v", err)
	}
}
