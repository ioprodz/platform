package cv_infra

import (
	"ioprodz/common/config"
	cv_models "ioprodz/cv/_models"
)

func CreateCVRepository() cv_models.CVRepository {
	isTest := config.Load().ENVIRONMENT == "test"

	if isTest {
		return CreateMemoryCVRepo()
	} else {
		return CreateMongoCVRepo()
	}
}
