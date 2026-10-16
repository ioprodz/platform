package auth_provider

import (
	"log"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

// ConfigureIoprodz registers the ioprodz account (Zitadel OIDC) as goth provider "ioprodz".
// Skipped when no client id is configured.
func ConfigureIoprodz(issuer string, clientId string, clientSecret string, baseUrl string) {
	if clientId == "" {
		return
	}
	provider, err := openidConnect.NewNamed("ioprodz", clientId, clientSecret,
		baseUrl+"/auth/ioprodz/callback",
		issuer+"/.well-known/openid-configuration",
		"openid", "profile", "email",
		// lets the user's token read their own grants (portal)
		"urn:zitadel:iam:org:project:id:zitadel:aud",
	)
	if err != nil {
		log.Printf("ioprodz login disabled: %v", err)
		return
	}
	// NewNamed registers "<name>-oidc"; the routes use /auth/ioprodz
	provider.SetName("ioprodz")
	goth.UseProviders(provider)
}
