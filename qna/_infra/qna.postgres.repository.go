package qna_infra

import (
	"ioprodz/common/db"
	"ioprodz/common/policies"
	qna_models "ioprodz/qna/_models"
)

type QNAPostgresRepository struct {
	base policies.BaseRepository[qna_models.QNA]
}

func (repo *QNAPostgresRepository) Create(qna qna_models.QNA) error {
	return repo.base.Create(qna)
}

func (repo *QNAPostgresRepository) List() ([]qna_models.QNA, error) {
	return repo.base.List()
}

func (repo *QNAPostgresRepository) Get(id string) (qna_models.QNA, error) {
	return repo.base.Get(id)
}

func (repo *QNAPostgresRepository) Update(entity qna_models.QNA) error {
	return repo.base.Update(entity)
}

func (repo *QNAPostgresRepository) Delete(id string) error {
	return repo.base.Delete(id)
}

func CreatePostgresQNARepo() *QNAPostgresRepository {
	repo := &QNAPostgresRepository{base: db.CreatePostgresRepo[qna_models.QNA]("qna")}

	return repo
}
