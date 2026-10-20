package auth_infra

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// skips the .env file, which lives at the repo root
	os.Setenv("APP_ENV", "production")
	os.Exit(m.Run())
}

func TestIoprodzTokenCookieIsEncryptedAndReadBack(t *testing.T) {
	token := IoprodzToken{AccessToken: "secret-access-token", UserId: "u1", LoginName: "ana@ioprodz.com", ExpiresAt: time.Now().Add(time.Hour)}
	recorder := httptest.NewRecorder()
	SetIoprodzToken(recorder, httptest.NewRequest("GET", "/auth/ioprodz/callback", nil), token)

	setCookie := recorder.Header().Get("Set-Cookie")
	if !strings.Contains(setCookie, "Path=/security") || !strings.Contains(setCookie, "HttpOnly") {
		t.Fatalf("cookie must be HttpOnly and limited to /security: %s", setCookie)
	}
	if strings.Contains(setCookie, "secret-access-token") {
		t.Fatal("the token must not be readable in the cookie")
	}

	request := httptest.NewRequest("GET", "/security", nil)
	request.Header.Set("Cookie", strings.Split(setCookie, ";")[0])
	got := GetIoprodzToken(request)
	if !got.Valid() || got.AccessToken != token.AccessToken || got.UserId != "u1" || got.LoginName != "ana@ioprodz.com" {
		t.Fatalf("token not read back: %+v", got)
	}
}

func TestExpiredIoprodzTokenIsNotKept(t *testing.T) {
	recorder := httptest.NewRecorder()
	SetIoprodzToken(recorder, httptest.NewRequest("GET", "/", nil), IoprodzToken{AccessToken: "x", ExpiresAt: time.Now().Add(-time.Minute)})
	if recorder.Header().Get("Set-Cookie") != "" {
		t.Fatal("an expired token must not be stored")
	}
	if (IoprodzToken{}).Valid() {
		t.Fatal("no token is not a valid token")
	}
}
