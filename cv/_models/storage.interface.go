package cv_models

import "ioprodz/common/policies"

type CVRepository interface {
	policies.BaseRepository[CV]
}
