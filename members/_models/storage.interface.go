package members_models

import "ioprodz/common/policies"

type MembersRepository interface {
	policies.BaseRepository[Member]
}
