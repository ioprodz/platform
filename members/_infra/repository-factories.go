package members_infra

import (
	"ioprodz/common/config"
	members_models "ioprodz/members/_models"
)

func CreateMembersRepository() members_models.MembersRepository {
	isTest := config.Load().ENVIRONMENT == "test"

	if isTest {
		return CreateMemoryMemberRepo()
	} else {
		return CreateMongoMemberRepo()
	}
}
