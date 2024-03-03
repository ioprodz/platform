package auth_models

import "ioprodz/common/policies"

type AccountRepository interface {
	policies.BaseRepository[Account]
	GetByProviderId(provider string, providerUserId string) (Account, error)
}

type SessionRepository interface {
	policies.BaseRepository[Session]
	GetByHash(hash string) (Session, error)
	GetByAccountId(accountId string) []Session
}
