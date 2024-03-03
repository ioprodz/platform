package auth_infra

import (
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/db"
	"ioprodz/common/policies"
)

type AccountMongoRepository struct {
	base db.BaseMongoRepository[auth_models.Account]
}

func (b *AccountMongoRepository) Create(entity auth_models.Account) error {
	return b.base.Create(entity)
}

func (b *AccountMongoRepository) List() ([]auth_models.Account, error) {
	return b.base.List()
}

func (b *AccountMongoRepository) Get(id string) (auth_models.Account, error) {
	return b.base.Get(id)
}

func (b *AccountMongoRepository) Update(entity auth_models.Account) error {
	return b.base.Update(entity)
}

func (b *AccountMongoRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (b *AccountMongoRepository) GetByProviderId(provider string, providerUserId string) (auth_models.Account, error) {

	list, _ := b.List()
	for _, acc := range list {
		if acc.Provider == provider && acc.ProviderUserId == providerUserId {
			return acc, nil
		}
	}

	return auth_models.Account{}, &policies.StorageError{Message: "Account '" + providerUserId + "' not found by for provider '" + provider + "'"}
}

func CreateMongoAccountRepo() *AccountMongoRepository {
	repo := &AccountMongoRepository{base: *db.CreateMongoRepo[auth_models.Account]("user_accounts")}
	return repo
}
