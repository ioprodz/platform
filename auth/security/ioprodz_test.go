package auth_security

import (
	"strings"
	"testing"
)

func TestAddLinksOpenZitadelLoginPagesForTheUser(t *testing.T) {
	links := addLinks("https://auth.example", "ana+x@ioprodz.com")

	want := map[string]string{
		links.Passkey:       "/ui/v2/login/passkey/set?",
		links.SecurityKey:   "/ui/v2/login/u2f/set?",
		links.Authenticator: "/ui/v2/login/otp/time-based/set?",
		links.Password:      "/ui/v2/login/password/change?",
	}
	for link, path := range want {
		if !strings.HasPrefix(link, "https://auth.example"+path) {
			t.Errorf("expected %s on %s", path, link)
		}
		if !strings.HasSuffix(link, "loginName=ana%2Bx%40ioprodz.com") {
			t.Errorf("login name not passed (escaped) in %s", link)
		}
	}
}
