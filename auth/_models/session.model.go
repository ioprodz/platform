package auth_models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id             string
	AccountId      string
	UaString       string
	Hash           string
	FirstCreatedAt string
	LastUsedAt     string
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

func NewSession(accountId string, uaString string, hash string) Session {
	return Session{
		Id:             uuid.NewString(),
		AccountId:      accountId,
		UaString:       uaString,
		Hash:           hash,
		LastUsedAt:     time.Now().Format(time.RFC3339),
		FirstCreatedAt: time.Now().Format(time.RFC3339)}
}
