package blog_models

import "ioprodz/common/policies"

type BlogRepository interface {
	policies.Repository[Blog]
	Update(blog Blog) error
}
