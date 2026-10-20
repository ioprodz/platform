package auth_infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
)

type IoprodzKey struct {
	Id   string
	Name string
}

// IoprodzSignInMethods is what an ioprodz account can sign in with.
type IoprodzSignInMethods struct {
	Password      bool
	Passkeys      []IoprodzKey
	SecurityKeys  []IoprodzKey
	Authenticator bool
}

const factorReady = "AUTH_FACTOR_STATE_READY"

// FetchIoprodzSignInMethods reads the user's own sign-in methods from Zitadel
// with the user's own access token.
func FetchIoprodzSignInMethods(issuer string, token IoprodzToken) (IoprodzSignInMethods, error) {
	methods := IoprodzSignInMethods{Passkeys: []IoprodzKey{}, SecurityKeys: []IoprodzKey{}}

	var types struct {
		AuthMethodTypes []string `json:"authMethodTypes"`
	}
	if err := ioprodzCall("GET", issuer+"/v2/users/"+url.PathEscape(token.UserId)+"/authentication_methods", token.AccessToken, nil, &types); err != nil {
		return methods, err
	}
	methods.Password = slices.Contains(types.AuthMethodTypes, "AUTHENTICATION_METHOD_TYPE_PASSWORD")

	var passkeys struct {
		Result []struct {
			Id    string `json:"id"`
			State string `json:"state"`
			Name  string `json:"name"`
		} `json:"result"`
	}
	if err := ioprodzCall("POST", issuer+"/auth/v1/users/me/passwordless/_search", token.AccessToken, struct{}{}, &passkeys); err != nil {
		return methods, err
	}
	for _, passkey := range passkeys.Result {
		if passkey.State == factorReady {
			methods.Passkeys = append(methods.Passkeys, IoprodzKey{Id: passkey.Id, Name: passkey.Name})
		}
	}

	var factors struct {
		Result []struct {
			State string    `json:"state"`
			Otp   *struct{} `json:"otp"`
			U2f   *struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"u2f"`
		} `json:"result"`
	}
	if err := ioprodzCall("POST", issuer+"/auth/v1/users/me/auth_factors/_search", token.AccessToken, struct{}{}, &factors); err != nil {
		return methods, err
	}
	for _, factor := range factors.Result {
		if factor.State != factorReady {
			continue
		}
		if factor.Otp != nil {
			methods.Authenticator = true
		}
		if factor.U2f != nil {
			methods.SecurityKeys = append(methods.SecurityKeys, IoprodzKey{Id: factor.U2f.Id, Name: factor.U2f.Name})
		}
	}
	return methods, nil
}

// RemoveIoprodzSignInMethod removes one of the user's own sign-in methods:
// kind "passkey" or "key" with its id, or "authenticator".
func RemoveIoprodzSignInMethod(issuer string, accessToken string, kind string, id string) error {
	paths := map[string]string{
		"passkey":       "/auth/v1/users/me/passwordless/" + url.PathEscape(id),
		"key":           "/auth/v1/users/me/auth_factors/u2f/" + url.PathEscape(id),
		"authenticator": "/auth/v1/users/me/auth_factors/otp",
	}
	path, ok := paths[kind]
	if !ok {
		return fmt.Errorf("unknown sign-in method %q", kind)
	}
	return ioprodzCall("DELETE", issuer+path, accessToken, nil, nil)
}

// ioprodzCall sends one request to Zitadel's API. Errors carry Zitadel's own
// message, which the page can show.
func ioprodzCall(method string, endpoint string, accessToken string, body interface{}, result interface{}) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, endpoint, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	// auth.ioprodz.com sits behind Cloudflare, which rejects some default user agents
	req.Header.Set("User-Agent", "ioprodz-platform/1.0")

	res, err := ioprodzAuthClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var failure struct {
			Message string `json:"message"`
		}
		json.NewDecoder(res.Body).Decode(&failure)
		if failure.Message == "" {
			failure.Message = fmt.Sprintf("status %d", res.StatusCode)
		}
		return fmt.Errorf("%s", failure.Message)
	}
	if result == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(result)
}
