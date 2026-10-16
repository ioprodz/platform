package auth_infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var ioprodzAuthClient = &http.Client{Timeout: 10 * time.Second}

// FetchIoprodzProjectIds returns the ids of the Zitadel projects the user
// holds at least one role in, using the user's own access token
// (needs scope urn:zitadel:iam:org:project:id:zitadel:aud).
func FetchIoprodzProjectIds(issuer string, accessToken string) ([]string, error) {
	body := bytes.NewBufferString(`{"query":{"limit":1000}}`)
	req, err := http.NewRequest("POST", issuer+"/auth/v1/usergrants/me/_search", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	// auth.ioprodz.com sits behind Cloudflare, which rejects some default user agents
	req.Header.Set("User-Agent", "ioprodz-platform/1.0")

	res, err := ioprodzAuthClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ioprodz account grants: status %d", res.StatusCode)
	}

	var payload struct {
		Result []struct {
			ProjectId string   `json:"projectId"`
			RoleKeys  []string `json:"roleKeys"`
		} `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	projectIds := []string{}
	for _, grant := range payload.Result {
		if len(grant.RoleKeys) > 0 {
			projectIds = append(projectIds, grant.ProjectId)
		}
	}
	return projectIds, nil
}
