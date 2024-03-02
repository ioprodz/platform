package auth_models

import "encoding/json"

type Account struct {
	Id             string
	Email          string
	Provider       string
	ProviderUserId string
}

func (a Account) GetId() string {
	return a.Id
}

func AccountFromJSON(jsonData []byte) Account {
	var account Account
	if err := json.Unmarshal(jsonData, &account); err != nil {
		panic("unable to parse account json")
	}
	return account
}
