package qna_infra

import (
	"ioprodz/common/db"
	qna_models "ioprodz/qna/_models"
)

type answerPostgresRepository struct {
	base db.BasePostgresRepository[qna_models.Answers]
}

func (repo *answerPostgresRepository) Create(qna qna_models.Answers) error {
	return repo.base.Create(qna)
}

func (repo *answerPostgresRepository) List() ([]qna_models.Answers, error) {
	return repo.base.List()
}

func (repo *answerPostgresRepository) Get(id string) (qna_models.Answers, error) {
	return repo.base.Get(id)
}

func (repo *answerPostgresRepository) Update(item qna_models.Answers) error {
	return repo.base.Update(item)
}

func (repo *answerPostgresRepository) Delete(id string) error {
	return repo.base.Delete(id)
}

func CreatePostgresAnswerRepo() *answerPostgresRepository {
	repo := &answerPostgresRepository{base: *db.CreatePostgresRepo[qna_models.Answers]("answers")}

	return repo
}
