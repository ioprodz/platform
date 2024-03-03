package qna_infra

import (
	"ioprodz/common/config"
	qna_models "ioprodz/qna/_models"
)

func CreateAnswersRepository() qna_models.AnswersRepository {
	isTest := config.Load().ENVIRONMENT == "test"

	if isTest {
		return CreateMemoryAnswerRepo()
	} else {
		return CreateMongoAnswerRepo()
	}
}

func CreateQNARepository() qna_models.QNARepository {
	isTest := config.Load().ENVIRONMENT == "test"

	if isTest {
		return CreateMemoryQNARepo()
	} else {
		return CreateMongoQNARepo()
	}
}
