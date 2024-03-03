package auth_infra

import (
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/db"
	"ioprodz/common/policies"
)

type AccountMemoryRepository struct {
	base db.BaseMemoryRepository[auth_models.Account]
}

func (b *AccountMemoryRepository) Create(entity auth_models.Account) error {
	return b.base.Create(entity)
}

func (b *AccountMemoryRepository) List() ([]auth_models.Account, error) {
	return b.base.List()
}

func (b *AccountMemoryRepository) Get(id string) (auth_models.Account, error) {
	return b.base.Get(id)
}

func (b *AccountMemoryRepository) Update(entity auth_models.Account) error {
	return b.base.Update(entity)
}

func (b *AccountMemoryRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (b *AccountMemoryRepository) GetByProviderId(provider string, providerUserId string) (auth_models.Account, error) {

	list, _ := b.List()
	for _, acc := range list {
		if acc.Provider == provider && acc.ProviderUserId == providerUserId {
			return acc, nil
		}
	}

	return auth_models.Account{}, &policies.StorageError{Message: "Account '" + providerUserId + "' not found by for provider '" + provider + "'"}
}

func CreateMemoryAccountRepo() *AccountMemoryRepository {

	repo := &AccountMemoryRepository{base: *db.CreateMemoryRepo[auth_models.Account]()}
	repo.seed()
	return repo
}

func (r *AccountMemoryRepository) seed() {

	r.Create(auth_models.AccountFromJSON([]byte(`{
		"id":"account-001",
		"email":"test@email.com",
		"provider":"github",
		"providerId":"123456"
	}`)))
}
