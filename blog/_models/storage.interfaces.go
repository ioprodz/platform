package blog_models

import "ioprodz/common/policies"

type BlogRepository interface {
	policies.BaseRepository[Blog]
}
