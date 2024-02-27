package qna_infra

import (
	"ioprodz/common/db"
	qna_models "ioprodz/qna/_models"
)

type answerMongoRepository struct {
	base db.BaseMongoRepository[qna_models.Answers]
}

func (repo *answerMongoRepository) Create(qna qna_models.Answers) error {
	return repo.base.Create(qna)
}

func (repo *answerMongoRepository) List() ([]qna_models.Answers, error) {
	return repo.base.List()
}

func (repo *answerMongoRepository) Get(id string) (qna_models.Answers, error) {
	return repo.base.Get(id)
}

func (repo *answerMongoRepository) Update(item qna_models.Answers) error {
	return repo.base.Update(item)
}

func (repo *answerMongoRepository) Delete(id string) error {
	return repo.base.Delete(id)
}

func CreateMongoAnswerRepo() *answerMongoRepository {
	repo := &answerMongoRepository{base: *db.CreateMongoRepo[qna_models.Answers]("answers")}

	return repo
}
