package auth_infra

import (
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/db"
	"ioprodz/common/policies"
)

type AccountPostgresRepository struct {
	base db.BasePostgresRepository[auth_models.Account]
}

func (b *AccountPostgresRepository) Create(entity auth_models.Account) error {
	return b.base.Create(entity)
}

func (b *AccountPostgresRepository) List() ([]auth_models.Account, error) {
	return b.base.List()
}

func (b *AccountPostgresRepository) Get(id string) (auth_models.Account, error) {
	return b.base.Get(id)
}

func (b *AccountPostgresRepository) Update(entity auth_models.Account) error {
	return b.base.Update(entity)
}

func (b *AccountPostgresRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (b *AccountPostgresRepository) GetByProviderId(provider string, providerUserId string) (auth_models.Account, error) {

	list, _ := b.List()
	for _, acc := range list {
		if acc.Provider == provider && acc.ProviderUserId == providerUserId {
			return acc, nil
		}
	}

	return auth_models.Account{}, &policies.StorageError{Message: "Account '" + providerUserId + "' not found by for provider '" + provider + "'"}
}

func CreatePostgresAccountRepo() *AccountPostgresRepository {
	repo := &AccountPostgresRepository{base: *db.CreatePostgresRepo[auth_models.Account]("user_accounts")}
	return repo
}
