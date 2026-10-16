package auth_models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id             string
	AccountId      string
	Email          string
	AvatarUrl      string
	Name           string
	UaString       string
	Hash           string
	FirstCreatedAt string
	LastUsedAt     string
	Provider       string
	// Zitadel project ids the user holds a role in, read at login
	ProjectIds []string
}

func (s Session) GetId() string {
	return s.Id
}

func (s *Session) SetLastUsedNow() {
	s.LastUsedAt = time.Now().Format(time.RFC3339)
}

func SessionFromJSON(jsonData []byte) Session {
	var session Session
	if err := json.Unmarshal(jsonData, &session); err != nil {
		panic("unable to parse session json")
	}
	return session
}

func NewSession(accountId string, uaString string, hash string, avatarUrl string, name string, email string, provider string) Session {
	return Session{
		Provider:       provider,
		Id:             uuid.NewString(),
		AccountId:      accountId,
		UaString:       uaString,
		Hash:           hash,
		AvatarUrl:      avatarUrl,
		Name:           name,
		Email:          email,
		LastUsedAt:     time.Now().Format(time.RFC3339),
		FirstCreatedAt: time.Now().Format(time.RFC3339)}
}
