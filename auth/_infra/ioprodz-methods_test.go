package auth_infra

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func fakeZitadel(t *testing.T, calls *[]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*calls = append(*calls, r.Method+" "+r.URL.Path)
		if r.Header.Get("Authorization") != "Bearer user-token" {
			t.Errorf("expected the user's own token, got %q", r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/v2/users/u1/authentication_methods":
			w.Write([]byte(`{"authMethodTypes":["AUTHENTICATION_METHOD_TYPE_PASSWORD","AUTHENTICATION_METHOD_TYPE_PASSKEY"]}`))
		case "/auth/v1/users/me/passwordless/_search":
			w.Write([]byte(`{"result":[{"id":"p1","state":"AUTH_FACTOR_STATE_READY","name":"MacBook"},{"id":"p2","state":"AUTH_FACTOR_STATE_NOT_READY","name":"half"}]}`))
		case "/auth/v1/users/me/auth_factors/_search":
			w.Write([]byte(`{"result":[{"state":"AUTH_FACTOR_STATE_READY","otp":{}},{"state":"AUTH_FACTOR_STATE_READY","u2f":{"id":"k1","name":"YubiKey"}}]}`))
		case "/auth/v1/users/me/auth_factors/otp":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code":9,"message":"At least one second factor is required"}`))
		default:
			w.Write([]byte(`{}`))
		}
	}))
}

func TestFetchIoprodzSignInMethodsListsReadyMethodsOnly(t *testing.T) {
	calls := []string{}
	zitadel := fakeZitadel(t, &calls)
	defer zitadel.Close()

	methods, err := FetchIoprodzSignInMethods(zitadel.URL, IoprodzToken{AccessToken: "user-token", UserId: "u1"})
	if err != nil {
		t.Fatal(err)
	}
	if !methods.Password || !methods.Authenticator {
		t.Errorf("expected password and authenticator, got %+v", methods)
	}
	if len(methods.Passkeys) != 1 || methods.Passkeys[0].Name != "MacBook" {
		t.Errorf("expected only the ready passkey, got %+v", methods.Passkeys)
	}
	if len(methods.SecurityKeys) != 1 || methods.SecurityKeys[0].Id != "k1" {
		t.Errorf("expected the security key, got %+v", methods.SecurityKeys)
	}
}

func TestRemoveIoprodzSignInMethodCallsTheRightEndpoint(t *testing.T) {
	calls := []string{}
	zitadel := fakeZitadel(t, &calls)
	defer zitadel.Close()

	RemoveIoprodzSignInMethod(zitadel.URL, "user-token", "passkey", "p1")
	RemoveIoprodzSignInMethod(zitadel.URL, "user-token", "key", "k1")
	want := []string{"DELETE /auth/v1/users/me/passwordless/p1", "DELETE /auth/v1/users/me/auth_factors/u2f/k1"}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, calls)
		}
	}

	err := RemoveIoprodzSignInMethod(zitadel.URL, "user-token", "authenticator", "otp")
	if err == nil || err.Error() != "At least one second factor is required" {
		t.Errorf("expected Zitadel's message, got %v", err)
	}
	if err := RemoveIoprodzSignInMethod(zitadel.URL, "user-token", "password", ""); err == nil {
		t.Error("an unknown kind must be refused")
	}
}
