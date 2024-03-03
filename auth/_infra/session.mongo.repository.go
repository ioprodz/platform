package auth_infra

import (
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/db"
	"ioprodz/common/policies"
)

type SessionMongoRepository struct {
	base db.BaseMongoRepository[auth_models.Session]
}

func (b *SessionMongoRepository) Create(entity auth_models.Session) error {
	return b.base.Create(entity)
}

func (b *SessionMongoRepository) List() ([]auth_models.Session, error) {
	return b.base.List()
}

func (b *SessionMongoRepository) Get(id string) (auth_models.Session, error) {
	return b.base.Get(id)
}

func (b *SessionMongoRepository) Update(entity auth_models.Session) error {
	return b.base.Update(entity)
}

func (b *SessionMongoRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (b *SessionMongoRepository) GetByHash(hash string) (auth_models.Session, error) {

	list, _ := b.List()
	for _, session := range list {
		if session.Hash == hash {
			return session, nil
		}
	}

	return auth_models.Session{}, &policies.StorageError{Message: "Session not found by hash '" + hash + "'"}
}

func (b *SessionMongoRepository) GetByAccountId(accountId string) []auth_models.Session {
	list, _ := b.List()

	accountSessions := []auth_models.Session{}
	for _, session := range list {
		if session.AccountId == accountId {
			accountSessions = append(accountSessions, session)
		}
	}
	return accountSessions
}

func CreateMongoSessionRepo() *SessionMongoRepository {
	repo := &SessionMongoRepository{base: *db.CreateMongoRepo[auth_models.Session]("sessions")}
	return repo
}
