package auth_infra

import (
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/db"
	"ioprodz/common/policies"
)

type SessionMemoryRepository struct {
	base db.BaseMemoryRepository[auth_models.Session]
}

func (b *SessionMemoryRepository) Create(entity auth_models.Session) error {
	return b.base.Create(entity)
}

func (b *SessionMemoryRepository) List() ([]auth_models.Session, error) {
	return b.base.List()
}

func (b *SessionMemoryRepository) Get(id string) (auth_models.Session, error) {
	return b.base.Get(id)
}

func (b *SessionMemoryRepository) Update(entity auth_models.Session) error {
	return b.base.Update(entity)
}

func (b *SessionMemoryRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (b *SessionMemoryRepository) GetByHash(hash string) (auth_models.Session, error) {

	list, _ := b.List()
	for _, session := range list {
		if session.Hash == hash {
			return session, nil
		}
	}

	return auth_models.Session{}, &policies.StorageError{Message: "Session not found by hash '" + hash + "'"}
}

func (b *SessionMemoryRepository) GetByAccountId(accountId string) []auth_models.Session {
	list, _ := b.List()

	accountSessions := []auth_models.Session{}
	for _, session := range list {
		if session.AccountId == accountId {
			accountSessions = append(accountSessions, session)
		}
	}
	return accountSessions
}

func CreateMemorySessionRepo() *SessionMemoryRepository {

	repo := &SessionMemoryRepository{base: *db.CreateMemoryRepo[auth_models.Session]()}
	repo.seed()
	return repo
}

func (r *SessionMemoryRepository) seed() {

	r.Create(auth_models.SessionFromJSON([]byte(`{
		"id":"session-001",
		"hash":"hash",
		"lastUsed":"2024-02-29T23:19:52.345"
	}`)))
}
