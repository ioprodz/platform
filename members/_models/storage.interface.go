package members_models

import "ioprodz/common/policies"

type MembersRepository interface {
	policies.BaseRepository[Member]
}

type UserProjectRepository interface {
	policies.BaseRepository[UserProject]
	GetByUserId(userId string) ([]UserProject, error)
}
