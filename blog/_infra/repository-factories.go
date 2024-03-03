package blog_infra

import (
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/config"
)

func CreateBlogRepository() blog_models.BlogRepository {
	isTest := config.Load().ENVIRONMENT == "test"

	if isTest {
		return CreateMemoryBlogRepo()
	} else {
		return CreateMongoBlogRepo()
	}
}
