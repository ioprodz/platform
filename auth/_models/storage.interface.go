package auth_models

import "ioprodz/common/policies"

type AccountRepository struct {
	policies.BaseRepository[Account]
}
