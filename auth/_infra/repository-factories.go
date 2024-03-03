package auth_infra

import (
	auth_models "ioprodz/auth/_models"
	"ioprodz/common/config"
)

func CreateAccountRepository() auth_models.AccountRepository {
	isTest := config.Load().ENVIRONMENT == "test"

	if isTest {
		return CreateMemoryAccountRepo()
	} else {
		return CreateMongoAccountRepo()
	}
}

func CreateSessionRepository() auth_models.SessionRepository {
	isTest := config.Load().ENVIRONMENT == "test"

	if isTest {
		return CreateMemorySessionRepo()
	} else {
		return CreateMongoSessionRepo()
	}
}
