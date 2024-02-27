package qna_infra

import (
	"ioprodz/common/db"
	"ioprodz/common/policies"
	qna_models "ioprodz/qna/_models"
)

type QNAMongoRepository struct {
	base policies.BaseRepository[qna_models.QNA]
}

func (repo *QNAMongoRepository) Create(qna qna_models.QNA) error {
	return repo.base.Create(qna)
}

func (repo *QNAMongoRepository) List() ([]qna_models.QNA, error) {
	return repo.base.List()
}

func (repo *QNAMongoRepository) Get(id string) (qna_models.QNA, error) {
	return repo.base.Get(id)
}

func (repo *QNAMongoRepository) Update(entity qna_models.QNA) error {
	return repo.base.Update(entity)
}

func (repo *QNAMongoRepository) Delete(id string) error {
	return repo.base.Delete(id)
}

func CreateMongoQNARepo() *QNAMongoRepository {
	repo := &QNAMongoRepository{base: db.CreateMongoRepo[qna_models.QNA]("qna")}

	return repo
}
